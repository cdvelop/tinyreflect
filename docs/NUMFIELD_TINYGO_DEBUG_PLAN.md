# NumField TinyGo Debug Plan - Phase 2

## Current Status

### What Works ✅
- **Stdlib compilation:** NumField returns 5 fields correctly
- **All tests pass with stdlib:** `go test ./...` passes
- **Build tags working:** Correct StructType selected for each compiler

### What Fails ❌
- **TinyGo WASM:** NumField returns 0 fields
- **Execution stops:** After calling NumField, no error but returns 0

## Problem Analysis - Deeper Investigation

### Test Results Comparison

**Standard Go (Working):**
```
Theme Hello, PWA!
main.js:2 === Testing tinyreflect ===
main.js:2 Found: 5 fields
main.js:2 === Testing tinyreflect ===
main.js:2 DEBUG: About to call ValueOf
main.js:2 DEBUG: ValueOf returned, v.typ_ = true
main.js:2 DEBUG: Type() returned, typ = true
main.js:2 DEBUG: About to call NumField
main.js:2 SUCCESS: Found 5 fields
main.js:2 Field 0: StringField = string:test string
...
```

**TinyGo (Failing):**
```
Theme Hello, PWA!
main.js:2 === Testing tinyreflect ===
main.js:2 Found: 0 fields
main.js:2 === Testing tinyreflect ===
main.js:2 DEBUG: About to call ValueOf
main.js:2 DEBUG: ValueOf returned, v.typ_ = true
main.js:2 DEBUG: Type() returned, typ = true
main.js:2 DEBUG: About to call NumField
main.js:2 [STOPS HERE - no output]
```

### Root Cause Analysis

#### Issue 1: Missing `underlying()` Call

**TinyGo's implementation:**
```go
// /usr/local/lib/tinygo/src/internal/reflectlite/type.go:603
func (t *RawType) NumField() int {
    if t.Kind() != Struct {
        panic(errTypeNumField)
    }
    return int((*structType)(unsafe.Pointer(t.underlying())).numField)  // ← CALLS underlying()!
}

// Line 237
func (t *RawType) underlying() *RawType {
    if t.isNamed() {
        return (*elemType)(unsafe.Pointer(t)).elem
    }
    return t
}
```

**Our current implementation (WRONG):**
```go
// TypeOf.go:117
func (t *Type) NumField() (int, error) {
    if t.Kind() != K.Struct {
        return 0, Err(ref, D.Numbers, D.Fields, D.Type, "Struct")
    }
    st := (*StructType)(unsafe.Pointer(t))  // ← MISSING underlying() call!
    return st.numFields(), nil
}
```

**Why this matters:**
- In TinyGo, when you create a named struct type (like `TestStruct`), the Type pointer points to a "named type" wrapper
- The named type has a different memory layout - it's an `elemType` that contains a pointer to the actual `structType`
- Without calling `underlying()`, we're reading memory from the wrong location
- The `numField` at that location is 0 (uninitialized or wrong offset)

#### Issue 2: Named Types vs Underlying Types

TinyGo has two type representations:

**Named Type (elemType):**
```go
type elemType struct {
    RawType
    numMethod uint16
    ptrTo     *RawType
    elem      *RawType  // ← Points to actual underlying type
    pkgpath   *byte
    name      [1]byte   // null terminated name
}
```

**Underlying Struct Type (structType):**
```go
type structType struct {
    RawType
    numMethod uint16
    ptrTo     *RawType
    pkgpath   *byte
    size      uint32
    numField  uint16    // ← The field count we need!
    fields    [1]structField
}
```

**The problem:**
- When we do `TypeOf(TestStruct{})`, TinyGo returns a pointer to `elemType` (named type)
- We cast it directly to `StructType` and try to read `numField`
- But at that memory location in `elemType`, there's no `numField` - it's in the wrong place!
- We need to follow the `elem` pointer first to get to the actual `structType`

#### Issue 3: Our StructType Layout Mismatch

