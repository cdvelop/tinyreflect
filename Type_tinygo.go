//go:build tinygo

package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Type is the runtime representation of a Go type (TinyGo version).
// This matches TinyGo's internal/reflectlite.RawType structure.
// NOTE: This is just 1 byte! The actual type information extends beyond.
type Type struct {
	meta uint8 // metadata byte, contains kind and flags
}

// TinyGo-specific type constants
// Based on /usr/local/lib/tinygo/src/internal/reflectlite/type.go
const (
	kindMask  = 31 // mask to apply to the meta byte to get the Kind value
	flagNamed = 32 // flag that is set if this is a named type
)

// elemType represents a named type in TinyGo.
// All types that have an element type: named, chan, slice, array, map
// Layout matches TinyGo's internal/reflectlite/type.go:167
type elemType struct {
	meta      uint8
	numMethod uint16
	ptrTo     *Type
	elem      *Type // Pointer to underlying type
}

// Kind returns the type's Kind.
func (t *Type) Kind() Kind {
	if t == nil {
		return K.Invalid
	}
	if tag := t.ptrtag(); tag != 0 {
		return K.Pointer
	}
	return Kind(t.meta & kindMask)
}

// underlying returns the underlying type.
// For named types, returns the elem pointer recursively.
// For unnamed types, returns self.
func (t *Type) underlying() *Type {
	// Loop until we find a non-named type
	for t.isNamed() {
		t = (*elemType)(unsafe.Pointer(t)).elem
		if t == nil {
			return nil
		}
	}
	return t
}

// isNamed checks if this is a named type.
func (t *Type) isNamed() bool {
	if t.ptrtag() != 0 {
		return false
	}
	return t.meta&flagNamed != 0
}

// ptrtag returns the pointer tag (last 2 bits of pointer address).
func (t *Type) ptrtag() uintptr {
	return uintptr(unsafe.Pointer(t)) & 0b11
}

// StructID returns 0 for TinyGo (not implemented yet)
func (t *Type) StructID() uint32 {
	return 0
}

// Name returns the type's name (simplified for TinyGo)
func (t *Type) Name() string {
	return t.Kind().String()
}

// IfaceIndir always returns true for TinyGo (simplified)
func (t *Type) IfaceIndir() bool {
	return true
}

// Size returns the size in bytes of the type.
// For TinyGo, we need to read from the underlying type structure.
func (t *Type) Size() uintptr {
	if t == nil {
		return 0
	}

	underlying := t.underlying()
	kind := underlying.Kind()

	// For basic types, return standard sizes
	switch kind {
	case K.Bool, K.Int8, K.Uint8:
		return 1
	case K.Int16, K.Uint16:
		return 2
	case K.Int32, K.Uint32, K.Float32:
		return 4
	case K.Int64, K.Uint64, K.Float64, K.Complex64:
		return 8
	case K.Complex128:
		return 16
	case K.Int, K.Uint, K.Uintptr, K.Pointer, K.UnsafePointer:
		return unsafe.Sizeof(uintptr(0))
	case K.String:
		return unsafe.Sizeof("")
	case K.Slice:
		return unsafe.Sizeof([]byte{})
	case K.Struct:
		// Read size from structType
		st := (*StructType)(unsafe.Pointer(underlying))
		return uintptr(st.size)
	case K.Array:
		// Read size from arrayType
		at := (*ArrayType)(unsafe.Pointer(underlying))
		return uintptr(at.Len) * at.Elem.Size()
	}

	return 0
}
