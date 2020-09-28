package serverHTTP

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
)

// Service represent service
type Service interface {
	ActivateUser(ctx context.Context, code string, userID string) (res *account.LoginResponse, err error)
	ActivateEmail(ctx context.Context, code string, userID string) error
}
