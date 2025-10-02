//go:build !tinygo

package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Type is the runtime representation of a Go type (stdlib version).
// This matches the stdlib's internal/abi.Type structure.
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
	GCData    *byte
	Str       NameOff // string form
	PtrToThis TypeOff // type for pointer to this type, may be zero
}

// Name returns the type's name within its package for a defined type.
// For other types, returns the kind name (e.g., "int", "string").
func (t *Type) Name() string {
	return t.Kind_.String()
}

// StructID returns a unique identifier for struct types based on runtime hash
// Returns 0 for non-struct types
func (t *Type) StructID() uint32 {
	if t.Kind() == K.Struct {
		return t.Hash
	}
	return 0
}

// Kind returns the type's Kind.
func (t *Type) Kind() Kind {
	return t.Kind_ & KindMask
}

// IfaceIndir reports whether t is stored indirectly in an interface value.
func (t *Type) IfaceIndir() bool {
	// Simplified for now - most types are stored indirectly
	kind := t.Kind()
	return kind != K.Pointer && kind != K.UnsafePointer
}
