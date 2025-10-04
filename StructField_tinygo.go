//go:build tinygo

package tinyreflect

import "unsafe"

// TinyGo's actual internal structField layout from internal/reflectlite/type.go
// This matches the exact layout used by TinyGo's reflect implementation
type tinygoStructField struct {
	fieldType *Type          // the type of the field
	data      unsafe.Pointer // pointer to packed data: flags, offset, name, tag
}

// Constants for the flags byte (from TinyGo's internal/reflectlite/type.go)
const (
	structFieldFlagAnonymous = 1 << iota
	structFieldFlagHasTag
	structFieldFlagIsExported
	structFieldFlagIsEmbedded
)

// Convert TinyGo's internal field to our StructField
// This parses the packed data format used by TinyGo
func (f *tinygoStructField) toStructField() *StructField {
	if f == nil || f.data == nil {
		return nil
	}

	data := f.data

	// Read flags byte
	flagsByte := *(*byte)(data)
	data = unsafe.Add(data, 1)

	// Read offset (uvarint32)
	offset, lenOffs := uvarint32(unsafe.Slice((*byte)(data), maxVarintLen32))
	data = unsafe.Add(data, lenOffs)

	// Read name (null-terminated string)
	namePtr := (*byte)(data)
	// Calculate name length by finding null terminator
	nameLen := 0
	for {
		b := *(*byte)(unsafe.Add(data, nameLen))
		if b == 0 {
			break
		}
		nameLen++
		if nameLen > 1000 { // safety check
			break
		}
	}

	return &StructField{
		Name: Name{Bytes: namePtr},
		Typ:  f.fieldType,
		Off:  uintptr(offset),
	}
}

// uvarint32 decodes a uint32 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0.
// This is compatible with TinyGo's implementation
func uvarint32(buf []byte) (uint32, int) {
	var x uint32
	var s uint
	for i, b := range buf {
		if i == maxVarintLen32 {
			return 0, -(i + 1) // overflow
		}
		if b < 0x80 {
			if i == maxVarintLen32-1 && b > 1 {
				return 0, -(i + 1) // overflow
			}
			return x | uint32(b)<<s, i + 1
		}
		x |= uint32(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

const maxVarintLen32 = 5 // maximum length of a varint32
