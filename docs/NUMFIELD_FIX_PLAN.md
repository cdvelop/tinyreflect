# NumField Fix Plan for TinyGo Compatibility

## Problem Analysis

When compiling with standard Go (stdlib), `NumField()` returns 5 fields correctly.
When compiling with TinyGo, `NumField()` returns 0 fields.

### Stdlib vs TinyGo Comparison

#### Standard Go (internal/abi/type.go)
```go
// /usr/local/go/src/internal/abi/type.go
type StructType struct {
	Type
	PkgPath Name
	Fields  []StructField  // Slice of fields
}

type StructField struct {
	Name   Name    // name is always non-empty
	Typ    *Type   // type of field
	Offset uintptr // byte offset of field
}
```

**Stdlib reflectlite NumField implementation:**
```go
// /usr/local/go/src/internal/reflectlite/type.go:341
func (t rtype) NumField() int {
	tt := t.Type.StructType()
	if tt == nil {
		panic("reflect: NumField of non-struct type")
	}
	return len(tt.Fields)  // Uses slice length
}
```

#### TinyGo (internal/reflectlite/type.go)
```go
// /usr/local/lib/tinygo/src/internal/reflectlite/type.go:219
type structType struct {
	RawType
	numMethod uint16
	ptrTo     *RawType
	pkgpath   *byte
	size      uint32
	numField  uint16      // Field count as uint16!
	fields    [1]structField  // Array that extends in memory
}

type structField struct {
	fieldType *RawType
	data      unsafe.Pointer // various bits of information, packed in a byte array
}
```

**TinyGo reflectlite NumField implementation:**
```go
// /usr/local/lib/tinygo/src/internal/reflectlite/type.go:603
func (t *RawType) NumField() int {
	if t.Kind() != Struct {
		panic(errTypeNumField)
	}
	return int((*structType)(unsafe.Pointer(t.underlying())).numField)  // Uses uint16 field
}
```

### Root Cause

**Current tinyreflect StructType:**
```go
type StructType struct {
	Type
	PkgPath Name
	Fields  []StructField  // Matches stdlib but WRONG memory layout for TinyGo
}
```

**Key Problems:**
1. **Stdlib uses slice:** `Fields []StructField` with `len(Fields)` 
2. **TinyGo uses fixed layout:** `numField uint16` + array `[1]structField` extending in memory
3. **Memory incompatibility:** When TinyGo runtime creates type info, it uses its own layout
4. **Current tinyreflect:** Assumes stdlib layout, reads wrong memory location in TinyGo

### Decision: Which Implementation to Follow?

**We must follow TinyGo's layout** because:
1. **TinyGo is the constraint:** It's the more limited environment
2. **Size optimization:** TinyGo's layout is more compact (uint16 + array vs slice header)
3. **Direct memory mapping:** TinyGo's layout maps directly to compiled type info
4. **Stdlib can adapt:** We can make stdlib work with TinyGo's layout using unsafe casts

## Solution Strategy - MINIMAL CHANGES

**Goal:** Update internal struct layout to match TinyGo while keeping SAME API.

**Approach:** Use TinyGo's memory layout universally with build-tag specific field access:
- Single `StructType` definition that works for both
- Methods use `unsafe` pointer arithmetic internally when needed
- API remains unchanged - users see no difference

## Implementation Plan

### Step 1: Update StructType.go (Single File)
Replace current struct with TinyGo-compatible layout:

```go
type StructType struct {
	Type
	numMethod uint16
	ptrTo     *Type
	pkgpath   *byte
	size      uint32
	numField  uint16
	fields    [1]StructField
}
```

**Note:** This layout works for both stdlib and TinyGo because we'll access fields through helper functions.

### Step 2: Add Internal Helper Function
Add single helper function to access fields by index:

```go
// getField returns pointer to field at index i
// Uses unsafe arithmetic because fields extends beyond [1]
func (st *StructType) getField(i int) *StructField {
	if i < 0 || i >= int(st.numField) {
		return nil
	}
	// Calculate field address: &st.fields[0] + i*sizeof(StructField)
	return (*StructField)(unsafe.Add(unsafe.Pointer(&st.fields[0]), i*unsafe.Sizeof(StructField{})))
}
```

