# Plan Simple: Solo Cambiar Slice a Array Fijo

## Situación Actual
- El código ya tiene la estructura `TinyReflect` con `structCache []structCacheEntry`
- Usa `make([]structCacheEntry, tr.maxStructs)` (slice dinámico)
- **PROBLEMA**: Esto no es compatible con TinyGo que prefiere arrays fijos

## Objetivo Único
Cambiar **SOLO** el slice dinámico a array fijo y medir performance. Si mejora, terminar ahí.

## Cambio Específico

**ACTUAL:**
```go
type TinyReflect struct {
    structCache []structCacheEntry  // ❌ Slice dinámico
    structCount int32
    cacheLock   int32
    log         func(msgs ...any)
    maxStructs  int32               // ❌ Ya no necesario
}
```

**PROPUESTO:**
```go
type TinyReflect struct {
    structCache [128]structCacheEntry  // ✅ Array fijo
    structCount int32
    cacheLock   int32  
    log         func(msgs ...any)
    // maxStructs ELIMINADO - array tiene tamaño fijo
}
```

## Cambios en Constructor

**ACTUAL:**
```go
func New(args ...any) *TinyReflect {
    tr := &TinyReflect{
        maxStructs: StructSize128,
        log:        func(...any) {},
    }
    // ...
    tr.structCache = make([]structCacheEntry, tr.maxStructs)  // ❌ Allocation
    return tr
}
```

**PROPUESTO:**
```go
func New(args ...any) *TinyReflect {
    tr := &TinyReflect{
        log: func(...any) {},
        // structCache ya inicializado como array (zero value)
    }
    // Solo procesar args, NO make()
    for _, arg := range args {
        switch v := arg.(type) {
        case func(...any):
            tr.log = v
        }
    }
    return tr
}
```

## Cambios en Funciones que Usan el Cache

Las funciones que actualmente usan `len(tr.structCache)` o `tr.maxStructs` necesitan ajuste:

**ACTUAL:**
```go
// En alguna función que busca en cache
for i := 0; i < len(tr.structCache); i++ {
    // ...
}
```

**PROPUESTO:**
```go
// Funciona igual - len() funciona con arrays
for i := 0; i < len(tr.structCache); i++ {
    // ...
}
```

## Beneficios de Este Cambio Simple

### **Rendimiento**
- ✅ Zero allocations en constructor (elimina `make()`)
- ✅ Layout de memoria predecible
- ✅ Mejor cache locality

### **Compatibilidad**  
- ✅ 100% compatible con TinyGo
- ✅ Sin allocaciones dinámicas
- ✅ WebAssembly optimizado

### **Simplicidad**
- ✅ Cambio mínimo, riesgo mínimo
- ✅ Mismo algoritmo de búsqueda
- ✅ Fácil de validar con benchmarks

## Implementación

**Paso único:** 
1. Cambiar `structCache []structCacheEntry` → `structCache [128]structCacheEntry`
2. Eliminar `maxStructs` field y `make()` en constructor
3. Ejecutar benchmarks para medir mejora

**Si este cambio simple mejora el performance → TERMINAR AHÍ.**

**Solo si no mejora lo suficiente → considerar optimizaciones adicionales.**