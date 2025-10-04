# **Technical Report on Struct Field Access in TinyGo WASM via Hybrid RTTI and Unsafe Memory Manipulation**

## **I. Executive Summary: The TinyGo Reflection Conundrum**

The objective of accessing struct field values and names reliably using the Go reflect package when compiled by TinyGo targeting WebAssembly (WASM), while maintaining a minimal binary size (approximately 50 KB), presents a fundamental conflict between Go's run-time dynamism and TinyGo's aggressive size optimization goals. The current implementation failure, characterized by corrupted field names (e.g., Int8Field�  Int1) and the critical RuntimeError: memory access out of bounds \[User Query\], is a direct consequence of TinyGo’s mandatory Run-Time Type Information (RTTI) stripping.1

Standard reflection methods such as reflect.Value.Field(i) cannot be guaranteed under TinyGo’s constraint model. The analysis confirms that TinyGo’s maximum size optimization (-opt=z) successfully achieves the small binary footprint 3 but simultaneously invalidates the memory pointers intended for non-critical string metadata, such as field names. The observed crash is the result of the reflection logic attempting to dereference a corrupted pointer that lies outside the WASM linear memory boundaries.

The definitive solution requires a low-level, size-conscious strategy: abandoning the reliance on dynamic reflection for value retrieval and replacing it with hybrid access. This approach leverages the *stable* structural metadata retained by TinyGo's RTTI—specifically the field count and the byte offset (reflect.StructField.Offset)—and combines it with unsafe package primitives (unsafe.Pointer, uintptr arithmetic) to manually calculate and safely access the precise memory location of each field value within the WASM memory segment.

## **II. Deep Diagnostic Analysis: RTTI Stripping and WASM Memory Boundary Violation**

The diagnostic analysis focuses on two primary failure modes: the corruption of field name metadata and the subsequent catastrophic memory access out of bounds error, both rooted in the compiler’s optimization targets.

### **2.1. The TinyGo Constraint Model (TGC-1) and RTTI Stripping**

TinyGo is explicitly designed for constrained environments, such as microcontrollers and minimal WASM runtimes. This mandate necessitates substantial deviations from the standard Go compiler (gc), particularly concerning the reflect package, which is known for introducing significant runtime overhead and complexity.5 TinyGo provides only limited support for

reflect; packages reliant on deep reflection, such as encoding/json, frequently panic at runtime or exhibit unexpected behavior.1

To achieve the user's size requirement of 50 KB, TinyGo requires aggressive optimization, defaulting to or mandating the \-opt=z flag.3 This flag instructs the LLVM backend to prioritize code size over execution speed, leading to the aggressive stripping of unnecessary data. In the context of RTTI, large data structures, particularly string literals representing field names, are prime targets for removal by the linker if the compiler determines they are not strictly required by the compiled code flow.

The system observation of successfully identifying the total number of fields (Found: 5 fields) confirms that the core type descriptor structure (the rtype) and the array of field definitions are fundamentally present and correctly indexed. However, the data *contained* within those field definitions, specifically the string pointers for names, has been nullified or overwritten during the linking process. This illustrates that TinyGo preserves the structural skeleton of the RTTI (counts, offsets, type indices) but discards the heavy string payload (names, package paths) to minimize binary size.

### **2.2. Analysis of Field Name Corruption and RTTI Misalignment**

The observed output trace provides empirical evidence of metadata corruption: DEBUG NameByIndex: name \= Int8Field�  Int1 \[User Query\]. This sequence of characters, followed by null terminators (�) and extraneous data, is characteristic of reading memory from an invalid or partial pointer. The pointer intended for the field name string likely points to a garbage address, an area of the WASM heap that was partially initialized, or a region immediately adjacent to another small, retained literal.

The mechanism of failure occurs within the reflection logic, specifically in functions like NameByIndex. The internal TinyGo reflect.StructField structure, which contains Name string, Offset uintptr, and Type Type 7, uses a pointer to an external memory segment for the string data. When the reflection mechanism attempts to read the

Name field’s underlying (ptr, len) tuple, the ptr component is invalid due to optimization stripping.

