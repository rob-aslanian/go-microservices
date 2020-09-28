package model

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"time"
)

type FriendshipStatus string

const FriendshipStatus_Requested = FriendshipStatus("Requested")
const FriendshipStatus_Approved = FriendshipStatus("Approved")
const FriendshipStatus_Denied = FriendshipStatus("Denied")
const FriendshipStatus_Ignored = FriendshipStatus("Ignored")

type NewFriendship struct {
	SenderId   string `json:"_from"`
	ReceiverId string `json:"_to"`

	Status      FriendshipStatus `json:"status"`
	Description string           `json:"description"`

	CreatedAt   time.Time `json:"created_at"`
	RespondedAt time.Time `json:"responded_at"`
}

type FriendshipFilter struct {
	Query     string   `validate:"omitempty,min=3,max=30"`
	Category  string   `validate:"max=50"`
	Letter    string   `validate:"max=1"`
	SortBy    string   `validate:"omitempty,oneof=first_name last_name recently_added"`
	Companies []string `validate:"max=5,dive,len=24,alphanum"`
}

func NewFriendshipFilterFromRPC(data *networkRPC.FriendshipFilter) *FriendshipFilter {
	return &FriendshipFilter{
		Query:     data.Query,
		Category:  data.Category,
		Letter:    data.Letter,
		SortBy:    data.SortBy,
		Companies: data.Companies,
	}
}

type FriendshipRequestFilter struct {
	Status string `validate:"omitempty,oneof=Requested Approved Denied Ignored"`
	Sent   bool
}

func NewFriendshipRequestFilterFromRPC(data *networkRPC.FriendRequestFilter) *FriendshipRequestFilter {
	return &FriendshipRequestFilter{
		Status: data.Status,
		Sent:   data.Sent,
	}
}

type Friendship struct {
	Id        string `json:"_key"`
	Friend    User   `json:"friend"`
	MyRequest bool   `json:"my_request"`

	Status      FriendshipStatus `json:"status"`
	Description string           `json:"description"`
	Categories  []string         `json:"categories"`
	Following   bool             `json:"following"`
	CreatedAt   time.Time        `json:"created_at"`
	RespondedAt time.Time        `json:"responded_at"`
}

func (f *Friendship) ToRPC() *networkRPC.Friendship {
	return &networkRPC.Friendship{
		Id:          f.Id,
		Status:      string(f.Status),
		MyRequest:   f.MyRequest,
		Description: f.Description,
		Categories:  f.Categories,
		Following:   f.Following,
		CreatedAt:   f.CreatedAt.Unix(),
		RespondedAt: f.RespondedAt.Unix(),
		Friend:      f.Friend.ToRPC(),
	}
}

type FriendshipArr []*Friendship

func (friendships FriendshipArr) ToRPC() *networkRPC.FriendshipArr {
	networkFriendships := make([]*networkRPC.Friendship, len(friendships))
	for i, friendship := range friendships {
		networkFriendships[i] = friendship.ToRPC()
	}
	return &networkRPC.FriendshipArr{Friendships: networkFriendships}
}
