package tinyreflect

import (
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// TFlag is used by a Type to signal what extra type information is
// available in the memory directly following the Type value.
type TFlag uint8

// NameOff is the Offset to a Name from moduledata.types.  See resolveNameOff in runtime.
type NameOff int32

// TypeOff is the Offset to a type from moduledata.types.  See resolveTypeOff in runtime.
type TypeOff int32

// Essential constants for type operations
const (
	KindDirectIface Kind = 1 << 5
	KindMask        Kind = (1 << 5) - 1
)

// Dictionary reference for reflection errors - "reflect" cannot be translated
const ref = "reflect"

// EmptyInterface describes the layout of a "any" or a "any."
type EmptyInterface struct {
	Type *Type
	Data unsafe.Pointer
}

// Type is defined in Type_stdlib.go and Type_tinygo.go with build tags
// to handle the different internal representations in stdlib vs TinyGo
// Methods like Kind(), Name(), StructID(), IfaceIndir() are defined in
// Type_stdlib.go and Type_tinygo.go with build-specific implementations.

func (t *Type) String() string {
	return t.Kind().String()
}

// StructType returns t cast to a *StructType, or nil if its tag does not match.
func (t *Type) StructType() *StructType {
	if t.Kind() != K.Struct {
		return nil
	}
	// Get underlying type first (handles named types in TinyGo)
	ut := t.underlying()
	return (*StructType)(unsafe.Pointer(ut))
}

// Field returns the i'th field of the struct type.
// It returns an error if the type is not a struct or the index is out of range.
func (t *Type) Field(i int) (StructField, error) {
	// Get underlying type first (handles named types)
	ut := t.underlying()
	if ut.Kind() != K.Struct {
		return StructField{}, Err(ref, D.Field, D.NotOfType, "Struct")
	}
	st := (*StructType)(unsafe.Pointer(ut))
	if st == nil {
		return StructField{}, Err(ref, D.Field, D.NotOfType, "Struct")
	}
	if i < 0 || i >= st.numFields() {
		return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)
	}
	f := st.getField(i)
	if f == nil {
		return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)
	}
	return *f, nil
}

// NumField returns the number of fields in the struct type.
// It returns an error if the type is not a struct.
func (t *Type) NumField() (int, error) {
	// Get underlying type first (handles named types)
	ut := t.underlying()

	// Convert to StructType and check if it has fields
	st := (*StructType)(unsafe.Pointer(ut))
	if st == nil {
		return 0, Err(ref, D.Numbers, D.Fields, D.NotOfType, "Struct")
	}

	n := st.numFields()
	if n == 0 {
		// Verify it's actually a struct type
		if ut.Kind() != K.Struct {
			return 0, Err(ref, D.Numbers, D.Fields, D.Type, "Struct")
		}
	}

	return n, nil
}

// PtrType represents a pointer type.
type PtrType struct {
	Type
	Elem *Type // pointer element type
}

// Name returns the name of a struct type's i'th field.
// It panics if the type's Kind is not Struct.
// It panics if i is out of range.
func (t *Type) NameByIndex(i int) (string, error) {
	// Get underlying type first (handles named types)
	ut := t.underlying()
	tt := (*StructType)(unsafe.Pointer(ut))

	if tt == nil {
		println("DEBUG NameByIndex: StructType is nil")
		return "", Err(ref, D.Type, D.NotOfType, "Struct")
	}

	numFields := tt.numFields()
	println("DEBUG NameByIndex: numFields =", numFields, "requested i =", i)

	// If no fields, verify it's actually supposed to be a struct
	if numFields == 0 {
		kind := ut.Kind()
		println("DEBUG NameByIndex: numFields=0, checking kind =", kind)
		if kind != K.Struct {
			return "", Err(ref, D.Type, D.NotOfType, "Struct")
		}
	}

	if i < 0 || i >= numFields {
		return "", Err(ref, D.Index, D.Out, D.Of, D.Range)
	}

	f := tt.getField(i)
	println("DEBUG NameByIndex: field pointer =", f != nil)
	if f == nil {
		return "", Err(ref, "field is nil")
	}

	println("DEBUG NameByIndex: f.Name.Bytes =", f.Name.Bytes != nil)
	if f.Name.Bytes == nil {
		println("DEBUG NameByIndex: Name.Bytes is nil, returning empty")
		return "", nil
	}

	// Try to safely read the name with error recovery
	defer func() {
		if r := recover(); r != nil {
			println("DEBUG NameByIndex: PANIC recovered:", r)
		}
	}()

	name := f.Name.Name()
	println("DEBUG NameByIndex: name =", name)
	return name, nil
}

// SliceType returns t cast to a *SliceType, or nil if its tag does not match.
func (t *Type) SliceType() *SliceType {
	if t.Kind() != K.Slice {
		return nil
	}
	return (*SliceType)(unsafe.Pointer(t))
}

// ArrayType returns t cast to a *ArrayType, or nil if its tag does not match.
func (t *Type) ArrayType() *ArrayType {
	if t.Kind() != K.Array {
		return nil
	}
	return (*ArrayType)(unsafe.Pointer(t))
}

// PtrType returns t cast to a *PtrType, or nil if its tag does not match.
func (t *Type) PtrType() *PtrType {
	if t.Kind() != K.Pointer {
		return nil
	}
	return (*PtrType)(unsafe.Pointer(t))
}

// Elem returns the element type for t if t is an array, channel, map, pointer, or slice, otherwise nil.
func (t *Type) Elem() *Type {
	switch t.Kind() {
	case K.Array:
		tt := (*ArrayType)(unsafe.Pointer(t))
		return tt.Elem
	case K.Pointer:
		tt := (*PtrType)(unsafe.Pointer(t))
		return tt.Elem
	case K.Slice:
		tt := (*SliceType)(unsafe.Pointer(t))
		return tt.Elem
	default:
		return nil
	}
}
