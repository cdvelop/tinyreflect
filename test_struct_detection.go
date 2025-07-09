package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestStructTypeDetection(t *testing.T) {
	// Test struct type detection
	type TestStruct struct {
		A, B, C, D, E int64
	}

	s := TestStruct{1, 2, 3, 4, 5}

	// Test Convert detection
	kind := Convert(s).GetKind()
	t.Logf("Convert detects struct as: %s", kind)

	// Test refValueOf detection
	v := refValueOf(s)
	t.Logf("refValueOf detects struct as: %s", v.refKind())

	// Test if it's consistent
	if kind != v.refKind() {
		t.Errorf("Inconsistent detection: Convert=%s, refValueOf=%s", kind, v.refKind())
	}
}
