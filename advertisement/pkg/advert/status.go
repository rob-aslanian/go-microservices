package advert

// Status represent status advert
type Status string

const (
	// StatusDraft draft
	StatusDraft Status = "draft"
	// StatusActive active status
	StatusActive Status = "active"
	// StatusInActive active status
	StatusInActive Status = "in_active"
	// StatusPaused paused
	StatusPaused Status = "paused"
	// StatusCompleted completed
	StatusCompleted Status = "completed"
	// StatusNotRunning not running. Can be if insufficient funds
	StatusNotRunning Status = "not_running"
	// StatusRejected rejected. Can be after moderation
	StatusRejected Status = "rejected"
	// StatusScheduled rejected.
	StatusScheduled Status = "scheduled"
)
