package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	companyadmin "gitlab.lan/Rightnao-site/microservices/search/internal/company-admin"
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

// GetIDsOfFriends ...
func (n Network) GetIDsOfFriends(ctx context.Context, userID string) ([]string, error) {
	passContext(&ctx)

	resp, err := n.networkClient.GetFriendIdsOf(ctx, &networkRPC.User{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetList(), nil
}

// GetIDsOfFollowingCompanies ...
func (n Network) GetIDsOfFollowingCompanies(ctx context.Context) ([]string, error) {
	passContext(&ctx)

	resp, err := n.networkClient.GetFilteredFollowingCompanies(ctx, &networkRPC.CompanyFilter{})
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(resp.GetFollows()))

	for _, f := range resp.GetFollows() {
		if f.GetCompany() != nil {
			ids = append(ids, f.GetCompany().GetId())
		}
	}

	return ids, nil
}

// GetBlockedIDs ...
func (n Network) GetBlockedIDs(ctx context.Context) ([]string, error) {
	passContext(&ctx)

	resp, err := n.networkClient.GetblockedUsersOrCompanies(ctx, &networkRPC.Empty{})
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(resp.GetList()))

	for _, f := range resp.GetList() {
		ids = append(ids, f.GetId())
	}

	return ids, nil
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
