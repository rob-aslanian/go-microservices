package service

import (
	"context"
)

// AuthRPC represents Auth gRPC client
type AuthRPC interface {
	GetUserID(ctx context.Context, token string) (string, error)
}

// NetworkRPC represents Network gRPC client
type NetworkRPC interface {
	IsAdmin(ctx context.Context, companyID string) (bool, error)
}
