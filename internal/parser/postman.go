package parser

import (
	"encoding/json"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
)

type PostmanCollection struct {
	Info struct {
		Name   string `json:"name"`
		Schema string `json:"schema"`
	} `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanItem struct {
	Name     string          `json:"name"`
	Item     []PostmanItem   `json:"item,omitempty"` // Dossiers imbriqués
	Request  *PostmanRequest `json:"request,omitempty"`
	Response []PostmanResponse `json:"response,omitempty"`
}

type PostmanRequest struct {
	Method string      `json:"method"`
	URL    PostmanURL  `json:"url"`
	Header []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"header"`
}

type PostmanURL struct {
	Raw  string   `json:"raw"`
	Path []string `json:"path"`
}

type PostmanResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Code   int    `json:"code"`
	Body   string `json:"body"`
}

type PostmanParser struct{}

func (p *PostmanParser) CanParse(filePath string, content []byte) bool {
	contentStr := string(content)
	return strings.Contains(contentStr, "schema.getpostman.com/json/collection") ||
		strings.Contains(contentStr, "\"_postman_id\"")
}

func (p *PostmanParser) Parse(filePath string, content []byte) ([]*domain.Endpoint, error) {
	var coll PostmanCollection
	if err := json.Unmarshal(content, &coll); err != nil {
		return nil, err
	}

	var endpoints []*domain.Endpoint
	p.extractItems(coll.Item, &endpoints)

	return endpoints, nil
}

func (p *PostmanParser) extractItems(items []PostmanItem, out *[]*domain.Endpoint) {
	for _, it := range items {
		if it.Item != nil && len(it.Item) > 0 {
			// Dossier récursif
			p.extractItems(it.Item, out)
			continue
		}

		if it.Request == nil {
			continue
		}

		path := "/" + strings.Join(it.Request.URL.Path, "/")
		cleanPath := CleanPath(path)
		method := domain.HTTPMethod(strings.ToUpper(it.Request.Method))
		if method == "" {
			method = domain.MethodGet
		}

		endpoint := &domain.Endpoint{
			Path:      cleanPath,
			Method:    method,
			Summary:   it.Name,
			Responses: make(map[int]*domain.Response),
		}

		// Récupérer les exemples de réponses définis dans Postman
		for _, resp := range it.Response {
			code := resp.Code
			if code == 0 {
				code = 200
			}

			var parsedBody interface{}
			if err := json.Unmarshal([]byte(resp.Body), &parsedBody); err != nil {
				parsedBody = resp.Body
			}

			endpoint.Responses[code] = &domain.Response{
				StatusCode: code,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       parsedBody,
			}
		}

		if len(endpoint.Responses) == 0 {
			endpoint.Responses[200] = &domain.Response{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}
		}

		*out = append(*out, endpoint)
	}
}
