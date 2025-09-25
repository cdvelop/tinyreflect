package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// sliceHeader is the runtime representation of a slice.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// stringHeader is the runtime representation of a string.
type stringHeader struct {
	Data unsafe.Pointer
	Len  int
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

// Indirect returns the value that v points to.
// If v is a nil pointer, Indirect returns a zero Value.
// If v is not a pointer, Indirect returns v.
func Indirect(v Value) Value {
	if v.kind() != K.Pointer {
		return v
	}
	elem, err := v.Elem()
	if err != nil {
		return Value{}
	}
	return elem
}

// MakeSlice creates a new zero-initialized slice value
// for the specified slice type, length, and capacity.
func MakeSlice(typ *Type, len, cap int) (Value, error) {
	if typ == nil {
		return Value{}, Err(D.Value, D.Type, D.Nil)
	}
	if typ.Kind().String() != "slice" {
		return Value{}, Err("MakeSlice of non-slice type")
	}
	if len < 0 || cap < 0 || len > cap {
		return Value{}, Err("invalid slice length or capacity")
	}

	sliceType := (*SliceType)(unsafe.Pointer(typ))
	elemType := sliceType.Elem
	if elemType == nil {
		return Value{}, Err("MakeSlice element type is nil")
	}

	var data unsafe.Pointer
	if elemType.Size != 0 {
		mem := make([]byte, uintptr(cap)*elemType.Size)
		data = unsafe.Pointer(&mem[0])
	}

	sliceHeader := &sliceHeader{Data: data, Len: len, Cap: cap}
	return Value{typ, unsafe.Pointer(sliceHeader), flagIndir | flag(K.Slice)}, nil
}

// NewValue returns a Value representing a pointer to a new zero value
// for the specified type.
func NewValue(typ *Type) Value {
	if typ == nil {
		return Value{}
	}

	// Create properly aligned memory for the target value
	size := typ.Size
	if size == 0 {
		size = 1 // Ensure minimum size for zero-sized types
	}

	// Allocate aligned memory for the value
	valuePtr := make([]byte, size+7) // Add padding for alignment
	alignedValuePtr := unsafe.Pointer((uintptr(unsafe.Pointer(&valuePtr[0])) + 7) &^ 7)

	// Allocate aligned memory for the pointer to the value
	ptrStorage := make([]uintptr, 1) // Use uintptr slice for proper alignment
	ptrStorage[0] = uintptr(alignedValuePtr)

	// Create a proper pointer type that won't cause alignment issues
	// We'll use typ.Elem() construction pattern to avoid creating local PtrType
	ptrTypeSize := unsafe.Sizeof(uintptr(0))

	// Create type info for pointer type - allocate it properly aligned
	ptrTypeStorage := make([]byte, unsafe.Sizeof(PtrType{})+7)
	alignedPtrType := (*PtrType)(unsafe.Pointer((uintptr(unsafe.Pointer(&ptrTypeStorage[0])) + 7) &^ 7))

	// Initialize the pointer type properly
	alignedPtrType.Type = Type{
		Kind_:    K.Pointer,
		Size:     ptrTypeSize,
		PtrBytes: ptrTypeSize,
		Hash:     typ.Hash ^ 0x12345678, // Simple hash derivation
	}
	alignedPtrType.Elem = typ

	return Value{(*Type)(unsafe.Pointer(alignedPtrType)), unsafe.Pointer(&ptrStorage[0]), flag(K.Pointer) | flagIndir}
}
