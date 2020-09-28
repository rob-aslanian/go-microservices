package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"google.golang.org/grpc"
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

// GetUserID returns user id
func (a Auth) GetUserID(ctx context.Context) (string, error) {
	passContext(&ctx)

	token := retriveToken(ctx)

	u, err := a.authClient.GetUser(ctx, &authRPC.Session{
		Token: token,
	})

	handleError(err)

	// ---------------

	return u.GetId(), nil
}
