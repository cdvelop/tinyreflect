package tinyreflect

import "unsafe"

// ABI types consolidated from internal/abi for TinyString JSON functionality
// This eliminates duplication with convert.go types and provides single source of truth

// kind represents the specific kind of type that a Type represents (private)
// Unified with convert.go kind, using tp prefix for TinyString naming convention
type kind uint8

const (
	KInvalid kind = iota
	KBool
	KInt
	KInt8
	KInt16
	KInt32
	KInt64
	KUint
	KUint8
	KUint16
	KUint32
	KUint64
	KUintptr
	KFloat32
	KFloat64
	KComplex64
	KComplex128
	KArray
	KChan
	KFunc
	KInterface
	KMap
	KPointer
	KSlice
	KString
	KStruct
	KUnsafePtr

	// TinyString specific types - separate values to avoid conflicts
	KSlice   // String slice type (separate from KSlice)
	KPointer // String pointer type (separate from KPointer)
	KErr     // Error type (separate from KInvalid)
)

const (
	kindDirectIface = 1 << 5
	kindGCProg      = 1 << 6 // Type.gc points to GC program
	kindMask        = (1 << 5) - 1
)

// String returns the name of k
func (k kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return kindNames[0]
}

var kindNames = []string{
	KInvalid:    "invalid",
	KBool:       "bool",
	KInt:        "int",
	KInt8:       "int8",
	KInt16:      "int16",
	KInt32:      "int32",
	KInt64:      "int64",
	KUint:       "uint",
	KUint8:      "uint8",
	KUint16:     "uint16",
	KUint32:     "uint32",
	KUint64:     "uint64",
	KUintptr:    "uintptr",
	KFloat32:    "float32",
	KFloat64:    "float64",
	KComplex64:  "complex64",
	KComplex128: "complex128",
	KArray:      "array",
	KChan:       "chan",
	KFunc:       "func",
	KInterface:  "interface",
	KMap:        "map",
	KPointer:    "ptr",
	KSlice:      "slice",
	KString:     "string",
	KStruct:     "struct",
	KUnsafePtr:  "unsafe.Pointer",
}

// tFlag is used by a Type to signal what extra type information is available
type tFlag uint8

// nameOff is the offset to a name from moduledata.types
type nameOff int32

// typeOff is the offset to a type from moduledata.types
type typeOff int32

// refType is the runtime representation of a Go type (adapted from internal/abi)
// Used for JSON struct inspection and field access
type refType struct {
	size       uintptr
	ptrBytes   uintptr // number of (prefix) bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tFlag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcData    *byte
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}

// refKind returns the Kind for this type (private version)
func (t *refType) refKind() kind {
	return t.Kind() // Delegate to the existing Kind() method
}

// refPtrType represents a pointer type
type refPtrType struct {
	refType
	elem *refType // pointer element (pointed at) type
}

// refFieldType contains information about a struct field for JSON operations
type refFieldType struct {
	name    string       // original field name (e.g., "BirthDate")
	refType *refType     // type of the field
	offset  uintptr      // byte offset in the struct
	index   int          // field index
	tag     refStructTag // field tag string (e.g., `json:"birth_date"`)
}

// refFieldMeta represents the original ABI field structure with refName
type refFieldMeta struct {
	name   refName  // encoded field name with tag info
	typ    *refType // type of the field
	offset uintptr  // byte offset in the struct
}

// refStructTag is the tag string in a struct field (similar to reflect.StructTag)
type refStructTag string

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
func (tag refStructTag) Get(key string) string {
	value, _ := tag.Lookup(key)
	return value
}

// Lookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string.
func (tag refStructTag) Lookup(key string) (value string, ok bool) {
	// Simplified implementation based on Go's reflect.StructTag
	for tag != "" {
		// Skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			// Unquote the value
			if len(qvalue) >= 2 && qvalue[0] == '"' && qvalue[len(qvalue)-1] == '"' {
				value = qvalue[1 : len(qvalue)-1]
				// Simple unescape for basic cases
				result := ""
				for j := 0; j < len(value); j++ {
					if value[j] == '\\' && j+1 < len(value) {
						switch value[j+1] {
						case 'n':
							result += "\n"
						case 't':
							result += "\t"
						case 'r':
							result += "\r"
						case '\\':
							result += "\\"
						case '"':
							result += "\""
						default:
							result += string(value[j])
							continue
						}
						j++ // skip the escaped character
					} else {
						result += string(value[j])
					}
				}
				return result, true
			}
			return qvalue, true
		}
	}
	return "", false
}

