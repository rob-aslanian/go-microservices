package company

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/client"
	"google.golang.org/grpc/metadata"
)

type CompanyClient struct {
	company companyRPC.CompanyServiceClient
}

func NewCompanyClient(addr string) (*CompanyClient, error) {
	companyCon, err := client.CreateGrpcConnection(addr)
	if err != nil {
		return nil, err
	}
	company := companyRPC.NewCompanyServiceClient(companyCon)

	return &CompanyClient{company}, nil
}

func (a *CompanyClient) GetProfilesByIDs(ctx context.Context, ids []string) (interface{}, error) {
	passContext(&ctx)

	profiles, err := a.company.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	return profiles.GetProfiles(), nil
}

func passContext(ctx *context.Context) {

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("error while passing context")
	}
}
