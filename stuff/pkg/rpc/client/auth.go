package clientRPC

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth represents client of Auth
type Auth struct {
	authClient authRPC.AuthServiceClient
}

// NewAuthClient crates new gRPC client of Auth
func NewAuthClient(settings Settings) Auth {
	a := Auth{}
	a.connect(settings.Address)
	return a
}

func (a *Auth) connect(address string) {
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
	a.authClient = authRPC.NewAuthServiceClient(conn)
}

// GetUser returns user id
func (a Auth) GetUser(ctx context.Context, token string) (string, error) {
	u, err := a.authClient.GetUser(ctx, &authRPC.Session{
		Token: token,
	})

	stat, isOk := status.FromError(err)
	if isOk {
		if stat.Code() != codes.OK {
			fmt.Println("Code:", stat.Code(), "\tMessage:", stat.Message())
			return "", errors.New(stat.Message())
		}
	}

	return u.GetId(), nil
}
