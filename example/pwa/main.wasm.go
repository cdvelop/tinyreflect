//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	. "github.com/cdvelop/tinyreflect"
	. "github.com/cdvelop/tinystring"
)

// NestedTestStruct is a simple nested struct for testing.
type NestedTestStruct struct {
	NestedString string
	NestedInt    int
}

// TestStruct is a helper struct for testing reflection operations.
type TestStruct struct {
	StringField            string
	BoolField              bool
	IntField               int
	Int8Field              int8
	Int16Field             int16
	Int32Field             int32
	Int64Field             int64
	UintField              uint
	Uint8Field             uint8
	Uint16Field            uint16
	Uint32Field            uint32
	Uint64Field            uint64
	Float32Field           float32
	Float64Field           float64
	StringSliceField       []string
	BoolSliceField         []bool
	ByteSliceField         []byte
	IntSliceField          []int
	Int8SliceField         []int8
	Int16SliceField        []int16
	Int32SliceField        []int32
	Int64SliceField        []int64
	UintSliceField         []uint
	Uint8SliceField        []uint8
	Uint16SliceField       []uint16
	Uint32SliceField       []uint32
	Uint64SliceField       []uint64
	Float32SliceField      []float32
	Float64SliceField      []float64
	StructField            NestedTestStruct
	StructSliceField       []NestedTestStruct
	StringIntMapField      map[string]int
	IntStringMapField      map[int]string
	StringIntMapSliceField []map[string]int
	StringPtrField         *string
	IntPtrField            *int
}

func main() {
	// Crear el elemento div
	dom := js.Global().Get("document").Call("createElement", "div")

	buf := Convert().
		Write("<h1>TinyReflect WebAssembly Test</h1>").
		Write("<h2>Testing Field Reflection:</h2>").
		Write("<div id='results'>Running tests...</div>")

	dom.Set("innerHTML", buf.String())

	// Obtener el body del documento y agregar el elemento
	body := js.Global().Get("document").Get("body")
	body.Call("appendChild", dom)

	logger := func(msg ...any) {
		js.Global().Get("console").Call("log", Translate(msg...).String())
	}

	logger("Starting TinyReflect WebAssembly test")

	// Initialize all fields with test values
	strVal := "test string"
	intVal := 42
	data := TestStruct{
		StringField:            strVal,
		BoolField:              true,
		IntField:               intVal,
		Int8Field:              8,
		Int16Field:             16,
		Int32Field:             32,
		Int64Field:             64,
		UintField:              1,
		Uint8Field:             8,
		Uint16Field:            16,
		Uint32Field:            32,
		Uint64Field:            64,
		Float32Field:           32.32,
		Float64Field:           64.64,
		StringSliceField:       []string{"a", "b"},
		BoolSliceField:         []bool{true, false},
		ByteSliceField:         []byte("bytes"),
		IntSliceField:          []int{1, 2, 3},
		Int8SliceField:         []int8{1, 2},
		Int16SliceField:        []int16{1, 2},
		Int32SliceField:        []int32{1, 2},
		Int64SliceField:        []int64{1, 2},
		UintSliceField:         []uint{1, 2},
		Uint8SliceField:        []uint8{1, 2},
		Uint16SliceField:       []uint16{1, 2},
		Uint32SliceField:       []uint32{1, 2},
		Uint64SliceField:       []uint64{1, 2},
		Float32SliceField:      []float32{1.1, 2.2},
		Float64SliceField:      []float64{1.1, 2.2},
		StructField:            NestedTestStruct{NestedString: "nested", NestedInt: 99},
		StructSliceField:       []NestedTestStruct{{NestedString: "n1", NestedInt: 1}, {NestedString: "n2", NestedInt: 2}},
		StringIntMapField:      map[string]int{"key": 100},
		IntStringMapField:      map[int]string{42: "value"},
		StringIntMapSliceField: []map[string]int{{"k1": 1}, {"k2": 2}},
		StringPtrField:         &strVal,
		IntPtrField:            &intVal,
	}

	logger("TestStruct initialized with", len([]interface{}{}), "fields")

	// Test reflection functionality
	v := ValueOf(data)
	typ := v.Type()

	numFields, err := typ.NumField()
	if err != nil {
		logger("ERROR: NumField failed:", err)
		return
	}

	logger("SUCCESS: Found", numFields, "fields")

	// Update DOM with results count
	resultsDiv := js.Global().Get("document").Call("getElementById", "results")
	resultsDiv.Set("innerHTML", Translate("Found ", numFields, " fields. Check console for details.").String())

	// Test field access and print field information
	for i := 0; i < numFields; i++ {
		field, err := v.Field(i)
		if err != nil {
			logger("ERROR: Field(", i, ") failed:", err)
			continue
		}

		fieldName, err := typ.NameByIndex(i)
		if err != nil {
			logger("ERROR: NameByIndex(", i, ") failed:", err)
			continue
		}

		// Get field value as string for logging
		var valueStr string
		kindStr := field.Kind().String()
		switch {
		case kindStr == "string":
			valueStr = "string:" + field.String()
		case kindStr == "bool":
			if b, err := field.Bool(); err == nil {
				valueStr = "bool:" + Translate(b).String()
			} else {
				valueStr = "bool:(error)"
			}
		case kindStr == "int", kindStr == "int8", kindStr == "int16", kindStr == "int32", kindStr == "int64":
			if i, err := field.Int(); err == nil {
				valueStr = "int:" + Translate(i).String()
			} else {
				valueStr = "int:(error)"
			}
		case kindStr == "uint", kindStr == "uint8", kindStr == "uint16", kindStr == "uint32", kindStr == "uint64":
			if u, err := field.Uint(); err == nil {
				valueStr = "uint:" + Translate(u).String()
			} else {
				valueStr = "uint:(error)"
			}
		case kindStr == "float32", kindStr == "float64":
			if f, err := field.Float(); err == nil {
				valueStr = "float:" + Translate(f).String()
			} else {
				valueStr = "float:(error)"
			}
		case kindStr == "slice":
			if length, err := field.Len(); err == nil {
				valueStr = "slice:len=" + Translate(length).String()
			} else {
				valueStr = "slice:(error)"
			}
		case kindStr == "map":
			if length, err := field.Len(); err == nil {
				valueStr = "map:len=" + Translate(length).String()
			} else {
				valueStr = "map:(error)"
			}
		case kindStr == "ptr":
			// For pointers, just show the type
			valueStr = "ptr:" + field.Type().String()
		case kindStr == "struct":
			valueStr = "struct:" + field.Type().String()
		default:
			valueStr = "other:" + kindStr
		}

		logger("Field", i, ":", fieldName, "=", valueStr)
	}

	logger("TinyReflect WebAssembly test completed successfully!")

	select {}
}
