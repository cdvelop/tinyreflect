package tinyreflect

import (
	"fmt"
	"testing"
)

func TestOriginalIssue(t *testing.T) {
	// Recreate the original failing scenario
	numbers := []int{10, 20, 30}
	v := refValueOf(numbers)

	if v.err != nil {
		t.Errorf("refValueOf error: %v", v.err)
		return
	}

	fmt.Printf("Slice Kind: %s\n", v.refKind())
	fmt.Printf("Slice Length: %d\n", v.refLen())

	// Test element access - this was failing before
	elem := v.refIndex(1)
	if elem.err != nil {
		t.Errorf("Element access error: %v", elem.err)
		return
	}

	fmt.Printf("Element at index 1: %d\n", elem.refInt())

	// Verify the value is correct
	if elem.refInt() != 20 {
		t.Errorf("Expected 20, got %d", elem.refInt())
	}
}
