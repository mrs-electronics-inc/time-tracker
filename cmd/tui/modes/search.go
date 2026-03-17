package modes

import (
	"strings"

	"time-tracker/models"
)

func entryMatchesSearchQuery(entry models.TimeEntry, query string) bool {
	if query == "" {
		return true
	}

	normalizedQuery := strings.ToLower(query)
	project := strings.ToLower(entry.Project)
	title := strings.ToLower(entry.Title)

	return strings.Contains(project, normalizedQuery) || strings.Contains(title, normalizedQuery)
}
