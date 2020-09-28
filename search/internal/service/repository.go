package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	companyadmin "gitlab.lan/Rightnao-site/microservices/search/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

// Repository ...
type Repository interface {
	UserSearch(ctx context.Context, data *requests.UserSearch, connections, blockedIDs []string) ([]string, int64, error)
	JobSearch(ctx context.Context, data *requests.JobSearch, ids []string, blockedIDs []string) ([]string, int64, error)
	CandidateSearch(ctx context.Context, data *requests.CandidateSearch, blockedIDs []string) ([]string, []string, int64, error)
	CompanySearch(ctx context.Context, data *requests.CompanySearch, blockedIDs []string) ([]string, int64, error)
	ServiceSearch(ctx context.Context, data *requests.ServiceSearch) ([]string, int64, error)
	ServiceRequestSearch(ctx context.Context, data *requests.ServiceRequest) ([]string, int64, error)
}

// AuthRPC represents auth service
type AuthRPC interface {
	GetUserID(ctx context.Context) (string, error)
}

// NetworkRPC represents network service
type NetworkRPC interface {
	GetAdminLevel(ctx context.Context, companyID string) (companyadmin.AdminLevel, error)
	GetIDsOfFriends(ctx context.Context, userID string) ([]string, error)
	GetIDsOfFollowingCompanies(ctx context.Context) ([]string, error)
	GetBlockedIDs(ctx context.Context) ([]string, error)
}

// UserRPC represents user service
type UserRPC interface {
	GetProfilesByID(ctx context.Context, ids []string) ([]*userRPC.Profile, error)
	GetMapProfilesByID(ctx context.Context, ids []string) (map[string]interface{}, error)
}

// CompanyRPC represents company service
type CompanyRPC interface {
	GetCompanyProfiles(ctx context.Context, ids []string) ([]*companyRPC.Profile, error)
}

// JobsRPC represents jobs service
type JobsRPC interface {
	GetJob(ctx context.Context, id string) (*jobsRPC.JobViewForUser, error)
	GetCareerInterestsByIDs(ctx context.Context, companyID string, ids []string) ([]*jobsRPC.CareerInterests, error)
}

// InfoRPC ...
type InfoRPC interface {
	GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error)
}

// FilterRepository represents filters service
type FilterRepository interface {
	SaveUserSearchFilter(ctx context.Context, data *requests.UserSearchFilter) error
	SaveJobSearchFilter(ctx context.Context, data *requests.JobSearchFilter) error
	// SaveCandidateSearchFilter(ctx context.Context, data *requests.CandidateSearchFilter) error
	SaveCompanySearchFilter(ctx context.Context, data *requests.CompanySearchFilter) error
	SaveServiceSearchFilter(ctx context.Context, data *requests.ServiceSearchFilter) error
	SaveServiceRequestSearchFilter(ctx context.Context, data *requests.ServiceRequestSearchFilter) error

	// GetFiltersByType(ctx context.Context, userID string, data requests.FilterType) ([]interface{}, error)
	GetAllFilters(ctx context.Context, userID string) ([]interface{}, error)
	RemoveFilter(ctx context.Context, filterID string) error

	SaveUserSearchFilterForCompany(ctx context.Context, data *requests.UserSearchFilter) error
	SaveJobSearchFilterForCompany(ctx context.Context, data *requests.JobSearchFilter) error
	SaveCandidateSearchFilterForCompany(ctx context.Context, data *requests.CandidateSearchFilter) error
	SaveCompanySearchFilterForCompany(ctx context.Context, data *requests.CompanySearchFilter) error

	// GetFiltersByTypeForCompany(ctx context.Context, filterID string, data requests.FilterType) ([]interface{}, error)
	GetAllFiltersForCompany(ctx context.Context, companyID string) ([]interface{}, error)
	RemoveFilterForCompany(ctx context.Context, filterID, companyID string) error
}
