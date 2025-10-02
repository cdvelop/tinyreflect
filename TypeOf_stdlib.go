//go:build !tinygo

package tinyreflect

// underlying returns the underlying type.
// In stdlib, types don't have the named/unnamed distinction
// the same way as TinyGo, so we just return self.
// The Type already points to the correct structure.
func (t *Type) underlying() *Type {
	return t
}
