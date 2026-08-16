package tests

import (
	"testing"

	"github.com/teranga-host/mimic/internal/store"
)

func TestMemoryStoreCRUD(t *testing.T) {
	s := store.NewMemoryStore()

	// 1. Create
	user := map[string]interface{}{
		"name":  "Mamadou Diop",
		"email": "mamadou@example.com",
	}
	created := s.Create("users", user)
	id, ok := created["id"].(string)
	if !ok || id == "" {
		t.Fatalf("Expected generated ID for created user")
	}

	// 2. Get
	fetched, found := s.Get("users", id)
	if !found {
		t.Fatalf("Expected to find user with ID %s", id)
	}
	if fetched["name"] != "Mamadou Diop" {
		t.Errorf("Expected name 'Mamadou Diop', got %v", fetched["name"])
	}

	// 3. Collection
	collection := s.GetCollection("users")
	if len(collection) != 1 {
		t.Fatalf("Expected 1 item in collection, got %d", len(collection))
	}

	// 4. Update
	updated, ok := s.Update("users", id, map[string]interface{}{"name": "Mamadou Updated"})
	if !ok || updated["name"] != "Mamadou Updated" {
		t.Errorf("Update failed or unexpected name: %v", updated)
	}

	// 5. Delete
	deleted := s.Delete("users", id)
	if !deleted {
		t.Fatalf("Expected Delete to return true")
	}

	_, foundAfter := s.Get("users", id)
	if foundAfter {
		t.Fatalf("User should not exist after deletion")
	}
}
