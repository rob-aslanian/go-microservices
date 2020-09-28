package servicerequest

// Action is the action that can be taken for each service/request
type Action string

const (
	// ActionPause ...
	ActionPause Action = "pause"
	// ActionEdit ...
	ActionEdit Action = "edit"
	// ActionDelete ...
	ActionDelete Action = "delete"
	//ActionShare ...
	ActionShare Action = "share"
)

