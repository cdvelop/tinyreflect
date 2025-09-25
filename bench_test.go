package tinyreflect_test

import (
	"reflect"
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// BenchmarkStruct is a struct for benchmarking
type BenchmarkStruct struct {
	Name   string
	Age    int
	Active bool
	Data   []byte
	ID     uint64
	Score  float64
}

func (BenchmarkStruct) StructName() string {
	return "BenchmarkStruct"
}

// Benchmark standard library reflect TypeOf
func BenchmarkStdReflect_TypeOf(b *testing.B) {
	s := BenchmarkStruct{Name: "test", Age: 25}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reflect.TypeOf(s)
	}
}

// Benchmark tinyreflect TypeOf (first call - no cache)
func BenchmarkTinyReflect_TypeOf_First(b *testing.B) {
	s := BenchmarkStruct{Name: "test", Age: 25}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new instance each time to avoid cache
		tr2 := tinyreflect.New()
		_ = tr2.TypeOf(s)
	}
}

// Benchmark tinyreflect TypeOf (cached)
func BenchmarkTinyReflect_TypeOf_Cached(b *testing.B) {
	tr := tinyreflect.New()
	s := BenchmarkStruct{Name: "test", Age: 25}
	// Warm up cache
	_ = tr.TypeOf(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.TypeOf(s)
	}
}

// Benchmark standard library reflect field access
func BenchmarkStdReflect_FieldAccess(b *testing.B) {
	s := BenchmarkStruct{Name: "test", Age: 25, Active: true, ID: 123}
	v := reflect.ValueOf(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Field(0).String()
		_ = v.Field(1).Int()
		_ = v.Field(2).Bool()
		_ = v.Field(4).Uint()
	}
}

// Benchmark tinyreflect field access (first call - no cache)
func BenchmarkTinyReflect_FieldAccess_First(b *testing.B) {
	s := BenchmarkStruct{Name: "test", Age: 25, Active: true, ID: 123}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new instance each time to avoid cache
		tr2 := tinyreflect.New()
		v := tr2.ValueOf(s)
		_, _ = v.Field(0)
		_, _ = v.Field(1)
		_, _ = v.Field(2)
		_, _ = v.Field(4)
	}
}

// Benchmark tinyreflect field access (cached)
func BenchmarkTinyReflect_FieldAccess_Cached(b *testing.B) {
	tr := tinyreflect.New()
	s := BenchmarkStruct{Name: "test", Age: 25, Active: true, ID: 123}
	// Warm up cache
	v := tr.ValueOf(s)
	_, _ = v.Field(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := tr.ValueOf(s)
		_, _ = v.Field(0)
		_, _ = v.Field(1)
		_, _ = v.Field(2)
		_, _ = v.Field(4)
	}
}

// Benchmark standard library reflect field iteration
func BenchmarkStdReflect_FieldIteration(b *testing.B) {
	s := BenchmarkStruct{Name: "test", Age: 25, Active: true, Data: []byte("data"), ID: 123, Score: 95.5}
	v := reflect.ValueOf(s)
	typ := v.Type()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < v.NumField(); j++ {
			_ = typ.Field(j).Name
			_ = v.Field(j).Interface()
		}
	}
}

// Benchmark tinyreflect field iteration (first call - no cache)
func BenchmarkTinyReflect_FieldIteration_First(b *testing.B) {
	s := BenchmarkStruct{Name: "test", Age: 25, Active: true, Data: []byte("data"), ID: 123, Score: 95.5}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create new instance each time to avoid cache
		tr := tinyreflect.New()
		v := tr.ValueOf(s)
		typ := v.Type()
		num, _ := typ.NumField()
		for j := 0; j < num; j++ {
			_, _ = typ.NameByIndex(j)
			_, _ = v.Field(j)
		}
	}
}

// Benchmark tinyreflect field iteration (cached)
func BenchmarkTinyReflect_FieldIteration_Cached(b *testing.B) {
	tr := tinyreflect.New()
	s := BenchmarkStruct{Name: "test", Age: 25, Active: true, Data: []byte("data"), ID: 123, Score: 95.5}
	// Warm up cache
	v := tr.ValueOf(s)
	typ := v.Type()
	_, _ = typ.NumField()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := tr.ValueOf(s)
		typ := v.Type()
		num, _ := typ.NumField()
		for j := 0; j < num; j++ {
			_, _ = typ.NameByIndex(j)
			_, _ = v.Field(j)
		}
	}
}
