package tinyreflect

import (
	"sync/atomic"
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

type flag uintptr

const (
	flagKindWidth      = 5 // there are 27 kinds
	flagKindMask  flag = 1<<flagKindWidth - 1
	flagStickyRO  flag = 1 << 5
	flagEmbedRO   flag = 1 << 6
	flagIndir     flag = 1 << 7
	flagAddr      flag = 1 << 8
	flagRO        flag = flagStickyRO | flagEmbedRO
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
	tr *TinyReflect
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

// Elem returns the value that the pointer v points to.
// It panics if v's Kind is not Ptr.
// It returns the zero Value if v is a nil pointer.
func (v Value) Elem() (Value, error) {
	k := v.kind()
	switch k {
	case K.Interface:
		// Interface handling - simplified version
		// For now we'll return an error for interfaces
		return Value{}, Err(ref, D.Not, D.Supported, "interface elem not yet implemented")

	case K.Pointer:
		ptr := v.ptr
		if v.flag&flagIndir != 0 {
			ptr = *(*unsafe.Pointer)(ptr)
		}
		// The returned value's address is v's value.
		if ptr == nil {
			return Value{}, nil
		}
		// Use the Type.Elem() method to get the element type
		typ := v.typ().Elem()
		if typ == nil {
			return Value{}, Err(ref, D.Value, D.Type, D.Nil)
		}
		fl := v.flag&flagRO | flagIndir | flagAddr
		fl |= flag(typ.Kind())
		return Value{v.tr, typ, ptr, fl}, nil
	}
	return Value{}, Err(ref, D.Value, D.NotOfType, D.Pointer)
}

func (v Value) NumField() (int, error) {
	if v.typ_ == nil {
		return 0, Err(ref, D.Value, D.Nil)
	}
	if v.kind() != K.Struct {
		return 0, Err(ref, D.Numbers, D.Fields, D.NotOfType, "Struct")
	}

	// Fast path: check cache
	if v.tr != nil {
		structID := v.typ_.StructID()
		count := atomic.LoadInt32(&v.tr.structCount)
		for i := int32(0); i < count; i++ {
			if v.tr.structCache[i].structID == structID {
				return int(v.tr.structCache[i].fieldCount), nil
			}
		}
	}

	// Slow path: reflect
	st := v.typ_.StructType()
	if st == nil {
		return 0, Err(ref, D.Numbers, D.Fields, D.NotOfType, "Struct")
	}
	return len(st.Fields), nil
}

// Field returns the i'th field of the struct v.
// Returns an error if v is not a struct or i is out of range.
func (v Value) Field(i int) (Value, error) {
	if v.kind() != K.Struct {
		return Value{}, Err(ref, D.Value, D.NotOfType, "Struct")
	}

	// Fast path: check cache
	if v.tr != nil {
		structID := v.typ_.StructID()
		count := atomic.LoadInt32(&v.tr.structCount)
		for j := int32(0); j < count; j++ {
			if v.tr.structCache[j].structID == structID {
				cachedStruct := &v.tr.structCache[j]
				if uint(i) >= uint(cachedStruct.fieldCount) {
					return Value{}, Err(ref, D.Value, D.Index, D.Out, D.Of, D.Range)
				}
				fieldSchema := &cachedStruct.fieldSchemas[i]
				typ := fieldSchema.typ
				fl := v.flag&(flagStickyRO|flagIndir|flagAddr) | flag(typ.Kind())
				// This part is tricky without full name info, so we assume public fields in cache for now.
				// A more robust implementation might cache the IsExported flag as well.
				ptr := add(v.ptr, uintptr(fieldSchema.offset), "same as non-reflect &v.field")
				return Value{v.tr, typ, ptr, fl}, nil
			}
		}
	}

	// Slow path: reflect
	tt := (*StructType)(unsafe.Pointer(v.typ()))
	if uint(i) >= uint(len(tt.Fields)) {
		return Value{}, Err(ref, D.Value, D.Index, D.Out, D.Of, D.Range)
	}
	field := &tt.Fields[i]
	typ := field.Typ

	fl := v.flag&(flagStickyRO|flagIndir|flagAddr) | flag(typ.Kind())
	if !field.Name.IsExported() {
		if field.Embedded() {
			fl |= flagEmbedRO
		} else {
			fl |= flagStickyRO
		}
	}
	ptr := add(v.ptr, field.Off, "same as non-reflect &v.field")
	return Value{v.tr, typ, ptr, fl}, nil
}

// Type returns v's type.
func (v Value) Type() *Type {
	if v.typ_ == nil {
		// This is where "value type nil" error would come from
		return nil
	}
	return v.typ()
}

// typ returns the *abi.Type stored in the Value. This method is fast,
// but it doesn't always return the correct type for the Value.
// See abiType and Type, which do return the correct type.
func (v Value) typ() *Type {
	// Types are either static (for compiler-created types) or
	// heap-allocated but always reachable (for reflection-created
	// types, held in the central map). So there is no need to
	// escape types. Direct cast is safe here.
	return v.typ_
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
//nolint:govet
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

func (f flag) kind() Kind {
	return Kind(f & flagKindMask)
}

// NoEscape function removed to avoid go vet warnings.
// The direct pointer access in typ() method is safe since types are
// either static or heap-allocated and always reachable.

// Kind returns the specific kind of this value.
func (v Value) Kind() Kind {
	return v.kind()
}

// CanAddr reports whether the value's address can be obtained with [Value.Addr].
// Such values are called addressable. A value is addressable if it is
// an element of a slice, an element of an addressable array,
// a field of an addressable struct, or the result of dereferencing a pointer.
// If CanAddr returns false, calling [Value.Addr] will panic.
func (v Value) CanAddr() bool {
	return v.flag&flagAddr != 0
}

// String returns the string v's underlying value, as a string.
// String is a special case because of Go's String method convention.
// Unlike the other getters, it does not panic if v's Kind is not String.
// Instead, it returns a string of the form "<Translate value>" where Translate is v's type.
func (v Value) String() string {
	// stringNonString is split out to keep String inlineable for string kinds.
	if v.kind() == K.String {
		return *(*string)(v.ptr)
	}
	return v.stringNonString()
}

func (v Value) stringNonString() string {
	if v.kind() == K.Invalid {
		return "<invalid Value>"
	}
	// If you call String on a reflect.Value of other type, it's better to
	// print something than to panic. Useful in debugging.
	return "<" + v.Type().String() + " Value>"
}

// Int returns v's underlying value, as an int64.
// It returns an error if v's Kind is not Int, Int8, Int16, Int32, or Int64.
func (v Value) Int() (int64, error) {
	k := v.kind()
	p := v.ptr
	switch k {
	case K.Int:
		return int64(*(*int)(p)), nil
	case K.Int8:
		return int64(*(*int8)(p)), nil
	case K.Int16:
		return int64(*(*int16)(p)), nil
	case K.Int32:
		return int64(*(*int32)(p)), nil
	case K.Int64:
		return *(*int64)(p), nil
	}
	return 0, Err(ref, D.Value, D.NotOfType, "int")
}

// Uint returns v's underlying value, as a uint64.
// It returns an error if v's Kind is not Uint, Uint8, Uint16, Uint32, or Uint64.
func (v Value) Uint() (uint64, error) {
	k := v.kind()
	p := v.ptr
	switch k {
	case K.Uint:
		return uint64(*(*uint)(p)), nil
	case K.Uint8:
		return uint64(*(*uint8)(p)), nil
	case K.Uint16:
		return uint64(*(*uint16)(p)), nil
	case K.Uint32:
		return uint64(*(*uint32)(p)), nil
	case K.Uint64:
		return *(*uint64)(p), nil
	case K.Uintptr:
		return uint64(*(*uintptr)(p)), nil
	}
	return 0, Err(ref, D.Value, D.NotOfType, "uint")
}

// Float returns v's underlying value, as a float64.
// It returns an error if v's Kind is not Float32 or Float64.
func (v Value) Float() (float64, error) {
	k := v.kind()
	switch k {
	case K.Float32:
		return float64(*(*float32)(v.ptr)), nil
	case K.Float64:
		return *(*float64)(v.ptr), nil
	}
	return 0, Err(ref, D.Value, D.NotOfType, "float")
}

// Bool returns v's underlying value.
// It returns an error if v's Kind is not Bool.
func (v Value) Bool() (bool, error) {
	if v.kind() != K.Bool {
		return false, Err(ref, D.Value, D.NotOfType, "bool")
	}
	return *(*bool)(v.ptr), nil
}

// IsZero reports whether v is the zero value for its type.
// It mirrors reflect.Value.IsZero() behavior for supported types.
func (v Value) IsZero() bool {
	// Handle nil Value (from ValueOf(nil))
	if v.typ_ == nil {
		return true
	}

	switch v.kind() {
	case K.String:
		return *(*string)(v.ptr) == ""
	case K.Bool:
		return !*(*bool)(v.ptr)
	case K.Int:
		return *(*int)(v.ptr) == 0
	case K.Int8:
		return *(*int8)(v.ptr) == 0
	case K.Int16:
		return *(*int16)(v.ptr) == 0
	case K.Int32:
		return *(*int32)(v.ptr) == 0
	case K.Int64:
		return *(*int64)(v.ptr) == 0
	case K.Uint:
		return *(*uint)(v.ptr) == 0
	case K.Uint8:
		return *(*uint8)(v.ptr) == 0
	case K.Uint16:
		return *(*uint16)(v.ptr) == 0
	case K.Uint32:
		return *(*uint32)(v.ptr) == 0
	case K.Uint64:
		return *(*uint64)(v.ptr) == 0
	case K.Uintptr:
		return *(*uintptr)(v.ptr) == 0
	case K.Float32:
		return *(*float32)(v.ptr) == 0
	case K.Float64:
		return *(*float64)(v.ptr) == 0
	case K.Pointer, K.Interface:
		return v.ptr == nil
	case K.Slice:
		// For slices, check if the data pointer is nil
		if v.ptr == nil {
			return true
		}
		// Slice header: {data uintptr, len int, cap int}
		// First field is the data pointer
		dataPtr := (*uintptr)(v.ptr)
		return *dataPtr == 0
	case K.Map:
		return v.ptr == nil
	case K.Struct:
		// Recursively check all fields
		num, err := v.NumField()
		if err != nil {
			return false
		}
		for i := 0; i < num; i++ {
			field, err := v.Field(i)
			if err != nil || !field.IsZero() {
				return false
			}
		}
		return true
	default:
		// For unsupported types, consider non-zero to be safe
		return false
	}
}

// InterfaceZeroAlloc sets v's current value to the target pointer without interface{} boxing.
// This method eliminates interface{} boxing allocations for primitive types by directly
// manipulating the interface{} structure to avoid the boxing that occurs when returning any.
//
// For primitive types (int, string, bool, float64, etc.), it assigns the actual value directly
// to the interface{} structure without creating boxing overhead.
//
// For complex types (slices, maps, structs, etc.), it falls back to the standard Interface()
// method which does create boxing, but this is unavoidable for complex types.
//
// This optimization is particularly beneficial for high-performance code that needs to extract
// primitive values from reflection without the boxing overhead.
func (v Value) InterfaceZeroAlloc(target *any) {
	if v.typ_ == nil {
		*target = nil
		return
	}

	k := v.kind()

	// For primitive types, use direct unsafe manipulation to avoid boxing
	switch k {
	case K.String, K.Int, K.Int8, K.Int16, K.Int32, K.Int64,
		K.Uint, K.Uint8, K.Uint16, K.Uint32, K.Uint64, K.Uintptr,
		K.Bool, K.Float32, K.Float64:

		// Use packEface technique but directly modify the target
		t := v.typ()
		e := (*EmptyInterface)(unsafe.Pointer(target))
		e.Type = t
		e.Data = v.ptr

	default:
		// For complex types (slice, map, struct, interface, etc.), use standard boxing
		// This is unavoidable for complex types
		if iface, err := v.Interface(); err == nil {
			*target = iface
		}
	}
}
