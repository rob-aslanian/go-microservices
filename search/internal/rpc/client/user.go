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

// GetProfilesByID ...
func (u User) GetProfilesByID(ctx context.Context, ids []string) ([]*userRPC.Profile, error) {
	passContext(&ctx)

	resp, err := u.userClient.GetProfilesByID(ctx, &userRPC.UserIDs{
		ID: ids,
		// Language: // TODO:
	})
	if err != nil {
		return nil, err
	}

	// sort result
	orderedProfiles := make([]*userRPC.Profile, len(ids))
	for c, id := range ids {
		for i := range resp.GetProfiles() {
			if (resp.GetProfiles())[i].GetID() == id {
				log.Println(id, resp.GetProfiles()[i].GetID())
				orderedProfiles[c] = resp.GetProfiles()[i]
			}
		}
	}

	return orderedProfiles, nil
}

// GetMapProfilesByID ...
func (u User) GetMapProfilesByID(ctx context.Context, ids []string) (map[string]interface{}, error) {
	passContext(&ctx)

	resp, err := u.userClient.GetMapProfilesByID(ctx, &userRPC.UserIDs{
		ID: ids,
		// Language: // TODO:
	})
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{}, len(resp.GetProfiles()))
	for key, value := range resp.GetProfiles() {
		m[key] = value
	}

	return m, nil
}
