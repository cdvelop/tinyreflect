package tinyreflect

import (
	. "github.com/cdvelop/tinystring"
)

// kind returns the Kind of the value.
func (v Value) kind() Kind {
	return Kind(v.flag & flagKindMask)
}

// mustBeAssignable checks if the value is assignable and returns an error if not.
// For tinyreflect, we simplify this to just checking if the value is valid.
func (v Value) mustBeAssignable() error {
	if v.flag == 0 {
		return Err(D.Value, D.Not, "assignable")
	}
	return nil
}

// mustBe checks if the value's kind is one of the expected kinds and returns an error if not.
func (v Value) mustBe(expected Kind) error {
	if k := v.kind(); k != expected {
		return Err(D.Call, D.Of, D.Method, k.String(), D.Value)
	}
	return nil
}

// SetString sets the string value to the field represented by Value.
// It uses unsafe to write the value to the memory location of the field.
func (v Value) SetString(x string) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}
	if err := v.mustBe(K.String); err != nil {
		return err
	}
	*(*string)(v.ptr) = x
	return nil
}

// SetBool sets the bool value to the field represented by Value.
func (v Value) SetBool(x bool) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}
	if err := v.mustBe(K.Bool); err != nil {
		return err
	}
	*(*bool)(v.ptr) = x
	return nil
}

// SetBytes sets the byte slice value to the field represented by Value.
func (v Value) SetBytes(x []byte) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}
	if err := v.mustBe(K.Slice); err != nil {
		return err
	}
	// For tinyreflect, we simplify and assume it's a []byte slice
	*(*[]byte)(v.ptr) = x
	return nil
}

// SetInt sets the int value to the field represented by Value.
func (v Value) SetInt(x int64) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}

	switch k := v.kind(); k {
	case K.Int:
		*(*int)(v.ptr) = int(x)
	case K.Int8:
		*(*int8)(v.ptr) = int8(x)
	case K.Int16:
		*(*int16)(v.ptr) = int16(x)
	case K.Int32:
		*(*int32)(v.ptr) = int32(x)
	case K.Int64:
		*(*int64)(v.ptr) = x
	default:
		return Err(D.Call, D.Of, "SetInt", D.Method, k.String(), D.Value)
	}
	return nil
}

// SetUint sets the uint value to the field represented by Value.
func (v Value) SetUint(x uint64) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}

	switch k := v.kind(); k {
	case K.Uint:
		*(*uint)(v.ptr) = uint(x)
	case K.Uint8:
		*(*uint8)(v.ptr) = uint8(x)
	case K.Uint16:
		*(*uint16)(v.ptr) = uint16(x)
	case K.Uint32:
		*(*uint32)(v.ptr) = uint32(x)
	case K.Uint64:
		*(*uint64)(v.ptr) = x
	case K.Uintptr:
		*(*uintptr)(v.ptr) = uintptr(x)
	default:
		return Err(D.Call, D.Of, "SetUint", D.Method, k.String(), D.Value)
	}
	return nil
}

// SetFloat sets the float value to the field represented by Value.
func (v Value) SetFloat(x float64) error {
	if err := v.mustBeAssignable(); err != nil {
		return err
	}

	switch k := v.kind(); k {
	case K.Float32:
		*(*float32)(v.ptr) = float32(x)
	case K.Float64:
		*(*float64)(v.ptr) = x
	default:
		return Err(D.Call, D.Of, "SetFloat", D.Method, k.String(), D.Value)
	}
	return nil
}
