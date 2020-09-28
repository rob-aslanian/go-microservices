package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"google.golang.org/grpc"
)

// Chat represents client of Chat
type Chat struct {
	chatClient chatRPC.ChatServiceClient
}

// NewChatClient crates new gRPC client of Network
func NewChatClient(settings Settings) Chat {
	n := Chat{}
	n.connect(settings.Address)
	return n
}

func (n *Chat) connect(address string) {
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
	n.chatClient = chatRPC.NewChatServiceClient(conn)
}

// IsLive ...
func (n Chat) IsLive(ctx context.Context, id string) (bool, error) {
	passContext(&ctx)

	value, err := n.chatClient.IsUserLive(ctx, &chatRPC.IsUserLiveRequest{
		UserId: id,
	})

	err = handleError(err)
	if err != nil {
		return false, err
	}

	return value.GetValue(), nil
}
