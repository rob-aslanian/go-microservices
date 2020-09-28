package model

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"time"
)

type FollowRequest struct {
	FollowerId  string `json:"_from"`
	FollowingId string `json:"_to"`

	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	Id        string `json:"_key"`
	User      User   `json:"user"`
	Followers int    `json:"followers"`
	Following bool   `json:"following"`
	IsFriend  bool   `json:"is_friend"`

	CreatedAt time.Time `json:"created_at"`
}

func (f *Follow) ToRPC() *networkRPC.FollowInfo {
	return &networkRPC.FollowInfo{
		User:      f.User.ToRPC(),
		Followers: int32(f.Followers),
		Following: f.Following,
		IsFriend:  f.IsFriend,
		CreatedAt: f.CreatedAt.Unix(),
	}
}

type FollowArr []*Follow

func (follows FollowArr) ToRPC() *networkRPC.FollowsArr {
	networkFollows := make([]*networkRPC.FollowInfo, len(follows))
	for i, follow := range follows {
		networkFollows[i] = follow.ToRPC()
	}
	return &networkRPC.FollowsArr{Follows: networkFollows}
}
