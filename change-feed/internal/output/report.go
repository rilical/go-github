package output

import "github.com/github/openapi-change-feed/internal/changes"

func renderReport(_ changes.ChangeBatch) string {
	return "# API changes\n"
}
