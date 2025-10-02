# Limitaciones de TinyGo con TinyReflect

## Problema Fundamental

**TinyGo NO proporciona metadatos completos de tipos en tiempo de ejecución**, a diferencia de Go estándar.

### ¿Qué significa esto?

Cuando usas `TypeOf()` o `ValueOf()` en Go estándar:
- El runtime proporciona información completa sobre la estructura del tipo
- Los campos de un struct están disponibles automáticamente
- Los metadatos incluyen nombres de campos, tipos, offsets, etc.

Cuando usas `TypeOf()` o `ValueOf()` en TinyGo:
- El runtime proporciona información MUY limitada
- **Los campos de structs NO están disponibles automáticamente**
- TinyGo elimina estos metadatos para reducir el tamaño del binario

### Impacto en TinyReflect

Las siguientes operaciones **NO funcionan** en TinyGo sin modificaciones adicionales:

❌ `Type.NumField()` - No puede obtener el número de campos
❌ `Type.Field(i)` - No puede obtener información de campos
❌ `Type.NameByIndex(i)` - No puede obtener nombres de campos
❌ `Value.Field(i)` - No puede acceder a campos por índice
❌ Cualquier operación que requiera metadatos de struct

### ¿Por qué falla silenciosamente?

En TinyGo, cuando hacemos:
```go
st := (*StructType)(unsafe.Pointer(t))
len(st.Fields) // Esto devuelve 0 o causa panic
```

El slice `Fields` está vacío porque TinyGo no lo inicializa con los metadatos del struct.

## Soluciones Posibles

### Opción 1: Generación de Código (Recomendada)

Usar un generador de código que cree funciones de reflection específicas para cada tipo:

```go
//go:generate tinyreflect-gen

type User struct {
    Name string
    Age  int
}

// El generador crearía automáticamente:
func (u User) TinyReflectFields() []FieldInfo {
    return []FieldInfo{
        {Name: "Name", Type: "string", Offset: 0},
        {Name: "Age", Type: "int", Offset: 16},
    }
}
```

### Opción 2: Registro Manual

Registrar manualmente los metadatos de cada tipo:

```go
func init() {
    tinyreflect.RegisterStruct(User{}, []FieldInfo{
        {Name: "Name", Type: TypeOf(""), Offset: 0},
        {Name: "Age", Type: TypeOf(0), Offset: 16},
    })
}
```

### Opción 3: Usar reflect estándar de TinyGo

TinyGo tiene su propio paquete `reflect` limitado que podría usarse como base.

## Estado Actual

🔴 **TinyReflect NO es compatible con TinyGo para operaciones de struct reflection**

Las siguientes operaciones SÍ funcionan:
✅ `TypeOf()` para tipos básicos (int, string, bool, etc.)
✅ `ValueOf()` para tipos básicos
✅ Operaciones sobre valores básicos (Int(), String(), Bool(), etc.)
✅ Operaciones sobre slices y arrays (si no contienen structs)

## Recomendación

Para usar reflection con TinyGo, considera:

1. **Usar el paquete `reflect` estándar de TinyGo** (limitado pero funcional)
2. **Implementar un generador de código** para tus tipos específicos
3. **Evitar reflection** y usar interfaces o type switches cuando sea posible

## Referencias

- [TinyGo Reflection Limitations](https://tinygo.org/docs/reference/lang-support/#reflection)
- [TinyGo reflect package](https://tinygo.org/docs/reference/lang-support/stdlib/)