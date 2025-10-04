# Plan: Unsafe Pointer Arithmetic for Field Access in TinyGo WASM

## Executive Summary

Based on the technical analysis in `TinyGo_WASM_Struct_Reflection.md`, we have confirmed:

✅ **Current Success:**
- `NumField()` returns correct count (5 fields)
- Field offsets are reliably preserved by TinyGo RTTI
- Struct layout metadata is intact

❌ **Current Failure:**
- Field names are corrupted/stripped by `-opt=z` optimization
- `Name.Name()` causes `RuntimeError: memory access out of bounds`
- Cannot rely on string pointers in RTTI

## Solution Architecture: Hybrid RTTI + Unsafe

The document recommends **abandoning high-level reflection** for field value access and implementing **direct memory manipulation** using:

1. **Reliable RTTI metadata:** `StructField.Offset` (preserved even with `-opt=z`)
2. **Unsafe pointer arithmetic:** Direct memory address calculation
3. **Type-specific handlers:** For primitive types, strings, slices

## Implementation Plan

### Phase 1: Core Unsafe Field Access (Primitive Types)

#### 1.1. Create New File: `ValueUnsafe_tinygo.go`

Build tag: `//go:build tinygo`

Implement core unsafe field access method:

```go
// filepath: /home/cesar/Dev/Pkg/Mine/tinyreflect/ValueUnsafe_tinygo.go
//go:build tinygo

package tinyreflect

import "unsafe"

// FieldByOffset returns a new Value representing the field at the given byte offset.
// This uses direct pointer arithmetic and is the reliable method for TinyGo WASM.
func (v Value) FieldByOffset(offset uintptr) Value {
    // Step 1: Get base address of struct
    baseAddr := uintptr(v.ptr_)
    
    // Step 2: Calculate field address
    fieldAddr := baseAddr + offset
    
    // Step 3: Convert back to unsafe.Pointer
    fieldPtr := unsafe.Pointer(fieldAddr)
    
    // Return new Value pointing to field
    return Value{
        ptr_: fieldPtr,
        typ_: nil, // Will be set by caller based on field type
    }
}

// UnsafeFieldValue reads a field value directly using pointer arithmetic.
// Returns the field as interface{} by type switching on the field's Kind.
func (v Value) UnsafeFieldValue(fieldIndex int) (interface{}, error) {
    ut := v.typ_.underlying()
    st := (*StructType)(unsafe.Pointer(ut))
    
    if fieldIndex < 0 || fieldIndex >= int(st.numField) {
        return nil, errors.New("field index out of bounds")
    }
    
    // Get field metadata
    field := st.getField(fieldIndex)
    offset := field.offsetEmbed & 0xFFFFFF // Lower 24 bits are offset
    
    // Get field type
    fieldType := field.typ
    
    // Calculate field address
    baseAddr := uintptr(v.ptr_)
    fieldAddr := unsafe.Pointer(baseAddr + uintptr(offset))
    
    // Type switch based on Kind
    switch fieldType.Kind() {
    case K.Bool:
        return *(*bool)(fieldAddr), nil
    case K.Int:
        return *(*int)(fieldAddr), nil
    case K.Int8:
        return *(*int8)(fieldAddr), nil
    case K.Int16:
        return *(*int16)(fieldAddr), nil
    case K.Int32:
        return *(*int32)(fieldAddr), nil
    case K.Int64:
        return *(*int64)(fieldAddr), nil
    case K.Uint:
        return *(*uint)(fieldAddr), nil
    case K.Uint8:
        return *(*uint8)(fieldAddr), nil
    case K.Uint16:
        return *(*uint16)(fieldAddr), nil
    case K.Uint32:
        return *(*uint32)(fieldAddr), nil
    case K.Uint64:
        return *(*uint64)(fieldAddr), nil
    case K.Float32:
        return *(*float32)(fieldAddr), nil
    case K.Float64:
        return *(*float64)(fieldAddr), nil
    case K.String:
        return unsafeReadString(fieldAddr), nil
    default:
        return nil, errors.New("unsupported field type")
    }
}
```

#### 1.2. Implement String Reading (TGC-5)

TinyGo uses `uintptr` for string length (not `int`):

```go
// TinyGo's internal string header (different from stdlib)
type tinygoStringHeader struct {
    Data uintptr  // Pointer to string data
    Len  uintptr  // Length in bytes (uintptr, not int!)
}

func unsafeReadString(fieldAddr unsafe.Pointer) string {
    // Cast to TinyGo's string header
    header := (*tinygoStringHeader)(fieldAddr)
    
    // Safety check: ensure pointer is not nil
    if header.Data == 0 || header.Len == 0 {
        return ""
    }
    
    // Create byte slice from raw memory
    dataPtr := unsafe.Pointer(header.Data)
    bytes := unsafe.Slice((*byte)(dataPtr), header.Len)
    
    // Convert to string
    return string(bytes)
}
```

#### 1.3. Update `Value.Field()` to Use Unsafe

Modify existing `Field()` method to use unsafe approach:

