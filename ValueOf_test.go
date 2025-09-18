package tinyreflect

import (
	"testing"
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
