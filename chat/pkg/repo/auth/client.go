package auth

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/client"
	"golang.org/x/net/context"
)

type AuthClient struct {
	auth authRPC.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	authCon, err := client.CreateGrpcConnection(addr)
	if err != nil {
		return nil, err
	}
	auth := authRPC.NewAuthServiceClient(authCon)

	return &AuthClient{auth}, nil
}

func (a *AuthClient) GetUserId(ctx context.Context, token string) (string, error) {
	user, err := a.auth.GetUser(ctx, &authRPC.Session{Token:token})
	if err != nil {
		return "", err
	}
	return user.Id, nil
}
