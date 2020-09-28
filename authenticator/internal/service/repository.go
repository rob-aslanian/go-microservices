package service

import (
	"context"
	"time"

	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/location"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/session"
)

// AuthRPC gets the user ID and returns it as string
// If we want to use any of these functions, we have to call thoough this interface
// This is so that different services won't be connectected direclty, thus making future chages easier
type AuthRPC interface {
	GetUserID(ctx context.Context) (string, error)
}

// AuthentificationRepository is interface for authentification.
// This is so that different services won't be connectected direclty, thus making future chages easier
type AuthentificationRepository interface {
	InsertSession(ctx context.Context, ses session.Session) error
	DeactivateSessionByToken(ctx context.Context, token string) error
	DeactivateSpecificSessionByToken(cx context.Context, userID, sessionID string) (string, error)
	UpdateActivityTime(ctx context.Context, token string) error
	GetListOfSessions(ctx context.Context, id string, first, after int32) ([]session.Session, error)
	GetAmountOfSessions(ctx context.Context, userID string) (int32, error)
	GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error)
	DeactivateSessions(ctx context.Context, tokenException string) (string, []string, error)
}

// IPDBRepository recieves the city as string and returns to us as location type.
// If we want to use any of these functions, we have to call thoough this interface
// This is so that different services won't be connectected direclty, thus making future chages easier
type IPDBRepository interface {
	GetCityByIP(ip string) (location.Location, error)
}

// TokensRepository is interface for tokens.
// If we want to use any of these functions, we have to call thoough this interface
// This is so that different services won't be connectected direclty, thus making future chages easier
type TokensRepository interface {
	SaveToken(token string, id string) error
	DeleteTokens(tokens []string) error
	GetUserByToken(token string) (string, error)
	DeleteToken(token string) error
	ExpireToken(token string) error
}
