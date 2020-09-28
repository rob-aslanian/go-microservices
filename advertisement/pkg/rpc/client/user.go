package clientRPC

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// User ...
type User struct {
	userClient userRPC.UserServiceClient
}

// NewUserClient crates new gRPC client of Auth
func NewUserClient(settings Settings) User {
	u := User{}
	u.connect(settings.Address)
	return u
}

func (u *User) connect(address string) {
	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	u.userClient = userRPC.NewUserServiceClient(conn)

}

// GetUsersForAdvert ...
func (u User) GetUsersForAdvert(ctx context.Context, data account.UserForAdvert) ([]string, error) {
	passThroughContext(&ctx)

	res, err := u.userClient.GetUsersForAdvert(ctx, &userRPC.GetUsersForAdvertRequest{
		OwnerID:   data.OwnerID,
		Gender:    data.Gender,
		AgeFrom:   data.AgeFrom,
		AgeTo:     data.AgeTo,
		Languages: data.Languages,
		Locations: data.Locations,
	})

	if err != nil {
		return nil, err
	}

	return res.GetIDs(), nil
}
