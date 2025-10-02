//go:build tinygo

package tinyreflect

// TinyGo's internal structField layout
// Trying different field orders to match TinyGo's actual layout
type tinygoStructField struct {
	offsetEmbed uintptr // byte offset + embedded flag
	name        *byte   // name
	typ         *Type   // type
}

// Convert TinyGo's internal field to our StructField
func (f *tinygoStructField) toStructField() *StructField {
	return &StructField{
		Name: Name{Bytes: f.name},
		Typ:  f.typ,
		Off:  f.offsetEmbed &^ 1, // Remove embedded flag from offset
	}
}
