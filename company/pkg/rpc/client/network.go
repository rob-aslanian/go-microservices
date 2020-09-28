package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
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

// MakeCompanyOwner ...
func (n Network) MakeCompanyOwner(ctx context.Context, userID string, companyID string) error {
	passContext(&ctx)

	_, err := n.networkClient.MakeCompanyOwner(ctx, &networkRPC.UserCompanyId{
		UserId:    userID,
		CompanyId: companyID,
	})

	err = handleError(err)
	if err != nil {
		return err
	}

	return nil
}

// MakeCompanyAdmin ...
func (n Network) MakeCompanyAdmin(ctx context.Context, userID string, companyID string, adminLevel account.AdminLevel) error {
	passContext(&ctx)

	_, err := n.networkClient.MakeCompanyAdmin(ctx, &networkRPC.MakeCompanyAdminRequest{
		UserId:    userID,
		CompanyId: companyID,
		Level:     accountAdminLevelToRPC(adminLevel),
	})

	err = handleError(err)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyFollowersNumber ...
func (n Network) GetCompanyFollowersNumber(ctx context.Context, companyID string) (uint32, error) {
	passContext(&ctx)

	resp, err := n.networkClient.GetNumberOfFollowersForCompany(ctx, &networkRPC.Company{
		Id: companyID,
	})
	err = handleError(err)
	if err != nil {
		return 0, err
	}

	return uint32(resp.GetValue()), nil
}

// GetAdminLevel ...
func (n Network) GetAdminLevel(ctx context.Context, companyID string) (account.AdminLevel, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return account.AdminLevelUnknown, nil
	}

	resp, err := n.networkClient.GetAdminObject(ctx, &networkRPC.Company{
		Id: companyID,
	})

	// TODO: refactor
	if err != nil && err.Error() == `rpc error: code = Unknown desc = {"type":"NotFountError","description":"User is not admin of the company","data":null}` {
		return account.AdminLevelUnknown, nil
	}

	err = handleError(err)
	if err != nil {
		return account.AdminLevelUnknown, err
	}

	return adminLevelRPCToAccountAdminLevel(resp.GetLevel()), nil
}

// AddCompanyAdmin ...
func (n Network) AddCompanyAdmin(ctx context.Context, companyID string, userID string, level account.AdminLevel) error {
	passContext(&ctx)

	_, err := n.networkClient.MakeCompanyAdmin(ctx, &networkRPC.MakeCompanyAdminRequest{
		CompanyId: companyID,
		UserId:    userID,
		Level:     accountAdminLevelToRPC(level),
	})
	err = handleError(err)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCompanyAdmin ...
func (n Network) DeleteCompanyAdmin(ctx context.Context, companyID string, userID string) error {
	passContext(&ctx)

	_, err := n.networkClient.DeleteCompanyAdmin(ctx, &networkRPC.UserCompanyId{
		CompanyId: companyID,
		UserId:    userID,
	})
	err = handleError(err)
	if err != nil {
		return err
	}
	return nil
}

// GetCompanyCountings ...
func (n Network) GetCompanyCountings(ctx context.Context, companyID string) (followings, followers, employees int32, err error) {
	counts, err := n.networkClient.GetCompanyCountings(ctx, &networkRPC.Company{Id: companyID})
	if err != nil {
		return
	}

	return counts.GetFollowings(), counts.GetFollowers(), counts.GetEmployees(), nil
}

// IsFollow ...
func (n Network) IsFollow(ctx context.Context, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsFollowingCompany(ctx, &networkRPC.Company{
		Id: companyID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// IsFollowForCompany ...
func (n Network) IsFollowForCompany(ctx context.Context, userID, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsFollowingCompanyForCompany(ctx, &networkRPC.UserCompanyId{
		CompanyId: companyID,
		UserId:    userID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// IsFavourite ...
func (n Network) IsFavourite(ctx context.Context, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsFavouriteCompany(ctx, &networkRPC.Company{
		Id: companyID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// IsBlockedCompany ...
func (n Network) IsBlockedCompany(ctx context.Context, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsBlockedCompany(ctx, &networkRPC.Company{
		Id: companyID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// IsBlockedCompanyForCompany ...
func (n Network) IsBlockedCompanyForCompany(ctx context.Context, userID, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsBlockedCompanyForCompany(ctx, &networkRPC.UserCompanyId{
		UserId:    userID,
		CompanyId: companyID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// IsBlockedCompanyByUser ...
func (n Network) IsBlockedCompanyByUser(ctx context.Context, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsBlockedCompanyByUser(ctx, &networkRPC.Company{
		Id: companyID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

// IsBlockedCompanyByCompany ...
func (n Network) IsBlockedCompanyByCompany(ctx context.Context, userID, companyID string) (bool, error) {
	passContext(&ctx)

	token := retriveToken(ctx)
	if token == "" {
		return false, nil
	}

	res, err := n.networkClient.IsBlockedCompanyByCompany(ctx, &networkRPC.UserCompanyId{
		CompanyId: companyID,
		UserId:    userID,
	})
	if err != nil {
		return false, err
	}

	return res.GetValue(), nil
}

func accountAdminLevelToRPC(adminLevel account.AdminLevel) networkRPC.AdminLevel {
	al := networkRPC.AdminLevel_Admin

	switch adminLevel {
	case account.AdminLevelAdmin:
		al = networkRPC.AdminLevel_Admin
	case account.AdminLevelJob:
		al = networkRPC.AdminLevel_JobAdmin
	case account.AdminLevelVShop:
		al = networkRPC.AdminLevel_VShopAdmin
	case account.AdminLevelVService:
		al = networkRPC.AdminLevel_VServiceAdmin
	case account.AdminLevelCommercial:
		al = networkRPC.AdminLevel_CommercialAdmin
	}

	return al
}

func adminLevelRPCToAccountAdminLevel(adminLevel networkRPC.AdminLevel) account.AdminLevel {
	al := account.AdminLevelUnknown

	switch adminLevel {
	case networkRPC.AdminLevel_Admin:
		al = account.AdminLevelAdmin
	case networkRPC.AdminLevel_JobAdmin:
		al = account.AdminLevelJob
	case networkRPC.AdminLevel_VShopAdmin:
		al = account.AdminLevelVShop
	case networkRPC.AdminLevel_VServiceAdmin:
		al = account.AdminLevelVService
	case networkRPC.AdminLevel_CommercialAdmin:
		al = account.AdminLevelCommercial
	}

	return al
}
