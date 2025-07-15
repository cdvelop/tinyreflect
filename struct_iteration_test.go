package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// Test_StructIteration recorre los campos de una estructura para verificar nombres y valores.
func Test_StructIteration(t *testing.T) {

	// Estructura de prueba con tipos soportados
	type sampleStruct struct {
		Name    string
		Age     int
		Active  bool
		Balance float64
	}

	// Instancia de la estructura con valores
	instance := sampleStruct{
		Name:    "Test Name",
		Age:     30,
		Active:  true,
		Balance: 123.45,
	}

	// Valor reflejado de la instancia
	v := tinyreflect.ValueOf(instance)

	// Mapa para verificar los resultados esperados
	expectedFields := map[string]interface{}{
		"Name":    "Test Name",
		"Age":     int(30),
		"Active":  true,
		"Balance": float64(123.45),
	}

	typ := v.Type()
	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField() devolvió un error inesperado: %v", err)
	}

	// Iterar sobre los campos de la estructura
	for i := 0; i < numFields; i++ {
		// Obtener el nombre del campo por índice
		fieldName, err := typ.NameByIndex(i)
		if err != nil {
			t.Errorf("NameByIndex(%d) devolvió un error inesperado: %v", i, err)
			continue
		}

		// Obtener el valor del campo de la estructura por índice
		fieldValue, err := v.Field(i)
		if err != nil {
			t.Errorf("Field(%d) devolvió un error inesperado: %v", i, err)
			continue
		}

		// Verificar que el campo existe en el mapa de valores esperados
		expectedValue, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Campo '%s' no esperado en la estructura de prueba", fieldName)
			continue
		}

		// Obtener el valor del campo y compararlo
		value, err := fieldValue.Interface()
		if err != nil {
			t.Errorf("field.Interface() devolvió un error inesperado para el campo '%s': %v", fieldName, err)
			continue
		}

		if value != expectedValue {
			t.Errorf("Valor incorrecto para el campo '%s'. Esperado: %v, Obtenido: %v", fieldName, expectedValue, value)
		}

		// Eliminar el campo del mapa para asegurar que todos los campos son evaluados
		delete(expectedFields, fieldName)
	}

	// Si el mapa no está vacío, significa que no todos los campos esperados fueron encontrados
	if len(expectedFields) > 0 {
		for name := range expectedFields {
			t.Errorf("Campo esperado '%s' no fue encontrado en la iteración", name)
		}
	}
}
