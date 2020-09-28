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
	n := Stuff{}
	log.Println("stuff connected to:", settings.Address)
	n.connect(settings.Address)
	return n
}

func (n *Stuff) connect(address string) {
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
	n.stuffClient = stuffRPC.NewStuffServiceClient(conn)
}

// ContactInvitationForWallet ... 
func (n Stuff) ContactInvitationForWallet(ctx context.Context , name string , email string , message string , coins int32) error{
	_ , err := n.stuffClient.ContactInvitationForWallet(ctx, &stuffRPC.InvitationWalletRequest{
			Name:name,
			Email:email,
			SilverCoins:coins,
			Message:message,
		})

	if err != nil {
		return err
	}

	return nil
}

// AddGoldCoinsToWallet ... 
func (n Stuff) AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32) error{
	_ , err := n.stuffClient.AddGoldCoinsToWallet(ctx, &stuffRPC.WalletAddGoldCoins{
			UserID:userID,
			Coins:coins,
		})

	if err != nil {
		return err
	}

	return nil
}

// CreateWalletAccount ... 
func (n Stuff) CreateWalletAccount(ctx context.Context , userID string) error {
	 _ , err := n.stuffClient.CreateWalletAccount(ctx, &stuffRPC.UserId{
		ID:userID,
	})

	if err != nil {
		return err
	}

	return nil
}