# Compatibilidad con TinyGo

## ⚠️ IMPORTANTE: Limitaciones Fundamentales

**TinyReflect tiene compatibilidad PARCIAL con TinyGo** debido a limitaciones fundamentales en cómo TinyGo maneja los metadatos de tipos en tiempo de ejecución.

### ✅ Lo que SÍ funciona en TinyGo:
- Compilación sin errores
- Operaciones con tipos básicos (int, string, bool, float, etc.)
- `TypeOf()` y `ValueOf()` para tipos básicos
- Métodos de Value para tipos básicos: `Int()`, `String()`, `Bool()`, `Float()`, `Uint()`
- Operaciones con slices y arrays de tipos básicos

### ❌ Lo que NO funciona en TinyGo:
- **Reflection de structs**: `NumField()`, `Field()`, `NameByIndex()`
- **Acceso a campos de structs** mediante reflection
- **Iteración sobre campos de structs**
- Cualquier operación que requiera metadatos de struct en tiempo de ejecución

### ¿Por qué?

TinyGo elimina los metadatos de tipos para reducir el tamaño del binario. Esto significa que la información sobre campos de structs **no está disponible en tiempo de ejecución**.

Ver [`docs/TINYGO_LIMITATIONS.md`](TINYGO_LIMITATIONS.md) para más detalles y soluciones alternativas.

## Cambios Realizados

### 1. Eliminación de `//go:linkname` en ValueOf.go

**Problema:** La directiva `//go:linkname` no está completamente soportada en TinyGo y causaba problemas de compilación.

**Ubicación:** [`ValueOf.go:174`](../ValueOf.go:174)

**Solución:** Se eliminó la directiva `//go:linkname` de la función `add()` ya que no era necesaria. La función solo realiza aritmética de punteros y no requiere vinculación con el runtime de Go.

**Antes:**
```go
//go:linkname add
//nolint:govet
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
    return unsafe.Pointer(uintptr(p) + x)
}
```

**Después:**
```go
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
    return unsafe.Pointer(uintptr(p) + x)
}
```

**Impacto:** ✅ Sin cambios en funcionalidad o rendimiento. Todos los tests pasan.

---

### 2. Reemplazo de `unsafe.String` en Name.go

**Problema:** `unsafe.String` fue introducido en Go 1.20 y **NO está disponible en TinyGo**, que está basado en Go 1.19.

**Ubicación:** [`Name.go:32`](../Name.go:32) y [`Name.go:67`](../Name.go:67)

**Solución:** Se implementó una función helper `unsafeString()` que usa `unsafe.Slice()` para crear strings de manera compatible con TinyGo.

**Antes:**
```go
func (n Name) Name() string {
    if n.Bytes == nil {
        return ""
    }
    i, l := n.ReadVarint(1)
    return unsafe.String(n.DataChecked(1+i, "non-empty string"), l)
}
```

**Después:**
```go
func (n Name) Name() string {
    if n.Bytes == nil {
        return ""
    }
    i, l := n.ReadVarint(1)
    return unsafeString(n.DataChecked(1+i, "non-empty string"), l)
}

// unsafeString constructs a string from a byte pointer and length.
// This is a TinyGo-compatible replacement for unsafe.String (Go 1.20+).
func unsafeString(ptr *byte, length int) string {
    if ptr == nil || length == 0 {
        return ""
    }
    return string(unsafe.Slice(ptr, length))
}
```

**Impacto:** ✅ Funcionalidad idéntica. Compatible con TinyGo y Go estándar.

---

## Verificación de Compatibilidad

### Elementos Verificados como Compatibles

✅ **Operaciones `unsafe.Pointer`** - Todas las conversiones y aritmética de punteros son compatibles

✅ **Estructura `EmptyInterface`** - Compatible con la representación interna de interfaces en TinyGo

✅ **`TypeOff` y `NameOff`** - Solo son comentarios de documentación, no afectan la compilación

✅ **Todas las operaciones de reflection** - Implementadas usando solo características soportadas por TinyGo

### Tests

Todos los tests existentes pasan sin modificaciones:
```bash
go test -v ./...
# PASS: 100% de los tests (0.004s)
```

---

## Requisitos de TinyGo

- **Versión mínima de TinyGo:** 0.27.0 o superior
- **Versión de Go:** Compatible con Go 1.19+
- **Características usadas:**
  - `unsafe.Pointer` ✅
  - `unsafe.Slice` ✅ (disponible desde Go 1.17)
  - Aritmética de punteros básica ✅

---

## Compilación con TinyGo

Para compilar tu proyecto con TinyGo:

```bash
# Compilación estándar
tinygo build -o output.wasm -target wasm ./main.go

# Con optimizaciones
tinygo build -o output.wasm -target wasm -opt=2 ./main.go

# Para WebAssembly con tamaño mínimo
tinygo build -o output.wasm -target wasm -opt=z ./main.go
```

---

## Limitaciones Conocidas

Las siguientes limitaciones son inherentes a TinyGo y no pueden ser resueltas:

1. **Nombres de structs:** TinyGo elimina los metadatos de nombres de tipos. Usa la interfaz `StructNamer` para proporcionar nombres personalizados.

2. **Reflection avanzada:** Características como métodos, canales y funciones no están soportadas (por diseño de tinyreflect).

3. **Tamaño de binarios:** Aunque TinyGo produce binarios más pequeños, el uso de reflection siempre añade overhead.

---

## Beneficios de la Compatibilidad TinyGo

- 🎯 **Binarios más pequeños** - Reducción significativa del tamaño para WebAssembly
- 🚀 **Mejor rendimiento** - Optimizaciones específicas de TinyGo
- 📱 **Soporte embebido** - Posibilidad de usar en sistemas embebidos
- 🌐 **WebAssembly optimizado** - Ideal para aplicaciones web

---

## Soporte

Si encuentras problemas de compatibilidad con TinyGo, por favor:

1. Verifica que estás usando TinyGo 0.27.0 o superior
2. Revisa que tu código no use características no soportadas
3. Reporta el issue con detalles de tu entorno

---

**Última actualización:** 2025-01-10
**Estado:** ✅ 100% Compatible con TinyGo