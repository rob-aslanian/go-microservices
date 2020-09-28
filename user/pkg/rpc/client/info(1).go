package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"google.golang.org/grpc"
)

// Info represents client of Info
type Info struct {
	infoClient infoRPC.InfoServiceClient
}

// NewInfoClient crates new gRPC client of Info
func NewInfoClient(settings Settings) Info {
	a := Info{}
	a.connect(settings.Address)
	return a
}

func (a *Info) connect(address string) {
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
	a.infoClient = infoRPC.NewInfoServiceClient(conn)
}
