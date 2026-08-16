package tests

import (
	"os"
	"testing"

	"github.com/teranga-host/mimic/internal/domain"
	"github.com/teranga-host/mimic/internal/parser"
)

func TestDeclarativeParser(t *testing.T) {
	yamlContent := `
endpoints:
  - path: /api/v1/users
    method: GET
    resource: users
    is_crud: true
  - path: /api/v1/users/{id}
    method: GET
    resource: users
    is_crud: true
`
	tmpFile, err := os.CreateTemp("", "mimic-test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, _ = tmpFile.WriteString(yamlContent)
	_ = tmpFile.Close()

	endpoints, err := parser.DetectAndParse(tmpFile.Name())
	if err != nil {
		t.Fatalf("DetectAndParse failed: %v", err)
	}

	if len(endpoints) != 2 {
		t.Fatalf("Expected 2 endpoints, got %d", len(endpoints))
	}

	if endpoints[0].Path != "/api/v1/users" || endpoints[0].Method != domain.MethodGet {
		t.Errorf("Unexpected endpoint 0: %+v", endpoints[0])
	}
}

func TestLaravelParser(t *testing.T) {
	phpContent := `<?php
use Illuminate\Support\Facades\Route;

Route::get('/orders', [OrderController::class, 'index']);
Route::post('/orders', [OrderController::class, 'store']);
Route::apiResource('products', ProductController::class);
`
	tmpFile, err := os.CreateTemp("", "api.php")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, _ = tmpFile.WriteString(phpContent)
	_ = tmpFile.Close()

	endpoints, err := parser.DetectAndParse(tmpFile.Name())
	if err != nil {
		t.Fatalf("DetectAndParse Laravel failed: %v", err)
	}

	// 2 standard routes + 5 apiResource routes = 7 endpoints
	if len(endpoints) != 7 {
		t.Fatalf("Expected 7 endpoints from Laravel routes, got %d", len(endpoints))
	}
}
