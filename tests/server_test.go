package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teranga-host/mimic/internal/domain"
	"github.com/teranga-host/mimic/internal/server"
	"github.com/teranga-host/mimic/internal/store"
	"github.com/teranga-host/mimic/internal/ui"
)

func TestMockHandlerHTTPFlow(t *testing.T) {
	endpoints := []*domain.Endpoint{
		{
			Path:     "/api/v1/articles",
			Method:   domain.MethodGet,
			Resource: "articles",
			IsCRUD:   true,
		},
		{
			Path:     "/api/v1/articles",
			Method:   domain.MethodPost,
			Resource: "articles",
			IsCRUD:   true,
		},
		{
			Path:     "/api/v1/articles/{id}",
			Method:   domain.MethodGet,
			Resource: "articles",
			IsCRUD:   true,
		},
	}

	router := server.NewRouter(endpoints)
	memStore := store.NewMemoryStore()
	logger := ui.NewLogger()
	handler := server.NewMockHandler(router, memStore, logger)

	// 1. Test POST /api/v1/articles
	postBody, _ := json.Marshal(map[string]interface{}{
		"title": "Introduction to Go and Mimic",
	})
	reqPost := httptest.NewRequest(http.MethodPost, "/api/v1/articles", bytes.NewReader(postBody))
	recPost := httptest.NewRecorder()
	handler.ServeHTTP(recPost, reqPost)

	if recPost.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", recPost.Code)
	}

	var createdArticle map[string]interface{}
	_ = json.Unmarshal(recPost.Body.Bytes(), &createdArticle)
	articleID, ok := createdArticle["id"].(string)
	if !ok || articleID == "" {
		t.Fatalf("Expected created article to contain 'id'")
	}

	// 2. Test GET /api/v1/articles/{id}
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+articleID, nil)
	recGet := httptest.NewRecorder()
	handler.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", recGet.Code)
	}

	var fetchedArticle map[string]interface{}
	_ = json.Unmarshal(recGet.Body.Bytes(), &fetchedArticle)
	if fetchedArticle["title"] != "Introduction to Go and Mimic" {
		t.Errorf("Expected title match, got %v", fetchedArticle["title"])
	}

	// 3. Test 404 for unknown route
	req404 := httptest.NewRequest(http.MethodGet, "/api/v1/unknown-endpoint", nil)
	rec404 := httptest.NewRecorder()
	handler.ServeHTTP(rec404, req404)

	if rec404.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", rec404.Code)
	}
}
