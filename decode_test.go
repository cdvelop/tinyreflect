package tinyreflect

import (
	"testing"
)

func TestDecodeScenario(t *testing.T) {
	// This replicates the exact scenario from decoder.go line 46
	type simpleStruct struct {
		Name      string
		Timestamp int64
		Payload   []byte
		Ssid      []uint32
	}

	// Test scenario like in decoder: pointer to struct
	s := &simpleStruct{}

	// This is exactly what decoder.go does
	rv := Indirect(ValueOf(s))

	// Check that rv has a valid type
	typ := rv.Type()
	if typ == nil {
		t.Error("rv.Type() returned nil - this is the 'value type nil' error")
	} else {
		t.Logf("rv.Type() returned %p, Kind: %v", typ, typ.Kind())
	}

	// Check CanAddr
	canAddr := rv.CanAddr()
	t.Logf("rv.CanAddr() = %v", canAddr)

	// Check if the original value has a type
	originalV := ValueOf(s)
	if originalV.Type() == nil {
		t.Error("ValueOf(s).Type() returned nil")
	} else {
		t.Logf("ValueOf(s).Type() returned %p, Kind: %v", originalV.Type(), originalV.Type().Kind())
	}

	// Check if Indirect is working properly
	if rv.typ_ == nil {
		t.Error("Indirect result has nil typ_")
	}

	// Compare with direct struct (not pointer)
	directStruct := simpleStruct{}
	directRv := ValueOf(directStruct)
	if directRv.Type() == nil {
		t.Error("Direct struct ValueOf returned nil type")
	} else {
		t.Logf("Direct struct Type() returned %p, Kind: %v", directRv.Type(), directRv.Type().Kind())
	}
}
