package tinyreflect_test

import (
	"testing"

	"github.com/cdvelop/tinyreflect"
)

// Define User type at package level for testing consistency
type User struct {
	Name string
	Age  int
}

// Test for unique struct identification using hash-based approach
func TestUniqueStructIdentification(t *testing.T) {
	tr := tinyreflect.New()

	type Product struct {
		Title string
		Price float64
	}

	// Same structure name but different fields (simulates different packages)
	type User2 struct {
		ID     int
		Email  string
		Active bool
	}

	tests := []struct {
		name         string
		value        interface{}
		expectedKind string
	}{
		{
			name:         "User struct (2 fields)",
			value:        User{},
			expectedKind: "struct",
		},
		{
			name:         "Product struct (2 fields)",
			value:        Product{},
			expectedKind: "struct",
		},
		{
			name:         "User2 struct (3 fields)",
			value:        User2{},
			expectedKind: "struct",
		},
	}

	seenHashes := make(map[uint32]string)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := tr.TypeOf(tt.value)

			// Test struct ID generation
			structID := typ.StructID()
			t.Logf("StructID: %d", structID)

			// Test name
			name := typ.Name()
			t.Logf("Name: %s", name)

			// Verify each struct has a unique hash
			hash := typ.Hash
			if prevName, exists := seenHashes[hash]; exists {
				t.Errorf("Hash collision! %s and %s have same hash: %d", prevName, tt.name, hash)
			} else {
				seenHashes[hash] = tt.name
				t.Logf("Unique hash %d for %s", hash, tt.name)
			}

			// Verify Kind is correct
			if typ.Kind().String() != tt.expectedKind {
				t.Errorf("Expected kind %s, got %s", tt.expectedKind, typ.Kind().String())
			}

			// Verify StructID is not empty for structs
			if typ.Kind().String() == "struct" && structID == 0 {
				t.Error("StructID should not be empty for struct types")
			}
		})
	}

	t.Logf("Total unique struct hashes: %d", len(seenHashes))
}

// CRITICAL TEST: Same struct type from different initialization locations
// must have the SAME hash - this validates Go's runtime hash consistency
func TestSameStructSameHash(t *testing.T) {
	tr := tinyreflect.New()

	// Initialize User struct from different places
	user1 := User{}
	user2 := User{Name: "John", Age: 30}
	user3 := createUser()
	user4 := createUserPointer()

	// Get types
	type1 := tr.TypeOf(user1)
	type2 := tr.TypeOf(user2)
	type3 := tr.TypeOf(user3)
	type4 := tr.TypeOf(*user4) // Dereference pointer

	// All should have the SAME hash
	hash1 := type1.Hash
	hash2 := type2.Hash
	hash3 := type3.Hash
	hash4 := type4.Hash

	t.Logf("Hash from empty init: %d", hash1)
	t.Logf("Hash from value init: %d", hash2)
	t.Logf("Hash from function: %d", hash3)
	t.Logf("Hash from pointer deref: %d", hash4)

	// Critical validation - all must be equal
	if hash1 != hash2 {
		t.Errorf("CRITICAL: Hash mismatch! Empty init (%d) != Value init (%d)", hash1, hash2)
	}
	if hash1 != hash3 {
		t.Errorf("CRITICAL: Hash mismatch! Empty init (%d) != Function (%d)", hash1, hash3)
	}
	if hash1 != hash4 {
		t.Errorf("CRITICAL: Hash mismatch! Empty init (%d) != Pointer deref (%d)", hash1, hash4)
	}

	// Test StructID consistency
	id1 := type1.StructID()
	id2 := type2.StructID()
	id3 := type3.StructID()
	id4 := type4.StructID()

	t.Logf("StructID 1: %d", id1)
	t.Logf("StructID 2: %d", id2)
	t.Logf("StructID 3: %d", id3)
	t.Logf("StructID 4: %d", id4)

	if id1 != id2 || id1 != id3 || id1 != id4 {
		t.Errorf("CRITICAL: StructID inconsistency! Same struct type must have same StructID")
		t.Errorf("Expected all to be: %d", id1)
		t.Errorf("Got: %d, %d, %d, %d", id1, id2, id3, id4)
	} else {
		t.Log("SUCCESS: Same struct type has consistent StructID across different initializations")
	}
}

// Helper functions for different initialization patterns
func createUser() User {
	return User{Name: "Helper", Age: 25}
}

func createUserPointer() *User {
	return &User{Name: "Pointer", Age: 40}
}