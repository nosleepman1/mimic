package parser

import (
	"regexp"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
)

type LaravelParser struct{}

func (p *LaravelParser) CanParse(filePath string, content []byte) bool {
	contentStr := string(content)
	return strings.HasSuffix(filePath, "api.php") ||
		strings.Contains(contentStr, "Route::") ||
		strings.Contains(contentStr, "use Illuminate\\Support\\Facades\\Route")
}

// Regex pour matcher Route::get('/users', ...), Route::post('users/{id}', ...)
var routeRegex = regexp.MustCompile(`Route::(get|post|put|patch|delete|options|any)\s*\(\s*['"]([^'"]+)['"]`)

// Regex pour Route::apiResource('users', UserController::class)
var resourceRegex = regexp.MustCompile(`Route::apiResource\s*\(\s*['"]([^'"]+)['"]`)

func (p *LaravelParser) Parse(filePath string, content []byte) ([]*domain.Endpoint, error) {
	contentStr := string(content)
	endpoints := make([]*domain.Endpoint, 0)

	// 1. Parser les routes standard (Route::get, Route::post, etc.)
	matches := routeRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			methodStr := strings.ToUpper(match[1])
			pathStr := match[2]

			if methodStr == "ANY" {
				for _, m := range []domain.HTTPMethod{domain.MethodGet, domain.MethodPost, domain.MethodPut, domain.MethodDelete} {
					endpoints = append(endpoints, createLaravelEndpoint(pathStr, m))
				}
			} else {
				endpoints = append(endpoints, createLaravelEndpoint(pathStr, domain.HTTPMethod(methodStr)))
			}
		}
	}

	// 2. Parser les apiResource (Route::apiResource('users', ...))
	resourceMatches := resourceRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range resourceMatches {
		if len(match) >= 2 {
			resourceName := match[1]
			basePath := "/" + strings.TrimPrefix(resourceName, "/")

			// GET /users (index)
			endpoints = append(endpoints, &domain.Endpoint{
				Path:      basePath,
				Method:    domain.MethodGet,
				Resource:  resourceName,
				IsCRUD:    true,
				Summary:   "List " + resourceName,
				Responses: map[int]*domain.Response{200: {StatusCode: 200}},
			})

			// POST /users (store)
			endpoints = append(endpoints, &domain.Endpoint{
				Path:      basePath,
				Method:    domain.MethodPost,
				Resource:  resourceName,
				IsCRUD:    true,
				Summary:   "Create " + resourceName,
				Responses: map[int]*domain.Response{201: {StatusCode: 201}},
			})

			// GET /users/{id} (show)
			endpoints = append(endpoints, &domain.Endpoint{
				Path:      basePath + "/{id}",
				Method:    domain.MethodGet,
				Resource:  resourceName,
				IsCRUD:    true,
				Summary:   "Get " + resourceName,
				Responses: map[int]*domain.Response{200: {StatusCode: 200}},
			})

			// PUT /users/{id} (update)
			endpoints = append(endpoints, &domain.Endpoint{
				Path:      basePath + "/{id}",
				Method:    domain.MethodPut,
				Resource:  resourceName,
				IsCRUD:    true,
				Summary:   "Update " + resourceName,
				Responses: map[int]*domain.Response{200: {StatusCode: 200}},
			})

			// DELETE /users/{id} (destroy)
			endpoints = append(endpoints, &domain.Endpoint{
				Path:      basePath + "/{id}",
				Method:    domain.MethodDelete,
				Resource:  resourceName,
				IsCRUD:    true,
				Summary:   "Delete " + resourceName,
				Responses: map[int]*domain.Response{204: {StatusCode: 204}},
			})
		}
	}

	return endpoints, nil
}

func createLaravelEndpoint(pathStr string, method domain.HTTPMethod) *domain.Endpoint {
	cleanPath := CleanPath(pathStr)
	resource := extractResourceName(cleanPath)

	statusCode := 200
	if method == domain.MethodPost {
		statusCode = 201
	} else if method == domain.MethodDelete {
		statusCode = 204
	}

	return &domain.Endpoint{
		Path:      cleanPath,
		Method:    method,
		Resource:  resource,
		IsCRUD:    resource != "",
		Responses: map[int]*domain.Response{statusCode: {StatusCode: statusCode}},
	}
}

func extractResourceName(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && !strings.HasPrefix(parts[0], "{") {
		return parts[0]
	}
	return ""
}
