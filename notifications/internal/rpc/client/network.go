package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

// ---

// IsAdmin ...
func (n Network) IsAdmin(ctx context.Context, companyID string) (bool, error) {
	passContext(&ctx)

	_, err := n.networkClient.GetAdminObject(ctx, &networkRPC.Company{
		Id: companyID,
	})
	log.Println(err)

	// handleError(err)
	if err != nil {
		return false, err
	}

	// ---------------

	return true, nil
}

func passContext(ctx *context.Context) {
	// span := s.tracer.MakeSpan(*ctx, "passContext")
	// defer span.Finish()

	// fmt.Printf("Before: %+v\n", *ctx)

	md, ok := metadata.FromIncomingContext(*ctx)
	if ok {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("can't pass context")
	}

	// fmt.Printf("After: %+v\n", *ctx)
}
