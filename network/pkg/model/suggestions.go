package model

import "gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"

type UserSuggestion struct {
	User      User  `json:"user"`
	Following bool  `json:"following"`
	Followers int32 `json:"followers"`
}

func (f *UserSuggestion) ToRPC() *networkRPC.UserSuggestion {
	return &networkRPC.UserSuggestion{
		Following: f.Following,
		User:      f.User.ToRPC(),
		Followers: f.Followers,
	}
}

type UserSuggestionArr []*UserSuggestion

func (suggestions UserSuggestionArr) ToRPC() *networkRPC.UserSuggestionArr {
	networkSuggestions := make([]*networkRPC.UserSuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		networkSuggestions[i] = suggestion.ToRPC()
	}
	return &networkRPC.UserSuggestionArr{Suggestions: networkSuggestions}
}

type CompanySuggestion struct {
	Company   Company `json:"company"`
	Followers int32   `json:"followers"`
}

func (f *CompanySuggestion) ToRPC() *networkRPC.CompanySuggestion {
	return &networkRPC.CompanySuggestion{
		Company:   f.Company.ToRPC(),
		Followers: f.Followers,
	}
}

type CompanySuggestionArr []*CompanySuggestion

func (suggestions CompanySuggestionArr) ToRPC() *networkRPC.CompanySuggestionArr {
	networkSuggestions := make([]*networkRPC.CompanySuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		networkSuggestions[i] = suggestion.ToRPC()
	}
	return &networkRPC.CompanySuggestionArr{Suggestions: networkSuggestions}
}
