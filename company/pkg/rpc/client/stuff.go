package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	"google.golang.org/grpc"
)

// Stuff represents client of Stuff
type Stuff struct {
	stuffClient stuffRPC.StuffServiceClient
}

// NewStuffClient crates new gRPC client of Stuff
func NewStuffClient(settings Settings) Stuff {
	a := Stuff{}
	a.connect(settings.Address)
	return a
}

func (a *Stuff) connect(address string) {
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
	a.stuffClient = stuffRPC.NewStuffServiceClient(conn)
}

// AddGoldCoinsToWallet ... 
func (a Stuff) AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32) error{
	_ , err := a.stuffClient.AddGoldCoinsToWallet(ctx, &stuffRPC.WalletAddGoldCoins{
			UserID:userID,
			Coins:coins,
		})

	if err != nil {
		return err
	}

	return nil
}

