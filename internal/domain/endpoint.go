package domain

// HTTPMethod type
type HTTPMethod string

const (
	MethodGet     HTTPMethod = "GET"
	MethodPost    HTTPMethod = "POST"
	MethodPut     HTTPMethod = "PUT"
	MethodPatch   HTTPMethod = "PATCH"
	MethodDelete  HTTPMethod = "DELETE"
	MethodOptions HTTPMethod = "OPTIONS"
	MethodHead    HTTPMethod = "HEAD"
)

// Parameter représente un paramètre d'URL (path ou query)
type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"` // "path", "query", "header"
	Required    bool   `json:"required" yaml:"required"`
	Type        string `json:"type" yaml:"type"` // "string", "integer", "uuid"
	Description string `json:"description" yaml:"description"`
}

// Endpoint définit une route d'API à simuler
type Endpoint struct {
	Path        string            `json:"path" yaml:"path"`               // Ex: "/api/v1/users/{id}"
	Method      HTTPMethod        `json:"method" yaml:"method"`           // "GET", "POST", etc.
	Summary     string            `json:"summary" yaml:"summary"`
	Tags        []string          `json:"tags" yaml:"tags"`
	Parameters  []Parameter       `json:"parameters" yaml:"parameters"`
	RequestBody *Schema           `json:"request_body" yaml:"request_body"`
	Responses   map[int]*Response `json:"responses" yaml:"responses"`     // Code HTTP -> Réponse
	Resource    string            `json:"resource" yaml:"resource"`       // Nom de la ressource CRUD (ex: "users")
	IsCRUD      bool              `json:"is_crud" yaml:"is_crud"`         // Active la persistance en mémoire
}
