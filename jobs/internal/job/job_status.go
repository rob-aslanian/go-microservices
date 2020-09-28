package job

// Status ...
type Status string

const (
	// StatusActive ...
	StatusActive Status = "Active"
	// StatusDraft ...
	StatusDraft Status = "Draft"
	// StatusPaused ...
	StatusPaused Status = "Paused"
	// ExpiredPaused ...
	ExpiredPaused Status = "Expired"
)
