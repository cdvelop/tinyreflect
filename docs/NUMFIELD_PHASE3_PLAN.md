# NumField TinyGo Debug Plan - Phase 3: Deep Dive

## Current Situation

After implementing `underlying()` method:
- ✅ Stdlib tests pass (5 fields found)
- ❌ TinyGo WASM still returns 0 fields
- ✅ No compilation errors
- ✅ No panics - execution completes

## Problem: Why underlying() Isn't Working

### Investigation Results

#### 1. Our Current Implementation

**TypeOf_tinygo.go:**
```go
func (t *Type) underlying() *Type {
    if t.isNamed() {
        return (*elemType)(unsafe.Pointer(t)).elem
    }
    return t
}

func (t *Type) isNamed() bool {
    if t.ptrtag() != 0 {
        return false
    }
    return t.Kind_&flagNamed != 0  // ← PROBLEM: Checking wrong field!
}
```

**The Issue:**
- We're checking `t.Kind_` for the `flagNamed` bit
- But `Kind_` is type `Kind` which is uint8
- In TinyGo's `RawType`, the field is called `meta` not `Kind_`
- Our `Type` struct doesn't have a `meta` field directly accessible

#### 2. Type Structure Mismatch

**TinyGo's RawType:**
```go
// /usr/local/lib/tinygo/src/internal/reflectlite/type.go:162
type RawType struct {
    meta uint8  // ← This contains kind AND flags
}
```

**Our Type:**
```go
// TypeOf.go
type Type struct {
    Size        uintptr
    PtrBytes    uintptr
    Hash        uint32
    TFlag       TFlag    // ← Different flag system!
    Align_      uint8
    FieldAlign_ uint8
    Kind_       Kind     // ← This is just the kind, not meta!
    Equal       func(unsafe.Pointer, unsafe.Pointer) bool
    GCData      *byte
    Str         NameOff
    PtrToThis   TypeOff
}
```

**The Core Problem:**
- Our `Type` structure is based on **stdlib's abi.Type**
- TinyGo uses a completely different **RawType** structure
- We can't cast TinyGo's RawType to our Type correctly!
- Memory layout is COMPLETELY different!

### Memory Layout Comparison

#### Stdlib abi.Type (what our Type matches):
```
Offset  Field           Size
0       Size            8 bytes
8       PtrBytes        8 bytes
16      Hash            4 bytes
20      TFlag           1 byte
21      Align_          1 byte
22      FieldAlign_     1 byte
23      Kind_           1 byte
...
Total: ~40+ bytes
```

#### TinyGo RawType (what TinyGo actually uses):
```
Offset  Field           Size
0       meta            1 byte
Total: 1 byte base struct!
```

**This means:**
- When TinyGo gives us a `*RawType`, it points to just 1 byte!
- We cast it to our `*Type` which expects 40+ bytes
- We're reading garbage memory!
- `Kind_` at offset 23 is reading random memory
- `underlying()` never works because `isNamed()` reads garbage

## Root Cause: Fundamental Architecture Incompatibility

We tried to use **stdlib's Type structure** for **TinyGo's runtime**.

This cannot work because:
1. Stdlib Type has ~40 byte header
2. TinyGo RawType has 1 byte header
3. They're incompatible at the binary level
4. No amount of `unsafe.Pointer` casting will fix this

## The Real Solution: Use TinyGo's RawType Directly

### Strategy: Wrapper Pattern

Instead of trying to match stdlib's Type, we need:

**For TinyGo:**
- Use TinyGo's `RawType` directly
- Don't redefine structures
- Import or recreate just what we need
- Cast directly to TinyGo's types

**For Stdlib:**
- Keep current implementation (works fine)

### Implementation: Two Completely Different Backends

#### Option A: Import TinyGo's reflectlite (BLOCKED)
```go
//go:build tinygo

import "internal/reflectlite"

type Type = reflectlite.RawType
```

**Problem:** Can't import `internal/` packages from external modules.

#### Option B: Recreate TinyGo's Core Types (RECOMMENDED)

**Create minimal TinyGo-compatible types in our codebase:**

