package services

import (
	"errors"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/utils"
	"gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
	"golang.org/x/net/context"
)

func (s *NetworkService) authenticateUser(ctx context.Context) string {
	md := utils.ExtractMetadata(ctx, "token", "user-id")
	userId, ok := md["user-id"]
	if ok {
		return userId
	}

	token, ok := md["token"]
	if !ok {
		panic(hc_errors.NOT_AUTHENTICATED_ERROR)
	}

	senderId, err := s.auth.GetUserId(ctx, token)
	if err != nil {
		panic(err)
	}
	utils.AddToIncomingMetadata(ctx, "user-id", senderId)

	return senderId
}

func (s *NetworkService) requireAdminLevelForCompany(ctx context.Context, companyKey string, levels ...model.AdminLevel) {
	admin, err := s.GetAdminObject(ctx, companyKey)
	if err != nil {
		panic(err)
	}

	for _, l := range levels {
		if admin.Level == l {
			return
		}
	}
	panic(errors.New("You don't have access to this operation"))
}
