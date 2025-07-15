package tinyreflect

// SetString sets the string value to the field represented by Value.
// It uses unsafe to write the value to the memory location of the field.
// It's the caller's responsibility to ensure that the Value represents a string field.
func (v Value) SetString(x string) {
	// v.ptr is a pointer to the field. We cast it to a *string
	// and then dereference it to set the value.
	*(*string)(v.ptr) = x
}
