package tinyreflect

// StructNamer interface allows structs to provide their own name for reflection.
// This is required for TinyGo compatibility since runtime type name resolution
// is not available in TinyGo's limited reflection support.
//
// Types implementing this interface will have their StructName() method called
// during TypeOf() to cache the struct name for later retrieval via Type.Name().
//
// Example:
//
//	type User struct { Name string }
//	func (User) StructName() string { return "User" }
type StructNamer interface {
	StructName() string
}
