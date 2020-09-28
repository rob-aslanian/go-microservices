package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"google.golang.org/grpc"
)

// Network represents client of Network
type Network struct {
	networkClient networkRPC.NetworkServiceClient
}

// NewNetworkClient crates new gRPC client of Network
func NewNetworkClient(settings Settings) Network {
	n := Network{}
	n.connect(settings.Address)
	return n
}

func (n *Network) connect(address string) {
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
	n.networkClient = networkRPC.NewNetworkServiceClient(conn)
}
