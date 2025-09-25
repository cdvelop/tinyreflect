package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestSetString(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{StringField: "initial"}
	const newName = "changed"

	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)

	if err := nameField.SetString(newName); err != nil {
		t.Fatalf("SetString failed unexpectedly: %v", err)
	}

	if data.StringField != newName {
		t.Errorf("SetString failed: expected StringField to be %q, got %q", newName, data.StringField)
	}
}

func TestSetBool(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{BoolField: false}
	const newIsActive = true

	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	isActiveField, _ := structVal.Field(1)

	if err := isActiveField.SetBool(newIsActive); err != nil {
		t.Fatalf("SetBool failed unexpectedly: %v", err)
	}

	if data.BoolField != newIsActive {
		t.Errorf("SetBool failed: expected BoolField to be %v, got %v", newIsActive, data.BoolField)
	}
}

func TestSetBytes(t *testing.T) {
	tr := tinyreflect.New()
	data := TestStruct{ByteSliceField: []byte("initial")}
	newData := []byte("changed")

	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	dataField, _ := structVal.Field(16)

	if err := dataField.SetBytes(newData); err != nil {
		t.Fatalf("SetBytes failed unexpectedly: %v", err)
	}

	if string(data.ByteSliceField) != string(newData) {
		t.Errorf("SetBytes failed: expected ByteSliceField to be %q, got %q", string(newData), string(data.ByteSliceField))
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
		{2, 45}, {3, 127}, {4, 32767}, {5, 2147483647}, {6, 9223372036854775807},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		if err := field.SetInt(tc.value); err != nil {
			t.Errorf("SetInt on field %d failed: %v", tc.fieldIndex, err)
		}
	}

	if data.IntField != 45 || data.Int8Field != 127 || data.Int16Field != 32767 || data.Int32Field != 2147483647 || data.Int64Field != 9223372036854775807 {
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
		{7, 20}, {8, 255}, {9, 65535}, {10, 4294967295}, {11, 18446744073709551615},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		if err := field.SetUint(tc.value); err != nil {
			t.Errorf("SetUint on field %d failed: %v", tc.fieldIndex, err)
		}
	}

	if data.UintField != 20 || data.Uint8Field != 255 || data.Uint16Field != 65535 || data.Uint32Field != 4294967295 || data.Uint64Field != 18446744073709551615 {
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
		{13, 678.90}, {12, 123.45},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		if err := field.SetFloat(tc.value); err != nil {
			t.Errorf("SetFloat on field %d failed: %v", tc.fieldIndex, err)
		}
	}

	if data.Float64Field != 678.90 || data.Float32Field != 123.45 {
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
	data := TestStruct{StringField: "hello", IntField: 30}
	v := tr.ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0) // String field
	ageField, _ := structVal.Field(2)  // Int field

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
