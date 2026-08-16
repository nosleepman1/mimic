package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryStore gère la persistance thread-safe des ressources CRUD en mémoire
type MemoryStore struct {
	mu          sync.RWMutex
	collections map[string]map[string]map[string]interface{} // resource -> id -> item
	order       map[string][]string                          // resource -> list of ids
}

// NewMemoryStore instancie le store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		collections: make(map[string]map[string]map[string]interface{}),
		order:       make(map[string][]string),
	}
}

// GetCollection retourne tous les éléments d'une ressource
func (s *MemoryStore) GetCollection(resource string) []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items, exists := s.collections[resource]
	if !exists {
		return []map[string]interface{}{}
	}

	orderList := s.order[resource]
	result := make([]map[string]interface{}, 0, len(items))

	for _, id := range orderList {
		if item, ok := items[id]; ok {
			result = append(result, item)
		}
	}

	return result
}

// Get récupère un élément spécifique par son ID
func (s *MemoryStore) Get(resource, id string) (map[string]interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items, exists := s.collections[resource]
	if !exists {
		return nil, false
	}

	item, ok := items[id]
	return item, ok
}

// Create insère un nouvel élément en générant un ID si absent
func (s *MemoryStore) Create(resource string, data map[string]interface{}) map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.collections[resource]; !exists {
		s.collections[resource] = make(map[string]map[string]interface{})
		s.order[resource] = make([]string, 0)
	}

	// Déterminer ou générer l'ID
	idVal, ok := data["id"]
	var idStr string
	if !ok || idVal == nil || fmt.Sprintf("%v", idVal) == "" {
		idStr = uuid.New().String()
		data["id"] = idStr
	} else {
		idStr = fmt.Sprintf("%v", idVal)
	}

	// Ajouter timestamps
	if _, ok := data["created_at"]; !ok {
		data["created_at"] = time.Now().UTC().Format(time.RFC3339)
	}
	if _, ok := data["updated_at"]; !ok {
		data["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	}

	s.collections[resource][idStr] = data
	s.order[resource] = append(s.order[resource], idStr)

	return data
}

// Update met à jour un élément existant
func (s *MemoryStore) Update(resource, id string, data map[string]interface{}) (map[string]interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, exists := s.collections[resource]
	if !exists {
		return nil, false
	}

	existing, ok := items[id]
	if !ok {
		return nil, false
	}

	for k, v := range data {
		existing[k] = v
	}
	existing["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	existing["id"] = id // Garantir que l'ID n'est pas altéré

	return existing, true
}

// Delete supprime un élément
func (s *MemoryStore) Delete(resource, id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, exists := s.collections[resource]
	if !exists {
		return false
	}

	if _, ok := items[id]; !ok {
		return false
	}

	delete(items, id)

	// Mettre à jour l'ordre
	newOrder := make([]string, 0, len(s.order[resource])-1)
	for _, existingID := range s.order[resource] {
		if existingID != id {
			newOrder = append(newOrder, existingID)
		}
	}
	s.order[resource] = newOrder

	return true
}

// Seed pré-remplit une collection avec des données initiales
func (s *MemoryStore) Seed(resource string, items []map[string]interface{}) {
	for _, item := range items {
		s.Create(resource, item)
	}
}
