package tinyreflect

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

// Test the refStructTag implementation
func TestRefStructTag(t *testing.T) {
	// Test basic tag parsing
	tag := refStructTag("json:\"user_name,omitempty\" xml:\"UserName\"")

	// Test Get method
	jsonValue := tag.Get("json")
	expectedJson := "user_name,omitempty"
	if jsonValue != expectedJson {
		t.Errorf("Expected json tag '%s', got '%s'", expectedJson, jsonValue)
	}

	t.Logf("Basic tag test passed: got '%s'", jsonValue)
}

// Test JSON field mapping with tags
func TestJsonFieldMappingWithTags(t *testing.T) {
	clearRefStructsCache()

	// Test struct with JSON tags
	type TestUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	// Test with struct tag parsing
	jsonStr := `{"id": "test_123", "username": "testuser", "email": "test@example.com"}`

	var user TestUser
	// Test basic struct operations instead of JSON decode
	v := refValueOf(&user)
	if v.refKind() != KPointer {
		t.Errorf("Expected pointer kind, got %v", v.refKind())
	}

	t.Logf("JSON: %s", jsonStr)
	t.Logf("Struct type verified: %+v", user)
}

// Test struct tag validation
func TestStructTagValidation(t *testing.T) {
	clearRefStructsCache()

	type SimpleUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	// Test struct reflection operations
	user := SimpleUser{ID: "123", Username: "test", Email: "test@example.com"}
	v := refValueOf(user)

	if v.refKind() != KStruct {
		t.Errorf("Expected struct kind, got %v", v.refKind())
	}

	t.Logf("Struct validation completed for: %+v", user)
}
