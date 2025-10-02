//go:build !tinygo

package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

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

	// Create the Value with proper flags for addressability
	// The pointer should be flagIndir, and the pointed-to value should have flagAddr
	ptrFlags := flag(K.Pointer) | flagIndir
	return Value{(*Type)(unsafe.Pointer(alignedPtrType)), unsafe.Pointer(&ptrStorage[0]), ptrFlags}
}

// makeSliceData allocates memory for slice data (stdlib version)
func makeSliceData(elemType *Type, cap int) unsafe.Pointer {
	var data unsafe.Pointer
	if elemType.Size != 0 {
		mem := make([]byte, uintptr(cap)*elemType.Size)
		data = unsafe.Pointer(&mem[0])
	}
	return data
}
