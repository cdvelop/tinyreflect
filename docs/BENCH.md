## Performance

TinyReflect provides **direct reflection operations** that are optimized for performance without caching overhead. Below are benchmark results comparing TinyReflect with Go's standard `reflect` package:

| Operation | Standard `reflect` | TinyReflect | Improvement |
|-----------|-------------------|-------------|-------------|
| `TypeOf()` | 1.119 ns/op | **0.2310 ns/op** | **4.8× faster** |
| Field Access (4 fields) | 8.762 ns/op | **8.808 ns/op** | **Equivalent performance** |
| Field Iteration (6 fields) | 156.6 ns/op | **33.68 ns/op** | **4.8× faster** |

**Key Insights:**
- **TypeOf() operations** are significantly faster than standard reflect due to optimized type detection
- **Field access** performance is essentially equivalent to standard reflect
- **Field iteration** shows dramatic improvements (4.8× faster) due to direct struct field access
- **Zero memory allocations** for all core operations - no GC pressure
- **Predictable performance** - no caching means consistent timing across all calls
- **Best fit**: All reflection scenarios, especially iteration-heavy workloads and applications requiring consistent low-latency performance

### Architecture Benefits (v2.0)

**Performance Improvements with Cache Removal:**
- **Eliminated memory allocations**: No cache structures, atomic operations, or lock contention
- **Simplified call path**: Direct reflection operations without cache lookups
- **Reduced binary size**: Removed complex caching infrastructure
- **Better TinyGo compatibility**: No sync primitives or dynamic allocations
- **Predictable latency**: Consistent performance without cache hit/miss variance

**Results**: TinyReflect now outperforms standard reflect in most operations while maintaining zero-allocation characteristics.

### Benchmark Environment

**Hardware**: Intel Core i7-11800H @ 2.30GHz
**Go Version**: 1.24.4
**Memory**: 0 B/op, 0 allocs/op for all core operations