// refStructMeta represents a struct type with runtime metadata
type refStructMeta struct {
	refType
	pkgPath refName
	fields  []refFieldMeta
}

// refStructType contains cached information about a struct type for JSON operations
type refStructType struct {
	name    string         // name of the type
	refType *refType       // reference to the type information
	fields  []refFieldType // cached field information
}

// refSliceType represents a slice type
type refSliceType struct {
	refType
	elem *refType // slice element type
}

// refName is an encoded type name with optional extra data
type refName struct {
	bytes *byte
}

// Kind returns the kind of type
func (t *refType) Kind() kind {
	return kind(t.kind & kindMask)
}

// Size returns the size of data with type t
func (t *refType) Size() uintptr {
	return t.size
}

// Elem returns the element type for pointer, array, channel, map, or slice types
func (t *refType) Elem() *refType {
	switch t.Kind() {
	case KPointer:
		pt := (*refPtrType)(unsafe.Pointer(t))
		return pt.elem
	case KArray:
		at := (*refArrayType)(unsafe.Pointer(t))
		return at.elem
	case KSlice:
		st := (*refSliceType)(unsafe.Pointer(t))
		return st.elem
	// Add other cases as needed
	default:
		return nil
	}
}

// refElem returns the element type for pointer, array, channel, map, or slice types (private version)
func (t *refType) refElem() *refType {
	return t.Elem() // Delegate to the existing Elem() method
}

// refArrayType represents an array type
type refArrayType struct {
	refType
	elem *refType // array element type
	len  uintptr
}

// NumField returns the number of fields in a struct meta
func (t *refStructMeta) NumField() int {
	return len(t.fields)
}

// Field returns the i'th field of the struct meta
func (t *refStructMeta) Field(i int) *refFieldMeta {
	if i < 0 || i >= len(t.fields) {
		panic("reflect: Field index out of range")
	}
	return &t.fields[i]
}

// Name returns the name string for refName
func (n refName) Name() string {
	if n.bytes == nil {
		return ""
	}
	i, l := n.readVarint(1)
	return unsafe.String(n.dataChecked(1+i, "non-empty string"), l)
}

// Tag returns the tag string associated with the name
func (n refName) Tag() string {
	if n.bytes == nil {
		return ""
	}
	// Tags are typically stored after the name data
	// This is a simplified implementation - in the real Go runtime,
	// tags are stored with more complex encoding
	i, l := n.readVarint(1)
	if l == 0 {
		return ""
	}
	// Skip the name string
	nameStart := 1 + i
	nameEnd := nameStart + l

	// Check if there's tag data after the name
	if nameEnd < 100 { // Safety check to prevent reading too far
		tagI, tagL := n.readVarint(nameEnd)
		if tagL > 0 {
			return unsafe.String(n.dataChecked(nameEnd+tagI, "tag string"), tagL)
		}
	}
	return ""
}

// IsExported returns whether the name is exported
func (n refName) IsExported() bool {
	return (*n.bytes)&(1<<0) != 0
}

// readVarint parses a varint as encoded by encoding/binary
func (n refName) readVarint(off int) (int, int) {
	v := 0
	for i := 0; ; i++ {
		x := *n.dataChecked(off+i, "read varint")
		v += int(x&0x7f) << (7 * i)
		if x&0x80 == 0 {
			return i + 1, v
		}
	}
}

// dataChecked does pointer arithmetic on n's bytes
func (n refName) dataChecked(off int, whySafe string) *byte {
	return (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(n.bytes)) + uintptr(off)))
}

// clearObjectCache clears the global object cache - useful for testing
func clearObjectCache() {
	// This function is deprecated, use clearRefStructsCache in reflect.go instead
}
