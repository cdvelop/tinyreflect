package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

// TestNewElemBehavior verifies that New() creates proper pointers
// and that Elem() correctly retrieves typed elements from those pointers.
// This test ensures the core pointer creation and dereferencing works correctly.
func TestNewElemBehavior(t *testing.T) {
	type InnerStruct struct {
		V int
		S string
	}

	// Test 1: Verify New() creates a non-nil pointer to zero value
	t.Run("NewCreatesNonNilPointer", func(t *testing.T) {
		sample := InnerStruct{}
		sampleType := TypeOf(sample)

		// Verify sample type is valid
		if sampleType == nil {
			t.Fatal("TypeOf returned nil")
		}
		if sampleType.Kind() != K.Struct {
			t.Errorf("Expected struct kind, got %v", sampleType.Kind())
		}

		// Create new pointer using New()
		newPtr := New(sampleType)
		if newPtr.Type() == nil {
			t.Fatal("New() returned value with nil type")
		}
		if newPtr.Kind() != K.Pointer {
			t.Errorf("Expected pointer kind, got %v", newPtr.Kind())
		}

		// Verify pointer is not nil (critical fix)
		isNil, err := newPtr.IsNil()
		if err != nil {
			t.Fatalf("IsNil() failed: %v", err)
		}
		if isNil {
			t.Error("New() should create non-nil pointer, but IsNil() returned true")
		}
	})

	// Test 2: Verify Elem() returns properly typed element
	t.Run("ElemReturnsTypedElement", func(t *testing.T) {
		sample := InnerStruct{}
		sampleType := TypeOf(sample)
		newPtr := New(sampleType)

		// Get element from pointer
		elem, err := newPtr.Elem()
		if err != nil {
			t.Fatalf("Elem() failed: %v", err)
		}

		// Verify element has proper type (this was the bug)
		if elem.Type() == nil {
			t.Fatal("Elem() returned value with nil type - this indicates a bug in New() or Elem()")
		}
		if elem.Type().Kind() != K.Struct {
			t.Errorf("Expected struct element, got kind %v", elem.Type().Kind())
		}

		// Verify element is zero value but typed
		if elem.Kind() != K.Struct {
			t.Errorf("Element should have struct kind, got %v", elem.Kind())
		}
	})

	// Test 3: Verify field access works on created elements
	t.Run("ElementFieldAccess", func(t *testing.T) {
		sample := InnerStruct{}
		sampleType := TypeOf(sample)
		newPtr := New(sampleType)
		elem, err := newPtr.Elem()
		if err != nil {
			t.Fatalf("Elem() failed: %v", err)
		}

		// Try to access fields - this should work if type is proper
		field0, err := elem.Field(0)
		if err != nil {
			t.Fatalf("Field(0) failed: %v", err)
		}
		if field0.Type() == nil {
			t.Error("Field(0) has nil type")
		}

		field1, err := elem.Field(1)
		if err != nil {
			t.Fatalf("Field(1) failed: %v", err)
		}
		if field1.Type() == nil {
			t.Error("Field(1) has nil type")
		}
	})

	// Test 4: Verify New() works with basic types
	t.Run("NewWithBasicTypes", func(t *testing.T) {
		basicTypes := []interface{}{
			int(0),
			string(""),
			bool(false),
		}

		for _, sample := range basicTypes {
			sampleType := TypeOf(sample)
			if sampleType == nil {
				t.Errorf("TypeOf returned nil for %Translate", sample)
				continue
			}

			newPtr := New(sampleType)
			if newPtr.Type() == nil {
				t.Errorf("New() returned value with nil type for %Translate", sample)
				continue
			}

			// Verify it's a pointer
			if newPtr.Kind() != K.Pointer {
				t.Errorf("Expected pointer kind for %Translate, got %v", sample, newPtr.Kind())
			}

			// Verify not nil
			isNil, err := newPtr.IsNil()
			if err != nil {
				t.Errorf("IsNil() failed for %Translate: %v", sample, err)
				continue
			}
			if isNil {
				t.Errorf("New() created nil pointer for %Translate", sample)
			}
		}
	})
}
