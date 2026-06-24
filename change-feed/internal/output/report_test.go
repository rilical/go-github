package output

import (
	"strings"
	"testing"
)

func TestReportIsByteStable(t *testing.T) {
	if renderReport(sampleBatch()) != renderReport(sampleBatch()) {
		t.Fatal("report.md must be deterministic")
	}
}

func TestReportShowsSummaryLine(t *testing.T) {
	if !strings.Contains(renderReport(sampleBatch()), "1 changes") {
		t.Error("report should state the change count")
	}
}
