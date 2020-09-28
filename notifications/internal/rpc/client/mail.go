package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/mailRPC"
	"google.golang.org/grpc"
)

// Mail represents client of Mail
type Mail struct {
	mailClient mailRPC.MailServiceClient
}

// NewMailClient crates new gRPC client of Mail
func NewMailClient(settings Settings) Mail {
	m := Mail{}
	m.connect(settings.Address)
	return m
}

func (m *Mail) connect(address string) {
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
	m.mailClient = mailRPC.NewMailServiceClient(conn)
}
