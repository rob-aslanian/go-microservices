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

// NewCompanyClient crates new gRPC client of User
func NewCompanyClient(settings Settings) Company {
	a := Company{}
	a.connect(settings.Address)
	return a
}

func (u *Company) connect(address string) {
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
	u.companyClient = companyRPC.NewCompanyServiceClient(conn)
}

// // CheckPassword check password of user
// func (u User) CheckPassword(ctx context.Context, password string) error {
// 	passContext(&ctx)
//
// 	_, err := u.userClient.CheckPassword(ctx, &userRPC.CheckPasswordRequest{
// 		Password: password,
// 	})
//
// 	// TODO: handle error
//
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// // GetProfilesByID ...
// func (u User) GetProfilesByID(ctx context.Context, ids []string) (interface{}, error) {
// 	passContext(&ctx)
//
// 	resp, err := u.userClient.GetProfilesByID(ctx, &userRPC.UserIDs{
// 		ID: ids,
// 		// Language: // TODO:
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// sort result
// 	orderedProfiles := make([]*userRPC.Profile, len(ids))
// 	for c, id := range ids {
// 		for i := range resp.GetProfiles() {
// 			if (resp.GetProfiles())[i].GetID() == id {
// 				log.Println(id, resp.GetProfiles()[i].GetID())
// 				orderedProfiles[c] = resp.GetProfiles()[i]
// 			}
// 		}
// 	}
//
// 	return orderedProfiles, nil
// }
//
// // GetMapProfilesByID ...
// func (u User) GetMapProfilesByID(ctx context.Context, ids []string) (map[string]interface{}, error) {
// 	passContext(&ctx)
//
// 	resp, err := u.userClient.GetMapProfilesByID(ctx, &userRPC.UserIDs{
// 		ID: ids,
// 		// Language: // TODO:
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	m := make(map[string]interface{}, len(resp.GetProfiles()))
// 	for key, value := range resp.GetProfiles() {
// 		m[key] = value
// 	}
//
// 	return m, nil
// }

// GetCompanyProfiles ...
func (u *Company) GetCompanyProfiles(ctx context.Context, ids []string) ([]*companyRPC.Profile, error) {
	passContext(&ctx)

	resp, err := u.companyClient.GetCompanyProfiles(ctx, &companyRPC.GetCompanyProfilesRequest{
		Ids: ids,
		// Language: // TODO:
	})
	if err != nil {
		return nil, err
	}

	// sort result
	orderedProfiles := make([]*companyRPC.Profile, len(ids))
	for c, id := range ids {
		for i := range resp.GetProfiles() {
			if (resp.GetProfiles())[i].GetId() == id {
				orderedProfiles[c] = resp.GetProfiles()[i]
			}
		}
	}

	return orderedProfiles, nil
}
