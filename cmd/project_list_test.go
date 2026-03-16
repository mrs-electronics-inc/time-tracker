package cmd

import (
	"bytes"
	"strings"
	"testing"

	"time-tracker/models"
)

type projectListTestStorage struct {
	projects []models.Project
}

func (s projectListTestStorage) LoadProjects() ([]models.Project, error) {
	copyOfProjects := make([]models.Project, len(s.projects))
	copy(copyOfProjects, s.projects)
	return copyOfProjects, nil
}

func TestListProjects_SortsCaseInsensitiveAndShowsColumns(t *testing.T) {
	storage := projectListTestStorage{
		projects: []models.Project{
			{Name: "zeta", Code: "Z", Category: "Backlog"},
			{Name: "alpha", Code: "A", Category: "Client"},
			{Name: "Beta", Code: "B", Category: "Internal"},
		},
	}

	var out bytes.Buffer
	if err := listProjects(storage, &out); err != nil {
		t.Fatalf("listProjects returned error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Name") || !strings.Contains(output, "Code") || !strings.Contains(output, "Category") {
		t.Fatalf("expected Name/Code/Category columns, got output:\n%s", output)
	}

	alphaIndex := strings.Index(output, "alpha")
	betaIndex := strings.Index(output, "Beta")
	zetaIndex := strings.Index(output, "zeta")

	if alphaIndex < 0 || betaIndex < 0 || zetaIndex < 0 {
		t.Fatalf("expected all project names in output, got:\n%s", output)
	}

	if !(alphaIndex < betaIndex && betaIndex < zetaIndex) {
		t.Fatalf("expected projects sorted by case-insensitive name, got:\n%s", output)
	}
}
