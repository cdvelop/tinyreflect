package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestInterface(t *testing.T) {
	// Test with a nil type
	v := Value{}
	_, err := v.Interface()
	if err == nil {
		t.Error("Interface with nil type: expected an error, but got nil")
	}

	// Test with an interface
	var i interface{} = 123
	v = ValueOf(i)
	v.flag = flag(K.Interface) // force kind to Interface
	_, err = v.Interface()
	if err == nil {
		t.Error("Interface on interface: expected an error, but got nil")
	}
}
