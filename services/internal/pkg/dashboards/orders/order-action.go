package order

// Action is the order action that can be taken for each service/request
type Action string

const (
	// ActionAccept ...
	ActionAccept Action = "accept"
	// ActionDecline ...
	ActionDecline Action = "decline"
)
