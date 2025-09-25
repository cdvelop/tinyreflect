# TinyReflect Cache System Proposal

## Executive Summary

Transform TinyReflect from global functions to instance-based architecture with transparent internal caching. Designed specifically for TinyGo/WebAssembly constraints using arrays and atomic operations instead of maps.

## ⚠️ **CRITICAL FACTOR: TinyGo Compatibility**

**IMPORTANT:** TinyGo has significant limitations with maps and sync primitives. This restriction fundamentally modifies the cache system design, favoring implementations based on arrays and simple atomic operations.


## Current Problems

**Global Variables Issues:**
```go
// ❌ ACTUAL: Problemático
var structNameCache []cacheEntry  // Global state
func ValueOf(i any) Value         // Uses global cache

// ✅ PROPUESTO: Instance-based  
tr := tinyreflect.New()           // Instance with internal cache
v := tr.ValueOf(i)               // Uses instance cache
```

**Problems eliminated:**
- Race conditions, testing conflicts, memory leaks
- TinyGo incompatibility, inability to have multiple instances

## Proposed Architecture

```go
type TinyReflect struct {
    structCache [256]structCacheEntry  // Fixed array for TinyGo compatibility  
    structCount int32                  // Atomic counter
    log func(msgs ...any)             // Optional logging
}

// Constructor with variadic args
func New(args ...any) *TinyReflect

// Instance methods (replace global functions)  
func (tr *TinyReflect) ValueOf(i any) Value
func (tr *TinyReflect) TypeOf(i any) *Type
```

**Usage Example:**
```go
// One-time initialization
tr := tinyreflect.New()           // Default: 128 structs, no logging
// tr := tinyreflect.New(StructSize256, log.Println) // 256 structs + logging

// Transparent caching - same API
type User struct { Name string; Age int }

u := User{"Alice", 42}
v := tr.ValueOf(u)               // First call: caches struct schema
typ := v.Type()                  // Standard tinyreflect API

u2 := User{"Bob", 30}            
v2 := tr.ValueOf(u2)            // Fast: uses cached schema
numFields, _ := typ.NumField()   // Fast: from cache
```

## Benefits vs Trade-offs

**✅ Benefits:**
- **Universal optimization:** All libraries using tinyreflect get automatic caching
- **TinyGo compatible:** Arrays + atomics only, no sync.Map
- **Predictable performance:** Fixed memory layout, no GC pressure  
- **Transparent API:** Same tinyreflect API, just faster internally

**⚠️ Trade-offs:**
- **Breaking change:** Requires migration from global functions to instance methods
- **Fixed capacity:** Limited number of cached structs (128-256 typically sufficient)
- **Upfront memory:** Pre-allocates cache arrays at startup

## Cache Design

**Analysis**: Based on TinyBin/StructSQL code analysis - only structs need complex caching. Basic types (int, string) use simple direct codecs.

**Cache Structure:**
```go  
type structCacheEntry struct {
    structID     uint32                // Hash-based key
    nameLen      uint8                 // Struct name length
    fieldCount   uint8                 // Number of fields  
    structName   [32]byte              // Struct name
    fieldSchemas [16]fieldSchema       // Complete field schemas
}

type fieldSchema struct {
    nameLen   uint8      // Field name length
    kind      Kind       // Field type kind
    offset    uint16     // Field offset in struct
    fieldName [20]byte   // Field name
}
```

**What this eliminates:**
- `typ.NumField()` calls (cached in `fieldCount`)  
- `typ.Field(i)` loops (cached in `fieldSchemas[]`)
- `field.Name.Name()` calls (cached in `fieldName`)
- `field.Typ.Kind()` calls (cached in `kind`)

**Memory vs Performance Trade-off:**
- **Without cache**: 37 bytes but still requires expensive `typ.NumField()` + `typ.Field(i)` loops
- **With full schema cache**: 392 bytes but eliminates all reflection loops
- **Result**: +955% memory for -100% reflection overhead

**Transparent Caching:**
```go
tr := tinyreflect.New()

// Same API - first call caches schema internally
typ := tr.TypeOf(myStruct)      
numFields, _ := typ.NumField()   // Expensive first time, cached

// Subsequent calls - same API but fast from cache  
for i := 0; i < numFields; i++ {
    field, _ := typ.Field(i)     // Fast: from cache
    name := field.Name.Name()    // Fast: from cache  
}
```

## Implementation Details

**Key Optimizations:**
- **Hash reuse:** Uses existing `Type.Hash` as cache keys (no additional hash calculation)
- **Transparent caching:** Same tinyreflect API, performance improvement is internal
- **Logging system:** Optional `log func(msgs ...any)` for debugging (no-op by default)

**Constructor Examples:**
```go
tr := tinyreflect.New()                        // Default: 128 structs, no logging
tr := tinyreflect.New(StructSize256)           // 256 structs  
tr := tinyreflect.New(log.Println)             // With logging
tr := tinyreflect.New(StructSize256, t.Logf)   // For tests
```

**Core Methods:**
```go
func (tr *TinyReflect) ValueOf(i any) Value    // Replaces global ValueOf
func (tr *TinyReflect) TypeOf(i any) *Type     // Replaces global TypeOf  
func (tr *TinyReflect) Indirect(v Value) Value // Replaces global Indirect
```

## Migration Strategy & Implementation Plan

**Migration Phases:**
1. **Legacy wrappers** (temporary): Global functions delegate to default instance
2. **Complete cutover**: Update all consumers to instance-based API  
3. **Cleanup**: Remove legacy functions in next major release

**Implementation Priorities:**
1. Remove global variables, implement `New(args...)` constructor, unified struct cache
2. TinyGo validation, collision handling, memory layout optimization
3. Documentation updates, benchmarking different cache sizes

**TinyGo Requirements:**
- Arrays + atomics only (no sync.Map)
- No dynamic allocation after initialization
- All state in TinyReflect struct (zero global variables)

## Final Design

```go
type TinyReflect struct {
    structCache []structCacheEntry    // Configurable size array
    structCount int32                 // Atomic counter  
    cacheLock   int32                 // Atomic lock
    log         func(msgs ...any)     // Optional logging
    maxStructs  int32                 // Capacity limit
}
```

**Performance Expectations:**
| Aspect | TinyReflect | Go reflect | TinyBin |
|--------|-------------|------------|---------|
| Struct analysis | O(1) after first | O(n) always | O(1) map |
| TinyGo compatible | ✅ Arrays+atomics | ✅ Compatible | ❌ sync.Map |  
| WebAssembly | ✅ Optimal | ❌ Slow | ⚠️ Overhead |

## Conclusion

Essential architectural evolution for tiny* ecosystem - TinyGo/WebAssembly optimized with:
- 100% TinyGo compatibility (arrays + atomics only)
- Transparent caching (same API, better performance)  
- Predictable memory layout (fixed allocations)
- Universal benefits for all tinyreflect consumers

---
**Status:** Ready for TinyGo-first implementation
