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

func (a *User) connect(address string) {
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
	a.userClient = userRPC.NewUserServiceClient(conn)
}


// GetUserByInvitedID ... 
func (a User) GetUserByInvitedID(ctx context.Context , userID string) (int32 , error) {
	count , err := a.userClient.GetUserByInvitedID(ctx, &userRPC.UserId{
		Id:userID,
	})

	if err != nil {
		return 0 ,  err
	}

	return count.GetCount() , nil
}
