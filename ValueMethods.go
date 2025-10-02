package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Index returns v's i'th element.
// It returns an error if v's Kind is not Array, Slice, or String or i is out of range.
func (v Value) Index(i int) (Value, error) {
	switch v.kind() {
	case K.Array:
		if v.typ_ == nil {
			return Value{}, Err(D.Value, D.Type, D.Nil)
		}
		arrayType := (*ArrayType)(unsafe.Pointer(v.typ_))
		if uint(i) >= uint(arrayType.Len) {
			return Value{}, Err(D.Index, D.Out, D.Of, D.Range)
		}

		elemType := arrayType.Elem
		elemAddr := unsafe.Pointer(uintptr(v.ptr) + uintptr(i)*getElemSize(elemType))
		fl := v.flag&^flagKindMask | flag(elemType.Kind()) | flagIndir | flagAddr
		return Value{elemType, elemAddr, fl}, nil

	case K.Slice:
		sliceHeader := (*sliceHeader)(v.ptr)
		if uint(i) >= uint(sliceHeader.Len) {
			return Value{}, Err(D.Index, D.Out, D.Of, D.Range)
		}

		if v.typ_ == nil {
			return Value{}, Err(D.Value, D.Type, D.Nil)
		}
		sliceType := (*SliceType)(unsafe.Pointer(v.typ_))
		elemType := sliceType.Elem
		elemAddr := unsafe.Pointer(uintptr(sliceHeader.Data) + uintptr(i)*getElemSize(elemType))
		fl := v.flag&^flagKindMask | flag(elemType.Kind()) | flagIndir | flagAddr
		return Value{elemType, elemAddr, fl}, nil

	case K.String:
		stringHeader := (*stringHeader)(v.ptr)
		if uint(i) >= uint(stringHeader.Len) {
			return Value{}, Err(D.Index, D.Out, D.Of, D.Range)
		}

		byteAddr := unsafe.Pointer(uintptr(stringHeader.Data) + uintptr(i))
		fl := flag(K.Uint8) | flagIndir | flagAddr
		uint8Type := TypeOf(uint8(0))
		return Value{uint8Type, byteAddr, fl}, nil
	}

	return Value{}, Err(D.Call, D.Of, "Index", D.Method, v.kind().String(), D.Value)
}

// Len returns v's length.
// It returns an error if v's Kind is not Array, Chan, Map, Slice, or String.
func (v Value) Len() (int, error) {
	switch v.kind() {
	case K.Array:
		if v.typ_ == nil {
			return 0, Err(D.Value, D.Type, D.Nil)
		}
		arrayType := (*ArrayType)(unsafe.Pointer(v.typ_))
		return int(arrayType.Len), nil

	case K.Slice:
		sliceHeader := (*sliceHeader)(v.ptr)
		return sliceHeader.Len, nil

	case K.String:
		stringHeader := (*stringHeader)(v.ptr)
		return stringHeader.Len, nil
	}

	return 0, Err(D.Call, D.Of, "Len", D.Method, v.kind().String(), D.Value)
}

// Cap returns v's capacity.
// It returns an error if v's Kind is not Array, Chan, Slice or pointer to Array.
func (v Value) Cap() (int, error) {
	switch v.kind() {
	case K.Array:
		if v.typ_ == nil {
			return 0, Err(D.Value, D.Type, D.Nil)
		}
		arrayType := (*ArrayType)(unsafe.Pointer(v.typ_))
		return int(arrayType.Len), nil

	case K.Slice:
		sliceHeader := (*sliceHeader)(v.ptr)
		return sliceHeader.Cap, nil
	}

	return 0, Err(D.Call, D.Of, "Cap", D.Method, v.kind().String(), D.Value)
}

// IsNil reports whether its argument v is nil.
// It returns an error if v's Kind is not Chan, Func, Interface, Map, Pointer, or Slice.
func (v Value) IsNil() (bool, error) {
	switch v.kind() {
	case K.Pointer:
		if v.flag&flagIndir != 0 {
			return *(*unsafe.Pointer)(v.ptr) == nil, nil
		}
		return v.ptr == nil, nil

	case K.Slice:
		sliceHeader := (*sliceHeader)(v.ptr)
		return sliceHeader.Data == nil, nil

	case K.Interface:
		return v.ptr == nil, nil
	}

	return false, Err(D.Call, D.Of, "IsNil", D.Method, v.kind().String(), D.Value)
}

// Addr method is implemented in ValueMethods_stdlib.go and ValueMethods_tinygo.go

// Set assigns x to the value v.
// It returns an error if CanSet would return false.
func (v Value) Set(x Value) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}

	if v.typ_ == nil || x.typ_ == nil {
		return Err(D.Value, D.Type, D.Nil)
	}

	// Allow assignment between compatible pointer types
	if v.kind() == K.Pointer && x.kind() == K.Pointer {
		// For pointers, allow assignment if pointing to compatible types
		// This is more permissive than strict type equality
	} else if v.typ_ != x.typ_ {
		return Err(D.Value, D.Of, D.Type, x.typ_.String(), D.Not, "assignable", D.Type, v.typ_.String())
	}

	size := getTypeSize(v.typ_)
	if size == 0 {
		return nil
	}

	var srcPtr, dstPtr unsafe.Pointer
	if x.flag&flagIndir != 0 {
		srcPtr = x.ptr
	} else {
		srcPtr = unsafe.Pointer(&x.ptr)
	}
	if v.flag&flagIndir != 0 {
		dstPtr = v.ptr
	} else {
		dstPtr = unsafe.Pointer(&v.ptr)
	}

	for i := uintptr(0); i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(dstPtr) + i)) = *(*byte)(unsafe.Pointer(uintptr(srcPtr) + i))
	}

	return nil
}
