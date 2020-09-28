package model

import "gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"

/********************** FOR USERS **********************/
type UserFilter struct {
	Query     string   `validate:"omitempty,min=3,max=30"`
	Category  string   `validate:"max=30"`
	Letter    string   `validate:"max=1"`
	SortBy    string   `validate:"omitempty,oneof=first_name last_name recently_added"`
	Companies []string `validate:"max=5,dive,len=24,alphanum"`
}

func NewUserFilterFromRPC(data *networkRPC.UserFilter) *UserFilter {
	if data == nil {
		return &UserFilter{}
	}
	return &UserFilter{
		Query:     data.Query,
		Category:  data.Category,
		Letter:    data.Letter,
		SortBy:    data.SortBy,
		Companies: data.Companies,
	}
}

type IdWithUserFilter struct {
	Id     string `validate:"len=24"`
	Filter *UserFilter
}

func NewIdWithUserFilterFromRPC(data *networkRPC.IdWithUserFilter) *IdWithUserFilter {
	return &IdWithUserFilter{
		Id:     data.Id,
		Filter: NewUserFilterFromRPC(data.Filter),
	}
}

/********************** FOR COMPANIES **********************/
type CompanyFilter struct {
	Query    string `validate:"omitempty,min=3,max=30"`
	Category string `validate:"max=30"`
	Letter   string `validate:"max=1"`
	SortBy   string `validate:"omitempty,oneof=name rating size number_of_followers recently_added"`
}

func NewCompanyFilterFromRPC(data *networkRPC.CompanyFilter) *CompanyFilter {
	if data == nil {
		return &CompanyFilter{}
	}
	return &CompanyFilter{
		Query:    data.Query,
		Category: data.Category,
		Letter:   data.Letter,
		SortBy:   data.SortBy,
	}
}

type IdWithCompanyFilter struct {
	Id     string `validate:"len=24"`
	Filter *CompanyFilter
}

func NewIdWithCompanyFilterFromRPC(data *networkRPC.IdWithCompanyFilter) *IdWithCompanyFilter {
	return &IdWithCompanyFilter{
		Id:     data.Id,
		Filter: NewCompanyFilterFromRPC(data.Filter),
	}
}
