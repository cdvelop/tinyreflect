package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Interface returns v's current value as an any.
// It is equivalent to:
//
//	var i any = (v's underlying value)
//
// For a Value created from a nil interface value, Interface returns nil.
func (v Value) Interface() (i any, err error) {
	if v.typ_ == nil {
		return nil, Err(ref, D.Value, D.Nil)
	}

	if v.kind() == K.Interface {
		return nil, Err(ref, D.Type, D.Not, D.Supported)
	}

	i = packEface(v)

	return
}

// packEface converts v to the empty interface.
func packEface(v Value) any {
	t := v.typ()
	var i any
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	e.Type = t
	e.Data = v.ptr
	return i
}