**Our TinyGo StructType:**
```go
// StructType_tinygo.go
type StructType struct {
    Type                    // Embedded Type
    numMethod uint16
    ptrTo     *Type
    pkgpath   *byte
    size      uint32
    numField  uint16        // ← We expect this here
    Fields    [1]StructField
}
```

**But if Type is a named type, the layout is actually:**
```go
type elemType struct {
    Type                    // Embedded Type (RawType in TinyGo)
    numMethod uint16
    ptrTo     *Type
    elem      *Type         // ← Pointer to underlying type!
    pkgpath   *byte
    name      [1]byte       // ← Name, not size!
}
```

**Memory layout mismatch:**
- Position of `numField` in `StructType`: offset ~40 bytes
- Position of `elem` pointer in `elemType`: offset ~20 bytes
- When we read `numField` from `elemType`, we're reading garbage memory!

## Solution Strategy

### Option A: Implement `underlying()` Method (RECOMMENDED)

Add TinyGo-compatible `underlying()` method to get the actual struct type:

**Steps:**
1. Add `underlying()` method to `Type` with build tags
2. Implement TinyGo version that checks `isNamed()` and follows `elem` pointer
3. Implement stdlib version that just returns self
4. Update all methods to call `underlying()` before casting to `StructType`

**Pros:**
- Matches TinyGo's exact behavior
- Handles named types correctly
- Minimal changes to existing code

**Cons:**
- Need to understand named type detection (`isNamed()`, `ptrtag()`)
- More complex implementation

### Option B: Always Use StructType() Helper

Update `StructType()` method to handle underlying type resolution:

**Steps:**
1. Update `Type.StructType()` to check if named and resolve
2. All code already uses `StructType()` so should just work
3. Build-tag specific implementations

**Pros:**
- Centralized fix in one place
- Cleaner API

**Cons:**
- Still need to implement `underlying()` logic
- May hide complexity

### Recommended: Option A + Option B Combined

1. Implement `underlying()` as internal helper
2. Update `StructType()` to use it
3. Ensure all field access goes through helpers

## Implementation Plan

### Step 1: Add `underlying()` Helper (TinyGo)

Create `TypeOf_tinygo.go`:
```go
//go:build tinygo

package tinyreflect

import "unsafe"

// elemType represents a named type in TinyGo
type elemType struct {
    Type
    numMethod uint16
    ptrTo     *Type
    elem      *Type  // Pointer to underlying type
    pkgpath   *byte
    // name follows as [1]byte but we don't need it
}

// underlying returns the underlying type.
// For named types, returns the elem pointer.
// For unnamed types, returns self.
func (t *Type) underlying() *Type {
    if t.isNamed() {
        return (*elemType)(unsafe.Pointer(t)).elem
    }
    return t
}

// isNamed checks if this is a named type
func (t *Type) isNamed() bool {
    // Check ptrtag first
    if t.ptrtag() != 0 {
        return false
    }
    // Check TFlagNamed flag
    return t.TFlag&tflagNamed != 0
}

// ptrtag returns the pointer tag (last 2 bits)
func (t *Type) ptrtag() uintptr {
    return uintptr(unsafe.Pointer(t)) & 0b11
}

// TFlag constants for TinyGo
const (
    tflagNamed TFlag = 1 << 2  // Type is named
)
```

### Step 2: Add `underlying()` Helper (Stdlib)

Create `TypeOf_stdlib.go`:
```go
//go:build !tinygo

package tinyreflect

// underlying returns the underlying type.
// In stdlib, types don't have the named/unnamed distinction
// the same way, so we just return self.
func (t *Type) underlying() *Type {
    return t
}
```

### Step 3: Update `StructType()` Method

In `TypeOf.go` (no build tags - works for both):
```go
// StructType returns t cast to a *StructType, or nil if its tag does not match.
func (t *Type) StructType() *StructType {
    if t.Kind() != K.Struct {
        return nil
    }
    // Get underlying type first (handles named types in TinyGo)
    ut := t.underlying()
    return (*StructType)(unsafe.Pointer(ut))
}
```

### Step 4: Update `NumField()` - ALREADY USES StructType()

