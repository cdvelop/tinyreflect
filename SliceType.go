package tinyreflect

// SliceType represents a slice type.
type SliceType struct {
	Type
	Elem *Type // slice element type
}

// ArrayType represents an array type.
type ArrayType struct {
	Type
	Elem  *Type   // array element type
	Slice *Type   // slice type
	Len   uintptr // array length
}

// Elem returns the element type of the slice
func (t *SliceType) Element() *Type {
	return t.Elem
}

// Elem returns the element type of the array
func (t *ArrayType) Element() *Type {
	return t.Elem
}

// Length returns the length of the array
func (t *ArrayType) Length() int {
	return int(t.Len)
}
