package changes

type Severity int

const (
	Info     Severity = iota + 1 // 1 — oasdiff INFO
	Warn                         // 2 — oasdiff WARN
	Breaking                     // 3 — oasdiff ERR; the trusted headline signal
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
