package chat

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/client"
	"google.golang.org/grpc/metadata"
)

type ChatClient struct {
	chat chatRPC.ChatServiceClient
}

func NewChatClient(addr string) (*ChatClient, error) {
	chatCon, err := client.CreateGrpcConnection(addr)
	if err != nil {
		return nil, err
	}
	chat := chatRPC.NewChatServiceClient(chatCon)

	return &ChatClient{chat}, nil
}

func (a *ChatClient) GetProfilesByIDs(ctx context.Context, senderID, targetID string, value bool) error {
	passContext(&ctx)

	_, err := a.chat.BlockConversetionByParticipants(ctx, &chatRPC.BlockRequest{
		SenderID: senderID,
		TargetID: targetID,
		Value:    value,
	})
	if err != nil {
		return err
	}

	return nil
}

func passContext(ctx *context.Context) {

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("error while passing context")
	}
}
