//go:build wasm
// +build wasm

package main

import (
	"reflect"
	"syscall/js"

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
		Write("<h1>Reflect WebAssembly Example</h1>").
		Write("<div id='results'>Running tests...</div>")

	dom.Set("innerHTML", buf.String())

	// Obtener el body del documento y agregar el elemento
	body := js.Global().Get("document").Get("body")
	body.Call("appendChild", dom)

	logger := func(msg ...any) {
		js.Global().Get("console").Call("log", Translate(msg...).String())
	}

	// Initialize struct with test values
	data := TestStruct{
		StringField: "test string",
		BoolField:   true,
		IntField:    42,
		Int8Field:   8,
		Int16Field:  16,
	}

	logger("=== Getting field names and values ===")

	v := reflect.ValueOf(data)
	t := v.Type()
	numFields := t.NumField()

	logger("Found", numFields, "fields:")

	// Build HTML output
	htmlOutput := Convert().Write("<ul>")

	for i := 0; i < numFields; i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		fieldValue := field.Interface()

		logger("Field", i, ":", fieldName, "=", fieldValue)

		htmlOutput.Write("<li>Field ").Write(i).Write(": ").Write(fieldName).Write(" = ").Write(fieldValue).Write("</li>")
	}

	htmlOutput.Write("</ul>")

	// Update DOM with results
	resultsDiv := js.Global().Get("document").Call("getElementById", "results")
	resultsDiv.Set("innerHTML", htmlOutput.String())

	logger("Test completed successfully!")

	select {}
}
