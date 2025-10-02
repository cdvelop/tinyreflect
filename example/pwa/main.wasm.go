//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	. "github.com/cdvelop/tinyreflect"
	. "github.com/cdvelop/tinystring"
)

// TestStruct is a helper struct for testing reflection operations.
type TestStruct struct {
	StringField string
	BoolField   bool
	IntField    int
	Int8Field   int8
	Int16Field  int16
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

	// Initialize all fields with test values
	strVal := "test string"
	intVal := 42
	data := TestStruct{
		StringField: strVal,
		BoolField:   true,
		IntField:    intVal,
		Int8Field:   8,
		Int16Field:  16,
	}

	// First test with library tinyreflect
	logger("=== Testing tinyreflect ===")
	stdV := ValueOf(data)
	stdT := stdV.Type()
	stdNumFields, err := stdT.NumField()
	if err != nil {
		//logger("ERROR: NumField failed:", err)
	}

	logger("Found:", stdNumFields, "fields")

	// Now test with tinyreflect
	logger("=== Testing tinyreflect ===")
	logger("DEBUG: About to call ValueOf")
	v := ValueOf(data)
	logger("DEBUG: ValueOf returned, v.typ_ =", v.Type() != nil)

	typ := v.Type()
	logger("DEBUG: Type() returned, typ =", typ != nil)

	if typ == nil {
		logger("ERROR: Type() returned nil")
		return
	}

	logger("DEBUG: About to call NumField")
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
	logger("DEBUG: Starting field loop, numFields =", numFields)
	for i := 0; i < numFields; i++ {
		logger("DEBUG: Getting field", i)
		field, err := v.Field(i)
		if err != nil {
			logger("ERROR: Field(", i, ") failed:", err)
			continue
		}
		logger("DEBUG: Got field", i, "successfully")

		logger("DEBUG: Getting field name for index", i)
		fieldName, err := typ.NameByIndex(i)
		if err != nil {
			logger("ERROR: NameByIndex(", i, ") failed:", err)
			continue
		}
		logger("DEBUG: Got field name:", fieldName)

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
