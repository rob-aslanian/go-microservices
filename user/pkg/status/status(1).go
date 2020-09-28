package status

// UserStatus ...
type UserStatus string

const (
	// UserStatusUnKnown ... 
	UserStatusUnKnown UserStatus = "UNKNOWN"
	// UserStatusNotActivated ...
	UserStatusNotActivated UserStatus = "NOT_ACTIVATED"

	// UserStatusActivated ...
	UserStatusActivated UserStatus = "ACTIVATED"

	// UserStatusDeactivated ...
	UserStatusDeactivated UserStatus = "DEACTIVATED"
)
