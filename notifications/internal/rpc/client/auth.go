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

//  ----------------------------------

// GetUserID returns user id
func (a Auth) GetUserID(ctx context.Context, token string) (string, error) {
	u, err := a.authClient.GetUser(ctx, &authRPC.Session{
		Token: token,
	})

	// handleError(err)
	if err != nil {
		return "", err
	}

	// ---------------

	return u.GetId(), nil
}

//
// // LoginUser creates session for user. Returns token.
// func (a Auth) LoginUser(ctx context.Context, userID string) (string, error) {
// 	result, err := a.authClient.LoginUser(ctx, &authRPC.User{
// 		Id: userID,
// 	})
//
// 	handleError(err)
//
// 	return result.GetToken(), nil
// }
//
// // SignOut ...
// func (a Auth) SignOut(ctx context.Context, token string) error {
// 	_, err := a.authClient.LogoutSession(ctx, &authRPC.Session{
// 		Token: token,
// 	})
// 	handleError(err)
// 	return err
// }
//
// // GetTimeOfLastActivity ...
// func (a Auth) GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error) {
// 	tm, err := a.authClient.GetTimeOfLastActivity(ctx, &authRPC.User{
// 		Id: id,
// 	})
//
// 	handleError(err)
// 	return time.Unix(0, tm.GetTime()), err
// }
