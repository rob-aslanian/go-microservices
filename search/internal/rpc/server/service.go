package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/search/internal/requests"
)

// Service define functions inside Service
type Service interface {
	UserSearch(ctx context.Context, data *requests.UserSearch) ([]*userRPC.Profile, int64, error)
	CompanySearch(ctx context.Context, data *requests.CompanySearch) (interface{}, int64, error)
	JobSearch(ctx context.Context, data *requests.JobSearch) (interface{}, int64, error)
	CandidateSearch(ctx context.Context, companyID string, data *requests.CandidateSearch) (interface{}, int64, error)
	ServiceSearch(ctx context.Context, data *requests.ServiceSearch) ([]string, int64, error)
	ServiceRequestSearch(ctx context.Context, data *requests.ServiceRequest) ([]string, int64, error)

	SaveUserSearchFilter(ctx context.Context, data *requests.UserSearchFilter) (string, error)
	SaveJobSearchFilter(ctx context.Context, data *requests.JobSearchFilter) (string, error)
	// SaveCandidateSearchFilter(ctx context.Context, data *requests.CandidateSearchFilter) (string, error)
	SaveCompanySearchFilter(ctx context.Context, data *requests.CompanySearchFilter) (string, error)
	SaveServiceSearchFilter(ctx context.Context, data *requests.ServiceSearchFilter) (string, error)
	SaveServiceRequestSearchFilter(ctx context.Context, data *requests.ServiceRequestSearchFilter) (string, error)


	SaveUserSearchFilterForCompany(ctx context.Context, companyID string, data *requests.UserSearchFilter) (string, error)
	SaveJobSearchFilterForCompany(ctx context.Context, companyID string, data *requests.JobSearchFilter) (string, error)
	SaveCandidateSearchFilterForCompany(ctx context.Context, companyID string, data *requests.CandidateSearchFilter) (string, error)
	SaveCompanySearchFilterForCompany(ctx context.Context, companyID string, data *requests.CompanySearchFilter) (string, error)

	// GetFiltersByType(ctx context.Context, filterType requests.FilterType) ([]interface{}, error)
	GetAllFilters(ctx context.Context) ([]interface{}, error)
	RemoveFilter(ctx context.Context, filterID string) error

	// GetFiltersByTypeForCompany(ctx context.Context, companyID string, filterType requests.FilterType) ([]interface{}, error)
	GetAllFiltersForCompany(ctx context.Context, companyID string) ([]interface{}, error)
	RemoveFilterForCompany(ctx context.Context, filterID string, companyID string) error
}
