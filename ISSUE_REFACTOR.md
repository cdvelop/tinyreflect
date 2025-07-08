# TinyReflect Refactor Instructions

## Project Context
TinyReflect is a minimal reflection package for small devices, based on Go's reflectlite, designed specifically for TinyGo/WebAssembly targets.

## Dependencies & Constraints
- **ONLY** depends on: tinystring, sync, unsafe
- **NO** standard library dependencies (fmt, strings, strconv, errors, reflect)
- Must use tinystring's Kind definitions from `kind.go`
- Target: TinyGo/WebAssembly for minimal binary size

## Error Handling Requirements
- **NO** custom error messages in tinyreflect
- **MUST** use tinystring's multilingual error system (D.* dictionary)
- Use `Err()` function from tinystring for error creation
- If missing error terms, add them to tinystring's `dictionary.go` first
- Follow pattern: `Err(D.Cannot, D.Set, D.Value)` instead of hardcoded strings

## Kind System Integration
- Use tinystring's Kind definitions (KString, KInt, KBool, etc.)
- Import: `. "github.com/cdvelop/tinystring"` for Kind access
- Remove any duplicate Kind definitions from tinyreflect
- Adapt all Kind references to use tinystring's constants

## Code Structure Rules
- Prefix all public types/functions with 'ref' to avoid API pollution
- Keep minimal interface - only essential reflection for JSON operations
- Use unsafe.Pointer for low-level memory operations
- Maintain thread safety with sync primitives where needed

## Current Issues to Fix
1. Replace all `errorType("message")` with `Err(D.*)` calls
2. Fix recursive Kind() method in abi.go (use t.kind field, not t.Kind())
3. Add missing initFromValue() method to refValue
4. Fix missing return statements and function signatures
5. Ensure all Kind references use tinystring's definitions
6. Add missing error dictionary terms to tinystring if needed

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
// Use reflection with tinystring error handling
val := refValueOf(data)
if err := val.someOperation(); err != nil {
    return Err(D.Cannot, D.Set, D.Value) // Not hardcoded strings
}
```

## Success Criteria
- Zero compilation errors
- All errors use tinystring's multilingual system
- Minimal binary size impact
- Full TinyGo compatibility
- Only depends on: tinystring, sync, unsafe

## Next Steps After Document Creation
1. Add missing dictionary terms to tinystring if needed
2. Fix Kind system integration
3. Replace all error messages with D.* patterns
4. Test compilation and basic functionality
5. Update README.md with Why TinyReflect? section matching tinystring's pattern
