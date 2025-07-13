# TinyReflect Implementation Plan

## Project Overview
TinyReflect is a minimal Go reflection package designed for WebAssembly and TinyGo compatibility. It aims to replace the standard library `reflect` package with ultra-minimal, focused implementations that dramatically reduce binary size while maintaining essential functionality.

## Architecture Goals
- **Maximum code reuse** with tinystring package
- **Minimal binary footprint** for WebAssembly deployment  
- **TinyGo compatibility** with no compilation issues
- **Focused API** supporting only essential operations for JSON-like data handling
- **Zero external dependencies** beyond standard interfaces

## Current Issues to Fix

### 1. Kind Type Conflicts (CRITICAL)
**Problem**: `Kind` type and `K` variable are redeclared in `kind.go` but already exist in tinystring
**Solution**: Remove duplicate declarations, use tinystring's Kind implementation
**Files to modify**: 
- Remove/refactor `kind.go` completely
- Update imports in files using Kind

### 2. Missing StructTag.Get Implementation (FAILING TEST)
**Problem**: `StructTag.Get()` method calls `Convert(tag).TagValue(key)` but this method doesn't exist in tinystring
**Solution**: Implement proper tag parsing using standard library approach
**Files to modify**:
- `StructTag.go` - implement proper Get() method

### 3. Missing Core Type System Integration
**Problem**: Types need proper integration with tinystring's type detection
**Solution**: Leverage tinystring's existing type system
**Files to check/update**:
- `TypeOf.go` - ensure proper Kind integration
- `StructType.go` - struct field enumeration

## Error Handling System (TinyString Integration)

TinyReflect integrates with TinyString's multilingual error system which provides translations for error messages in 9 languages.

### How It Works
1. **Dictionary Constants**: Use predefined constants from `D` struct (e.g., `D.Field`, `D.Type`, `D.Struct`)
2. **Ref Constant**: The word "reflect" cannot be translated, so it's defined as `const ref = "reflect"`
3. **Error Creation**: Use `Err(ref, D.Constant1, D.Constant2, ...)` to compose multilingual error messages
4. **Available Constants**: Check `dictionary.go` in tinystring for all available translated terms

### Example Usage
```go
// Error: field index out of range
return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)

// Error: field not of type struct  
return StructField{}, Err(ref, D.Field, D.Type, D.Struct)
```

### Available Dictionary Terms for Reflection
Key terms available in `D` struct:
- `D.Field`, `D.Fields` - field/fields
- `D.Type` - type
- `D.Struct` - struct
- `D.Value` - value
- `D.Range` - range
- `D.Out`, `D.Of` - out, of
- `D.Numbers` - numbers
- `D.Cannot` - cannot
- `D.Empty` - empty

## Implementation Strategy
1. **Remove kind.go conflicts**
   - Delete redundant Kind declarations  
   - Use tinystring.K constants
   - Fix import dependencies

2. **Implement StructTag.Get() properly**
   - Parse tag strings following Go standard format: `key:"value"`
   - Handle multiple tags separated by spaces
   - Support quoted values with escape sequences
   - Return empty string for missing keys

3. **Fix test failures**
   - Ensure `TestGetFieldName` passes completely
   - Validate struct field name, type, and tag extraction

### Phase 2: Core Functionality Implementation
1. **TypeOf() function**
   - Extract type information from runtime
   - Create Type structures with proper Kind values
   - Support for all basic types listed in README

2. **Struct field enumeration**
   - Implement `NumField()` method
   - Implement `Field(i int)` method  
   - Extract field names, types, tags, and offsets

3. **Type.String() method**
   - Return proper type names using tinystring
   - Integrate with existing naming conventions

### Phase 3: Extended Type Support  
1. **Slice and Map support**
   - Element type detection
   - Length/capacity information where applicable

2. **Pointer support**
   - Pointer dereferencing
   - Nil pointer handling

3. **Value operations** (if needed)
   - Basic value extraction
   - Type conversion utilities

## File Architecture

### Core Files (Keep/Modify)
- `TypeOf.go` - Main type detection entry point
- `StructField.go` - Struct field representation  
- `StructTag.go` - Tag parsing and extraction
- `StructType.go` - Struct type operations
- `Name.go` - Name representation and tag storage
- `ValueOf.go` - Value operations (minimal)