The persistence of the structural skeleton (correct field count and traversal) alongside the corruption of the name pointers dictates a key architectural conclusion: the system should rely solely on the structural data elements that are essential for memory layout, namely the field offset and the field type identifier, and must not trust any associated string pointers managed by the minimal RTTI runtime (TGC-4).

### **2.3. Deconstructing the RuntimeError: memory access out of bounds (TGC-2)**

The ultimate failure, RuntimeError: memory access out of bounds, is a definitive WebAssembly security boundary violation. WASM mandates a sandbox environment where all memory access is confined to a finite, linear memory array.8 When a Go pointer (

uintptr) is dereferenced, the runtime verifies that the offset is within the allocated memory bounds.

The crash occurs during the attempted name retrieval for field 4, specifically when the implementation attempts to access the memory pointed to by the corrupted f.Name.Bytes \[User Query\]. If the TinyGo linker zeroed out the string pointer component of StructField (or set it to a small, invalid offset), attempting to read bytes from this location (which often falls outside the initialized heap, or into a reserved low address) triggers the bounds error.9

This confirms that the existing implementation of tinyreflect is fundamentally flawed not in its logic for iterating fields, but in its dependency on the integrity of RTTI string pointers preserved by the TinyGo compiler under maximal optimization. To proceed, this dependency must be completely eliminated, necessitating a shift to direct memory addressing using the retained Offset value.

## **III. The Engineering Imperative: A Hybrid RTTI \+ unsafe Solution**

To achieve the twin goals of functional reflection (accessing values) and binary size minimization (50 KB), the architectural approach must transition from high-level, dynamic reflection to controlled, low-level memory manipulation guided by static RTTI metadata.

### **3.1. The Unsafe Mandate for Minimal Binaries (TGC-3)**

Standard reflection is inherently slower because it necessitates runtime type inspection and method resolution, bypassing compile-time optimizations.5 This performance and size cost is unacceptable in a constrained WASM environment. The

unsafe package provides the necessary means to bypass Go's compile-time type safety checks, allowing direct manipulation of memory addresses.10

The compiled output of unsafe operations—primarily pointer arithmetic—is extremely minimal and highly efficient, translating directly into optimized low-level instructions in the WASM binary. By using unsafe, the implementation retrieves field values without relying on the sophisticated, heavy, and often incomplete TinyGo reflection runtime methods (like Value.Field(i).Interface()), thereby maintaining the required minimal size footprint.

### **3.2. Go unsafe Primitives in the WASM Context**

In the WASM execution model, a Go uintptr is an offset into the module's linear memory, not a traditional host OS memory address.8 Therefore, pointer arithmetic performed using

uintptr operates directly within the WASM sandbox boundaries. Precision in using unsafe is paramount to avoid new memory corruption or boundary errors.

The following table details the essential unsafe primitives required for this hybrid approach:

Go unsafe Primitives for Struct Access in TinyGo WASM

| Primitive | Purpose in Reflection Bypass | WASM Context Consideration |
| :---- | :---- | :---- |
| unsafe.Pointer | Serves as the universal bridge between typed pointers (e.g., \*MyStruct) and the numerical uintptr memory offset.10 | Must be used only for immediate type casting to satisfy Go’s rules regarding pointer conversions. |
| uintptr | The unsigned integer type that holds the memory offset. This is the only type on which pointer arithmetic (addition of the field offset) can be legally performed.10 | Represents the direct offset into the WASM linear memory segment. |
| unsafe.Offsetof | A compiler intrinsic used to determine the exact byte offset of a field from the start of its containing struct instance.11 | This intrinsic correctly accounts for the compiler’s decisions regarding struct padding and alignment, which are critical for structural integrity.13 |
| unsafe.Add (Go 1.17+) | A utility function that performs safe pointer addition: unsafe.Pointer(uintptr(ptr) \+ uintptr(len)).14 | Used to calculate the target field's memory address by adding the reliable RTTI offset to the struct's base pointer. |

## **IV. Implementation Strategy I: Reliable Field Name Management**

Since the analysis confirms that field names cannot be reliably accessed via reflect.StructField.Name due to RTTI stripping, field name management must be decoupled from the runtime reflection process.

### **4.1. Decoupling Names for Size Minimization (TGC-4)**

The pursuit of a 50 KB binary inherently means accepting the loss of embedded string metadata. If the application requires field names (e.g., for generating JSON or displaying debugging information), the most robust and size-minimal method involves managing these names externally to the reflection runtime:

