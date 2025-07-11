# TinyString JSON API Integration Instructions

## Goal
Implement minimal, dependency-free JSON encoding and decoding for structs in TinyString, using the new reflection system. The API must be compatible with TinyGo and avoid all dependencies (including encoding/json).

## Core Philosophy
This implementation must adhere to TinyString's core principles as outlined in WHY.md:
- **🏆 Smallest possible binary Size** - Minimize WebAssembly footprint
- **📦 Zero dependencies** - No imports beyond standard interfaces
- **🔧 Maximum code reuse** - Leverage existing code in TinyString
- **✅ Full TinyGo compatibility** - No compilation issues with TinyGo

## API Design
- Encoding: `Convert(&struct{}).JsonEncode(w io.Writer) error`
- Decoding: `Convert(r io.Reader).JsonDecode(&struct{}) error` (only pointer for decoding)

## Buffer Architecture Requirements
- **All internal JSON methods must use TinyString's buffer system** (`buffOut`, `buffWork`, `buffErr`)
- **Private methods must NOT return `(string, error)` pairs** - use buffer destinations instead
- **Reuse existing buffer methods**: `wrString()`, `wrBytes()`, `getString()`, `rstBuffer()`
- **Follow existing patterns**: See `quote.go`, `num_float.go` for examples of proper buffer usage

## API Exposure Rules
- **ONLY** expose the two JSON methods: `JsonEncode` and `JsonDecode`
- **ALL** internal reflection types and methods must be PRIVATE (lowercase)
- **NO** public exposure of struct tags, field access, or any reflection internals
- Keep TinyString's API clean and focused

## Writer/Reader Types
- Use standard Go interfaces: `io.Writer` and `io.Reader` for maximum compatibility.

## JSON Standard
- Support basic struct tags for field naming (e.g., `json:"field_name"`).
- Only exported fields are encoded/decoded.
- No support for omitempty or advanced tag options.
- Minimal implementation: just enough for basic JSON compatibility and TinyGo support.

## Error Handling
- All errors must use the multilingual error system (`Err(D.Type, ...)`) as in the rest of the library.

## Dependencies
- 100% custom implementation. **No use of encoding/json or any external package.**
- **STRICTLY FORBIDDEN**: No use of the standard `reflect` package - this defeats the purpose of TinyString
- Must use only the custom reflection implementation from tinyreflect for minimal binary Size

## File Structure
- Implement encoding in `jsonencode.go` and decoding in `jsondecode.go`.
- Place tests in `jsonencode_test.go` and `jsondecode_test.go`.

## Implementation Strategy
1. **Maximum Code Reuse**: Reuse existing TinyString functionality wherever possible:
   - Use existing buffer system (`buffOut`, `buffWork`, `buffErr`) from `memory.go`
   - Reuse existing string escaping logic from `quote.go`
   - Reuse existing numeric parsing from `num_float.go`, `num_int.go`
   - Reuse existing type detection from `Kind_.go`
   - Reuse existing error handling from `error.go`
2. **Buffer-Based Architecture**: All internal JSON methods must use buffer destinations, not return `(string, error)` pairs
3. **Private API Design**: All JSON helper methods must be private and follow TinyString's buffer-first architecture
4. **Minimal Reflection**: Add only the minimum reflection needed for struct field access to existing `reflect.go`

## Supported Types (as per README.md)
JSON functionality must support the following types:
- **Basic types**: `string`, `bool`
- **Numeric types**: All int/uint variants, float32, float64
- **All basic slices**: `[]string`, `[]bool`, `[]byte`, etc.
- **Structs**: Only with supported field types
- **Maps with string keys**: `map[string]string`, `map[string]int`, etc.
- **Pointers**: Only to supported types above

## Test Coverage Requirements
- Create comprehensive tests that cover 100% of the JSON API
- Test all supported types listed in README.md
- Include test cases for:
  - Simple structs with primitive fields
  - Nested structs
  - Slices of structs
  - Maps with string keys
  - Error handling cases (invalid JSON, type mismatches)
  - TinyGo compatibility

## README Update
- Add a usage example for the new JSON API to the README **after implementation is complete**.

## Implementation Requirements
- **Buffer-first architecture**: All JSON helper methods must use buffer destinations
- **Code reuse**: Leverage existing functionality from `quote.go`, `num_float.go`, `num_int.go`, etc.
- **Private methods only**: No public exposure of JSON internals beyond `JsonEncode` and `JsonDecode`
- **Minimal struct field access**: Only the minimum reflection needed for basic struct operations
- **Error handling**: Use existing multilingual error system with buffer destinations
- **TinyGo compatibility**: Must compile and run correctly with TinyGo

## Code Reuse Examples
- String escaping: Use patterns from `Quote()` method in `quote.go`
- Number formatting: Use `wrFloat64()`, `wrIntBase()` from existing numeric modules
- Buffer management: Use `wrString()`, `rstBuffer()`, `getString()` from `memory.go`
- Error handling: Use `wrErr()` pattern with dictionary terms
