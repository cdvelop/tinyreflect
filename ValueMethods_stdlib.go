//go:build !tinygo

package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Addr returns a pointer value representing the address of v.
// It returns an error if CanAddr returns false.
func (v Value) Addr() (Value, error) {
	if v.flag&flagAddr == 0 {
		return Value{}, Err(D.Value, D.Not, "addressable")
	}

	if v.typ_ == nil {
		return Value{}, Err(D.Value, D.Type, D.Nil)
	}

	ptrType := &PtrType{
		Type: Type{Kind_: K.Pointer, Size: unsafe.Sizeof(uintptr(0))},
		Elem: v.typ_,
	}

	fl := (v.flag & flagRO) | flag(K.Pointer)
	return Value{&ptrType.Type, v.ptr, fl}, nil
}

// getElemSize returns the size of an element type (stdlib version)
func getElemSize(elemType *Type) uintptr {
	return elemType.Size
}

// getTypeSize returns the size of a type (stdlib version)
func getTypeSize(typ *Type) uintptr {
	return typ.Size
}
