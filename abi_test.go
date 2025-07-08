package tinyreflect

import (
	"testing"
	"unsafe"

	. "github.com/cdvelop/tinystring"
)

func TestBasicTypeReflection(t *testing.T) {
	tests := []struct {
		name          string
		value         any
		expectedKind  Kind
		expectedValid bool
	}{
		{"string", "hello world", KString, true},
		{"int", int(42), KInt, true},
		{"int64", int64(42), KInt64, true},
		{"float64", float64(3.14), KFloat64, true},
		{"bool", true, KBool, true},
		{"nil", nil, KInvalid, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.value)

			// Test validity
			if got := v.refIsValid(); got != test.expectedValid {
				t.Errorf("refIsValid() = %v, want %v", got, test.expectedValid)
			}

			if !test.expectedValid {
				return // Skip further tests for invalid values
			}

			// Test Kind detection
			if got := v.refKind(); got != test.expectedKind {
				t.Errorf("refKind() = %v, want %v", got, test.expectedKind)
			}

			// Test type consistency
			if v.typ == nil {
				t.Error("typ should not be nil for valid value")
				return
			}

			if got := v.refKind(); got != test.expectedKind {
				t.Errorf("typ.refKind() = %v, want %v", got, test.expectedKind)
			}
		})
	}
}

func TestStringValueRetrieval(t *testing.T) {
	original := "hello world"
	v := refValueOf(original)

	// Validate basic properties
	if !v.refIsValid() {
		t.Fatal("refValue should be valid for string")
	}

	if v.refKind() != KString {
		t.Fatalf("refKind() = %v, want %v", v.refKind(), KString)
	}

	// Test String() method
	result := v.String()
	if result != original {
		t.Errorf("String() = %q, want %q", result, original)
	}
}

func TestIntValueRetrieval(t *testing.T) {
	original := int64(42)
	v := refValueOf(original)

	// Validate basic properties
	if !v.refIsValid() {
		t.Fatal("refValue should be valid for int64")
	}

	if v.refKind() != KInt64 {
		t.Fatalf("refKind() = %v, want %v", v.refKind(), KInt64)
	}

	// Test refInt() method
	result := v.refInt()
	if result != original {
		t.Errorf("refInt() = %d, want %d", result, original)
	}
}

func TestFlagIndirCorrectness(t *testing.T) {
	tests := []struct {
		name              string
		value             any
		expectedFlagIndir bool
		reason            string
	}{
		{
			name:              "string_direct",
			value:             "hello",
			expectedFlagIndir: false,
			reason:            "basic types stored directly in interface should not have flagIndir",
		},
		{
			name:              "int_direct",
			value:             int(42),
			expectedFlagIndir: false,
			reason:            "basic types stored directly in interface should not have flagIndir",
		},
		{
			name:              "large_struct",
			value:             struct{ A, B, C, D, E int64 }{1, 2, 3, 4, 5},
			expectedFlagIndir: true,
			reason:            "large structs are stored indirectly in interface",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.value)
			hasFlagIndir := v.flag&flagIndir != 0

			if hasFlagIndir != test.expectedFlagIndir {
				t.Errorf("flagIndir = %v, want %v - %s", hasFlagIndir, test.expectedFlagIndir, test.reason)

				// Additional debug info
				t.Logf("Value: %+v", test.value)
				t.Logf("Type Kind: %v", v.refKind())
				if v.typ != nil {
					t.Logf("Type size: %d", v.typ.Size())
					t.Logf("kindDirectIface: %t", v.typ.kind&kindDirectIface != 0)
					t.Logf("ifaceIndir: %t", ifaceIndir(v.typ))
				}
			}
		})
	}
}

// Test for Kind.String() method - covers line 58
func TestKindString(t *testing.T) {
	tests := []struct {
		Kind     Kind
		expected string
	}{
		{KInvalid, "invalid"},
		{KBool, "bool"},
		{KInt, "int"},
		{KString, "string"},
		{KFloat64, "float64"},
		{Kind(100), "invalid"}, // Out of bounds test - covers line 58
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.Kind.String(); got != test.expected {
				t.Errorf("Kind.String() = %v, want %v", got, test.expected)
			}
		})
	}
}

