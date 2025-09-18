package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

// BenchmarkValue_InterfaceZeroAlloc benchmarks the new InterfaceZeroAlloc method
// that takes a pointer parameter to avoid return boxing
func BenchmarkValue_InterfaceZeroAlloc(b *testing.B) {
	// Test struct with different primitive types
	type TestStruct struct {
		IntField    int
		StringField string
		BoolField   bool
		FloatField  float64
		SliceField  []int // complex type for comparison
	}

	ts := TestStruct{
		IntField:    42,
		StringField: "benchmark_test_string",
		BoolField:   true,
		FloatField:  3.14159,
		SliceField:  []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := ValueOf(ts)

		// Benchmark InterfaceZeroAlloc - should have minimal allocations
		var target any
		if field0, err := v.Field(0); err == nil {
			field0.InterfaceZeroAlloc(&target)
		}
		if field1, err := v.Field(1); err == nil {
			field1.InterfaceZeroAlloc(&target)
		}
		if field2, err := v.Field(2); err == nil {
			field2.InterfaceZeroAlloc(&target)
		}
		if field3, err := v.Field(3); err == nil {
			field3.InterfaceZeroAlloc(&target)
		}
		if field4, err := v.Field(4); err == nil {
			field4.InterfaceZeroAlloc(&target) // complex type
		}
	}
}

// BenchmarkValue_Interface benchmarks the original Interface method for comparison
func BenchmarkValue_Interface(b *testing.B) {
	// Test struct with different primitive types
	type TestStruct struct {
		IntField    int
		StringField string
		BoolField   bool
		FloatField  float64
		SliceField  []int // complex type for comparison
	}

	ts := TestStruct{
		IntField:    42,
		StringField: "benchmark_test_string",
		BoolField:   true,
		FloatField:  3.14159,
		SliceField:  []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := ValueOf(ts)

		// Benchmark original Interface - will allocate for all types
		if field0, err := v.Field(0); err == nil {
			_, _ = field0.Interface()
		}
		if field1, err := v.Field(1); err == nil {
			_, _ = field1.Interface()
		}
		if field2, err := v.Field(2); err == nil {
			_, _ = field2.Interface()
		}
		if field3, err := v.Field(3); err == nil {
			_, _ = field3.Interface()
		}
		if field4, err := v.Field(4); err == nil {
			_, _ = field4.Interface()
		}
	}
}

// TestValue_InterfaceZeroAlloc tests the new InterfaceZeroAlloc method
func TestValue_InterfaceZeroAlloc(t *testing.T) {
	type TestStruct struct {
		IntField    int
		StringField string
		BoolField   bool
		FloatField  float64
		SliceField  []int
	}

	ts := TestStruct{
		IntField:    42,
		StringField: "test_string",
		BoolField:   true,
		FloatField:  3.14,
		SliceField:  []int{1, 2, 3},
	}

	v := ValueOf(ts)

	// Test int field
	var intTarget any
	if field0, err := v.Field(0); err == nil {
		field0.InterfaceZeroAlloc(&intTarget)
		if intTarget != 42 {
			t.Errorf("IntField: expected 42, got %v", intTarget)
		}
	}

	// Test string field
	var stringTarget any
	if field1, err := v.Field(1); err == nil {
		field1.InterfaceZeroAlloc(&stringTarget)
		if stringTarget != "test_string" {
			t.Errorf("StringField: expected 'test_string', got %v", stringTarget)
		}
	}

	// Test bool field
	var boolTarget any
	if field2, err := v.Field(2); err == nil {
		field2.InterfaceZeroAlloc(&boolTarget)
		if boolTarget != true {
			t.Errorf("BoolField: expected true, got %v", boolTarget)
		}
	}

	// Test float field
	var floatTarget any
	if field3, err := v.Field(3); err == nil {
		field3.InterfaceZeroAlloc(&floatTarget)
		if floatTarget != 3.14 {
			t.Errorf("FloatField: expected 3.14, got %v", floatTarget)
		}
	}

	// Test slice field (complex type - should still work)
	var sliceTarget any
	if field4, err := v.Field(4); err == nil {
		field4.InterfaceZeroAlloc(&sliceTarget)
		// For complex types, we can't easily compare, just check it's not nil
		if sliceTarget == nil {
			t.Errorf("SliceField: expected non-nil value")
		}
	}
}

func TestKindAndCanAddr(t *testing.T) {
	// Test with a non-addressable value
	v := ValueOf(123)
	if v.Kind() != K.Int {
		t.Errorf("Kind for int: expected Int, got %s", v.Kind())
	}
	if v.CanAddr() {
		t.Error("CanAddr for non-addressable value: expected false, got true")
	}

	// Test with an addressable value
	i := 123
	v = ValueOf(&i)
	elem, _ := v.Elem()
	if elem.Kind() != K.Int {
		t.Errorf("Kind for addressable int: expected Int, got %s", elem.Kind())
	}
	if !elem.CanAddr() {
		t.Error("CanAddr for addressable value: expected true, got false")
	}
}

