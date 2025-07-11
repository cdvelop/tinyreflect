package tinyreflect

type StructType struct {
	Type
	PkgPath Name
	Fields  []StructField
}
