package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"google.golang.org/grpc"
)

// User represents client of User
type User struct {
	userClient userRPC.UserServiceClient
}

// NewUserClient crates new gRPC client of User
func NewUserClient(settings Settings) User {
	a := User{}
	a.connect(settings.Address)
	return a
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

// CheckPassword check password of user
func (u User) CheckPassword(ctx context.Context, password string) error {
	passContext(&ctx)

	_, err := u.userClient.CheckPassword(ctx, &userRPC.CheckPasswordRequest{
		Password: password,
	})

	// TODO: handle error

	if err != nil {
		return err
	}

	return nil
}
