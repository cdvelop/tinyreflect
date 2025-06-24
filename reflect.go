package tinyreflect

import (
	"sync"
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Minimal reflectlite integration for TinyString JSON functionality
// This file contains essential reflection capabilities adapted from internal/reflectlite
// All functions and types are prefixed with 'ref' to avoid API pollution

// Import unified types from abi.go - no more duplication
// kind is now defined in abi.go with tp prefix

type refValue struct {
	// PRIMARY: Reflection fields integrated from refValue
	typ  *refType       // Reflection type information
	ptr  unsafe.Pointer // Pointer to the actual data
	flag refFlag        // Reflection flags for memory layout

	// ESSENTIAL: Core operation fields only
	kind         kind      // Type cache for performance (redundant with flag but kept for compatibility)
	roundDown    bool      // Operation flags
	separator    string    // String operations
	tmpStr       string    // String cache for performance
	lastConvType kind      // Cache validation
	err          errorType // Error handling

	// SPECIAL CASES: Complex types that need direct storage
	stringSliceVal []string // Slice operations
	stringPtrVal   *string  // Pointer operations
}

// refEface is the header for an interface{} value
type refEface struct {
	typ  *refType
	data unsafe.Pointer
}

// refFlag holds metadata about the value
type refFlag uintptr

const (
	flagKindWidth           = 5 // there are 27 kinds
	flagKindMask    refFlag = 1<<flagKindWidth - 1
	flagStickyRO    refFlag = 1 << 5
	flagEmbedRO     refFlag = 1 << 6
	flagIndir       refFlag = 1 << 7
	flagAddr        refFlag = 1 << 8
	flagMethod      refFlag = 1 << 9
	flagMethodShift         = 10
	flagKindShift           = flagMethodShift + 10 // room for method index
	flagRO          refFlag = flagStickyRO | flagEmbedRO
)

// refValueOf returns a new refValue initialized to the concrete value stored in i
// This replaces the old refValue-based function
func refValueOf(i any) *refValue {
	c := &refValue{separator: "_"}
	if i == nil {
		return c
	}
	c.initFromValue(i)
	return c
}

// ifaceIndir reports whether t is stored indirectly in an interface value
func ifaceIndir(t *refType) bool {
	return t.kind&kindDirectIface == 0
}

// Type returns the type of v
func (c *refValue) Type() *refType {
	return c.typ
}

// refElem returns the value that the interface c contains or that the pointer c points to
func (c *refValue) refElem() *refValue {
	k := c.refKind()
	switch k {
	case KInterface:
		var eface refEface
		if c.typ.kind&kindDirectIface != 0 {
			eface = refEface{typ: nil, data: c.ptr}
		} else {
			eface = *(*refEface)(c.ptr)
		}
		if eface.typ == nil {
			return &refValue{}
		}
		result := &refValue{separator: "_"}
		result.typ = eface.typ
		result.ptr = eface.data
		result.flag = refFlag(eface.typ.Kind())
		if ifaceIndir(eface.typ) {
			result.flag |= flagIndir
		}
		return result
	case KPointer:
		// Handle pointer dereferencing
		var ptr unsafe.Pointer
		if c.flag&flagIndir != 0 {
			// This is a pointer field from a struct - need to dereference to get the actual pointer
			ptr = *(*unsafe.Pointer)(c.ptr)
		} else {
			// This is a direct pointer from interface{}
			// c.ptr contains the pointer value itself (the address it points to)
			ptr = c.ptr
		}

		if ptr == nil {
			return &refValue{}
		}

		elemType := c.typ.Elem()
		if elemType == nil {
			return &refValue{}
		}

		// Create proper flags for the element
		// The element is addressable since we're dereferencing a pointer
		fl := c.flag&flagRO | flagAddr | refFlag(elemType.Kind())

		// For elements accessed through pointers, we don't need flagIndir
		// because ptr already points to the actual data
		result := &refValue{separator: "_"}
		result.typ = elemType
		result.ptr = ptr
		result.flag = fl
		return result
	}
	panic("reflect: call of reflect.Value.Elem on " + c.Type().Kind().String() + " value")
}

// refNumField returns the number of fields in the struct c
func (c *refValue) refNumField() int {
	c.mustBe(KStruct)
	tt := (*refStructMeta)(unsafe.Pointer(c.typ))
	return len(tt.fields)
}

// refField returns the i'th field of the struct c
func (c *refValue) refField(i int) *refValue {
	if c.refKind() != KStruct {
		panic("reflect: call of reflect.Value.Field on " + c.refKind().String() + " value")
	}
	tt := (*refStructMeta)(unsafe.Pointer(c.typ))
	if uint(i) >= uint(len(tt.fields)) {
		panic("reflect: Field index out of range")
	}
	field := &tt.fields[i]
	ptr := add(c.ptr, field.offset, "same as non-reflect &v.field")
	// Inherit read-only flags from parent, but allow assignment if parent allows it
	fl := c.flag&(flagRO) | refFlag(field.typ.Kind()) | flagAddr
	// For struct fields, flagIndir is needed only for pointer fields
	// because ptr points to the field location containing the pointer.
	// For other field types, ptr points directly to the field value.
	if field.typ.Kind() == KPointer {
		fl |= flagIndir
	}

	result := &refValue{separator: "_"}
	result.typ = field.typ
	result.ptr = ptr
	result.flag = fl
	return result
}

// refSetString sets c's underlying value to x
func (c *refValue) refSetString(x string) {
	c.mustBeAssignable()
	c.mustBe(KString)
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	*(*string)(ptr) = x
}

// refSetInt sets c's underlying value to x
func (c *refValue) refSetInt(x int64) {
	c.mustBeAssignable()
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch c.refKind() {
	case KInt:
		*(*int)(ptr) = int(x)
	case KInt8:
		*(*int8)(ptr) = int8(x)
	case KInt16:
		*(*int16)(ptr) = int16(x)
	case KInt32:
		*(*int32)(ptr) = int32(x)
	case KInt64:
		*(*int64)(ptr) = x
	default:
		panic("reflect: call of reflect.Value.SetInt on " + c.refKind().String() + " value")
	}
}

// refSetUint sets c's underlying value to x
func (c *refValue) refSetUint(x uint64) {
	c.mustBeAssignable()
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch c.refKind() {
	case KUint:
		*(*uint)(ptr) = uint(x)
	case KUint8:
		*(*uint8)(ptr) = uint8(x)
	case KUint16:
		*(*uint16)(ptr) = uint16(x)
	case KUint32:
		*(*uint32)(ptr) = uint32(x)
	case KUint64:
		*(*uint64)(ptr) = x
	case KUintptr:
		*(*uintptr)(ptr) = uintptr(x)
	default:
		c.err = errorType("reflect: call of reflect.Value.SetUint on " + c.refKind().String() + " value")
	}
}

// refSetFloat sets c's underlying value to x
func (c *refValue) refSetFloat(x float64) {
	c.mustBeAssignable()
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	switch c.refKind() {
	case KFloat32:
		*(*float32)(ptr) = float32(x)
	case KFloat64:
		*(*float64)(ptr) = x
	default:
		c.err = errorType("reflect: call of reflect.Value.SetFloat on " + c.refKind().String() + " value")
	}
}

// refSetBool sets c's underlying value to x
func (c *refValue) refSetBool(x bool) {
	c.mustBeAssignable()
	c.mustBe(KBool)
	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	*(*bool)(ptr) = x
}

// refSet assigns x to the value c
// c must be addressable and must not have been obtained by accessing unexported struct fields
func (c *refValue) refSet(x *refValue) {
	c.mustBeAssignable()
	if c.err != "" {
		return
	}
	x.mustBeExported() // do not let unexported x leak
	if x.err != "" {
		c.err = x.err
		return
	}

	// For pointer types, we need to copy the pointer value itself
	if c.refKind() == KPointer {
		// c.ptr points to the pointer variable
		// We need to set the pointer variable to the value that x represents
		if x.refKind() == KPointer {
			// Copy pointer value from x to c
			*(*unsafe.Pointer)(c.ptr) = *(*unsafe.Pointer)(x.ptr)
		} else {
			// x is not a pointer, this shouldn't happen in normal cases
			typedmemmove(c.typ, c.ptr, x.ptr)
		}
	} else {
		// For non-pointer types, copy the value
		typedmemmove(c.typ, c.ptr, x.ptr)
	}
}

// refZero returns a refValue representing the zero value for the specified type
func refZero(typ *refType) *refValue {
	if typ == nil {
		return &refValue{err: errorType("reflect: Zero(nil)")}
	}

	c := &refValue{separator: "_"}

	// For pointer types, zero value is nil pointer
	if typ.Kind() == KPointer {
		var nilPtr unsafe.Pointer // This is nil
		c.typ = typ
		c.ptr = unsafe.Pointer(&nilPtr)
		c.flag = refFlag(KPointer)
		return c
	}

	// For struct and other types, allocate memory for the zero value
	size := typ.Size()

	// Safety check: prevent huge allocations that could cause out of memory
	const maxSafeSize = 1024 * 1024 // 1MB limit
	if size > maxSafeSize {
		// For very large types, use a fixed small buffer
		size = 512
	}

	ptr := unsafe.Pointer(&make([]byte, size)[0])

	// Zero out the memory
	memclr(ptr, size)

	// Return the zero value with correct type and kind
	c.typ = typ
	c.ptr = ptr
	c.flag = refFlag(typ.Kind()) | flagAddr

	return c
}

// mustBeExported sets error if c was obtained using an unexported field
func (c *refValue) mustBeExported() {
	if c.err != "" {
		return
	}
	if c.flag&flagRO != 0 {
		c.err = errorType("reflect: use of unexported field")
	}
}

// mustBeAssignable sets error if c is not assignable
func (c *refValue) mustBeAssignable() {
	if c.err != "" {
		return
	}
	if c.flag&flagRO != 0 {
		c.err = errorType("reflect: cannot set value")
		return
	}
	if c.flag&flagAddr == 0 {
		c.err = errorType("reflect: cannot assign to value")
		return
	}
}

// mustBe sets error if c's kind is not expected
func (c *refValue) mustBe(expected kind) {
	if c.err != "" {
		return
	}
	if c.refKind() != expected {
		c.err = errorType("reflect: call of reflect.Value method on " + expected.String() + " value")
	}
}

// refKind returns the Kind without the flags
func (c *refValue) refKind() kind {
	return kind(c.flag & flagKindMask)
}

// typedmemmove copies a value of type t to dst from src
func typedmemmove(t *refType, dst, src unsafe.Pointer) {
	// Simplified version - just copy the bytes
	// This should use the actual Go runtime typedmemmove for safety
	// but for our purposes, a simple memory copy works
	memmove(dst, src, t.size)
}

// memmove copies n bytes from src to dst
func memmove(dst, src unsafe.Pointer, size uintptr) {
	// Simplified byte-by-byte copy
	// In real implementation, this would use runtime memmove
	dstBytes := (*[1 << 30]byte)(dst)
	srcBytes := (*[1 << 30]byte)(src)
	for i := uintptr(0); i < size; i++ {
		dstBytes[i] = srcBytes[i]
	}
}

// refIsValid reports whether c represents a value
func (c *refValue) refIsValid() bool {
	return c.flag != 0
}

// refInt returns c's underlying value, as an int64
func (c *refValue) refInt() int64 {
	if c.err != "" {
		return 0
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := c.refKind(); k {
	case KInt:
		return int64(*(*int)(ptr))
	case KInt8:
		return int64(*(*int8)(ptr))
	case KInt16:
		return int64(*(*int16)(ptr))
	case KInt32:
		return int64(*(*int32)(ptr))
	case KInt64:
		return *(*int64)(ptr)
	default:
		c.err = errorType("reflect: call of reflect.Value.Int on " + c.refKind().String() + " value")
		return 0
	}
}

// refUint returns c's underlying value, as a uint64
func (c *refValue) refUint() uint64 {
	if c.err != "" {
		return 0
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := c.refKind(); k {
	case KUint:
		return uint64(*(*uint)(ptr))
	case KUint8:
		return uint64(*(*uint8)(ptr))
	case KUint16:
		return uint64(*(*uint16)(ptr))
	case KUint32:
		return uint64(*(*uint32)(ptr))
	case KUint64:
		return *(*uint64)(ptr)
	case KUintptr:
		return uint64(*(*uintptr)(ptr))
	default:
		c.err = errorType("reflect: call of reflect.Value.Uint on " + c.refKind().String() + " value")
		return 0
	}
}

// refFloat returns c's underlying value, as a float64
func (c *refValue) refFloat() float64 {
	if c.err != "" {
		return 0
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}

	switch k := c.refKind(); k {
	case KFloat32:
		return float64(*(*float32)(ptr))
	case KFloat64:
		return *(*float64)(ptr)
	default:
		c.err = errorType("reflect: call of reflect.Value.Float on " + c.refKind().String() + " value")
		return 0
	}
}

// refBool returns c's underlying value
func (c *refValue) refBool() bool {
	if c.err != "" {
		return false
	}

	c.mustBe(KBool)
	if c.err != "" {
		return false
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	return *(*bool)(ptr)
}

// refString returns c's underlying value, as a string
func (c *refValue) refString() string {
	if c.err != "" {
		return ""
	}

	if !c.refIsValid() {
		return ""
	}

	// Don't enforce mustBe() - allow reading strings from struct fields
	if c.refKind() != KString {
		return ""
	}

	ptr := c.ptr
	if c.flag&flagIndir != 0 {
		ptr = *(*unsafe.Pointer)(ptr)
	}
	return *(*string)(ptr)
}

// Interface returns c's current value as an interface{}
func (c *refValue) Interface() any {
	if c.err != "" {
		return nil
	}

	if !c.refIsValid() {
		return nil
	}

	switch c.refKind() {
	case KString:
		return c.refString()
	case KInt:
		return int(c.refInt())
	case KInt8:
		return int8(c.refInt())
	case KInt16:
		return int16(c.refInt())
	case KInt32:
		return int32(c.refInt())
	case KInt64:
		return c.refInt()
	case KUint:
		return uint(c.refUint())
	case KUint8:
		return uint8(c.refUint())
	case KUint16:
		return uint16(c.refUint())
	case KUint32:
		return uint32(c.refUint())
	case KUint64:
		return c.refUint()
	case KUintptr:
		return uintptr(c.refUint())
	case KFloat32:
		return float32(c.refFloat())
	case KFloat64:
		return c.refFloat()
	case KBool:
		return c.refBool()
	case KInterface:
		// For interface{} types, extract the contained value directly
		var eface refEface
		if c.typ.kind&kindDirectIface != 0 {
			eface = refEface{typ: nil, data: c.ptr}
		} else {
			eface = *(*refEface)(c.ptr)
		}
		if eface.typ == nil {
			return nil
		}

		// Create a new interface{} with the contained value
		return *(*any)(unsafe.Pointer(&eface))
	case KStruct: // For struct types, create an interface{} with the struct value
		// The struct data is stored at c.ptr
		var eface refEface
		eface.typ = c.typ
		eface.data = c.ptr
		return *(*any)(unsafe.Pointer(&eface))
	default:
		// For complex types, return nil for now
		return nil
	}
}

// add returns p+x
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// Global cache for struct type information
// Using slice instead of map for TinyGo compatibility
var refStructsTypes []refStructType
var refStructsTypesMutex sync.RWMutex

// getStructType fills struct information if not cached, assigns to provided pointer
func getStructType(typ *refType, out *refStructType) {
	if typ.Kind() != KStruct {
		return
	}

	// Get unique type name for caching
	typeName := getTypeName(typ)

	// First try read lock to check cache
	refStructsTypesMutex.RLock()
	for i := range refStructsTypes {
		if refStructsTypes[i].name == typeName {
			*out = refStructsTypes[i]
			refStructsTypesMutex.RUnlock()
			return
		}
	}
	refStructsTypesMutex.RUnlock()

	// Not in cache, need write lock to add new entry
	refStructsTypesMutex.Lock()
	defer refStructsTypesMutex.Unlock()

	// Double-check pattern: another goroutine might have added it
	for i := range refStructsTypes {
		if refStructsTypes[i].name == typeName {
			*out = refStructsTypes[i]
			return
		}
	}

	// Create new struct info
	structType := (*refStructMeta)(unsafe.Pointer(typ))
	fields := make([]refFieldType, len(structType.fields))
	for i, f := range structType.fields {
		fieldName := f.name.Name()
		fieldTag := f.name.Tag() // Get the tag string
		fields[i] = refFieldType{
			name:    fieldName,
			refType: f.typ,
			offset:  f.offset,
			index:   i,
			tag:     refStructTag(fieldTag),
		}
	}

	// Create new struct info
	newInfo := refStructType{
		name:    typeName,
		refType: typ,
		fields:  fields,
	}

	// Add to cache
	refStructsTypes = append(refStructsTypes, newInfo)

	// Assign to output
	*out = newInfo
}

// clearRefStructsCache clears the global struct cache - useful for testing
func clearRefStructsCache() {
	refStructsTypesMutex.Lock()
	defer refStructsTypesMutex.Unlock()
	refStructsTypes = refStructsTypes[:0] // Clear slice while preserving capacity
}

func getTypeName(typ *refType) string {
	if typ == nil {
		return "nil"
	}

	// Use type pointer and size to create unique identifier
	// Convert uintptr to string manually since Convert() doesn't handle uintptr
	ptr := uintptr(unsafe.Pointer(typ))
	ptrStr := ""
	if ptr != 0 {
		// Convert uintptr to base-10 string manually
		temp := ptr
		if temp == 0 {
			ptrStr = "0"
		} else {
			digits := ""
			for temp > 0 {
				digit := temp % 10
				digits = string(rune('0'+digit)) + digits
				temp /= 10
			}
			ptrStr = digits
		}
	}

	sizeStr := Convert(int64(typ.size)).String()
	kindStr := typ.Kind().String()
	return kindStr + "_" + sizeStr + "_" + ptrStr
}

// refLen returns the length of c
// It panics if c's Kind is not Slice
func (c *refValue) refLen() int {
	if c.err != "" {
		return 0
	}
	k := c.refKind()
	switch k {
	case KSlice:
		// For slices, the length is stored in the slice header
		return (*sliceHeader)(c.ptr).Len
	default:
		c.err = errorType("reflect: call of reflect.Value.Len on " + k.String() + " value")
		return 0
	}
}

// refIndex returns c's i'th element
// It panics if c's Kind is not Slice or if i is out of range
func (c *refValue) refIndex(i int) *refValue {
	if c.err != "" {
		return &refValue{err: c.err}
	}
	k := c.refKind()
	switch k {
	case KSlice:
		s := (*sliceHeader)(c.ptr)
		if i < 0 || i >= s.Len {
			c.err = errorType("reflect: slice index out of range")
			return &refValue{err: c.err}
		}

		// Get element type
		elemType := c.typ.Elem()
		if elemType == nil {
			return &refValue{err: errorType("reflect: slice element type is nil")}
		}

		elemSize := elemType.Size()

		// Calculate pointer to element
		elemPtr := unsafe.Pointer(uintptr(s.Data) + uintptr(i)*elemSize)
		// Create new refValue for the element
		result := &refValue{separator: "_"}
		result.typ = elemType
		result.ptr = elemPtr
		result.flag = refFlag(elemType.Kind())

		// If element is stored indirectly, set the flag
		// Note: strings should never be indirect in slices
		if elemType.Kind() != KString && elemType.kind&kindDirectIface == 0 {
			result.flag |= flagIndir
		}

		return result
	default:
		c.err = errorType("reflect: call of reflect.Value.Index on " + k.String() + " value")
		return &refValue{err: c.err}
	}
}

// refMakeSlice creates a new slice with the given type, length, and capacity
func refMakeSlice(typ *refType, len, cap int) *refValue {
	if typ.Kind() != KSlice {
		panic("refMakeSlice called on non-slice type")
	}

	elemType := typ.Elem()
	elemSize := elemType.Size()

	// Allocate memory for the slice data
	var dataPtr unsafe.Pointer
	if cap > 0 {
		dataPtr = mallocSliceData(elemSize, cap)
		// Initialize the memory to zero
		memclr(dataPtr, uintptr(cap)*elemSize)
	}

	// Create slice header
	slice := &sliceHeader{
		Data: dataPtr,
		Len:  len,
		Cap:  cap,
	}

	// Create a refValue pointing to the slice header
	return &refValue{
		separator: "_",
		typ:       typ,
		ptr:       unsafe.Pointer(slice),
		flag:      refFlag(typ.Kind()),
	}
}

// Memory management utilities

// sliceHeader represents the header of a slice for low-level operations
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// memclr clears memory at the given pointer
func memclr(ptr unsafe.Pointer, size uintptr) {
	if ptr == nil || size == 0 {
		return
	}
	slice := (*[1 << 30]byte)(ptr)[:size:size]
	for i := range slice {
		slice[i] = 0
	}
}
