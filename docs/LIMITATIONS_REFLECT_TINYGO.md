# Limitations of Struct Name Reflection in TinyGo and the Go Runtime

TinyGo is designed to generate extremely small binaries for WebAssembly and embedded systems. To achieve this, it aggressively removes metadata from the final binary, including struct type names and other debug symbols that are present in the standard Go runtime.

## Why can't you get a struct's name in TinyGo?

- **The standard Go runtime** stores detailed type information, including struct names, to enable advanced reflection operations (`reflect.TypeOf(x).Name()`).
- **TinyGo**, to reduce binary size, removes these metadata. As a result, there is no reliable way to obtain the textual name of a struct at runtime.
- Any attempt to replicate this functionality would require re-implementing complex and undocumented parts of the Go runtime, breaking compatibility and increasing binary size.

## Recommended Alternative

Instead of relying on type names, use unique identifiers like `StructID()` (based on the type hash), which are available and compatible with TinyGo.

## Summary

- **It is not possible** to reliably implement `Type.Name()` for structs in TinyGo.
- Using unique identifiers (`StructID`) is the safe and compatible alternative.
- Keeping the reflection API minimal and focused on truly supported use cases is key for compatibility and small binary size.

---

For more details, see the documentation and README of this project.
