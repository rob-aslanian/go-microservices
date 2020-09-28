package service

import "context"

// AuthRPC represents a Auth gRPC client
type AuthRPC interface {
	GetUser(ctx context.Context, token string) (string, error)
}


// UserRPC ... 
type UserRPC interface {
	GetUserByInvitedID(ctx context.Context , userID string) (int32 , error)
}

