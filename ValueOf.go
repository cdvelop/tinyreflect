package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

type flag uintptr

const (
	flagKindWidth      = 5 // there are 27 kinds
	flagKindMask  flag = 1<<flagKindWidth - 1
	flagIndir     flag = 1 << 7
	flagStickyRO  flag = 1 << 5
	flagEmbedRO   flag = 1 << 6
	flagAddr      flag = 1 << 8
)

// TinyReflect - Minimal reflection library optimized for TinyGo/WebAssembly
// This package provides a thin API layer over tinystring's core reflection functionality
// ALL TYPE DETECTION AND CORE LOGIC IS DELEGATED TO TINYSTRING FOR MAXIMUM CODE REUSE

// Minimal reflectlite integration for TinyString JSON functionality
// This file contains essential reflection capabilities adapted from internal/reflectlite
// All functions and types are prefixed with 'ref' to avoid API pollution

// Import unified types from reflectlite/value.go - no more duplication
// Kind is now defined in tinystring/Kind.go as an anonymous struct for clean API

type Value struct {
	// typ_ holds the type of the value represented by a Value.
	// Access using the Typ method to avoid escape of v.
	typ_ *Type

	// Pointer-valued data or, if flagIndir is set, pointer to data.
	// Valid when either flagIndir is set or Typ.pointers() is true.
	ptr unsafe.Pointer

	// flag holds metadata about the value.
	// The lowest bits are flag bits:
	//	- flagStickyRO: obtained via unexported not embedded field, so read-only
	//	- flagEmbedRO: obtained via unexported embedded field, so read-only
	//	- flagIndir: val holds a pointer to the data
	//	- flagAddr: v.CanAddr is true (implies flagIndir)
	// Value cannot represent method values.
	// The next five bits give the Kind of the value.
	// This repeats Typ.Kind() except for method values.
	// The remaining 23+ bits give a method number for method values.
	// If flag.kind() != Func, code can assume that flagMethod is unset.
	// If ifaceIndir(Typ), code can assume that flagIndir is set.
	flag

	// A method value represents a curried method invocation
	// like r.Read for some receiver r. The Typ+val+flag bits describe
	// the receiver r, but the flag's Kind bits say Func (methods are
	// functions), and the top bits of the flag give the method number
	// in r's type's method table.
}

// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i. ValueOf(nil) returns the zero Value.
func ValueOf(i any) Value {
	if i == nil {
		return Value{}
	}
	return unpackEface(i)
}

// unpackEface converts the empty interface i to a Value.
func unpackEface(i any) Value {
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	// NOTE: don't read e.word until we know whether it is really a pointer or not.
	t := e.Type
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if t.IfaceIndir() {
		f |= flagIndir
	}
	return Value{t, e.Data, f}
}

// Elem returns the value that the pointer v points to.
// It panics if v's Kind is not Ptr.
// It returns the zero Value if v is a nil pointer.
func (v Value) Elem() (Value, error) {
	if v.kind() != K.Pointer {
		return Value{}, Err(ref, D.Value, D.NotOfType, D.Pointer)
	}
	if v.ptr == nil {
		return Value{}, nil
	}

	tt := (*PtrType)(unsafe.Pointer(v.typ()))
	typ := tt.Elem
	fl := v.flag&^flagKindMask | flag(typ.Kind())
	fl &^= flagStickyRO | flagEmbedRO
	fl |= flagIndir | flagAddr
	ptr := v.ptr
	return Value{typ, ptr, fl}, nil
}

func (v Value) NumField() (int, error) {
	if v.typ_ == nil {
		return 0, Err(ref, D.Value, D.Nil)
	}
	st := v.typ_.StructType()
	if st == nil {
		return 0, Err(ref, D.Numbers, D.Fields, D.NotOfType, D.Struct)
	}
	return len(st.Fields), nil
}

// Field returns the i'th field of the struct v.
// Returns an error if v is not a struct or i is out of range.
func (v Value) Field(i int) (Value, error) {
	if v.kind() != K.Struct {
		return Value{}, Err(ref, D.Value, D.NotOfType, D.Struct)
	}
	tt := (*StructType)(unsafe.Pointer(v.typ()))
	if uint(i) >= uint(len(tt.Fields)) {
		return Value{}, Err(ref, D.Value, D.Index, D.Out, D.Of, D.Range)
	}
	field := &tt.Fields[i]
	typ := field.Typ

	// Inherit permission bits from v, but clear flagEmbedRO.
	fl := v.flag&(flagStickyRO|flagIndir|flagAddr) | flag(typ.Kind())
	// Using an unexported field forces flagRO.
	if !field.Name.IsExported() {
		if field.Embedded() {
			fl |= flagEmbedRO
		} else {
			fl |= flagStickyRO
		}
	}
	// Either flagIndir is set and v.ptr points at struct,
	// or flagIndir is not set and v.ptr is the actual struct data.
	// In the former case, we want v.ptr + offset.
	// In the latter case, we must have field.offset = 0,
	// so v.ptr + field.offset is still the correct address.
	ptr := add(v.ptr, field.Off, "same as non-reflect &v.field")
	return Value{typ, ptr, fl}, nil
}

// Type returns v's type.
func (v Value) Type() *Type {
	return v.typ()
}

// typ returns the *abi.Type stored in the Value. This method is fast,
// but it doesn't always return the correct type for the Value.
// See abiType and Type, which do return the correct type.
func (v Value) typ() *Type {
	// Types are either static (for compiler-created types) or
	// heap-allocated but always reachable (for reflection-created
	// types, held in the central map). So there is no need to
	// escape types. noescape here help avoid unnecessary escape
	// of v.
	return (*Type)(NoEscape(unsafe.Pointer(v.typ_)))
}

// add returns p+x.
//
// The whySafe string is ignored, so that the function still inlines
// as efficiently as p+x, but all call sites should use the string to
// record why the addition is safe, which is to say why the addition
// does not cause x to advance to the very end of p's allocation
// and therefore point incorrectly at the next block in memory.
//
// add should be an internal detail (and is trivially copyable),
// but widely used packages access it using linkname.
// Notable members of the hall of shame include:
//   - github.com/pinpoint-apm/pinpoint-go-agent
//   - github.com/vmware/govmomi
//
// Do not remove or change the type signature.
// See go.dev/issue/67401.
//
//go:linkname add
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func (f flag) kind() Kind {
	return Kind(f & flagKindMask)
}

// NoEscape hides the pointer p from escape analysis, preventing it
// from escaping to the heap. It compiles down to nothing.
//
// WARNING: This is very subtle to use correctly. The caller must
// ensure that it's truly safe for p to not escape to the heap by
// maintaining runtime pointer invariants (for example, that globals
// and the heap may not generally point into a stack).
//
//go:nosplit
//go:nocheckptr
func NoEscape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x)
}
