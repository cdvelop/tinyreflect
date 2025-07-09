package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestSliceInt(t *testing.T) {
	// Test []int
	numbers := []int{10, 20, 30}
	v := refValueOf(numbers)

	if v.err != nil {
		t.Errorf("Error: %v", v.err)
		return
	}

	if v.refKind() != KSlice {
		t.Errorf("Expected KSlice, got %v", v.refKind())
	}

	if v.refLen() != 3 {
		t.Errorf("Expected length 3, got %d", v.refLen())
	}

	// Test element access
	elem := v.refIndex(1)
	if elem.err != nil {
		t.Errorf("Index error: %v", elem.err)
		return
	}

	if elem.refInt() != 20 {
		t.Errorf("Expected 20, got %d", elem.refInt())
	}
}

func TestSliceByte(t *testing.T) {
	// Test []byte
	data := []byte{65, 66, 67} // "ABC"
	v := refValueOf(data)

	if v.err != nil {
		t.Errorf("Error: %v", v.err)
		return
	}

	if v.refKind() != KByte {
		t.Errorf("Expected KByte, got %v", v.refKind())
	}

	if v.refLen() != 3 {
		t.Errorf("Expected length 3, got %d", v.refLen())
	}

	// Test element access
	elem := v.refIndex(0)
	if elem.err != nil {
		t.Errorf("Index error: %v", elem.err)
		return
	}

	if elem.refUint() != 65 {
		t.Errorf("Expected 65, got %d", elem.refUint())
	}
}

func TestSliceBool(t *testing.T) {
	// Test []bool
	flags := []bool{true, false, true}
	v := refValueOf(flags)

	if v.err != nil {
		t.Errorf("Error: %v", v.err)
		return
	}

	if v.refKind() != KSlice {
		t.Errorf("Expected KSlice, got %v", v.refKind())
	}

	if v.refLen() != 3 {
		t.Errorf("Expected length 3, got %d", v.refLen())
	}

	// Test element access
	elem := v.refIndex(0)
	if elem.err != nil {
		t.Errorf("Index error: %v", elem.err)
		return
	}

	if !elem.refBool() {
		t.Errorf("Expected true, got %v", elem.refBool())
	}
}

func TestSliceString(t *testing.T) {
	// Test []string (should use KSliceStr)
	names := []string{"Alice", "Bob", "Charlie"}
	v := refValueOf(names)

	if v.err != nil {
		t.Errorf("Error: %v", v.err)
		return
	}

	if v.refKind() != KSliceStr {
		t.Errorf("Expected KSliceStr, got %v", v.refKind())
	}

	if v.refLen() != 3 {
		t.Errorf("Expected length 3, got %d", v.refLen())
	}

	// Test element access
	elem := v.refIndex(1)
	if elem.err != nil {
		t.Errorf("Index error: %v", elem.err)
		return
	}

	if elem.refString() != "Bob" {
		t.Errorf("Expected 'Bob', got '%s'", elem.refString())
	}
}