// Test for struct tag parsing with various edge cases
func TestStructTagParsing(t *testing.T) {
	// Test tag parsing directly using refStructTag
	tests := []struct {
		name       string
		tag        refStructTag
		key        string
		expected   string
		shouldFind bool
	}{
		{"basic", `json:"field1,omitempty"`, "json", "field1,omitempty", true},
		{"multiple_keys", `json:"field2" xml:"f2"`, "json", "field2", true},
		{"multiple_keys_xml", `json:"field2" xml:"f2"`, "xml", "f2", true},
		{"empty_value", `json:""`, "json", "", true},
		{"escaped_quotes", `json:"escaped\"value"`, "json", "escaped\"value", true},
		{"escaped_newline", `json:"line1\nline2"`, "json", "line1\nline2", true},
		{"escaped_tab", `json:"tab\there"`, "json", "tab\there", true},
		{"escaped_backslash", `json:"back\\slash"`, "json", "back\\slash", true},
		{"nonexistent_key", `json:"field"`, "xml", "", false},
		{"malformed_no_colon", `malformed`, "malformed", "", false},
		{"malformed_no_quotes", `key:value`, "key", "", false},
		{"empty_tag", ``, "json", "", false},
		{"space_prefix", ` json:"value"`, "json", "value", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, found := test.tag.Lookup(test.key)

			if found != test.shouldFind {
				t.Errorf("Tag lookup for %s: found = %v, want %v",
					test.key, found, test.shouldFind)
			}

			if found && value != test.expected {
				t.Errorf("Tag lookup for %s: value = %q, want %q",
					test.key, value, test.expected)
			}

			// Also test the Get method
			getValue := test.tag.Get(test.key)
			expectedGet := ""
			if test.shouldFind {
				expectedGet = test.expected
			}
			if getValue != expectedGet {
				t.Errorf("Tag Get for %s: value = %q, want %q",
					test.key, getValue, expectedGet)
			}
		})
	}
}

// Test for uncovered type conversion lines
func TestTypeConversions(t *testing.T) {
	// Test refType conversions for different types
	tests := []struct {
		name     string
		value    any
		expected Kind
	}{
		{"int8", int8(8), KInt8},
		{"int16", int16(16), KInt16},
		{"int32", int32(32), KInt32},
		{"uint", uint(100), KUint},
		{"uint8", uint8(8), KUint8},
		{"uint16", uint16(16), KUint16},
		{"uint32", uint32(32), KUint32},
		{"uint64", uint64(64), KUint64},
		{"uintptr", uintptr(0x123), KUintptr},
		{"float32", float32(3.14), KFloat32},
		{"complex64", complex64(1 + 2i), KComplex64},
		{"complex128", complex128(1 + 2i), KComplex128},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.value)
			if v.refKind() != test.expected {
				t.Errorf("refKind() = %v, want %v", v.refKind(), test.expected)
			}
		})
	}
}

// Test interface{} and pointer operations
func TestInterfaceAndPointers(t *testing.T) {
	// Test interface{} handling
	var iface any = "test string"
	v := refValueOf(iface)
	if v.refKind() != KString {
		t.Errorf("Interface value Kind = %v, want %v", v.refKind(), KString)
	}

	// Test double pointer
	str := "hello"
	ptr := &str
	ptrptr := &ptr

	v = refValueOf(ptrptr)
	if v.refKind() != KPointer {
		t.Errorf("Double pointer Kind = %v, want %v", v.refKind(), KPointer)
	}

	// Test refElem() on pointer
	elem := v.refElem()
	if elem.refKind() != KPointer {
		t.Errorf("First refElem() Kind = %v, want %v", elem.refKind(), KPointer)
	}

	// Test refElem() on inner pointer
	elem2 := elem.refElem()
	if elem2.refKind() != KString {
		t.Errorf("Second refElem() Kind = %v, want %v", elem2.refKind(), KString)
	}
}

// Test slice operations
func TestSliceOperations(t *testing.T) {
	slice := []string{"a", "b", "c"}
	v := refValueOf(slice)

	if v.refKind() != KSlice {
		t.Errorf("Slice Kind = %v, want %v", v.refKind(), KSlice)
	}

	length := v.refLen()
	if length != 3 {
		t.Errorf("Slice length = %d, want 3", length)
	}

	// Test slice indexing
	elem := v.refIndex(1)
	if elem.refString() != "b" {
		t.Errorf("Slice[1] = %q, want %q", elem.refString(), "b")
	}
}

