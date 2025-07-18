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
		// If there's an error (e.g., nil pointer), return zero value
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

	// For now, return a simplified implementation
	// TODO: Implement proper pointer creation
	return Value{}
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
