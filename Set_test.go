package tinyreflect

import (
	"testing"
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
	// The struct we want to modify
	data := TestStruct{Name: "initial", Age: 30}
	const newName = "changed"

	// Get a Value of the struct. Since we want to modify the original 'data'
	// variable, we need to work with a pointer to it.
	v := ValueOf(&data)

	// The ValueOf(&data) returns a Ptr. We need to get the element it points to
	// before we can access its fields. This is what Elem() does.
	structVal, err := v.Elem()
	if err != nil {
		t.Fatalf("Failed to get element from pointer value: %v", err)
	}

	// Find the field we want to modify.
	// For this test, we know "Name" is the first field (index 0).
	nameField, err := structVal.Field(0)
	if err != nil {
		t.Fatalf("Failed to get field 'Name': %v", err)
	}

	// Now, set the new string value on the 'Name' field.
	nameField.SetString(newName)

	// Finally, verify that the original struct's field has been updated.
	if data.Name != newName {
		t.Errorf("SetString failed: expected Name to be %q, but got %q", newName, data.Name)
	}
}

func TestSetBool(t *testing.T) {
	// The struct we want to modify
	data := TestStruct{IsActive: false}
	const newIsActive = true

	// Get a Value of the struct.
	v := ValueOf(&data)

	// Get the element it points to.
	structVal, err := v.Elem()
	if err != nil {
		t.Fatalf("Failed to get element from pointer value: %v", err)
	}

	// Find the field we want to modify.
	// For this test, we know "IsActive" is the third field (index 2).
	isActiveField, err := structVal.Field(2)
	if err != nil {
		t.Fatalf("Failed to get field 'IsActive': %v", err)
	}

	// Now, set the new boolean value on the 'IsActive' field.
	isActiveField.SetBool(newIsActive)

	// Finally, verify that the original struct's field has been updated.
	if data.IsActive != newIsActive {
		t.Errorf("SetBool failed: expected IsActive to be %v, but got %v", newIsActive, data.IsActive)
	}
}

func TestSetBytes(t *testing.T) {
	// The struct we want to modify
	data := TestStruct{Data: []byte("initial")}
	newData := []byte("changed")

	// Get a Value of the struct.
	v := ValueOf(&data)

	// Get the element it points to.
	structVal, err := v.Elem()
	if err != nil {
		t.Fatalf("Failed to get element from pointer value: %v", err)
	}

	// Find the field we want to modify.
	// For this test, we know "Data" is the fourth field (index 3).
	dataField, err := structVal.Field(3)
	if err != nil {
		t.Fatalf("Failed to get field 'Data': %v", err)
	}

	// Now, set the new byte slice value on the 'Data' field.
	dataField.SetBytes(newData)

	// Finally, verify that the original struct's field has been updated.
	if string(data.Data) != string(newData) {
		t.Errorf("SetBytes failed: expected Data to be %q, but got %q", string(newData), string(data.Data))
	}
}

func TestSetInt(t *testing.T) {
	data := TestStruct{}
	v := ValueOf(&data)
	structVal, _ := v.Elem()

	testCases := []struct {
		fieldIndex int
		value      int64
		expected   interface{}
	}{
		{1, 45, int(45)},
		{6, 127, int8(127)},
		{7, 32767, int16(32767)},
		{8, 2147483647, int32(2147483647)},
		{9, 9223372036854775807, int64(9223372036854775807)},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		field.SetInt(tc.value)
	}

	if data.Age != testCases[0].expected {
		t.Errorf("SetInt failed for Age: expected %d, got %d", testCases[0].expected, data.Age)
	}
	if data.Int8Val != testCases[1].expected {
		t.Errorf("SetInt failed for Int8Val: expected %d, got %d", testCases[1].expected, data.Int8Val)
	}
	if data.Int16Val != testCases[2].expected {
		t.Errorf("SetInt failed for Int16Val: expected %d, got %d", testCases[2].expected, data.Int16Val)
	}
	if data.Int32Val != testCases[3].expected {
		t.Errorf("SetInt failed for Int32Val: expected %d, got %d", testCases[3].expected, data.Int32Val)
	}
	if data.Int64Val != testCases[4].expected {
		t.Errorf("SetInt failed for Int64Val: expected %d, got %d", testCases[4].expected, data.Int64Val)
	}
}

