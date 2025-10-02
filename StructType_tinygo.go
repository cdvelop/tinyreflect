//go:build tinygo

package tinyreflect

import "unsafe"

// StructType represents a struct type.
// Layout matches TinyGo's internal structType for compatibility.
type StructType struct {
	Type
	numMethod uint16
	ptrTo     *Type
	pkgpath   *byte
	size      uint32
	numField  uint16
	Fields    [1]tinygoStructField // Array extends in memory beyond [1] - using TinyGo's internal layout
}

// getField returns a pointer to the field at index i.
// Uses unsafe pointer arithmetic because fields array extends beyond [1].
func (st *StructType) getField(i int) *StructField {
	if i < 0 || i >= int(st.numField) {
		return nil
	}
	// Calculate field address using TinyGo's internal field size
	fieldSize := unsafe.Sizeof(tinygoStructField{})
	println("DEBUG getField: i =", i, "fieldSize =", fieldSize)

	offset := uintptr(i) * fieldSize
	tinyField := (*tinygoStructField)(unsafe.Add(unsafe.Pointer(&st.Fields[0]), offset))

	// Convert to our StructField format
	return tinyField.toStructField()
}

// numFields returns the number of fields in the struct.
// For TinyGo, we use the numField uint16 field.
func (st *StructType) numFields() int {
	return int(st.numField)
}
