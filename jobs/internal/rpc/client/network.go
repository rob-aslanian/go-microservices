package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	companyadmin "gitlab.lan/Rightnao-site/microservices/jobs/internal/company-admin"
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

// GetAdminLevel ...
func (n Network) GetAdminLevel(ctx context.Context, companyID string) (companyadmin.AdminLevel, error) {
	passContext(&ctx)

	resp, err := n.networkClient.GetAdminObject(ctx, &networkRPC.Company{
		Id: companyID,
	})

	// TODO: refactor
	if err != nil && err.Error() == `rpc error: code = Unknown desc = {"type":"NotFountError","description":"User is not admin of the company","data":null}` {
		return companyadmin.AdminLevelUnknown, nil
	}

	// TODO: handle error
	if err != nil {
		return companyadmin.AdminLevelUnknown, err
	}

	return adminLevelRPCToAccountAdminLevel(resp.GetLevel()), nil
}

func accountAdminLevelToRPC(adminLevel companyadmin.AdminLevel) networkRPC.AdminLevel {
	al := networkRPC.AdminLevel_Admin

	switch adminLevel {
	case companyadmin.AdminLevelAdmin:
		al = networkRPC.AdminLevel_Admin
	case companyadmin.AdminLevelJob:
		al = networkRPC.AdminLevel_JobAdmin
	case companyadmin.AdminLevelVShop:
		al = networkRPC.AdminLevel_VShopAdmin
	case companyadmin.AdminLevelVService:
		al = networkRPC.AdminLevel_VServiceAdmin
	case companyadmin.AdminLevelCommercial:
		al = networkRPC.AdminLevel_CommercialAdmin
		// default:
		// 	al = networkRPC.AdminLevel_Unknown
	}

	return al
}

func adminLevelRPCToAccountAdminLevel(adminLevel networkRPC.AdminLevel) companyadmin.AdminLevel {
	al := companyadmin.AdminLevelUnknown

	switch adminLevel {
	case networkRPC.AdminLevel_Admin:
		al = companyadmin.AdminLevelAdmin
	case networkRPC.AdminLevel_JobAdmin:
		al = companyadmin.AdminLevelJob
	case networkRPC.AdminLevel_VShopAdmin:
		al = companyadmin.AdminLevelVShop
	case networkRPC.AdminLevel_VServiceAdmin:
		al = companyadmin.AdminLevelVService
	case networkRPC.AdminLevel_CommercialAdmin:
		al = companyadmin.AdminLevelCommercial
	}

	return al
}
