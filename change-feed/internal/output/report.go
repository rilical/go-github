package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/github/openapi-change-feed/internal/changes"
)

func renderReport(b changes.ChangeBatch) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# API changes %s..%s\n\n", b.BaseSpec.ETag, b.HeadSpec.ETag)
	fmt.Fprintf(&sb, "%d changes - %d breaking, %d warn, %d info (%d added, %d removed)\n\n",
		b.Summary.Total, b.Summary.Breaking, b.Summary.Warn, b.Summary.Info, b.Summary.Added, b.Summary.Removed)

	recs := append([]changes.ChangeRecord(nil), b.Changes...)
	sort.SliceStable(recs, func(i, j int) bool {
		if recs[i].Severity != recs[j].Severity {
			return recs[i].Severity > recs[j].Severity
		}
		if recs[i].Path != recs[j].Path {
			return recs[i].Path < recs[j].Path
		}
		return recs[i].ID < recs[j].ID
	})

	for _, bucket := range []changes.Severity{changes.Breaking, changes.Warn, changes.Info} {
		var lines []string
		for _, r := range recs {
			if r.Severity != bucket {
				continue
			}
			lines = append(lines, fmt.Sprintf("- `%s %s` %s", r.Method, r.Path, r.Text))
		}
		if len(lines) == 0 {
			continue
		}
		fmt.Fprintf(&sb, "## %s\n\n%s\n\n", strings.ToUpper(bucket.String()), strings.Join(lines, "\n"))
	}

	return sb.String()
}