Current code already uses `StructType()`:
```go
func (t *Type) NumField() (int, error) {
    if t.Kind() != K.Struct {
        return 0, Err(ref, D.Numbers, D.Fields, D.Type, "Struct")
    }
    st := (*StructType)(unsafe.Pointer(t))  // ← Change this
    return st.numFields(), nil
}
```

**Fix:**
```go
func (t *Type) NumField() (int, error) {
    if t.Kind() != K.Struct {
        return 0, Err(ref, D.Numbers, D.Fields, D.Type, "Struct")
    }
    st := t.StructType()  // ← Use helper instead of direct cast
    if st == nil {
        return 0, Err(ref, D.Numbers, D.Fields, D.NotOfType, "Struct")
    }
    return st.numFields(), nil
}
```

### Step 5: Fix `Field()` Method in TypeOf.go

Currently still uses direct array access:
```go
func (t *Type) Field(i int) (StructField, error) {
    if t.Kind() != K.Struct {
        return StructField{}, Err(ref, D.Field, D.NotOfType, "Struct")
    }
    st := (*StructType)(unsafe.Pointer(t))
    if i < 0 || i >= len(st.Fields) {  // ← WRONG: uses Fields directly
        return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)
    }
    return st.Fields[i], nil  // ← WRONG: uses Fields directly
}
```

**Fix:**
```go
func (t *Type) Field(i int) (StructField, error) {
    if t.Kind() != K.Struct {
        return StructField{}, Err(ref, D.Field, D.NotOfType, "Struct")
    }
    st := t.StructType()  // Use helper
    if st == nil {
        return StructField{}, Err(ref, D.Field, D.NotOfType, "Struct")
    }
    if i < 0 || i >= st.numFields() {  // Use helper
        return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)
    }
    f := st.getField(i)  // Use helper
    if f == nil {
        return StructField{}, Err(ref, D.Field, D.Out, D.Of, D.Range)
    }
    return *f, nil
}
```

### Step 6: Search for All Direct StructType Casts

Need to find and fix all places that cast directly to StructType:
```bash
grep -n "(*StructType)(unsafe.Pointer" *.go
```

Replace all with calls to `t.StructType()` or `t.underlying()`.

## Files to Create/Modify

### New Files
1. **`TypeOf_tinygo.go`** - TinyGo-specific underlying() implementation
2. **`TypeOf_stdlib.go`** - Stdlib-specific underlying() implementation (trivial)

### Files to Modify
1. **`TypeOf.go`** - Update StructType(), NumField(), NameByIndex(), Field()
2. **`ValueOf.go`** - Update all direct StructType casts

### Files Already Fixed
- ✅ `StructType_stdlib.go` - Stdlib layout
- ✅ `StructType_tinygo.go` - TinyGo layout

## Verification Steps

1. **Verify underlying() works:**
   - Add debug logging in underlying()
   - Check if isNamed() returns true for TestStruct
   - Check if elem pointer is valid

2. **Test with TinyGo:**
   ```bash
   cd example/pwa
   tinygo build -target=wasm -o main.wasm main.wasm.go
   # Test in browser - should show 5 fields
   ```

3. **Test with stdlib:**
   ```bash
   go test ./...  # Should still pass
   ```

## Success Criteria

- [ ] TinyGo WASM shows "Found 5 fields" in console
- [ ] All field values display correctly
- [ ] Stdlib tests still pass
- [ ] No panics or errors in either environment
- [ ] Binary size remains under 50KB

## Expected Issues & Solutions

### Issue: Don't know tflagNamed value
**Solution:** Check TinyGo source or use runtime inspection

### Issue: elemType struct layout unclear
**Solution:** Add debug logging to verify memory layout

### Issue: Breaking stdlib
**Solution:** Ensure underlying() returns self in stdlib

## Debug Strategy

Add temporary logging to verify:
1. Is type named? (isNamed())
2. What is elem pointer value?
3. What is numField value before/after underlying()?
4. Memory addresses and offsets

```go
// Debug helper
func (t *Type) debugType() string {
    return fmt.Sprintf("Type@%p Kind=%v Named=%v", t, t.Kind(), t.isNamed())
}
```
