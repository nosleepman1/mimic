package domain

// Request capture les détails de la requête reçue par le mock
type Request struct {
	Method      HTTPMethod             `json:"method"`
	Path        string                 `json:"path"`
	PathParams  map[string]string      `json:"path_params"`
	QueryParams map[string][]string    `json:"query_params"`
	Headers     map[string]string      `json:"headers"`
	Body        map[string]interface{} `json:"body"`
}

// Response définit la réponse simulée par le mock
type Response struct {
	StatusCode int                    `json:"status_code" yaml:"status_code"`
	Headers    map[string]string      `json:"headers" yaml:"headers"`
	Body       interface{}            `json:"body" yaml:"body"`
	Schema     *Schema                `json:"schema" yaml:"schema"`
}
