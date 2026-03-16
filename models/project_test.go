package models

import (
	"encoding/json"
	"testing"
)

func TestProject_JSONFields(t *testing.T) {
	project := Project{
		Name:     "Auth Refactor",
		Code:     "12572",
		Category: "Infrastructure",
	}

	b, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("marshal project: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal project json: %v", err)
	}

	if got["name"] != project.Name {
		t.Fatalf("expected name %q, got %q", project.Name, got["name"])
	}

	if got["code"] != project.Code {
		t.Fatalf("expected code %q, got %q", project.Code, got["code"])
	}

	if got["category"] != project.Category {
		t.Fatalf("expected category %q, got %q", project.Category, got["category"])
	}
}
