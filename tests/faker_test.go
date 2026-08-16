package tests

import (
	"strings"
	"testing"

	"github.com/teranga-host/mimic/internal/domain"
	"github.com/teranga-host/mimic/internal/faker"
)

func TestFakerGenerator(t *testing.T) {
	gen := faker.NewGenerator()

	// 1. Phone number test
	phone := faker.RandomPhoneSN()
	if !strings.HasPrefix(phone, "+221") {
		t.Errorf("Expected phone to start with +221, got %s", phone)
	}

	// 2. Full Name test
	name := faker.RandomFullName()
	if len(name) < 3 || !strings.Contains(name, " ") {
		t.Errorf("Unexpected name format: %s", name)
	}

	// 3. Schema generation
	schema := &domain.Schema{
		Type: domain.TypeObject,
		Properties: map[string]*domain.Schema{
			"email":  {Type: domain.TypeEmail},
			"phone":  {Type: domain.TypePhoneSN},
			"amount": {Type: domain.TypeXOF},
			"active": {Type: domain.TypeBoolean},
		},
	}

	res := gen.GenerateFromSchema(schema, "user")
	resMap, ok := res.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result from schema generator, got %T", res)
	}

	if _, ok := resMap["email"].(string); !ok {
		t.Errorf("Expected string email, got %T", resMap["email"])
	}
	if _, ok := resMap["active"].(bool); !ok {
		t.Errorf("Expected bool active, got %T", resMap["active"])
	}
}
