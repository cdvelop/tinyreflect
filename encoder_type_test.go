package tinyreflect_test

import (
	"github.com/cdvelop/tinyreflect"
	"testing"
)

func TestEncoderTypeMethod(t *testing.T) {
	t.Run("rv.Type() should not return nil", func(t *testing.T) {
		testCases := []interface{}{
			42,
			"hello",
			true,
			[]int{1, 2, 3},
			struct{ X int }{42},
			&struct{ X int }{42},
		}

		for _, tc := range testCases {
			// This replicates the exact code from encoder.go that was failing
			rv := tinyreflect.Indirect(tinyreflect.ValueOf(tc))
			typ := rv.Type()

			if typ == nil {
				t.Errorf("rv.Type() returned nil for value %v", tc)
			} else {
				t.Logf("rv.Type() returned %p for value %v", typ, tc)
			}
		}
	})

	t.Run("Test specific encoder case", func(t *testing.T) {
		// Test the specific case that was failing in the encoder
		x := struct {
			Name string
			Age  int
		}{"John", 25}

		// This is exactly what the encoder does
		rv := tinyreflect.Indirect(tinyreflect.ValueOf(x))
		typ := rv.Type()

		if typ == nil {
			t.Error("Encoder case: rv.Type() returned nil")
		} else {
			t.Logf("Encoder case: rv.Type() returned %p", typ)
		}
	})
}