func TestSetUint(t *testing.T) {
	data := TestStruct{}
	v := ValueOf(&data)
	structVal, _ := v.Elem()

	testCases := []struct {
		fieldIndex int
		value      uint64
		expected   interface{}
	}{
		{4, 20, uint(20)},
		{10, 255, uint8(255)},
		{11, 65535, uint16(65535)},
		{12, 4294967295, uint32(4294967295)},
		{13, 18446744073709551615, uint64(18446744073709551615)},
		{14, 12345, uintptr(12345)},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		field.SetUint(tc.value)
	}

	if data.ID != testCases[0].expected {
		t.Errorf("SetUint failed for ID: expected %d, got %d", testCases[0].expected, data.ID)
	}
	if data.Uint8Val != testCases[1].expected {
		t.Errorf("SetUint failed for Uint8Val: expected %d, got %d", testCases[1].expected, data.Uint8Val)
	}
	if data.Uint16Val != testCases[2].expected {
		t.Errorf("SetUint failed for Uint16Val: expected %d, got %d", testCases[2].expected, data.Uint16Val)
	}
	if data.Uint32Val != testCases[3].expected {
		t.Errorf("SetUint failed for Uint32Val: expected %d, got %d", testCases[3].expected, data.Uint32Val)
	}
	if data.Uint64Val != testCases[4].expected {
		t.Errorf("SetUint failed for Uint64Val: expected %d, got %d", testCases[4].expected, data.Uint64Val)
	}
	if data.UintptrVal != testCases[5].expected {
		t.Errorf("SetUint failed for UintptrVal: expected %d, got %d", testCases[5].expected, data.UintptrVal)
	}
}

func TestSetFloat(t *testing.T) {
	data := TestStruct{}
	v := ValueOf(&data)
	structVal, _ := v.Elem()

	testCases := []struct {
		fieldIndex int
		value      float64
		expected   interface{}
	}{
		{5, 678.90, float64(678.90)},
		{15, 123.45, float32(123.45)},
	}

	for _, tc := range testCases {
		field, _ := structVal.Field(tc.fieldIndex)
		field.SetFloat(tc.value)
	}

	if data.Balance != testCases[0].expected {
		t.Errorf("SetFloat failed for Balance: expected %f, got %f", testCases[0].expected, data.Balance)
	}
	if data.Float32Val != testCases[1].expected {
		t.Errorf("SetFloat failed for Float32Val: expected %f, got %f", testCases[1].expected, data.Float32Val)
	}
}

func TestSetOnZeroValue(t *testing.T) {
	var v Value
	if err := v.SetString("hello"); err == nil {
		t.Error("expected an error when setting string on a zero value, but got nil")
	}
	if err := v.SetBool(true); err == nil {
		t.Error("expected an error when setting bool on a zero value, but got nil")
	}
	if err := v.SetBytes([]byte("hello")); err == nil {
		t.Error("expected an error when setting bytes on a zero value, but got nil")
	}
	if err := v.SetInt(123); err == nil {
		t.Error("expected an error when setting int on a zero value, but got nil")
	}
	if err := v.SetUint(123); err == nil {
		t.Error("expected an error when setting uint on a zero value, but got nil")
	}
	if err := v.SetFloat(123.45); err == nil {
		t.Error("expected an error when setting float on a zero value, but got nil")
	}
}

func TestSetStringFailsOnInt(t *testing.T) {
	data := TestStruct{Age: 30}
	v := ValueOf(&data)
	structVal, _ := v.Elem()
	ageField, _ := structVal.Field(1)
	err := ageField.SetString("hello")
	if err == nil {
		t.Error("expected an error when setting a string on an int field, but got nil")
	}
}

func TestSetIntFailsOnString(t *testing.T) {
	data := TestStruct{Name: "hello"}
	v := ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)
	err := nameField.SetInt(123)
	if err == nil {
		t.Error("expected an error when setting an int on a string field, but got nil")
	}
}

func TestSetUintFailsOnString(t *testing.T) {
	data := TestStruct{Name: "hello"}
	v := ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)
	err := nameField.SetUint(123)
	if err == nil {
		t.Error("expected an error when setting a uint on a string field, but got nil")
	}
}

func TestSetFloatFailsOnString(t *testing.T) {
	data := TestStruct{Name: "hello"}
	v := ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)
	err := nameField.SetFloat(123.45)
	if err == nil {
		t.Error("expected an error when setting a float on a string field, but got nil")
	}
}

func TestSetBoolFailsOnString(t *testing.T) {
	data := TestStruct{Name: "hello"}
	v := ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)
	err := nameField.SetBool(true)
	if err == nil {
		t.Error("expected an error when setting a bool on a string field, but got nil")
	}
}

func TestSetBytesFailsOnString(t *testing.T) {
	data := TestStruct{Name: "hello"}
	v := ValueOf(&data)
	structVal, _ := v.Elem()
	nameField, _ := structVal.Field(0)
	err := nameField.SetBytes([]byte("hello"))
	if err == nil {
		t.Error("expected an error when setting a bytes slice on a string field, but got nil")
	}
}
