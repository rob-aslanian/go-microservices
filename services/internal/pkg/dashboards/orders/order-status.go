package order

// Status is used to view the progress of each order
type Status string

const (
	// StatusnNew ...
	StatusnNew Status = "new"
	// StatusInProgress ...
	StatusInProgress Status = "in_progress"
	// StatusOutOfSchedule ...
	StatusOutOfSchedule Status = "out_of_schedule"
	// Delivered ..
	Delivered Status = "delivered"
	// Completed ...
	Completed Status = "completed"
	// Disputed ...
	Disputed Status = "disputed"
	// Canceled ...
	Canceled Status = "canceled"
)