func TestAllGetters(t *testing.T) {
	testCases := []struct {
		value    interface{}
		kind     Kind
		intVal   int64
		uintVal  uint64
		floatVal float64
		boolVal  bool
		strVal   string
		wantErr  bool
	}{
		{int(1), K.Int, 1, 0, 0, false, "", false},
		{int8(2), K.Int8, 2, 0, 0, false, "", false},
		{int16(3), K.Int16, 3, 0, 0, false, "", false},
		{int32(4), K.Int32, 4, 0, 0, false, "", false},
		{int64(5), K.Int64, 5, 0, 0, false, "", false},
		{uint(6), K.Uint, 0, 6, 0, false, "", false},
		{uint8(7), K.Uint8, 0, 7, 0, false, "", false},
		{uint16(8), K.Uint16, 0, 8, 0, false, "", false},
		{uint32(9), K.Uint32, 0, 9, 0, false, "", false},
		{uint64(10), K.Uint64, 0, 10, 0, false, "", false},
		{uintptr(11), K.Uintptr, 0, 11, 0, false, "", false},
		{float32(12.0), K.Float32, 0, 0, 12.0, false, "", false},
		{float64(13.0), K.Float64, 0, 0, 13.0, false, "", false},
		{true, K.Bool, 0, 0, 0, true, "", false},
		{"hello", K.String, 0, 0, 0, false, "hello", false},
	}

	for _, tc := range testCases {
		v := ValueOf(tc.value)

		// Test Int()
		i, err := v.Int()
		if tc.kind >= K.Int && tc.kind <= K.Int64 {
			if err != nil {
				t.Errorf("Int() on %v: unexpected error: %v", tc.value, err)
			}
			if i != tc.intVal {
				t.Errorf("Int() on %v: expected %d, got %d", tc.value, tc.intVal, i)
			}
		} else {
			if err == nil {
				t.Errorf("Int() on %v: expected error, got nil", tc.value)
			}
		}

		// Test Uint()
		u, err := v.Uint()
		if tc.kind >= K.Uint && tc.kind <= K.Uintptr {
			if err != nil {
				t.Errorf("Uint() on %v: unexpected error: %v", tc.value, err)
			}
			if u != tc.uintVal {
				t.Errorf("Uint() on %v: expected %d, got %d", tc.value, tc.uintVal, u)
			}
		} else {
			if err == nil {
				t.Errorf("Uint() on %v: expected error, got nil", tc.value)
			}
		}

		// Test Float()
		f, err := v.Float()
		if tc.kind >= K.Float32 && tc.kind <= K.Float64 {
			if err != nil {
				t.Errorf("Float() on %v: unexpected error: %v", tc.value, err)
			}
			if f != tc.floatVal {
				t.Errorf("Float() on %v: expected %f, got %f", tc.value, tc.floatVal, f)
			}
		} else {
			if err == nil {
				t.Errorf("Float() on %v: expected error, got nil", tc.value)
			}
		}

		// Test Bool()
		b, err := v.Bool()
		if tc.kind == K.Bool {
			if err != nil {
				t.Errorf("Bool() on %v: unexpected error: %v", tc.value, err)
			}
			if b != tc.boolVal {
				t.Errorf("Bool() on %v: expected %v, got %v", tc.value, tc.boolVal, b)
			}
		} else {
			if err == nil {
				t.Errorf("Bool() on %v: expected error, got nil", tc.value)
			}
		}

		// Test String()
		s := v.String()
		if tc.kind == K.String {
			if s != tc.strVal {
				t.Errorf("String() on %v: expected %s, got %s", tc.value, tc.strVal, s)
			}
		} else {
			if s == tc.strVal {
				t.Errorf("String() on %v: expected non-empty string, got empty", tc.value)
			}
		}
	}
}

func TestStringNonString(t *testing.T) {
	// Test with invalid value
	vInvalid := Value{}
	if vInvalid.stringNonString() != "<invalid Value>" {
		t.Errorf("stringNonString for invalid value: expected '<invalid Value>', got %s", vInvalid.stringNonString())
	}
}

func TestElemErrors(t *testing.T) {
	// Test with an interface
	var i interface{} = 123
	v := ValueOf(i)
	_, err := v.Elem()
	if err == nil {
		t.Error("Elem on interface: expected an error, but got nil")
	}
}

func TestNumField(t *testing.T) {
	type S struct {
		A int
		B string
	}
	s := S{}
	v := ValueOf(s)
	n, err := v.NumField()
	if err != nil {
		t.Errorf("NumField on struct: unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("NumField on struct: expected 2, got %d", n)
	}

	// Test with a non-struct
	i := 123
	v = ValueOf(i)
	_, err = v.NumField()
	if err == nil {
		t.Error("NumField on non-struct: expected an error, but got nil")
	}

	// Test with nil type
	v = ValueOf(s)
	v.typ_ = nil
	_, err = v.NumField()
	if err == nil {
		t.Error("NumField with nil type: expected an error, but got nil")
	}
}

func TestField(t *testing.T) {
	type S struct {
		A int
		B string
	}
	s := S{}
	v := ValueOf(s)

	// Test with valid index
	f, err := v.Field(0)
	if err != nil {
		t.Errorf("Field with valid index: unexpected error: %v", err)
	}
	if f.Kind() != K.Int {
		t.Errorf("Field with valid index: expected kind Int, got %s", f.Kind())
	}

	// Test with invalid index
	_, err = v.Field(2)
	if err == nil {
		t.Error("Field with invalid index: expected an error, but got nil")
	}

	// Test on a non-struct
	i := 123
	v = ValueOf(i)
	_, err = v.Field(0)
	if err == nil {
		t.Error("Field on non-struct: expected an error, but got nil")
	}
}

func TestFieldUnexported(t *testing.T) {
	type E struct {
		e int
	}
	type S struct {
		E
		s int
	}
	s := S{}
	v := ValueOf(s)

	// Test with unexported embedded field
	f, err := v.Field(0)
	if err != nil {
		t.Errorf("Field with unexported embedded field: unexpected error: %v", err)
	}
	if f.Kind() != K.Struct {
		t.Errorf("Field with unexported embedded field: expected kind Struct, got %s", f.Kind())
	}

	// Test with unexported field
	f, err = v.Field(1)
	if err != nil {
		t.Errorf("Field with unexported field: unexpected error: %v", err)
	}
	if f.Kind() != K.Int {
		t.Errorf("Field with unexported field: expected kind Int, got %s", f.Kind())
	}
}
