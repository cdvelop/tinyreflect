package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestNewElemBehavior verifies that NewValue() creates proper pointers
// and that Elem() correctly retrieves typed elements from those pointers.
func TestNewElemBehavior(t *testing.T) {
	tr := tinyreflect.New()
	type InnerStruct struct {
		V int
		S string
	}

	t.Run("NewCreatesNonNilPointer", func(t *testing.T) {
		sample := InnerStruct{}
		sampleType := tr.TypeOf(sample)

		if sampleType == nil || sampleType.Kind().String() != "struct" {
			t.Fatalf("TypeOf failed for struct, got kind %v", sampleType.Kind())
		}

		newPtr := tr.NewValue(sampleType)
		if newPtr.Type() == nil || newPtr.Kind().String() != "ptr" {
			t.Fatalf("NewValue returned incorrect type, got kind %v", newPtr.Kind())
		}

		isNil, err := newPtr.IsNil()
		if err != nil {
			t.Fatalf("IsNil() failed: %v", err)
		}
		if isNil {
			t.Error("NewValue() should create non-nil pointer, but IsNil() returned true")
		}
	})

	t.Run("ElemReturnsTypedElement", func(t *testing.T) {
		sampleType := tr.TypeOf(InnerStruct{})
		newPtr := tr.NewValue(sampleType)

		elem, err := newPtr.Elem()
		if err != nil {
			t.Fatalf("Elem() failed: %v", err)
		}
		if elem.Type() == nil || elem.Kind().String() != "struct" {
			t.Fatalf("Elem() returned incorrect type, got kind %v", elem.Kind())
		}
	})

	t.Run("ElementFieldAccess", func(t *testing.T) {
		sampleType := tr.TypeOf(InnerStruct{})
		newPtr := tr.NewValue(sampleType)
		elem, _ := newPtr.Elem()

		_, err := elem.Field(0)
		if err != nil {
			t.Errorf("Field(0) failed: %v", err)
		}
		_, err = elem.Field(1)
		if err != nil {
			t.Errorf("Field(1) failed: %v", err)
		}
	})

	t.Run("NewWithBasicTypes", func(t *testing.T) {
		basicTypes := []interface{}{int(0), string(""), bool(false)}

		for _, sample := range basicTypes {
			sampleType := tr.TypeOf(sample)
			if sampleType == nil {
				t.Errorf("TypeOf returned nil for %T", sample)
				continue
			}

			newPtr := tr.NewValue(sampleType)
			if newPtr.Kind().String() != "ptr" {
				t.Errorf("Expected pointer kind for %T, got %v", sample, newPtr.Kind())
			}

			isNil, _ := newPtr.IsNil()
			if isNil {
				t.Errorf("NewValue() created nil pointer for %T", sample)
			}
		}
	})
}