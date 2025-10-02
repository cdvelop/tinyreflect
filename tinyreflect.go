package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

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

// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) *Type {
	if i == nil {
		return nil
	}
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	return e.Type
}

// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i. ValueOf(nil) returns the zero Value.
func ValueOf(i any) Value {
	if i == nil {
		return Value{}
	}
	return unpackEface(i)
}

// unpackEface converts the empty interface i to a Value.
func unpackEface(i any) Value {
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	t := e.Type
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if t.IfaceIndir() {
		f |= flagIndir
	}
	return Value{t, e.Data, f}
}

// Indirect returns the value that v points to.
// If v is a nil pointer, Indirect returns a zero Value.
// If v is not a pointer, Indirect returns v.
func Indirect(v Value) Value {
	if v.kind() != K.Pointer {
		return v
	}
	elem, err := v.Elem()
	if err != nil {
		return Value{}
	}
	return elem
}

// MakeSlice creates a new zero-initialized slice value
// for the specified slice type, length, and capacity.
func MakeSlice(typ *Type, len, cap int) (Value, error) {
	if typ == nil {
		return Value{}, Err(D.Value, D.Type, D.Nil)
	}
	if typ.Kind().String() != "slice" {
		return Value{}, Err("MakeSlice of non-slice type")
	}
	if len < 0 || cap < 0 || len > cap {
		return Value{}, Err("invalid slice length or capacity")
	}

	sliceType := (*SliceType)(unsafe.Pointer(typ))
	elemType := sliceType.Elem
	if elemType == nil {
		return Value{}, Err("MakeSlice element type is nil")
	}

	data := makeSliceData(elemType, cap)

	sliceHeader := &sliceHeader{Data: data, Len: len, Cap: cap}
	return Value{typ, unsafe.Pointer(sliceHeader), flagIndir | flag(K.Slice)}, nil
}

// NewValue is implemented in tinyreflect_stdlib.go and tinyreflect_tinygo.go
