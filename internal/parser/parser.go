package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/teranga-host/mimic/internal/domain"
)

// Parser interface universelle pour l'ingestion de contrats d'API
type Parser interface {
	CanParse(filePath string, content []byte) bool
	Parse(filePath string, content []byte) ([]*domain.Endpoint, error)
}

var parsers = []Parser{
	&DeclarativeParser{},
	&OpenAPIParser{},
	&PostmanParser{},
	&LaravelParser{},
}

// DetectAndParse détecte automatiquement le format du fichier et extrait tous les endpoints
func DetectAndParse(filePath string) ([]*domain.Endpoint, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidContractFile, err)
	}

	for _, p := range parsers {
		if p.CanParse(filePath, content) {
			return p.Parse(filePath, content)
		}
	}

	ext := filepath.Ext(filePath)
	return nil, fmt.Errorf("%w: impossible de parser le fichier '%s' (extension %s)", domain.ErrUnsupportedFormat, filePath, ext)
}

// CleanPath standardise les chemins d'URL (ex: /users/:id -> /users/{id})
func CleanPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + strings.TrimPrefix(part, ":") + "}"
		}
	}

	return strings.Join(parts, "/")
}
