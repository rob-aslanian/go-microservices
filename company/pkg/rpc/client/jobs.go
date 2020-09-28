package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"google.golang.org/grpc"
)

// Jobs represents client of Jobs
type Jobs struct {
	jobsClient jobsRPC.JobsServiceClient
}

// NewJobsClient crates new gRPC client of Jobs
func NewJobsClient(settings Settings) Jobs {
	a := Jobs{}
	a.connect(settings.Address)
	return a
}

func (a *Jobs) connect(address string) {
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
	a.jobsClient = jobsRPC.NewJobsServiceClient(conn)
}

// GetAmountOfActiveJobsOfCompany ...
func (a Jobs) GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error) {
	passContext(&ctx)

	res, err := a.jobsClient.GetAmountOfActiveJobsOfCompany(ctx, &jobsRPC.ID{
		Id: companyID,
	})

	handleError(err)

	// ---------------

	return res.GetAmount(), nil
}