### Step 3: Update NumField in TypeOf.go
Change from `len(st.Fields)` to `int(st.numField)`:

```go
func (t *Type) NumField() (int, error) {
	if t.Kind() != K.Struct {
		return 0, Err(ref, D.Numbers, D.Fields, D.Type, "Struct")
	}
	st := (*StructType)(unsafe.Pointer(t))
	return int(st.numField), nil  // Changed from len(st.Fields)
}
```

### Step 4: Update NameByIndex in TypeOf.go
Use helper function instead of direct slice access:

```go
func (t *Type) NameByIndex(i int) (string, error) {
	if t.Kind() != K.Struct {
		return "", Err(ref, D.Type, D.NotOfType, "Struct")
	}
	tt := (*StructType)(unsafe.Pointer(t))
	if i < 0 || i >= int(tt.numField) {
		return "", Err(ref, D.Index, D.Out, D.Of, D.Range)
	}
	f := tt.getField(i)  // Use helper instead of tt.Fields[i]
	return f.Name.Name(), nil
}
```

### Step 5: Update Field() in ValueOf.go
Use helper function for field access:

```go
func (v Value) Field(i int) (Value, error) {
	if v.kind() != K.Struct {
		return Value{}, Err(ref, D.Value, D.NotOfType, "Struct")
	}
	tt := (*StructType)(unsafe.Pointer(v.typ()))
	if uint(i) >= uint(tt.numField) {
		return Value{}, Err(ref, D.Value, D.Index, D.Out, D.Of, D.Range)
	}
	field := tt.getField(i)  // Use helper instead of tt.Fields[i]
	// ... rest remains the same
}
```

### Step 6: Update NumField in ValueOf.go
Similar change:

```go
func (v Value) NumField() (int, error) {
	// ... validation code ...
	st := v.typ_.StructType()
	if st == nil {
		return 0, Err(ref, D.Numbers, D.Fields, D.NotOfType, "Struct")
	}
	return int(st.numField), nil  // Changed from len(st.Fields)
}
```

## Files to Modify (NO NEW FILES)

1. **StructType.go** - Update struct definition
2. **TypeOf.go** - Update NumField() and NameByIndex()
3. **ValueOf.go** - Update Field() and NumField()

**No file renames, no duplications, no build tags needed!**

## Why This Works for Both

1. **TinyGo:** Reads correct memory layout directly
2. **Stdlib:** When stdlib creates the type info, it populates the same memory layout
3. **Helper function:** Works universally using pointer arithmetic
4. **Size:** No code duplication = minimal binary size

## Verification Steps

1. Test with stdlib: `go test ./...`
2. Test with TinyGo: `tinygo test`
3. Build WASM with TinyGo: `tinygo build -o main.wasm -target=wasm example/pwa/main.wasm.go`
4. Verify PWA shows "Found 5 fields" in browser console
5. Check binary size remains under 50KB

## Success Criteria

- [ ] NumField() returns 5 for both stdlib and TinyGo
- [ ] Field iteration works correctly
- [ ] All existing tests pass (stdlib and TinyGo)
- [ ] PWA example displays correct field count and values
- [ ] Binary size stays under 50KB
- [ ] No API changes - existing code continues to work
- [ ] Single implementation - no code duplication


### Supported Types (Minimalist Approach)
TinyReflect **intentionally** supports only a minimal set of types to keep binary size small:

**✅ Supported Types:**
- **Basic types**: `string`, `bool`
- **All numeric types**: `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`
- **All basic slices**: `[]string`, `[]bool`, `[]byte`, `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`, `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`, `[]float32`, `[]float64`
- **Structs**: Only with supported field types
- **Struct slices**: `[]struct{...}` where all fields are supported types
- **Maps**: `map[K]V` where K and V are supported types only
- **Map slices**: `[]map[K]V` where K and V are supported types only
- **Pointers**: Only to supported types above

**❌ Unsupported Types:**
- `any`, `chan`, `func`
- `complex64`, `complex128`
- `uintptr`, `unsafe.Pointer` (used internally only)
- Arrays (different from slices)
- Nested complex types beyond supported scope

This focused approach ensures minimal code size while covering the most common JSON-like data operations including simple structs.