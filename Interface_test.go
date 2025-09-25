package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestInterface(t *testing.T) {

	// Test with a zero Value (should return an error)
	v := tinyreflect.Value{}
	_, err := v.Interface()
	if err == nil {
		t.Error("Interface with zero Value: expected an error, but got nil")
	}

	// Test with a valid value
	i := 123
	v = tinyreflect.ValueOf(i)
	iface, err := v.Interface()
	if err != nil {
		t.Errorf("Interface on valid value: unexpected error: %v", err)
	}
	if val, ok := iface.(int); !ok || val != 123 {
		t.Errorf("Interface on valid value: expected 123, got %v", iface)
	}

	// Test with a nil interface value
	var nilIface any = nil
	v = tinyreflect.ValueOf(nilIface)
	iface, err = v.Interface()
	if err == nil {
		t.Error("Interface on nil interface value: expected an error")
	}
	if iface != nil {
		t.Errorf("Interface on nil interface value: expected nil, got %v", iface)
	}
}