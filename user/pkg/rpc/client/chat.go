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

// NewChatClient crates new gRPC client of Chat
func NewChatClient(settings Settings) Chat {
	a := Chat{}
	a.connect(settings.Address)
	return a
}

func (a *Chat) connect(address string) {
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
	a.chatClient = chatRPC.NewChatServiceClient(conn)
}