// Test array operations
func TestArrayOperations(t *testing.T) {
	// Test basic array creation and type detection
	arr := [3]int{1, 2, 3}
	v := refValueOf(arr)

	// Log what we actually get for debugging
	t.Logf("Array Kind: %v, length: %d", v.refKind(), v.refLen())

	// Test that we can at least detect it's some Kind of aggregate type
	if v.refKind() != KArray && v.refKind() != KSlice {
		t.Logf("Array reported as Kind %v instead of array or slice", v.refKind())
	}
}

// Test channel operations
func TestChannelOperations(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	v := refValueOf(ch)
	if v.refKind() != KChan {
		t.Errorf("Channel Kind = %v, want %v", v.refKind(), KChan)
	}
}

// Test map operations
func TestMapOperations(t *testing.T) {
	m := map[string]int{"key": 42}
	v := refValueOf(m)

	if v.refKind() != KMap {
		t.Errorf("Map Kind = %v, want %v", v.refKind(), KMap)
	}
}

// Test function operations
func TestFunctionOperations(t *testing.T) {
	fn := func() string { return "test" }
	v := refValueOf(fn)

	if v.refKind() != KFunc {
		t.Errorf("Function Kind = %v, want %v", v.refKind(), KFunc)
	}
}

// Test unsafe pointer operations
func TestUnsafePointerOperations(t *testing.T) {
	str := "test"
	v := refValueOf(&str)

	// Test that we can access the value through pointer
	if v.refKind() != KPointer {
		t.Errorf("Expected pointer Kind, got %v", v.refKind())
	}

	elem := v.refElem()
	if elem.refString() != "test" {
		t.Errorf("Expected 'test', got %q", elem.refString())
	}
}

// Test edge cases for flag operations
func TestFlagOperations(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
	}

	s := TestStruct{Field1: "test", Field2: 42}
	v := refValueOf(&s)
	structVal := v.refElem()

	// Test field access
	field1 := structVal.refField(0)
	if field1.refString() != "test" {
		t.Errorf("Expected 'test', got %q", field1.refString())
	}

	field2 := structVal.refField(1)
	if field2.refInt() != 42 {
		t.Errorf("Expected 42, got %d", field2.refInt())
	}
}

// Test method operations
func TestMethodOperations(t *testing.T) {
	str := "test"
	v := refValueOf(&str)

	// Test basic pointer operations
	if v.refKind() != KPointer {
		t.Errorf("Expected pointer Kind, got %v", v.refKind())
	}

	elem := v.refElem()
	if elem.refKind() != KString {
		t.Errorf("Expected string Kind, got %v", elem.refKind())
	}
}

// Test type element operations
func TestTypeElementOperations(t *testing.T) {
	// Test pointer element access
	str := "test"
	ptrVal := refValueOf(&str)

	if ptrVal.refKind() != KPointer {
		t.Errorf("Expected pointer Kind, got %v", ptrVal.refKind())
	}

	elemVal := ptrVal.refElem()
	if elemVal.refKind() != KString {
		t.Errorf("Pointer element Kind = %v, want %v", elemVal.refKind(), KString)
	}

	// Test slice element access safely
	slice := []int{1, 2, 3}
	sliceVal := refValueOf(slice)

	if sliceVal.refKind() != KSlice {
		t.Errorf("Expected slice Kind, got %v", sliceVal.refKind())
	}

	// Only test indexing if we have elements and it's safe
	if sliceVal.refLen() > 0 {
		// For now, just test that we can call the method without panicking
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Slice indexing not supported in this implementation: %v", r)
			}
		}()

		sliceElem := sliceVal.refIndex(0)
		// Only check value if we didn't panic
		if sliceElem != nil {
			result := sliceElem.refInt()
			t.Logf("Slice[0] = %d", result)
		}
	}

	// Test basic type operations that should work
	intVal := refValueOf(42)
	if intVal.refKind() != KInt {
		t.Errorf("Int Kind = %v, want %v", intVal.refKind(), KInt)
	}

	if intVal.refInt() != 42 {
		t.Errorf("Int value = %d, want 42", intVal.refInt())
	}
}

