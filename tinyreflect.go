package tinyreflect

import (
	"sync/atomic"
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

// Default cache sizes
const (
	StructSize128 = 128
	StructSize256 = 256
)

// sliceHeader is the runtime representation of a slice.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// stringHeader is the runtime representation of a string.
type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

// TinyReflect provides an instance-based reflection API with an internal cache
// optimized for TinyGo and WebAssembly environments.
type TinyReflect struct {
	structCache []structCacheEntry // The cache for struct schemas.
	structCount int32              // Atomic counter for the number of cached structs.
	cacheLock   int32              // Atomic lock for cache access.
	log         func(msgs ...any)  // Optional logger function.
	maxStructs  int32              // The maximum number of structs that can be cached.
}

// New creates a new TinyReflect instance with optional configuration.
// You can provide a cache size (e.g., StructSize256) and a logger function.
// If no arguments are provided, it defaults to a cache size of 128 and no logging.
func New(args ...any) *TinyReflect {
	tr := &TinyReflect{
		maxStructs: StructSize128, // Default cache size
		log:        func(...any) {}, // Default no-op logger
	}

	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			tr.maxStructs = int32(v)
		case func(...any):
			tr.log = v
		}
	}

	tr.structCache = make([]structCacheEntry, tr.maxStructs)

	return tr
}

// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func (tr *TinyReflect) TypeOf(i any) *Type {
	if i == nil {
		return nil
	}
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	typ := e.Type

	underlying := typ
	for underlying != nil && underlying.Kind().String() == "ptr" {
		underlying = underlying.Elem()
	}

	if underlying != nil && underlying.Kind().String() == "struct" {
		tr.cacheStructSchema(i, underlying)
	}

	return typ
}

// ValueOf returns a new Value initialized to the concrete value
// stored in the interface i. ValueOf(nil) returns the zero Value.
func (tr *TinyReflect) ValueOf(i any) Value {
	if i == nil {
		return Value{tr: tr}
	}
	return tr.unpackEface(i)
}

// unpackEface converts the empty interface i to a Value.
func (tr *TinyReflect) unpackEface(i any) Value {
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	t := e.Type
	if t == nil {
		return Value{tr: tr}
	}
	f := flag(t.Kind())
	if t.IfaceIndir() {
		f |= flagIndir
	}
	return Value{tr, t, e.Data, f}
}

// Indirect returns the value that v points to.
// If v is a nil pointer, Indirect returns a zero Value.
// If v is not a pointer, Indirect returns v.
func (tr *TinyReflect) Indirect(v Value) Value {
	if v.kind().String() != "ptr" {
		return v
	}
	elem, err := v.Elem()
	if err != nil {
		return Value{tr: tr}
	}
	return elem
}

// MakeSlice creates a new zero-initialized slice value
// for the specified slice type, length, and capacity.
func (tr *TinyReflect) MakeSlice(typ *Type, len, cap int) (Value, error) {
	if typ == nil {
		return Value{}, Err(D.Value, D.Type, D.Nil)
	}
	if typ.Kind().String() != "slice" {
		return Value{}, Err("MakeSlice of non-slice type")
	}
	if len < 0 || cap < 0 || len > cap {
		return Value{}, Err("invalid slice length or capacity")
	}

	sliceType := (*SliceType)(unsafe.Pointer(typ))
	elemType := sliceType.Elem
	if elemType == nil {
		return Value{}, Err("MakeSlice element type is nil")
	}

	var data unsafe.Pointer
	if elemType.Size != 0 {
		mem := make([]byte, uintptr(cap)*elemType.Size)
		data = unsafe.Pointer(&mem[0])
	}

	sliceHeader := &sliceHeader{Data: data, Len: len, Cap: cap}
	return Value{tr, typ, unsafe.Pointer(sliceHeader), flagIndir | flag(K.Slice)}, nil
}

// NewValue returns a Value representing a pointer to a new zero value
// for the specified type.
func (tr *TinyReflect) NewValue(typ *Type) Value {
	if typ == nil {
		return Value{tr: tr}
	}

	ptrType := PtrType{Type: Type{Kind_: K.Pointer, Size: unsafe.Sizeof(uintptr(0))}, Elem: typ}
	ptrToValue := make([]byte, typ.Size)
	ptr := make([]byte, unsafe.Sizeof(uintptr(0)))
	*(*unsafe.Pointer)(unsafe.Pointer(&ptr[0])) = unsafe.Pointer(&ptrToValue[0])

	return Value{tr, (*Type)(unsafe.Pointer(&ptrType)), unsafe.Pointer(&ptr[0]), flag(K.Pointer) | flagIndir}
}

func (tr *TinyReflect) cacheStructSchema(i any, typ *Type) {
	structID := typ.StructID()
	if structID == 0 {
		return // Not a struct or invalid type
	}

	count := atomic.LoadInt32(&tr.structCount)
	for j := int32(0); j < count; j++ {
		if tr.structCache[j].structID == structID {
			return // Already cached
		}
	}

	for !atomic.CompareAndSwapInt32(&tr.cacheLock, 0, 1) {
		// spin
	}
	defer atomic.StoreInt32(&tr.cacheLock, 0)

	count = atomic.LoadInt32(&tr.structCount)
	for j := int32(0); j < count; j++ {
		if tr.structCache[j].structID == structID {
			return
		}
	}

	if count >= tr.maxStructs {
		tr.log("tinyreflect: struct cache is full")
		return
	}

	var entry structCacheEntry
	entry.structID = structID

	// Use Indirect to correctly handle pointers when checking for StructNamer.
	v := tr.Indirect(tr.ValueOf(i))
	iface, err := v.Interface()
	if err != nil {
		iface = i // Fallback for safety
	}

	var structName string
	if sn, ok := iface.(StructNamer); ok {
		structName = sn.StructName()
	} else {
		structName = typ.Name()
	}
	entry.nameLen = uint8(copy(entry.structName[:], structName))

	numFields, _ := typ.NumField()
	if numFields > len(entry.fieldSchemas) {
		numFields = len(entry.fieldSchemas)
	}
	entry.fieldCount = uint8(numFields)

	st := typ.StructType()
	if st == nil {
		return
	}

	for k := 0; k < numFields; k++ {
		field := st.Fields[k]
		entry.fieldSchemas[k].nameLen = uint8(copy(entry.fieldSchemas[k].fieldName[:], field.Name.Name()))
		entry.fieldSchemas[k].kind = field.Typ.Kind()
		entry.fieldSchemas[k].offset = uint16(field.Off)
		entry.fieldSchemas[k].typ = field.Typ
	}

	newIndex := atomic.AddInt32(&tr.structCount, 1) - 1
	if newIndex < tr.maxStructs {
		tr.structCache[newIndex] = entry
		tr.log("tinyreflect: cached schema for struct", structName)
	}
}