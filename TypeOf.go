package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// TFlag is used by a Type to signal what extra type information is
// available in the memory directly following the Type value.
type TFlag uint8

// NameOff is the Offset to a Name from moduledata.types.  See resolveNameOff in runtime.
type NameOff int32

// TypeOff is the Offset to a type from moduledata.types.  See resolveTypeOff in runtime.
type TypeOff int32

// Essential constants for type operations
const (
	KindDirectIface Kind = 1 << 5
	KindMask        Kind = (1 << 5) - 1
)

// Dictionary reference for reflection errors - "reflect" cannot be translated
const ref = "reflect"

// EmptyInterface describes the layout of a "interface{}" or a "any."
type EmptyInterface struct {
	Type *Type
	Data unsafe.Pointer
}

// Type is the runtime representation of a Go type.
//
// Be careful about accessing this type at build time, as the version
// of this type in the compiler/linker may not have the same layout
// as the version in the target binary, due to pointer width
// differences and any experiments. Use cmd/compile/internal/rttype
// or the functions in compiletype.go to access this type instead.
// (TODO: this admonition applies to every type in this package.
// Put it in some shared location?)

type Type struct {
	Size        uintptr
	PtrBytes    uintptr // number of (prefix) bytes in the type that can contain pointers
	Hash        uint32  // Hash of type; avoids computation in Hash tables
	TFlag       TFlag   // extra type information flags
	Align_      uint8   // alignment of variable with this type
	FieldAlign_ uint8   // alignment of struct field with this type
	Kind_       Kind    // enumeration for C
	// function for comparing objects of this type
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
	// GCData stores the GC type data for the garbage collector.
	// Normally, GCData points to a bitmask that describes the
	// ptr/nonptr Fields of the type. The bitmask will have at
	// least PtrBytes/ptrSize bits.
	// If the TFlagGCMaskOnDemand bit is set, GCData is instead a
	// **byte and the pointer to the bitmask is one dereference away.
	// The runtime will build the bitmask if needed.
	// (See runtime/type.go:getGCMask.)
	// Note: multiple types may have the same value of GCData,
	// including when TFlagGCMaskOnDemand is set. The types will, of course,
	// have the same pointer layout (but not necessarily the same size).
	GCData    *byte
	Str       NameOff // string form
	PtrToThis TypeOff // type for pointer to this type, may be zero
}

// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) *Type {
	if i == nil {
		return nil
	}
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	return e.Type
}

func (t *Type) String() string {
	return t.Kind().String()
}

// Name returns the type's name within its package for a defined type.
// For TinyGo compatibility, struct names are resolved using the StructDictionary
func (t *Type) Name() string {
	if t.Kind_ == K.Struct {
		// The original implementation caused a stack overflow due to infinite recursion.
		// The fix is to resolve the name via the type's name offset (Str),
		// which is how the standard reflect package does it internally.
		return t.nameOff(t.Str).Name()
	}

	// For other types, return empty string (following Go's behavior for unnamed types)
	return t.Kind_.String()
}

// nameOff resolves a name offset from a base pointer.
func (t *Type) nameOff(off NameOff) Name {
	// This is a simplified version of what the runtime does.
	// It assumes the name data is accessible relative to the type pointer.
	// In a real scenario, this might involve looking into module data.
	return Name{} // Placeholder, real implementation needed based on tinystring
}

// StructID returns a unique identifier for struct types based on runtime hash
// Returns 0 for non-struct types
func (t *Type) StructID() uint32 {
	if t.Kind() == K.Struct {
		return t.Hash
	}
	return 0
}

// StructType returns t cast to a *StructType, or nil if its tag does not match.
func (t *Type) StructType() *StructType {
	if t.Kind() != K.Struct {
		return nil
	}
	return (*StructType)(unsafe.Pointer(t))
}

func (t *Type) Kind() Kind { return t.Kind_ & KindMask }

// IfaceIndir reports whether t is stored indirectly in an interface value.
func (t *Type) IfaceIndir() bool {
	return t.Kind_&KindDirectIface == 0
}

// Field returns the i'th field of the struct type.
// It returns an error if the type is not a struct or the index is out of range.
func (t *Type) Field(i int) (StructField, error) {
	if t.Kind() != K.Struct {
		return StructField{}, Err(ref, D.Field, D.NotOfType, D.Struct)
	}
	st := (*StructType)(unsafe.Pointer(t))
	if i < 0 || i >= len(st.Fields) {
		return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)
	}
	return st.Fields[i], nil
}

// NumField returns the number of fields in the struct type.
// It returns an error if the type is not a struct.
func (t *Type) NumField() (int, error) {
	if t.Kind() != K.Struct {
		return 0, Err(ref, D.Numbers, D.Fields, D.Type, D.Struct)
	}
	st := (*StructType)(unsafe.Pointer(t))
	return len(st.Fields), nil
}

// PtrType represents a pointer type.
type PtrType struct {
	Type
	Elem *Type // pointer element type
}

// Name returns the name of a struct type's i'th field.
// It panics if the type's Kind is not Struct.
// It panics if i is out of range.
func (t *Type) NameByIndex(i int) (string, error) {
	if t.Kind() != K.Struct {
		return "", Err(ref, D.Type, D.NotOfType, D.Struct)
	}
	tt := (*StructType)(unsafe.Pointer(t))
	if i < 0 || i >= len(tt.Fields) {
		return "", Err(ref, D.Index, D.Out, D.Of, D.Range)
	}
	f := &tt.Fields[i]
	return f.Name.Name(), nil
}

// SliceType returns t cast to a *SliceType, or nil if its tag does not match.
func (t *Type) SliceType() *SliceType {
	if t.Kind() != K.Slice {
		return nil
	}
	return (*SliceType)(unsafe.Pointer(t))
}

// ArrayType returns t cast to a *ArrayType, or nil if its tag does not match.
func (t *Type) ArrayType() *ArrayType {
	if t.Kind() != K.Array {
		return nil
	}
	return (*ArrayType)(unsafe.Pointer(t))
}

// PtrType returns t cast to a *PtrType, or nil if its tag does not match.
func (t *Type) PtrType() *PtrType {
	if t.Kind() != K.Pointer {
		return nil
	}
	return (*PtrType)(unsafe.Pointer(t))
}

// Elem returns the element type for t if t is an array, channel, map, pointer, or slice, otherwise nil.
func (t *Type) Elem() *Type {
	switch t.Kind() {
	case K.Array:
		tt := (*ArrayType)(unsafe.Pointer(t))
		return tt.Elem
	case K.Pointer:
		tt := (*PtrType)(unsafe.Pointer(t))
		return tt.Elem
	case K.Slice:
		tt := (*SliceType)(unsafe.Pointer(t))
		return tt.Elem
	default:
		return nil
	}
}
