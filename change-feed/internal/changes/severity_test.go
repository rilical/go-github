package changes

import "testing"

func TestSeverityFromOasdiffLevelAndBreaking(t *testing.T) {
	for _, c := range []struct {
		level int
		want  Severity
		brk   bool
	}{{1, Info, false}, {2, Warn, false}, {3, Breaking, true}} {
		got := SeverityFromOasdiffLevel(c.level)
		if got != c.want || got.IsBreaking() != c.brk {
			t.Errorf("level %d -> %v (breaking=%v), want %v (breaking=%v)",
				c.level, got, got.IsBreaking(), c.want, c.brk)
		}
	}
}
