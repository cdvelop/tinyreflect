package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Package-level functions required by codecs.go

// Indirect returns the value that v points to.
// If v is a nil pointer, Indirect returns a zero Value.
// If v is not a pointer, Indirect returns v.
func Indirect(v Value) Value {
	if v.kind() != K.Pointer {
		return v
	}

	// Use the existing Elem() method to get the pointed-to value
	elem, err := v.Elem()
	if err != nil {
		// If there's an error, return zero value - this preserves standard library behavior
		return Value{}
	}

	return elem
}

// MakeSlice creates a new zero-initialized slice value
// for the specified slice type, length, and capacity.
// Implementation adapted from /usr/local/go/src/reflect/value.go:2904
func MakeSlice(typ *Type, len, cap int) (Value, error) {
	if typ == nil {
		return Value{}, Err(D.Value, D.Type, D.Nil)
	}

	if typ.Kind().String() != "slice" {
		return Value{}, Err("MakeSlice", D.Of, "non-slice", D.Type)
	}
	if len < 0 {
		return Value{}, Err("MakeSlice", "negative", "len")
	}
	if cap < 0 {
		return Value{}, Err("MakeSlice", "negative", "cap")
	}
	if len > cap {
		return Value{}, Err("MakeSlice", "len", ">", "cap")
	}

	// Get the slice type information
	sliceType := (*SliceType)(unsafe.Pointer(typ))
	elemType := sliceType.Elem

	if elemType == nil {
		return Value{}, Err("MakeSlice", "element", D.Type, D.Nil)
	}

	// Simple implementation: create a real Go slice using make()
	// and extract its internal representation
	elemSize := elemType.Size

	if elemSize == 0 {
		// For zero-sized elements like struct{}, create a special case
		sliceHeader := &sliceHeader{
			Data: unsafe.Pointer(uintptr(1)), // Non-nil dummy pointer
			Len:  len,
			Cap:  cap,
		}
		return Value{typ, unsafe.Pointer(sliceHeader), flagIndir | flag(K.Slice)}, nil
	}

	// For regular elements, allocate memory and create slice header
	// Calculate total memory needed
	totalSize := uintptr(cap) * elemSize

	// Allocate raw memory using make([]byte, ...)
	// This ensures the memory is properly managed by Go's GC
	mem := make([]byte, totalSize)

	// Create the slice header pointing to our allocated memory
	sliceHeader := &sliceHeader{
		Data: unsafe.Pointer(&mem[0]),
		Len:  len,
		Cap:  cap,
	}

	// Return a Value that wraps the slice header
	return Value{typ, unsafe.Pointer(sliceHeader), flagIndir | flag(K.Slice)}, nil
}

// New returns a Value representing a pointer to a new zero value
// for the specified type. That is, the returned Value's Type is PtrTo(typ).
func New(typ *Type) Value {
	if typ == nil {
		return Value{}
	}

	// Step 1: Allocate memory for the pointed-to value
	size := typ.Size
	if size == 0 {
		size = 1 // Minimum allocation size
	}

	// Allocate memory for the actual value
	valueData := make([]byte, size)
	valuePtr := unsafe.Pointer(&valueData[0])

	// Step 2: Allocate memory for the pointer itself
	ptrData := make([]byte, unsafe.Sizeof(uintptr(0)))
	ptrPtr := unsafe.Pointer(&ptrData[0])

	// Step 3: Store the pointer to the value in the pointer memory
	*(*unsafe.Pointer)(ptrPtr) = valuePtr

	// Step 4: Create a basic PtrType structure
	ptrType := &PtrType{
		Type: Type{
			Size:        unsafe.Sizeof(uintptr(0)), // Size of a pointer
			PtrBytes:    unsafe.Sizeof(uintptr(0)), // All bytes are pointer bytes
			Hash:        typ.Hash ^ 0x12345678,     // Simple hash variation for pointer type
			TFlag:       0,                         // No special flags
			Align_:      8,                         // Pointer alignment (64-bit)
			FieldAlign_: 8,                         // Field alignment for pointers
			Kind_:       K.Pointer,                 // This is a pointer type
			Equal:       nil,                       // Simplified
			GCData:      nil,                       // Simplified
			Str:         0,                         // No name for generated pointer types
			PtrToThis:   0,                         // Not implemented
		},
		Elem: typ, // Points to the given type
	}

	// Return Value with pointer flags
	fl := flag(K.Pointer) | flagIndir
	return Value{&ptrType.Type, ptrPtr, fl}
}

// Zero returns a Value representing the zero value for the specified type.
func Zero(typ *Type) Value {
	if typ == nil {
		return Value{}
	}

	// For now, return a simplified implementation
	// TODO: Implement proper zero value creation
	return Value{}
}
