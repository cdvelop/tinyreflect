package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

func TestPointerTypeElemChain(t *testing.T) {
	tr := tinyreflect.New()
	// Test simple pointer
	var p *int
	typ1 := tr.TypeOf(p)
	t.Logf("Type of *int: %p, Kind: %v", typ1, typ1.Kind())

	elem1 := typ1.Elem()
	if elem1 == nil {
		t.Fatal("Elem of *int should not be nil")
	}
	t.Logf("Elem of *int: %p, Kind: %v", elem1, elem1.Kind())

	// Test pointer to pointer
	var pp **int
	typ2 := tr.TypeOf(pp)
	t.Logf("Type of **int: %p, Kind: %v", typ2, typ2.Kind())

	elem2 := typ2.Elem()
	if elem2 == nil {
		t.Fatal("Elem of **int is nil!")
	}
	t.Logf("Elem of **int: %p, Kind: %v", elem2, elem2.Kind())

	// Test elem of elem
	elem2elem := elem2.Elem()
	if elem2elem == nil {
		t.Fatal("Elem of elem of **int is nil!")
	}
	t.Logf("Elem of Elem of **int: %p, Kind: %v", elem2elem, elem2elem.Kind())
}