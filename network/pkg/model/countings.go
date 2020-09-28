package model

import "gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"

type UserCountings struct {
	NumOfConnections             int32 `json:"num_of_friends"`
	NumOfFollowings              int32 `json:"num_of_followings"`
	NumOfFollowers               int32 `json:"num_of_followers"`
	NumOfReceivedRecommendations int32 `json:"num_of_received_recommendations"`
}

func (this *UserCountings) ToRPC() *networkRPC.UserCountings {
	return &networkRPC.UserCountings{
		NumOfConnections:             this.NumOfConnections,
		NumOfFollowers:               this.NumOfFollowers,
		NumOfFollowings:              this.NumOfFollowings,
		NumOfReceivedRecommendations: this.NumOfReceivedRecommendations,
	}
}

type CompanyCountings struct {
	AmountOfFollowings int32 `json:"amount_of_followings"`
	AmountOfFollowers  int32 `json:"amount_of_followers"`
	AmountOfEmployees  int32 `json:"amount_of_employees"`
}
