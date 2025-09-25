package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// BenchmarkValue_InterfaceZeroAlloc benchmarks the new InterfaceZeroAlloc method
// that takes a pointer parameter to avoid return boxing.
func BenchmarkValue_InterfaceZeroAlloc(b *testing.B) {
	type TestStruct struct {
		IntField    int
		StringField string
		BoolField   bool
	}
	ts := TestStruct{IntField: 42, StringField: "benchmark", BoolField: true}
	tr := tinyreflect.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := tr.ValueOf(ts)
		var target any
		if field, err := v.Field(0); err == nil {
			field.InterfaceZeroAlloc(&target)
		}
		if field, err := v.Field(1); err == nil {
			field.InterfaceZeroAlloc(&target)
		}
		if field, err := v.Field(2); err == nil {
			field.InterfaceZeroAlloc(&target)
		}
	}
}

// TestValue_InterfaceZeroAlloc tests the new InterfaceZeroAlloc method.
func TestValue_InterfaceZeroAlloc(t *testing.T) {
	type TestStruct struct {
		IntField    int
		StringField string
	}
	ts := TestStruct{IntField: 42, StringField: "test"}
	tr := tinyreflect.New()
	v := tr.ValueOf(ts)

	var target any
	if field, err := v.Field(0); err == nil {
		field.InterfaceZeroAlloc(&target)
		if target != 42 {
			t.Errorf("IntField: expected 42, got %v", target)
		}
	}
	if field, err := v.Field(1); err == nil {
		field.InterfaceZeroAlloc(&target)
		if target != "test" {
			t.Errorf("StringField: expected 'test', got %v", target)
		}
	}
}

func TestKindAndCanAddr(t *testing.T) {
	tr := tinyreflect.New()
	v := tr.ValueOf(123)
	if v.Kind().String() != "int" {
		t.Errorf("Kind for int: expected 'int', got '%s'", v.Kind())
	}
	if v.CanAddr() {
		t.Error("CanAddr for non-addressable value: expected false")
	}

	i := 123
	v = tr.ValueOf(&i)
	elem, _ := v.Elem()
	if elem.Kind().String() != "int" {
		t.Errorf("Kind for addressable int: expected 'int', got '%s'", elem.Kind())
	}
	if !elem.CanAddr() {
		t.Error("CanAddr for addressable value: expected true")
	}
}

func TestAllGetters(t *testing.T) {
	tr := tinyreflect.New()
	testCases := []struct {
		name     string
		value    interface{}
		kind     string
		intVal   int64
		uintVal  uint64
		floatVal float64
		boolVal  bool
		strVal   string
	}{
		{"Int", int(1), "int", 1, 0, 0, false, ""},
		{"Uint", uint(6), "uint", 0, 6, 0, false, ""},
		{"Float32", float32(12.0), "float32", 0, 0, 12.0, false, ""},
		{"Bool", true, "bool", 0, 0, 0, true, ""},
		{"String", "hello", "string", 0, 0, 0, false, "hello"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := tr.ValueOf(tc.value)
			if v.Kind().String() != tc.kind {
				t.Fatalf("Kind mismatch: expected %s, got %s", tc.kind, v.Kind().String())
			}

			switch tc.kind {
			case "int":
				i, err := v.Int()
				if err != nil || i != tc.intVal {
					t.Errorf("Int() got %d, %v; want %d, nil", i, err, tc.intVal)
				}
			case "uint":
				u, err := v.Uint()
				if err != nil || u != tc.uintVal {
					t.Errorf("Uint() got %d, %v; want %d, nil", u, err, tc.uintVal)
				}
			case "float32":
				f, err := v.Float()
				if err != nil || f != tc.floatVal {
					t.Errorf("Float() got %f, %v; want %f, nil", f, err, tc.floatVal)
				}
			case "bool":
				b, err := v.Bool()
				if err != nil || b != tc.boolVal {
					t.Errorf("Bool() got %v, %v; want %v, nil", b, err, tc.boolVal)
				}
			case "string":
				s := v.String()
				if s != tc.strVal {
					t.Errorf("String() got %s; want %s", s, tc.strVal)
				}
			}
		})
	}
}

func TestNumFieldAndField(t *testing.T) {
	tr := tinyreflect.New()
	type S struct{ A int; B string }
	s := S{}
	v := tr.ValueOf(s)

	n, err := v.NumField()
	if err != nil || n != 2 {
		t.Errorf("NumField on struct: got %d, %v; want 2, nil", n, err)
	}

	f, err := v.Field(0)
	if err != nil || f.Kind().String() != "int" {
		t.Errorf("Field(0): got kind %s, %v; want 'int', nil", f.Kind(), err)
	}
}

func TestFieldUnexported(t *testing.T) {
	tr := tinyreflect.New()
	type E struct{ e int }
	type S struct {
		E
		s int
	}
	s := S{}
	v := tr.ValueOf(s)

	f, err := v.Field(0)
	if err != nil || f.Kind().String() != "struct" {
		t.Errorf("Field(0) on unexported embedded: got %s, %v; want 'struct', nil", f.Kind(), err)
	}

	f, err = v.Field(1)
	if err != nil || f.Kind().String() != "int" {
		t.Errorf("Field(1) on unexported: got %s, %v; want 'int', nil", f.Kind(), err)
	}
}