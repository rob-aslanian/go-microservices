package serverHTTP

import "context"

// Service represent service
type Service interface {
	// ActivateUser(ctx context.Context, code string, userID string) error
	// ActivateEmail(ctx context.Context, code string, userID string) error
}

// AuthRPC ...
type AuthRPC interface {
	GetUserID(ctx context.Context, token string) (string, error)
}

// NetworkRPC ...
type NetworkRPC interface {
	IsAdmin(ctx context.Context, companyID string) (bool, error)
}
