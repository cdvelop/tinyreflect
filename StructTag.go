package tinyreflect

import (
	. "github.com/cdvelop/tinystring"
)

// StructTag is the tag string in a struct field (similar to reflect.StructTag)
type StructTag string

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
func (tag StructTag) Get(key string) string {
	value, _ := tag.Lookup(key)
	return value
}

// Lookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string.
func (tag StructTag) Lookup(key string) (value string, ok bool) {
	value, err := Convert(tag).KV(key)
	if err == nil {
		return value, true
	}
	return "", false
}
