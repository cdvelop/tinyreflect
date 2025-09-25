# Race Conditions in TinyReflect

## Summary of the Issue
- `TinyReflect` caches struct schemas in a preallocated slice guarded by hand-rolled atomics.
- `structCount` is incremented before the new cache slot is fully materialized.
- Concurrent readers observe the incremented count and read a partially initialized entry, which triggers the data race reported by Go's race detector.
- The race is observable via the dedicated test `TestStructCacheRace` when running `go test -race`.

## Reproduction Steps
1. Ensure Go's race detector is enabled.
2. Run:
   ```bash
   go test -race -run TestStructCacheRace
   ```
3. The race detector reports conflicting accesses inside `cacheStructSchema`.

## Library Philosophy Alignment
To stay true to TinyReflect's mission (see `README.md`):
- **Minimal API surface** – avoid leaking complex synchronization primitives to users.
- **Predictable performance** – keep steady behaviour with and without TinyGo/WebAssembly builds.
- **TinyGo friendliness** – prefer patterns that compile without warnings under TinyGo.
- **Small binary footprint** – any fix should preserve or minimally impact code size.

## Chosen Mitigation — Defer `structCount` Publication (Two-Phase Commit)
**Idea:** Fill `structCache[newIndex]` completely, then publish the slot with an atomic store (e.g., `atomic.StoreInt32(&structCount, count+1)`). This ensures readers only observe fully initialized cache entries.

- **Pros**
  - Keeps the current data layout and avoids extra allocations.
  - Minimal runtime cost after the first cache miss.
  - Preserves TinyGo friendliness (only uses `sync/atomic`).
  - Matches TinyReflect's philosophy of a tiny, predictable surface area.
- **Cons**
  - Requires careful memory-order reasoning to guarantee happens-before semantics.
  - Slight refactor of initialization flow; mistakes could reintroduce subtle bugs.

## Next Steps
- Implement the two-phase publish flow inside `cacheStructSchema`.
- Re-run `go test -race ./...` (and TinyGo builds) to confirm the race is resolved.
- Document the change in the changelog/release notes if necessary.