```go
//go:build tinygo

package tinyreflect

import "unsafe"

// RawType matches TinyGo's internal/reflectlite.RawType exactly
type RawType struct {
    meta uint8
}

// Our public Type is just an alias
type Type = RawType

const (
    kindMask  = 31
    flagNamed = 32
)

func (t *RawType) Kind() Kind {
    if t == nil {
        return Invalid
    }
    if tag := t.ptrtag(); tag != 0 {
        return Pointer
    }
    return Kind(t.meta & kindMask)
}

func (t *RawType) underlying() *RawType {
    if t.isNamed() {
        return (*elemType)(unsafe.Pointer(t)).elem
    }
    return t
}

func (t *RawType) isNamed() bool {
    if t.ptrtag() != 0 {
        return false
    }
    return t.meta&flagNamed != 0  // ← Now reads correct field!
}

// elemType matches TinyGo's layout
type elemType struct {
    meta      uint8
    numMethod uint16
    ptrTo     *RawType
    elem      *RawType
}

// structType matches TinyGo's layout
type structType struct {
    meta      uint8
    numMethod uint16
    ptrTo     *RawType
    pkgpath   *byte
    size      uint32
    numField  uint16
    fields    [1]structField
}
```

#### Option C: Dynamic Type Detection (COMPLEX)

Detect at runtime which structure we're dealing with and adapt.

**Cons:** Very complex, error-prone, defeats the purpose.

## Recommended Implementation Plan

### Phase 1: Create TinyGo Type System

1. **Create `Type_tinygo.go`** with TinyGo's RawType as base
2. **Create `Type_stdlib.go`** with current Type (works fine)
3. **Update all `Type` references** to work with both

### Phase 2: Recreate TinyGo Struct Handling

1. **Redefine `structType`** to match TinyGo exactly
2. **Implement `NumField()`** using TinyGo's approach
3. **Fix field access** methods

### Phase 3: Handle Common Interface

Create interface that both implementations satisfy:
```go
type Typer interface {
    Kind() Kind
    NumField() (int, error)
    Field(i int) (StructField, error)
    // ... other methods
}
```

## Critical Realization

**We can't have a unified Type structure.**

TinyGo and Stdlib use completely different internal representations.
We need **two completely separate implementations** that present the same API.

## Files Structure

```
tinyreflect/
├── Type_stdlib.go        // Current Type struct (works)
├── Type_tinygo.go        // New: RawType-based (to fix)
├── StructType_stdlib.go  // Keep as is
├── StructType_tinygo.go  // Redefine to match TinyGo
├── TypeOf_stdlib.go      // Keep as is
├── TypeOf_tinygo.go      // Rewrite using RawType
├── ValueOf_stdlib.go     // Split from ValueOf.go
├── ValueOf_tinygo.go     // TinyGo-specific Value handling
└── tinyreflect.go        // Common interface/API only
```

## Next Steps

1. **Verify memory layout:** 
   - Add debug to print sizeof(Type) in both compilers
   - Confirm our Type doesn't match TinyGo's RawType

2. **Create minimal TinyGo backend:**
   - Start with just RawType
   - Implement Kind() correctly
   - Test that Kind() returns Struct

3. **Build up from there:**
   - Add underlying()
   - Add NumField()
   - Add Field access

## Why Previous Approach Failed

❌ **What we did:** Tried to create one Type structure for both
❌ **Why it failed:** TinyGo and Stdlib have incompatible memory layouts
✅ **What we need:** Two separate implementations with same API

## Success Criteria

- [ ] TinyGo: `sizeof(Type) == 1` (just meta byte)
- [ ] TinyGo: `Kind()` reads from `meta` field
- [ ] TinyGo: `isNamed()` correctly detects named types
- [ ] TinyGo: `underlying()` resolves to structType
- [ ] TinyGo: `NumField()` returns 5 fields
- [ ] Stdlib: Everything still works (don't break it!)
- [ ] Same API for both implementations

## Key Insight

The entire approach needs to change:
- **Don't try to unify the Type structure**
- **Create TinyGo-specific implementation from scratch**
- **Keep stdlib implementation as-is**
- **Only unify the public API**
