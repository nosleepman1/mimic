package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
	"gopkg.in/yaml.v3"
)

type OpenAPIParser struct{}

type RawOpenAPISpec struct {
	OpenAPI string                            `json:"openapi" yaml:"openapi"`
	Swagger string                            `json:"swagger" yaml:"swagger"`
	Paths   map[string]map[string]RawOperation `json:"paths" yaml:"paths"`
}

type RawOperation struct {
	Summary     string                 `json:"summary" yaml:"summary"`
	Description string                 `json:"description" yaml:"description"`
	Tags        []string               `json:"tags" yaml:"tags"`
	Parameters  []RawParameter         `json:"parameters" yaml:"parameters"`
	Responses   map[string]RawResponse `json:"responses" yaml:"responses"`
}

type RawParameter struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Required    bool   `json:"required" yaml:"required"`
	Description string `json:"description" yaml:"description"`
}

type RawResponse struct {
	Description string                  `json:"description" yaml:"description"`
	Content     map[string]RawMediaType `json:"content" yaml:"content"`
	Schema      *RawSchema              `json:"schema" yaml:"schema"` // Swagger 2.0 fallback
}

type RawMediaType struct {
	Schema  *RawSchema  `json:"schema" yaml:"schema"`
	Example interface{} `json:"example" yaml:"example"`
}

type RawSchema struct {
	Type        string                `json:"type" yaml:"type"`
	Format      string                `json:"format" yaml:"format"`
	Properties  map[string]*RawSchema `json:"properties" yaml:"properties"`
	Items       *RawSchema            `json:"items" yaml:"items"`
	Example     interface{}           `json:"example" yaml:"example"`
	Description string                `json:"description" yaml:"description"`
}

func (p *OpenAPIParser) CanParse(filePath string, content []byte) bool {
	contentStr := string(content)
	return strings.Contains(contentStr, "openapi:") ||
		strings.Contains(contentStr, "\"openapi\"") ||
		strings.Contains(contentStr, "swagger:") ||
		strings.Contains(contentStr, "\"swagger\"")
}

func (p *OpenAPIParser) Parse(filePath string, content []byte) ([]*domain.Endpoint, error) {
	var spec RawOpenAPISpec

	if err := yaml.Unmarshal(content, &spec); err != nil {
		if errJson := json.Unmarshal(content, &spec); errJson != nil {
			return nil, fmt.Errorf("erreur syntaxe OpenAPI/Swagger: %w", err)
		}
	}

	endpoints := make([]*domain.Endpoint, 0)

	for pathStr, methods := range spec.Paths {
		cleanPath := CleanPath(pathStr)

		for methodStr, op := range methods {
			methodUpper := strings.ToUpper(methodStr)
			if !isStandardHTTPMethod(methodUpper) {
				continue
			}

			endpoint := &domain.Endpoint{
				Path:       cleanPath,
				Method:     domain.HTTPMethod(methodUpper),
				Summary:    op.Summary,
				Tags:       op.Tags,
				Parameters: make([]domain.Parameter, 0),
				Responses:  make(map[int]*domain.Response),
			}

			for _, param := range op.Parameters {
				endpoint.Parameters = append(endpoint.Parameters, domain.Parameter{
					Name:        param.Name,
					In:          param.In,
					Required:    param.Required,
					Description: param.Description,
				})
			}

			// Parsing des réponses
			for codeStr, respData := range op.Responses {
				code, err := strconv.Atoi(codeStr)
				if err != nil {
					code = 200 // default fallback
				}

				resp := &domain.Response{
					StatusCode: code,
					Headers:    map[string]string{"Content-Type": "application/json"},
				}

				// OpenAPI 3.x
				if jsonMedia, ok := respData.Content["application/json"]; ok {
					if jsonMedia.Example != nil {
						resp.Body = jsonMedia.Example
					} else if jsonMedia.Schema != nil {
						resp.Schema = convertSchema(jsonMedia.Schema)
					}
				} else if respData.Schema != nil {
					// Swagger 2.0
					resp.Schema = convertSchema(respData.Schema)
				}

				endpoint.Responses[code] = resp
			}

			if len(endpoint.Responses) == 0 {
				endpoint.Responses[200] = &domain.Response{
					StatusCode: 200,
					Headers:    map[string]string{"Content-Type": "application/json"},
				}
			}

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

func convertSchema(raw *RawSchema) *domain.Schema {
	if raw == nil {
		return nil
	}

	s := &domain.Schema{
		Type:        domain.FieldType(raw.Type),
		Format:      raw.Format,
		Example:     raw.Example,
		Description: raw.Description,
	}

	if raw.Properties != nil {
		s.Properties = make(map[string]*domain.Schema)
		for k, v := range raw.Properties {
			s.Properties[k] = convertSchema(v)
		}
	}

	if raw.Items != nil {
		s.Items = convertSchema(raw.Items)
	}

	return s
}

func isStandardHTTPMethod(m string) bool {
	switch m {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD":
		return true
	default:
		return false
	}
}
