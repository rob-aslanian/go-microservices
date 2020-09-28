package network

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/client"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/utils"
	"google.golang.org/grpc/metadata"
)

type NetworkClient struct {
	network networkRPC.NetworkServiceClient
}

func NewAdminLevelProvider(addr string) (*NetworkClient, error) {
	netConn, err := client.CreateGrpcConnection(addr)
	if err != nil {
		return nil, err
	}

	net := networkRPC.NewNetworkServiceClient(netConn)
	return &NetworkClient{network: net}, nil
}

func (this *NetworkClient) GetAdminLevelFor(ctx context.Context, companyId string) (string, error) {
	res, err := this.network.GetAdminObject(utils.ToOutContext(ctx), &networkRPC.Company{Id: companyId})
	if err != nil {
		return "", err
	}
	return res.Level.String(), nil
}

func (n *NetworkClient) GetAllFriendshipsID(ctx context.Context) ([]string, error) {
	ids, err := n.network.GetAllFriendshipsID(utils.ToOutContext(ctx), &networkRPC.Empty{})
	if err != nil {
		return []string{}, nil
	}

	return ids.GetList(), nil
}

// GetBlockedIDs ...
func (n NetworkClient) GetBlockedIDs(ctx context.Context) ([]string, error) {
	passContext(&ctx)

	resp, err := n.network.GetblockedUsersOrCompanies(ctx, &networkRPC.Empty{})
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(resp.GetList()))

	for _, f := range resp.GetList() {
		ids = append(ids, f.GetId())
	}

	return ids, nil
}

// IsBlockedByUser checks if user with UserID is blocked by requestor.
func (n NetworkClient) IsBlockedByUser(ctx context.Context, userID string) (bool, error) {
	passContext(&ctx)

	value, err := n.network.IsBlockedByUser(ctx, &networkRPC.User{
		Id: userID,
	})

	if err != nil {
		return false, err
	}

	return value.GetValue(), nil
}

// IsBlockedCompanyByUser checks if company with companyID is blocked by requestor.
func (n NetworkClient) IsBlockedCompanyByUser(ctx context.Context, companyID string) (bool, error) {
	passContext(&ctx)

	value, err := n.network.IsBlockedCompanyByUser(ctx, &networkRPC.Company{
		Id: companyID,
	})

	if err != nil {
		return false, err
	}

	return value.GetValue(), nil
}

func passContext(ctx *context.Context) {

	md, b := metadata.FromIncomingContext(*ctx)
	if b {
		*ctx = metadata.NewOutgoingContext(*ctx, md)
	} else {
		log.Println("error while passing context")
	}
}
