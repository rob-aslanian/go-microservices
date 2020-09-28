package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"google.golang.org/grpc"
)

// Company represents client of Company
type Company struct {
	companyClient companyRPC.CompanyServiceClient
}

// NewCompanyClient crates new gRPC client of Company
func NewCompanyClient(settings Settings) Company {
	a := Company{}
	a.connect(settings.Address)
	return a
}

func (a *Company) connect(address string) {
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
	a.companyClient = companyRPC.NewCompanyServiceClient(conn)
}
