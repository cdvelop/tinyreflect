package tinyreflect

// A StructField describes a single field in a struct.
type StructField struct {
	Name Name    // name is always non-empty
	Typ  *Type   // type of field
	Off  uintptr // offset within struct, in bytes
}

func (f *StructField) Embedded() bool {
	return f.Name.IsEmbedded()
}

// Tag returns the field's tag as a StructTag.
func (f StructField) Tag() StructTag {
	return StructTag(f.Name.Tag())
}
