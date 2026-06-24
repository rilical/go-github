package changes

import (
	"encoding/json"
	"fmt"
)

type Severity int

const (
	Info     Severity = iota + 1 // 1: oasdiff INFO
	Warn                         // 2: oasdiff WARN
	Breaking                     // 3: oasdiff ERR; the trusted headline signal
)

func (s Severity) String() string {
	switch s {
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Breaking:
		return "breaking"
	default:
		return "unknown"
	}
}

// MarshalJSON emits the severity as a self-describing string ("info"/"warn"/"breaking")
// so the change feed is readable by any consumer without the Go enum mapping.
func (s Severity) MarshalJSON() ([]byte, error) {
	if s < Info || s > Breaking {
		return nil, fmt.Errorf("changes: cannot marshal unknown severity %d", int(s))
	}
	return json.Marshal(s.String())
}

func (s *Severity) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	switch str {
	case "info":
		*s = Info
	case "warn":
		*s = Warn
	case "breaking":
		*s = Breaking
	default:
		return fmt.Errorf("changes: unknown severity %q", str)
	}
	return nil
}

func (s Severity) IsBreaking() bool { return s == Breaking }

func SeverityFromOasdiffLevel(level int) Severity {
	switch level {
	case 3:
		return Breaking
	case 2:
		return Warn
	default:
		return Info
	}
}
