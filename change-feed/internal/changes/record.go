package changes

import "time"

type ChangeRecord struct {
	ID          string     `json:"id"`
	Fingerprint string     `json:"fingerprint,omitempty"`
	Kind        ChangeKind `json:"kind"`
	Severity    Severity   `json:"severity"`
	Method      string     `json:"method,omitempty"`
	Path        string     `json:"path,omitempty"`
	OperationID string     `json:"operation_id,omitempty"`
	Section     string     `json:"section,omitempty"`
	Text        string     `json:"text"`
	BaseSpec    string     `json:"base_spec,omitempty"`
	HeadSpec    string     `json:"head_spec,omitempty"`
}

type Summary struct {
	Total    int `json:"total"`
	Breaking int `json:"breaking"`
	Warn     int `json:"warn"`
	Info     int `json:"info"`
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
}

type SpecRef struct {
	URL       string    `json:"url"`
	ETag      string    `json:"etag"`
	FetchedAt time.Time `json:"fetched_at"`
}

type ChangeBatch struct {
	GeneratedAt time.Time      `json:"generated_at"`
	BaseSpec    SpecRef        `json:"base_spec"`
	HeadSpec    SpecRef        `json:"head_spec"`
	Changes     []ChangeRecord `json:"changes"`
	Summary     Summary        `json:"summary"`
}
