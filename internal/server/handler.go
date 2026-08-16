package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
	"github.com/teranga-host/mimic/internal/faker"
	"github.com/teranga-host/mimic/internal/store"
	"github.com/teranga-host/mimic/internal/ui"
)

type MockHandler struct {
	router *Router
	store  *store.MemoryStore
	faker  *faker.Generator
	logger *ui.Logger
}

func NewMockHandler(router *Router, store *store.MemoryStore, logger *ui.Logger) *MockHandler {
	return &MockHandler{
		router: router,
		store:  store,
		faker:  faker.NewGenerator(),
		logger: logger,
	}
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := domain.HTTPMethod(r.Method)
	path := r.URL.Path

	match, found := h.router.Match(method, path)
	if !found {
		h.logger.LogRequest(r.Method, path, http.StatusNotFound, "Route not found in contract")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "Not Found",
			"message": "Endpoint is not defined in the API contract",
			"path":    path,
			"method":  r.Method,
		})
		return
	}

	endpoint := match.Endpoint

	// Gérer CRUD Stateful si la ressource est définie
	if endpoint.IsCRUD && endpoint.Resource != "" {
		if handled := h.handleCRUD(w, r, endpoint, match.PathParams); handled {
			return
		}
	}

	// Réponse standard basée sur le contrat et Faker
	h.handleStaticOrFaker(w, r, endpoint)
}

func (h *MockHandler) handleCRUD(w http.ResponseWriter, r *http.Request, ep *domain.Endpoint, pathParams map[string]string) bool {
	res := ep.Resource
	id, hasID := pathParams["id"]

	switch r.Method {
	case "GET":
		if hasID {
			// GET /resource/{id}
			item, found := h.store.Get(res, id)
			if !found {
				// Générer un item à la volée et le stocker
				generated := h.generateItem(ep, id)
				h.store.Create(res, generated)
				item = generated
			}
			h.respondJSON(w, r, http.StatusOK, item)
			return true
		} else {
			// GET /resource (collection)
			collection := h.store.GetCollection(res)
			if len(collection) == 0 {
				// Pré-remplir 3 éléments
				for i := 0; i < 3; i++ {
					h.store.Create(res, h.generateItem(ep, ""))
				}
				collection = h.store.GetCollection(res)
			}
			h.respondJSON(w, r, http.StatusOK, collection)
			return true
		}

	case "POST":
		// POST /resource
		bodyData := h.parseBody(r)
		created := h.store.Create(res, bodyData)
		h.respondJSON(w, r, http.StatusCreated, created)
		return true

	case "PUT", "PATCH":
		// PUT/PATCH /resource/{id}
		if hasID {
			bodyData := h.parseBody(r)
			updated, ok := h.store.Update(res, id, bodyData)
			if !ok {
				updated = h.store.Create(res, bodyData)
			}
			h.respondJSON(w, r, http.StatusOK, updated)
			return true
		}

	case "DELETE":
		// DELETE /resource/{id}
		if hasID {
			_ = h.store.Delete(res, id)
			h.logger.LogRequest(r.Method, r.URL.Path, http.StatusNoContent, "Deleted")
			w.WriteHeader(http.StatusNoContent)
			return true
		}
	}

	return false
}

func (h *MockHandler) handleStaticOrFaker(w http.ResponseWriter, r *http.Request, ep *domain.Endpoint) {
	statusCode := 200
	var respDef *domain.Response

	// Trouver la réponse 200/201 en priorité
	for code, def := range ep.Responses {
		if code >= 200 && code < 300 {
			statusCode = code
			respDef = def
			break
		}
	}

	if respDef == nil {
		for code, def := range ep.Responses {
			statusCode = code
			respDef = def
			break
		}
	}

	// Si un body statique / exemple est présent
	if respDef != nil && respDef.Body != nil {
		h.respondJSON(w, r, statusCode, respDef.Body)
		return
	}

	// Si un schéma est présent
	if respDef != nil && respDef.Schema != nil {
		fakeData := h.faker.GenerateFromSchema(respDef.Schema, ep.Resource)
		h.respondJSON(w, r, statusCode, fakeData)
		return
	}

	// Fallback : générer un objet réaliste déduit du chemin
	fallbackData := h.generateItem(ep, "")
	h.respondJSON(w, r, statusCode, fallbackData)
}

func (h *MockHandler) generateItem(ep *domain.Endpoint, idHint string) map[string]interface{} {
	item := make(map[string]interface{})
	if idHint != "" {
		item["id"] = idHint
	} else {
		item["id"] = h.faker.GuessFromFieldName("id")
	}

	fieldNames := []string{"name", "email", "phone_sn", "status", "created_at"}
	for _, f := range fieldNames {
		item[f] = h.faker.GuessFromFieldName(f)
	}

	return item
}

func (h *MockHandler) parseBody(r *http.Request) map[string]interface{} {
	bodyData := make(map[string]interface{})
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil && len(bodyBytes) > 0 {
			_ = json.Unmarshal(bodyBytes, &bodyData)
		}
	}
	return bodyData
}

func (h *MockHandler) respondJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}

	summary := ""
	if str, ok := data.(string); ok && len(str) < 50 {
		summary = strings.TrimSpace(str)
	}
	h.logger.LogRequest(r.Method, r.URL.Path, status, summary)
}
