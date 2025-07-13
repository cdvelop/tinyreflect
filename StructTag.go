package tinyreflect

import (
	. "github.com/cdvelop/tinystring"
)

// StructTag is the tag string in a struct field (similar to reflect.StructTag)
type StructTag string

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
func (tag StructTag) Get(key string) string {
	out, _ := Convert(string(tag)).TagValue(key)
	return out
}
