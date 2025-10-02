//go:build tinygo

package tinyreflect

import (
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

	// For TinyGo, we can't create Type literals
	// We need to use TypeOf to get a proper pointer type
	ptrType := &PtrType{
		Type: Type{meta: uint8(K.Pointer)},
		Elem: v.typ_,
	}

	fl := (v.flag & flagRO) | flag(K.Pointer)
	return Value{&ptrType.Type, v.ptr, fl}, nil
}

// getElemSize returns the size of an element type (TinyGo version)
func getElemSize(elemType *Type) uintptr {
	return elemType.Size()
}

// getTypeSize returns the size of a type (TinyGo version)
func getTypeSize(typ *Type) uintptr {
	return typ.Size()
}
