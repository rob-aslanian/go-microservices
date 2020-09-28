package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
)

// AuthRPC represents a Auth gRPC client
type AuthRPC interface {
	GetUser(ctx context.Context, token string) (string, error)
}

// UserRPC ...
type UserRPC interface {
	GetUsersForAdvert(ctx context.Context, data account.UserForAdvert) ([]string, error)
}