```go
// filepath: /home/cesar/Dev/Pkg/Mine/tinyreflect/ValueMethods_tinygo.go
//go:build tinygo

func (v Value) Field(i int) (Value, error) {
    ut := v.typ_.underlying()
    st := (*StructType)(unsafe.Pointer(ut))
    
    numFields := int(st.numField)
    if i < 0 || i >= numFields {
        return Value{}, errors.New("field index out of range")
    }
    
    // Get field using unsafe pointer arithmetic
    field := st.getField(i)
    offset := field.offsetEmbed & 0xFFFFFF
    
    baseAddr := uintptr(v.ptr_)
    fieldAddr := unsafe.Pointer(baseAddr + uintptr(offset))
    
    return Value{
        ptr_: fieldAddr,
        typ_: field.typ,
    }, nil
}
```

### Phase 2: Field Name Management (TGC-4)

Since RTTI string pointers are unreliable, implement external name mapping:

#### 2.1. Create Compile-Time Field Name Registry

```go
// filepath: /home/cesar/Dev/Pkg/Mine/tinyreflect/FieldNameRegistry.go

package tinyreflect

// FieldNames stores field names for a struct type, indexed by field position.
type FieldNames []string

// Global registry mapping type names to field name arrays
var fieldNameRegistry = make(map[string]FieldNames)

// RegisterFieldNames registers field names for a struct type at compile time.
// This should be called in init() functions or via code generation.
func RegisterFieldNames(typeName string, names FieldNames) {
    fieldNameRegistry[typeName] = names
}

// GetFieldName retrieves the name of field i for the given type.
// Returns empty string if not registered.
func (t *Type) GetFieldName(i int) string {
    typeName := t.Name()
    if typeName == "" {
        return ""
    }
    
    names, exists := fieldNameRegistry[typeName]
    if !exists || i < 0 || i >= len(names) {
        return ""
    }
    
    return names[i]
}
```

#### 2.2. Code Generation Tool (Optional)

Create `cmd/tinyreflect-gen/main.go` to auto-generate registrations:

```go
// Usage: go run cmd/tinyreflect-gen/main.go -type=TestStruct -package=main

// Output: generated_field_names.go
func init() {
    tinyreflect.RegisterFieldNames("main.TestStruct", tinyreflect.FieldNames{
        "StringField",
        "BoolField",
        "IntField",
        "Int8Field",
        "Int16Field",
    })
}
```

#### 2.3. Manual Registration (Immediate Solution)

For now, users can manually register in their code:

```go
// In user's main.go
func init() {
    tinyreflect.RegisterFieldNames("main.TestStruct", tinyreflect.FieldNames{
        "StringField",
        "BoolField",
        "IntField",
        "Int8Field",
        "Int16Field",
    })
}
```

### Phase 3: Advanced Type Support (TGC-5)

#### 3.1. Slice Reading

```go
type tinygoSliceHeader struct {
    Data uintptr
    Len  uintptr
    Cap  uintptr
}

func unsafeReadSlice(fieldAddr unsafe.Pointer, elemType *Type) interface{} {
    header := (*tinygoSliceHeader)(fieldAddr)
    
    if header.Data == 0 || header.Len == 0 {
        return nil // Empty slice
    }
    
    // Create slice based on element type
    // This requires type-specific handlers...
    // TODO: Implement based on elemType.Kind()
}
```

#### 3.2. Pointer Reading

```go
func unsafeReadPointer(fieldAddr unsafe.Pointer) unsafe.Pointer {
    // A pointer field contains the address of the pointed-to value
    return *(*unsafe.Pointer)(fieldAddr)
}
```

### Phase 4: Testing & Validation

#### 4.1. Update Test Code

Modify `example/pwa/main.wasm.go` to use new unsafe API:

```go
for i := 0; i < numFields; i++ {
    // Get field value using unsafe
    fieldValue, err := v.UnsafeFieldValue(i)
    if err != nil {
        logger("ERROR: Field", i, "error:", err)
        continue
    }
    
    // Get field name from registry
    fieldName := typ.GetFieldName(i)
    if fieldName == "" {
        fieldName = fmt.Sprintf("Field%d", i)
    }
    
    logger(fmt.Sprintf("Field %d: %s = %v", i, fieldName, fieldValue))
}
```

#### 4.2. Test with Primitive Types First

```go
type TestPrimitives struct {
    IntField    int
    Int8Field   int8
    Int16Field  int16
    BoolField   bool
    Float32Field float32
}
```

#### 4.3. Then Test Strings

```go
type TestStrings struct {
    Name        string
    Description string
}
```

### Phase 5: Documentation & Examples

#### 5.1. Update README.md

Document the TinyGo limitations and unsafe approach:

```markdown
## TinyGo WASM Support

### Known Limitations

- Field names require manual registration (RTTI stripping)
- Uses `unsafe` pointer arithmetic for field access
- Binary size optimized: ~50KB with full struct reflection

### Usage with TinyGo

```go
// Register field names at init time
func init() {
    tinyreflect.RegisterFieldNames("MyStruct", tinyreflect.FieldNames{
        "Field1", "Field2", "Field3",
    })
}

// Access fields by index (reliable)
v := tinyreflect.ValueOf(myStruct)
for i := 0; i < v.NumField(); i++ {
    value, _ := v.UnsafeFieldValue(i)
    name := v.Type().GetFieldName(i)
    fmt.Printf("%s = %v\n", name, value)
}
```