### Files to Remove/Refactor
- `kind.go` - Remove completely, use tinystring.Kind

### Test Files
- `StructTag_test.go` - Current failing test, must pass
- Additional tests for TypeOf, Field enumeration

## Dependencies and Integration

### TinyString Integration Points
1. **Kind constants**: Use `tinystring.K.*` for all type kinds
2. **Type conversion**: Leverage `tinystring.Convert()` where applicable  
3. **String operations**: Use tinystring methods for string manipulation
4. **Error handling**: Use tinystring's error system if available

### Standard Library Compatibility
- Runtime type information access (`unsafe` package)
- Memory layout compatibility with Go's internal types
- Tag parsing following `reflect.StructTag` behavior

## Supported Types (Per README)
**✅ Implemented Priority:**
- Basic types: `string`, `bool`
- Numeric types: all int/uint variants, float32/float64  
- Basic slices: `[]string`, `[]bool`, `[]byte`, etc.
- Structs: with supported field types only
- Maps: `map[K]V` where K,V are supported types
- Pointers: to supported types only

**❌ Explicitly Unsupported:**
- `interface{}`, `chan`, `func`
- `complex64`, `complex128`  
- `uintptr`, `unsafe.Pointer` (internal use only)
- Arrays (different from slices)
- Nested complex types beyond scope

## Testing Strategy
1. **Fix current test**: `TestGetFieldName` must pass completely
2. **Add comprehensive tests**: All supported types from README
3. **TinyGo compatibility**: Ensure compilation without warnings
4. **Binary size validation**: Measure WebAssembly output size

## Implementation Notes
- Use `unsafe` package carefully for runtime type access
- Follow Go's internal type layout for compatibility
- Minimize allocations for better performance  
- Keep API surface minimal and focused
- Document any TinyGo-specific optimizations

## Progress Tracking
- [x] Phase 1: Fix Kind conflicts and StructTag.Get()
- [x] Phase 1: Make TestGetFieldName pass
- [x] Phase 1: Implement proper error handling with multilingual support
- [ ] Phase 2: Complete TypeOf() implementation  
- [ ] Phase 2: Full struct field enumeration
- [ ] Phase 3: Extended type support
- [ ] Testing: Comprehensive test coverage
- [ ] Documentation: Update README with examples

## Current Files Status
- `kind.go` - ✅ DELETED (conflicts resolved)
- `StructTag.go` - ✅ FIXED (Get method working correctly)  
- `StructTag_test.go` - ✅ PASSING (all assertions pass)
- `TypeOf.go` - ✅ UPDATED (Kind constants added)
- `StructField.go` - ✅ LOOKS OK
- `Name.go` - ✅ LOOKS OK  
- `StructType.go` - ⚠️ NEEDS REVIEW
- `ValueOf.go` - ✅ CLEANED UP (duplicates removed)

## Recent Changes Made (Phase 1 - COMPLETED ✅)
1. **Removed kind.go completely** - Eliminated duplicate Kind type definitions
2. **Fixed StructTag.Get()** - Added explicit string conversion: `Convert(string(tag)).TagValue(key)`
3. **Cleaned up ValueOf.go** - Removed duplicate constants and EmptyInterface
4. **Added constants to TypeOf.go** - Moved KindDirectIface, KindMask, EmptyInterface to proper location
5. **Implemented proper error handling** - Uses TinyString's multilingual system with `ref` constant
6. **All tests now pass** - TestGetFieldName validates struct field name, type, and tag extraction

## Error System Integration ✅
- **Constant**: `const ref = "reflect"` (non-translatable technical term)
- **Usage**: `Err(ref, D.Field, D.Out, D.Of, D.Range)` for "reflect field out of range"
- **Multilingual**: Automatically translates using TinyString's 9-language dictionary
- **Available Terms**: D.Field, D.Type, D.Struct, D.Value, D.Range, D.Numbers, etc.

Last Updated: Phase 1 COMPLETED - All core functionality working with proper multilingual error support
Next Action: Begin Phase 2 - Review StructType.go and extend type support
