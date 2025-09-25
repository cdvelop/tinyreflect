## Performance

TinyReflect provides **transparent caching** that dramatically improves performance for repeated operations on the same struct types. Below are benchmark results comparing TinyReflect with Go's standard `reflect` package:

| Operation | Standard `reflect` | TinyReflect (First Call) | TinyReflect (Cached) | Improvement |
|-----------|-------------------|--------------------------|----------------------|-------------|
| `TypeOf()` | 1.18 ns/op | 8,737 ns/op | **32.35 ns/op** | Standard faster (~27× slower cached) |
| Field Access (4 fields) | 8.96 ns/op | 8,823 ns/op | 13.05 ns/op | Within ~1.5× of standard |
| Field Iteration (6 fields) | 139.4 ns/op | 36.66 ns/op | **40.54 ns/op** | **3.4× faster** |

**Key Insights:**
- **First-time operations** remain orders of magnitude slower while the struct schema is materialized.
- **Cache hits** eliminate repeated allocations and bring most operations close to standard reflect performance (field access) or faster (iteration-heavy workloads).
- **Iteration workloads** benefit the most, showing ~3.4× speedups once cached.
- **Go 1.24.4 improvements**: Field access performance improved slightly (~1.6× → ~1.5× slower than standard).
- **Best fit**: Scenarios that reuse struct types repeatedly (serialization, patching, diffing) where predictable latency matters more than absolute first-hit speed.

### Recent Optimizations (v1.1)

**Performance Improvements Made:**
- **Fast-path cache checking**: `TypeOf()` now checks cache before expensive operations
- **Eliminated string comparisons**: Replaced `Kind().String() == "ptr"` with direct `Kind() == K.Pointer` comparisons  
- **Optimized struct caching**: Reduced lock contention and expensive `Interface()` calls
- **Smart pointer handling**: More efficient dereferencing for `StructNamer` interface detection
- **Cached field operations**: Added `NumFieldCached()` and `NameByIndexCached()` methods for better performance

**Results**: ~20% improvement in `TypeOf()` performance (32× → 27× slower than standard reflect).

### Go 1.24.4 Performance Update

**Additional improvements with Go 1.24.4:**
- **Field Access performance**: Slight improvement from ~1.6× to ~1.5× slower than standard reflect
- **Consistent iteration performance**: Maintains ~3.4× faster field iteration compared to standard reflect
- **Overall stability**: Performance remains consistent across multiple runs with improved compiler optimizations

> Benchmarks run on Intel Core i7-11800H, Go 1.24.4. Results may vary by hardware and Go version.