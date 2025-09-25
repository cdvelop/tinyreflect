package tinyreflect

import . "github.com/cdvelop/tinystring"

// structCacheEntry holds the cached schema of a struct.
// This design is optimized for TinyGo by using fixed-size arrays
// to avoid dynamic allocations and map usage, which are limited in TinyGo.
type structCacheEntry struct {
	structID   uint32 // Hash-based key for quick lookups.
	nameLen    uint8  // Length of the struct name.
	fieldCount uint8  // Number of fields in the struct.

	// Fixed-size array for the struct name. 32 bytes should be enough for most names.
	structName [32]byte

	// Fixed-size array for field schemas. 16 fields is a reasonable limit for many structs.
	fieldSchemas [16]fieldSchema
}

// fieldSchema holds the cached metadata for a single struct field.
// Like structCacheEntry, it uses fixed-size arrays for TinyGo compatibility.
type fieldSchema struct {
	nameLen uint8 // Length of the field name.
	kind    Kind  // The type kind of the field.
	offset  uint16 // The field's offset within the struct.
	typ     *Type  // The type of the field.

	// Fixed-size array for the field name. 20 bytes is a reasonable limit.
	fieldName [20]byte
}