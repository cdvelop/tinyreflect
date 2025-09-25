package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
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

func TestTestStructAllFields(t *testing.T) {

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

	v := tinyreflect.ValueOf(data)
	typ := v.Type()

	numFields, err := typ.NumField()
	if err != nil {
		t.Fatalf("NumField failed: %v", err)
	}

	if numFields != 36 {
		t.Errorf("Expected 36 fields, got %d", numFields)
	}

	// Verify we can access all fields without errors
	for i := 0; i < numFields; i++ {
		field, err := v.Field(i)
		if err != nil {
			t.Errorf("Field(%d) failed: %v", i, err)
		}

		fieldName, err := typ.NameByIndex(i)
		if err != nil {
			t.Errorf("NameByIndex(%d) failed: %v", i, err)
		}

		t.Logf("Field %d: %s, Kind: %v", i, fieldName, field.Kind())
	}
}
