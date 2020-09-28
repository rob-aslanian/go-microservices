package model

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"time"
)

type BlockRequest struct {
	Blocker string `json:"_from"`
	Blocked string `json:"_to"`

	CreatedAt time.Time `json:"created_at"`
}

type BlockedUser struct {
	User User `json:"user"`

	BlockedAt time.Time `json:"created_at"`
}

type BlockedUserOrCompany struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	IsCompany bool   `json:"is_company"`
}

func (b *BlockedUserOrCompany) ToRPC() *networkRPC.BlockedUserOrCompany {
	return &networkRPC.BlockedUserOrCompany{
		Id:        b.Id,
		Name:      b.Name,
		Avatar:    b.Avatar,
		IsCompany: b.IsCompany,
	}
}

type BlockerUserOrCompanyArr []*BlockedUserOrCompany

func (arr BlockerUserOrCompanyArr) ToRPC() *networkRPC.BlockedUserOrCompanyArr {
	list := make([]*networkRPC.BlockedUserOrCompany, len(arr))
	for i, item := range arr {
		list[i] = item.ToRPC()
	}
	return &networkRPC.BlockedUserOrCompanyArr{List: list}
}
