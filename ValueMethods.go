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
		// For arrays, we need to check bounds and calculate element address
		if v.typ_ == nil {
			return Value{}, Err(D.Value, D.Type, D.Nil)
		}
		arrayType := (*ArrayType)(unsafe.Pointer(v.typ_))
		if uint(i) >= uint(arrayType.Len) {
			return Value{}, Err(D.Index, D.Out, D.Of, D.Range)
		}

		// Calculate element address
		elemType := arrayType.Elem
		elemSize := elemType.Size
		elemAddr := unsafe.Pointer(uintptr(v.ptr) + uintptr(i)*elemSize)

		// Create new value for the element
		fl := v.flag&^flagKindMask | flag(elemType.Kind())
		fl |= flagIndir | flagAddr

		return Value{elemType, elemAddr, fl}, nil

	case K.Slice:
		// For slices, we need to get the slice header and check bounds
		sliceHeader := (*sliceHeader)(v.ptr)
		if uint(i) >= uint(sliceHeader.Len) {
			return Value{}, Err(D.Index, D.Out, D.Of, D.Range)
		}

		// Get slice element type
		if v.typ_ == nil {
			return Value{}, Err(D.Value, D.Type, D.Nil)
		}
		sliceType := (*SliceType)(unsafe.Pointer(v.typ_))
		elemType := sliceType.Elem
		elemSize := elemType.Size

		// Calculate element address
		elemAddr := unsafe.Pointer(uintptr(sliceHeader.Data) + uintptr(i)*elemSize)

		// Create new value for the element
		fl := v.flag&^flagKindMask | flag(elemType.Kind())
		fl |= flagIndir | flagAddr

		return Value{elemType, elemAddr, fl}, nil

	case K.String:
		// For strings, we need to get the string header and return byte
		stringHeader := (*stringHeader)(v.ptr)
		if uint(i) >= uint(stringHeader.Len) {
			return Value{}, Err(D.Index, D.Out, D.Of, D.Range)
		}

		// Get the byte at position i
		byteAddr := unsafe.Pointer(uintptr(stringHeader.Data) + uintptr(i))

		// Create new value for the byte (uint8)
		fl := flag(K.Uint8) | flagIndir | flagAddr
		uint8Type := TypeOf(uint8(0))

		return Value{uint8Type, byteAddr, fl}, nil
	}

	return Value{}, Err(D.Call, D.Of, "Index", D.Method, v.kind().String(), D.Value)
}

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

// Len returns v's length.
// It returns an error if v's Kind is not Array, Chan, Map, Slice, or String.
func (v Value) Len() (int, error) {
	switch v.kind() {
	case K.Array:
		// For arrays, we need to get the length from the type
		if v.typ_ == nil {
			return 0, Err(D.Value, D.Type, D.Nil)
		}
		arrayType := (*ArrayType)(unsafe.Pointer(v.typ_))
		return int(arrayType.Len), nil

	case K.Slice:
		// For slices, get length from slice header
		sliceHeader := (*sliceHeader)(v.ptr)
		return sliceHeader.Len, nil

	case K.String:
		// For strings, get length from string header
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
		// For arrays, capacity equals length
		if v.typ_ == nil {
			return 0, Err(D.Value, D.Type, D.Nil)
		}
		arrayType := (*ArrayType)(unsafe.Pointer(v.typ_))
		return int(arrayType.Len), nil

	case K.Slice:
		// For slices, get capacity from slice header
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
		// For pointers, check if ptr is nil
		if v.flag&flagIndir != 0 {
			return *(*unsafe.Pointer)(v.ptr) == nil, nil
		}
		return v.ptr == nil, nil

	case K.Slice:
		// For slices, check if Data is nil
		sliceHeader := (*sliceHeader)(v.ptr)
		return sliceHeader.Data == nil, nil

	case K.Interface:
		// For interfaces, check if the interface is nil
		return v.ptr == nil, nil
	}

	return false, Err(D.Call, D.Of, "IsNil", D.Method, v.kind().String(), D.Value)
}

// Addr returns a pointer value representing the address of v.
// It returns an error if CanAddr would return false.
func (v Value) Addr() (Value, error) {
	if v.flag&flagAddr == 0 {
		return Value{}, Err(D.Value, D.Not, "addressable")
	}

	if v.typ_ == nil {
		return Value{}, Err(D.Value, D.Type, D.Nil)
	}

	// Create a pointer type for the value's type
	// This is a simplified version - ideally we'd use a proper PtrTo function
	ptrType := &PtrType{
		Type: Type{
			Size:        unsafe.Sizeof(uintptr(0)),
			PtrBytes:    unsafe.Sizeof(uintptr(0)),
			Hash:        v.typ_.Hash ^ 0x87654321, // Hash variation for pointer
			TFlag:       v.typ_.TFlag,
			Align_:      uint8(unsafe.Alignof(uintptr(0))),
			FieldAlign_: uint8(unsafe.Alignof(uintptr(0))),
			Kind_:       K.Pointer,
			Equal:       nil,
			GCData:      nil,
			Str:         0,
			PtrToThis:   0,
		},
		Elem: v.typ_,
	}

	// Following Go's pattern: use v.ptr directly and preserve flagRO
	fl := (v.flag & flagRO) | flag(K.Pointer)
	return Value{&ptrType.Type, v.ptr, fl}, nil
}

// Set assigns x to the value v.
// It returns an error if CanSet would return false.
func (v Value) Set(x Value) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}

	// Check for zero values
	if v.typ_ == nil {
		return Err(D.Value, D.Type, D.Nil)
	}
	if x.typ_ == nil {
		return Err(D.Value, D.Type, D.Nil)
	}

	// Check that the types are compatible
	if v.typ_ != x.typ_ {
		// For slices, check if they have the same element type and structure
		if v.typ_.Kind() == K.Slice && x.typ_.Kind() == K.Slice {
			// Allow assignment if the underlying structure is compatible
			vElem := v.typ_.Elem()
			xElem := x.typ_.Elem()
			if vElem != nil && xElem != nil && vElem == xElem {
				// Types are slice-compatible, proceed with assignment
			} else {
				return Err(D.Value, D.Of, D.Type, x.typ_.String(), D.Not, "assignable", D.Type, v.typ_.String())
			}
		} else if v.typ_.Kind() == K.Pointer && x.typ_.Kind() == K.Pointer {
			// For pointers, check if they point to compatible types
			vElem := v.typ_.Elem()
			xElem := x.typ_.Elem()
			if vElem != nil && xElem != nil && vElem == xElem {
				// Pointer types with same element type are compatible
			} else {
				return Err(D.Value, D.Of, D.Type, x.typ_.String(), D.Not, "assignable", D.Type, v.typ_.String())
			}
		} else {
			return Err(D.Value, D.Of, D.Type, x.typ_.String(), D.Not, "assignable", D.Type, v.typ_.String())
		}
	}

	// Get the size of the type for copying
	size := v.typ_.Size
	if size == 0 {
		return nil // Nothing to copy
	}

	// Copy the data from x to v
	var srcPtr unsafe.Pointer
	if x.flag&flagIndir != 0 {
		srcPtr = x.ptr
	} else {
		srcPtr = unsafe.Pointer(&x.ptr)
	}

	var dstPtr unsafe.Pointer
	if v.flag&flagIndir != 0 {
		dstPtr = v.ptr
	} else {
		dstPtr = unsafe.Pointer(&v.ptr)
	}

	// Use unsafe to copy the bytes
	for i := uintptr(0); i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(dstPtr) + i)) = *(*byte)(unsafe.Pointer(uintptr(srcPtr) + i))
	}

	return nil
}