1. **Compile-Time Index Mapping:** This is the recommended practice for minimal binaries. Instead of relying on the stripped RTTI data, the developer should use automated tools (like Go code generation via go generate) to create a separate, statically compiled map. This map would link the target struct type to a fixed, compile-time array of string literals representing the field names in their correct declaration order.15 The  
   tinyreflect library would then access fields solely by index, and use this external map to retrieve the corresponding name. This method ensures name retention while preserving aggressive RTTI stripping for maximum size reduction.  
2. **Struct Tags:** An alternative, though potentially heavier, approach is to store the names within Go struct tags (e.g., Field string \\tn:"Field"\`). If the TinyGo compiler preserves the StructField.Tag\` string literal, the name can be parsed from this field.7 However, this method adds string data back into the binary, increasing size, and still depends on the compiler retaining the tag data, which is not guaranteed under maximum optimization.

For a minimal binary targeting 50 KB, the reliance on indices (obtained reliably from t.NumField() and t.Field(i)) combined with an external, statically generated name map provides the necessary stability and control over data inclusion.

## **V. Implementation Strategy II: Guaranteed Field Value Access via Pointer Arithmetic**

The core technical solution involves using the reliable RTTI offset data to perform direct pointer arithmetic, allowing the memory address of the field value to be calculated and cast to the correct type.

### **5.1. Obtaining the Base Address**

To access any field, the absolute starting address of the struct instance within the WASM linear memory must be established. If the target struct instance is defined as s, the procedure is as follows:

Step 1: Obtain the Base Pointer. The address of the struct s is acquired as an untyped pointer:

Step 2: Convert to Calculable Offset. The unsafe.Pointer must be converted to uintptr, the numerical type that allows mathematical operations (pointer arithmetic) in Go:

The baseAddr now represents the starting byte offset of the struct instance within the WASM memory.10

### **5.2. Calculating the Field Address**

The next stage requires adding the field's compile-time offset to the base address. This offset is reliably retained within the RTTI structure, even when other RTTI metadata is stripped.

Step 3: Retrieve the Reliable Offset. The reflect package, while unreliable for string names, reliably provides the field offset for index i:

Step 4: Calculate the Target Field Address. The field's absolute memory offset is calculated by summing the base address and the field offset:

**Step 5: Convert Back to Pointer.** The numerical address must be immediately converted back into an unsafe.Pointer to ensure subsequent operations adhere to Go’s pointer safety rules, particularly regarding potential garbage collection interference.10

This targetPtr is guaranteed to point precisely to the starting byte of the field value within the WASM linear memory.11

### **5.3. Reading Primitive Field Values**

Once targetPtr is acquired, the final operation is a type cast followed by dereferencing. This technique completely bypasses the complex internal logic of reflect.Value.Field(i) and directly accesses the memory block.

For instance, to read a field known to be an 8-bit integer (int8):

This memory access pattern ensures that the correct value is retrieved, provided the RTTI correctly reported the field type (which can be derived from ).

### **5.4. Advanced Type Handling: Strings, Slices, and Pointers (TGC-5)**

Handling composite types like strings and slices requires careful attention to TinyGo’s specific memory representation, which deviates from standard Go for compatibility with constrained architectures.1

In TinyGo WASM, a string or slice is represented by a header struct that contains a data pointer and a length/capacity value. Critically, while standard Go uses int for length and capacity, TinyGo uses uintptr for these fields (even when compiling for WASM) to maintain consistency with its support for architectures like AVR:

Go

type StringHeader struct {  
    Data uintptr  
    Len  uintptr // TinyGo uses uintptr here, not int  
}

1

**Reading a String Field (Protocol for Value Retrieval):**

1. **Locate the Header:** Use the calculated targetPtr (from Section 5.2) which points to the start of the StringHeader structure.  
2. Cast to Header Type: Cast the targetPtr to the specific TinyGo structure pointer:

3. **Read String Data:** Access the fields of the header. The strHeaderPtr.Data provides the absolute WASM memory offset where the string bytes begin, and strHeaderPtr.Len provides the number of bytes.8  
4. **Reconstruct the String:** The unsafe.Slice function should be used to create a byte slice (byte) from the Data pointer and Len value, which can then be converted to a string. This methodology, requiring knowledge of TinyGo’s internal StringHeader layout, is essential to avoid secondary corruption or memory bounds violations when dealing with complex types.

By integrating the reliable offset RTTI metadata with precise pointer arithmetic, the system gains stable and size-efficient access to both primitive and complex field values, circumventing the integrity issues associated with higher-level reflection API calls in a constrained WASM environment.

## **VI. System Hardening and Compiler Configuration**

Achieving stability and the required size target demands a specific, immutable set of compiler flags and supplementary measures to manage the WASM environment.

### **6.1. Definitive TinyGo Build Options for Minimalism (TGC-1)**

The user’s current binary size of 50 KB implies the utilization of maximum optimization. It is critical to enforce the most aggressive configuration to ensure RTTI stripping remains consistent and predictable.

Recommended TinyGo Build Options for Minimal WASM

| Flag | Value/Setting | Purpose |
| :---- | :---- | :---- |
| \-target | wasm or wasip1 | Selects the correct WebAssembly target environment and associated runtime hooks.16 |
| \-opt | z (aggressive size optimization) | This flag is mandatory for achieving the 50 KB size goal. It maximally strips unused code and metadata, enforcing the RTTI string stripping that necessitates the unsafe approach.3 |
| \--no-debug | Enabled | Removes DWARF debug symbols and other linked debug information, drastically reducing the final binary size.16 |
| \-scheduler | none | Eliminates the entire Go scheduler runtime if the application does not utilize goroutines or concurrency, contributing to a smaller binary.18 |

### **6.2. Mitigation of Memory Boundary Issues**

The RuntimeError: memory access out of bounds can, in peripheral cases, be related to the underlying configuration of the WASM runtime’s linear memory.9 While the primary cause in this instance is an invalid pointer derived from RTTI corruption, general system stability requires mitigating potential memory allocation failures.

In environments where the WASM module is embedded (e.g., using Wazero), ensuring the host application allows for adequate initial memory or dynamic memory growth can prevent secondary failures.9 However, within the Go code compiled to WASM, the most critical safeguard is the precise handling of

unsafe.Pointer conversions. Pointers derived from uintptr arithmetic must be immediately type-cast and dereferenced. This discipline minimizes the window during which the TinyGo Garbage Collector (GC) might move or free the underlying struct s while the manual address offset is being manipulated, which could otherwise lead to memory corruption or bounds errors when the offset is later accessed.10

## **VII. Conclusion and Refactoring Summary**

The challenge of implementing dynamic struct field access within a highly constrained TinyGo WASM environment, especially when size optimization (-opt=z) is paramount, fundamentally reveals the trade-off between binary size and runtime reflection fidelity. The system analysis confirmed that TinyGo’s optimization regime renders RTTI metadata related to field names unreliable, leading directly to memory boundary violations.

### **Refactoring Recommendations for tinyreflect**

To successfully achieve the stated objective of accessing field values while maintaining the 50 KB binary size, the tinyreflect library must be refactored from a traditional reflection wrapper into a **metadata-guided memory access tool**.

1. **Eliminate Reliance on Name Pointers:** All code paths attempting to read reflect.StructField.Name must be removed. Field identification must rely on static indices or an external compile-time map.  
2. Implement Hybrid Access for Values: High-level calls like reflect.Value.Field(i).Interface() must be replaced entirely with the low-level unsafe pointer arithmetic detailed in Section V. The mechanism should be structured to perform the following generalized calculation within the library functions:

3. **Account for TinyGo Type Differences:** Specific handlers must be implemented for complex types (strings, slices) that respect TinyGo’s use of uintptr for length and capacity fields within their memory headers (TGC-5).

By focusing on the structural RTTI elements (offset, type) and performing manual memory dereferencing via unsafe primitives, the implementation achieves robust field value access, minimizes the runtime footprint, and successfully navigates the limitations imposed by the TinyGo WASM environment under maximal size optimization. The resultant implementation sacrifices the standard Go developer experience of full, dynamic reflection in favor of deterministic, high-performance, and size-efficient code critical for embedded and WebAssembly targets.

### [**SOURCES CITED**](SOURCES_CITED.md)