package tinyreflect

import (
	"testing"
)

func TestTypeMapKey(t *testing.T) {
	// Test how Type pointers behave as map keys
	type TestStruct struct {
		Name string
		Age  int
	}

	// Get the same type twice
	s1 := TestStruct{}
	s2 := TestStruct{}

	v1 := ValueOf(s1)
	v2 := ValueOf(s2)

	typ1 := v1.Type()
	typ2 := v2.Type()

	t.Logf("Type1: %p", typ1)
	t.Logf("Type2: %p", typ2)

	// Test if they're the same pointer
	if typ1 == typ2 {
		t.Log("✅ Same type instances have same pointer - map key will work")
	} else {
		t.Error("❌ Same type instances have different pointers - map key will fail!")
	}

	// Test with map
	testMap := make(map[*Type]string)
	testMap[typ1] = "first"

	// Try to retrieve with typ2
	value, ok := testMap[typ2]
	if ok {
		t.Logf("✅ Map lookup succeeded: %s", value)
	} else {
		t.Error("❌ Map lookup failed - different pointer addresses")
	}

	// Test with different structs
	type DifferentStruct struct {
		Name string
		Age  int
	}

	d1 := DifferentStruct{}
	vd1 := ValueOf(d1)
	typd1 := vd1.Type()

	_, ok = testMap[typd1]
	if !ok {
		t.Log("✅ Different struct type correctly not found in map")
	} else {
		t.Error("❌ Different struct type incorrectly found in map")
	}
}