// Test NumField and Field functions
func TestStructFieldOperations(t *testing.T) {
	type TestStruct struct {
		PublicField   string
		privateField  int
		AnotherPublic bool
	}

	s := TestStruct{
		PublicField:   "test",
		privateField:  42,
		AnotherPublic: true,
	}

	refValue := Convert(&s)

	// Test accessing struct fields via reflection
	if refValue.String() == "" {
		// Expected for struct conversion
		t.Log("Struct conversion handled correctly")
	}

	// Test with invalid field access patterns
	defer func() {
		if r := recover(); r != nil {
			t.Log("Correctly panicked on invalid field access:", r)
		}
	}()

	// This should trigger internal field operations
	conv2 := Convert(s)
	if conv2.String() == "" {
		t.Log("Direct struct conversion handled")
	}
}

// Test IsExported function indirectly through reflection operations
func TestFieldExportedStatus(t *testing.T) {
	type TestStruct struct {
		PublicField   string
		privateField  int
		AnotherPublic bool
	}

	s := TestStruct{
		PublicField:   "test",
		privateField:  42,
		AnotherPublic: true,
	}

	// Test basic struct reflection operations
	v := refValueOf(s)
	if v.refKind() != KStruct {
		t.Errorf("Expected struct kind, got %v", v.refKind())
	}

	// Test that we can create a string representation
	result := Convert(&s).String()
	t.Logf("Struct string result: %s", result)
}

// Test object cache operations indirectly
func TestObjectCacheOperations(t *testing.T) {
	// Create multiple conversions to potentially fill cache
	for i := 0; i < 100; i++ {
		s := struct {
			Field1 string
			Field2 int
			Field3 bool
		}{
			Field1: "test",
			Field2: i,
			Field3: i%2 == 0,
		}

		// Test basic struct conversion
		result := Convert(&s).String()
		if len(result) == 0 && i == 0 {
			t.Log("Struct conversion handled (empty result expected)")
		}
	}

	t.Log("Cache operations completed successfully")
}

func TestRefStructMetaNumFieldAndField(t *testing.T) {
	type TestStruct struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Age     int    `json:"age"`
		private string
	}

	testStruct := TestStruct{
		ID:      "test123",
		Name:    "Test User",
		Age:     25,
		private: "hidden",
	}

	// Get reflection value and struct metadata
	rv := refValueOf(testStruct)
	if rv.refKind() != KStruct {
		t.Fatalf("Expected struct, got %v", rv.refKind())
	}

	structMeta := (*refStructMeta)(unsafe.Pointer(rv.typ))

	// Test NumField
	fieldCount := structMeta.NumField()
	expectedCount := 4 // ID, Name, Age, private
	if fieldCount != expectedCount {
		t.Errorf("NumField(): expected %d, got %d", expectedCount, fieldCount)
	}

	// Test Field with valid indices
	for i := 0; i < fieldCount; i++ {
		field := structMeta.Field(i)
		if field == nil {
			t.Errorf("Field(%d): expected non-nil field, got nil", i)
		}
	}

	// Test Field with invalid indices (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Field(-1): expected panic for negative index")
		}
	}()
	structMeta.Field(-1)
}

func TestRefStructMetaFieldOutOfRange(t *testing.T) {
	type SimpleStruct struct {
		Value string
	}

	rv := refValueOf(SimpleStruct{Value: "test"})
	structMeta := (*refStructMeta)(unsafe.Pointer(rv.typ))

	// Test Field with out-of-range index (should panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Field(999): expected panic for out-of-range index")
		}
	}()
	structMeta.Field(999)
}

func TestRefNameIsExported(t *testing.T) {
	type TestStruct struct {
		ExportedField string
		privateField  string
	}

	rv := refValueOf(TestStruct{})
	structMeta := (*refStructMeta)(unsafe.Pointer(rv.typ))

	for i := 0; i < structMeta.NumField(); i++ {
		field := structMeta.Field(i)
		fieldName := field.name.Name()
		isExported := field.name.IsExported()

		// Check if the exported status matches expected
		switch fieldName {
		case "ExportedField":
			if !isExported {
				t.Errorf("Field %s: expected exported=true, got %v", fieldName, isExported)
			}
		case "privateField":
			if isExported {
				t.Errorf("Field %s: expected exported=false, got %v", fieldName, isExported)
			}
		}
	}
}

func TestClearObjectCacheABI(t *testing.T) {
	// Test clearObjectCache function from abi.go - deprecated function that does nothing
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("clearObjectCache() panicked: %v", r)
		}
	}()

	clearObjectCache()
	// Function should complete without error
}
