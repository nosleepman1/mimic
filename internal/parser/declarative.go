package parser

import (
	"encoding/json"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
	"gopkg.in/yaml.v3"
)

type DeclarativeFile struct {
	Endpoints []DeclarativeEndpoint `json:"endpoints" yaml:"endpoints"`
}

type DeclarativeEndpoint struct {
	Path        string                   `json:"path" yaml:"path"`
	Method      string                   `json:"method" yaml:"method"`
	Summary     string                   `json:"summary" yaml:"summary"`
	Resource    string                   `json:"resource" yaml:"resource"`
	IsCRUD      bool                     `json:"is_crud" yaml:"is_crud"`
	StatusCode  int                      `json:"status" yaml:"status"`
	Headers     map[string]string        `json:"headers" yaml:"headers"`
	Response    interface{}              `json:"response" yaml:"response"`
	Schema      *domain.Schema           `json:"schema" yaml:"schema"`
}

type DeclarativeParser struct{}

func (p *DeclarativeParser) CanParse(filePath string, content []byte) bool {
	contentStr := string(content)
	return strings.Contains(contentStr, "endpoints:") || strings.Contains(contentStr, "\"endpoints\"")
}

func (p *DeclarativeParser) Parse(filePath string, content []byte) ([]*domain.Endpoint, error) {
	var decFile DeclarativeFile

	if err := yaml.Unmarshal(content, &decFile); err != nil {
		if errJson := json.Unmarshal(content, &decFile); errJson != nil {
			return nil, err
		}
	}

	result := make([]*domain.Endpoint, 0, len(decFile.Endpoints))

	for _, dep := range decFile.Endpoints {
		method := domain.HTTPMethod(strings.ToUpper(dep.Method))
		if method == "" {
			method = domain.MethodGet
		}

		statusCode := dep.StatusCode
		if statusCode == 0 {
			if method == domain.MethodPost {
				statusCode = 201
			} else {
				statusCode = 200
			}
		}

		resp := &domain.Response{
			StatusCode: statusCode,
			Headers:    dep.Headers,
			Body:       dep.Response,
			Schema:     dep.Schema,
		}

		endpoint := &domain.Endpoint{
			Path:      CleanPath(dep.Path),
			Method:    method,
			Summary:   dep.Summary,
			Resource:  dep.Resource,
			IsCRUD:    dep.IsCRUD,
			Responses: map[int]*domain.Response{statusCode: resp},
		}

		result = append(result, endpoint)
	}

	return result, nil
}
