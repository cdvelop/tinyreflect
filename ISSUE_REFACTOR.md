# TinyReflect Refactor Instructions

## Project Context
TinyReflect is a minimal reflection package for small devices, based on Go's reflectlite, designed specifically for TinyGo/WebAssembly targets.

## Dependencies & Constraints
- **ONLY** depends on: tinystring, sync, unsafe
- **NO** standard library dependencies (fmt, strings, strconv, errors, reflect)
- Must use tinystring's Kind definitions from `kind.go`
- Target: TinyGo/WebAssembly for minimal binary size
- **CRITICAL**: Minimalist type support to reduce code size

## Supported Types (Strict Limitations)
TinyReflect must support ONLY these types to maintain minimal code:
- **Basic**: `string`, `bool`
- **Numeric**: All int/uint variants, float32, float64
- **All basic slices**: `[]string`, `[]bool`, `[]byte`, `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`, `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`, `[]float32`, `[]float64`
- **Structs**: Only with fields of supported types
- **Struct slices**: `[]struct{...}` where all fields are supported
- **Maps**: `map[K]V` where K and V are supported types only
- **Map slices**: `[]map[K]V` where K and V are supported types only
- **Pointers**: Only to supported types above
- **Unsupported**: `interface{}`, `chan`, `func`, `complex64`, `complex128`, `uintptr`, `unsafe.Pointer`, arrays

## Error Handling Requirements
- **NO** custom error messages in tinyreflect
- **NO** panic() calls - use error returns with tinystring's multilingual system
- **MUST** use tinystring's multilingual error system (D.* dictionary)
- Use `Err()` function from tinystring for error creation
- Pattern: `Err(D.Type, D.Not, D.Supported)` for unsupported types
- If missing error terms, add them to tinystring's `dictionary.go` first

## Kind System Integration
- Use tinystring's Kind definitions (KString, KInt, KBool, etc.)
- Import: `. "github.com/cdvelop/tinystring"` for Kind access
- Remove any duplicate Kind definitions from tinyreflect
- Adapt all Kind references to use tinystring's constants

## Code Structure Rules
- Prefix all public types/functions with 'ref' to avoid API pollution
- Keep minimal interface - only essential reflection for supported types
- Use unsafe.Pointer for low-level memory operations
- Maintain thread safety with sync primitives where needed
- **NO** panic() calls - return errors using tinystring's system
- **PRIORITY**: Reuse tinystring's Convert() type detection to minimize binary size
- Use Convert().GetKind() instead of duplicating type detection logic
- Reject unsupported types with `Err(D.Type, D.Not, D.Supported)`

## Memory Optimization Strategy (NEW)
- **ELIMINATE ptrValue**: No more interface{} allocations for type storage
- **USE refType directly**: Extract refType via unsafe from any value
- **MINIMAL conv struct**: Only Kind + refType + unsafe.Pointer for data access
- **NO redundant storage**: refType contains all type metadata needed
- **ZERO interface boxing**: Direct unsafe pointer access to data
- **TinyGo/WASM optimized**: Unsafe pointers more stable than interfaces

## Current Issues to Fix
1. Replace all `panic()` calls with `Err(D.*)` returns
2. Remove support for complex types (structs, interfaces, channels, functions)
3. Simplify type detection to only handle supported types
4. Replace all `errorType("message")` with `Err(D.*)` calls
5. Use Convert().GetKind() for all type detection
6. Remove code for unsupported operations and types
7. Add missing error dictionary terms to tinystring if needed
8. **OPTIMIZE**: Eliminate ptrValue from conv struct to prevent unnecessary allocations
9. **REFACTOR**: Use refType directly instead of interface{} storage
10. **MINIMIZE**: Replace interface{} with unsafe.Pointer for data access

## Error Message Migration Strategy
1. Identify all hardcoded error strings
2. Map them to appropriate D.* dictionary combinations
3. If terms missing, add to tinystring/dictionary.go first
4. Replace errorType() calls with Err() calls
5. Test compilation and functionality

## File Responsibilities
- `abi.go`: Type definitions, Kind system integration
- `reflect.go`: Core reflection operations, error handling via tinystring
- `tinyreflect.go`: Public API interface
- Follow tinystring's error patterns from README.md and TRANSLATE.md

## Target Usage Pattern
```go
import . "github.com/cdvelop/tinystring"

// OPTIMIZED: Extract refType directly without interface{} allocation
val := refValueOf(data)
refType := extractRefType(data) // Direct unsafe extraction
kind := refType.Kind()

// Only for supported types
if kind == KStruct {
    return Err(D.Type, D.Not, D.Supported) // Reject unsupported types
}

// Safe operations for supported types only - no interface{} boxing
if kind == KString {
    str := accessStringData(unsafe.Pointer(&data))
}

// NO ptrValue storage - direct unsafe access
if kind == KSlice {
    length := getSliceLength(unsafe.Pointer(&data))
}
```

## Success Criteria
- Zero compilation errors
- All errors use tinystring's multilingual system  
- **NO** panic() calls - graceful error handling
- Minimal binary size impact through maximum code reuse
- Full TinyGo compatibility
- Only depends on: tinystring, sync, unsafe
- **CRITICAL**: Use Convert().GetKind() to eliminate type detection duplication
- Support only minimal, essential types for JSON-like operations
- **OPTIMIZED**: Zero interface{} allocations for type storage
- **EFFICIENT**: Direct unsafe.Pointer access instead of ptrValue
- **MINIMAL**: refType-only approach eliminates memory overhead

## Next Steps After Document Creation
1. Add missing dictionary terms to tinystring if needed
2. Fix Kind system integration
3. Replace all error messages with D.* patterns
4. Test compilation and basic functionality
5. Update README.md with Why TinyReflect? section matching tinystring's pattern
6. **IMPLEMENT**: Eliminate ptrValue from conv struct in tinystring
7. **OPTIMIZE**: Replace interface{} storage with refType + unsafe.Pointer
8. **REFACTOR**: Update anyToBuff to extract refType directly
9. **TEST**: Validate memory optimization reduces allocations
10. **BENCHMARK**: Confirm binary size reduction achieved
