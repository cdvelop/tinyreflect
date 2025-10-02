//go:build !tinygo

package tinyreflect

// StructType represents a struct type.
// Layout matches stdlib's abi.StructType for compatibility.
type StructType struct {
	Type
	PkgPath Name
	Fields  []StructField
}

// getField returns a pointer to the field at index i.
// For stdlib, Fields is a slice so we access it directly.
func (st *StructType) getField(i int) *StructField {
	if i < 0 || i >= len(st.Fields) {
		return nil
	}
	return &st.Fields[i]
}

// numFields returns the number of fields in the struct.
// For stdlib, we use len(Fields).
func (st *StructType) numFields() int {
	return len(st.Fields)
}
