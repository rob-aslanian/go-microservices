package serverRPC

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/session"
)

// Service ...
type Service interface {
	LoginUser(ctx context.Context, id string) (string, error)
	GetUser(ctx context.Context, session string) (string, error)
	LogoutUser(ctx context.Context, session string) error
	LogoutOtherSession(ctx context.Context, sessionID string) error
	SignOutFromAll(ctx context.Context, session string) error
	GetListOfSessions(ctx context.Context, first, after int32) ([]session.Session, error)
	GetAmountOfSessions(ctx context.Context) (int32, error)
	GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error)
}
