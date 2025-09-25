package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// TestStruct is a helper struct for testing reflection operations.
type TestStruct struct {
	Name       string
	Age        int
	IsActive   bool
	Data       []byte
	ID         uint
	Balance    float64
	Int8Val    int8
	Int16Val   int16
	Int32Val   int32
	Int64Val   int64
	Uint8Val   uint8
	Uint16Val  uint16
	Uint32Val  uint32
	Uint64Val  uint64
	UintptrVal uintptr
	Float32Val float32
}

func TestSetString(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{Name: "initial"}
	const newName = "changed"

	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)

	if err := nameField.SetString(newName); err != nil {
		t.Fatalf("SetString failed unexpectedly: %v", err)
	}

	if data.Name != newName {
		t.Errorf("SetString failed: expected Name to be %q, got %q", newName, data.Name)
	}
}

func TestSetBool(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{IsActive: false}
	const newIsActive = true

	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	isActiveField, _ := structVal.Field(2)

	if err := isActiveField.SetBool(newIsActive); err != nil {
		t.Fatalf("SetBool failed unexpectedly: %v", err)
	}

	if data.IsActive != newIsActive {
		t.Errorf("SetBool failed: expected IsActive to be %v, got %v", newIsActive, data.IsActive)
	}
}

func TestSetBytes(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{Data: []byte("initial")}
	newData := []byte("changed")

	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	dataField, _ := structVal.Field(3)

	if err := dataField.SetBytes(newData); err != nil {
		t.Fatalf("SetBytes failed unexpectedly: %v", err)
	}

	if string(data.Data) != string(newData) {
		t.Errorf("SetBytes failed: expected Data to be %q, got %q", string(newData), string(data.Data))
	}
}

func TestSetInt(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{}
	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()

	testCases := []struct {
		fieldIndex int
		value      int64
	}{
		{1, 45}, {6, 127}, {7, 32767}, {8, 2147483647}, {9, 9223372036854775807},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		if err := field.SetInt(tc.value); err != nil {
			t.Errorf("SetInt on field %d failed: %v", tc.fieldIndex, err)
		}
	}

	if data.Age != 45 || data.Int8Val != 127 || data.Int16Val != 32767 || data.Int32Val != 2147483647 || data.Int64Val != 9223372036854775807 {
		t.Errorf("SetInt did not update struct values correctly. Got: %+v", data)
	}
}

func TestSetUint(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{}
	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()

	testCases := []struct {
		fieldIndex int
		value      uint64
	}{
		{4, 20}, {10, 255}, {11, 65535}, {12, 4294967295}, {13, 18446744073709551615}, {14, 12345},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		if err := field.SetUint(tc.value); err != nil {
			t.Errorf("SetUint on field %d failed: %v", tc.fieldIndex, err)
		}
	}

	if data.ID != 20 || data.Uint8Val != 255 || data.Uint16Val != 65535 || data.Uint32Val != 4294967295 || data.Uint64Val != 18446744073709551615 || data.UintptrVal != 12345 {
		t.Errorf("SetUint did not update struct values correctly. Got: %+v", data)
	}
}

func TestSetFloat(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{}
	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()

	testCases := []struct {
		fieldIndex int
		value      float64
	}{
		{5, 678.90}, {15, 123.45},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		if err := field.SetFloat(tc.value); err != nil {
			t.Errorf("SetFloat on field %d failed: %v", tc.fieldIndex, err)
		}
	}

	if data.Balance != 678.90 || data.Float32Val != 123.45 {
		t.Errorf("SetFloat did not update struct values correctly. Got: %+v", data)
	}
}

func TestSetOnZeroValue(t *testing.T) {
	var v tinyreflect.Value
	if err := v.SetString("hello"); err == nil {
		t.Error("expected error when setting string on a zero value")
	}
	if err := v.SetBool(true); err == nil {
		t.Error("expected error when setting bool on a zero value")
	}
	if err := v.SetBytes([]byte("hello")); err == nil {
		t.Error("expected error when setting bytes on a zero value")
	}
	if err := v.SetInt(123); err == nil {
		t.Error("expected error when setting int on a zero value")
	}
	if err := v.SetUint(123); err == nil {
		t.Error("expected error when setting uint on a zero value")
	}
	if err := v.SetFloat(123.45); err == nil {
		t.Error("expected error when setting float on a zero value")
	}
}

func TestSetTypeMismatchErrors(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{Name: "hello", Age: 30}
	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0) // String field
	ageField, _ := structVal.Field(1)  // Int field

	testCases := []struct {
		name      string
		action    func() error
		expectErr bool
	}{
		{"SetString on Int", func() error { return ageField.SetString("fail") }, true},
		{"SetInt on String", func() error { return nameField.SetInt(123) }, true},
		{"SetUint on String", func() error { return nameField.SetUint(123) }, true},
		{"SetFloat on String", func() error { return nameField.SetFloat(123.45) }, true},
		{"SetBool on String", func() error { return nameField.SetBool(true) }, true},
		{"SetBytes on String", func() error { return nameField.SetBytes([]byte("fail")) }, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.action()
			if (err != nil) != tc.expectErr {
				t.Errorf("expected error: %v, got: %v", tc.expectErr, err)
			}
		})
	}
}