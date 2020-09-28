package user

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/client"
	"google.golang.org/grpc/metadata"
)

type UserClient struct {
	user userRPC.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	userCon, err := client.CreateGrpcConnection(addr)
	if err != nil {
		return nil, err
	}
	user := userRPC.NewUserServiceClient(userCon)

	return &UserClient{user}, nil
}

func (a *UserClient) GetProfilesByIDs(ctx context.Context, ids []string) (interface{}, error) {
	passContext(&ctx)

	profiles, err := a.user.GetProfilesByID(ctx, &userRPC.UserIDs{
		ID: ids,
	})
	if err != nil {
		return nil, err
	}

	return profiles.GetProfiles(), nil
}

// GetConectionsPrivacy ...
func (a *UserClient) GetConectionsPrivacy(ctx context.Context, userID string) (string, error) {
	passContext(&ctx)

	res, err := a.user.GetConectionsPrivacy(ctx, &userRPC.ID{
		ID: userID,
	})
	if err != nil {
		return "", err
	}

	s := res.GetType()

	return s.String(), nil
}

func passContext(ctx *context.Context) {

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("error while passing context")
	}
}
