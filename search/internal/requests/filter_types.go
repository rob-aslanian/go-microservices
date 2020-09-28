package requests

// FilterType ...
type FilterType string

const (
	// TypeUserFilterType filter in user search
	TypeUserFilterType = "user"
	// TypeCompanyFilterType filter in company search
	TypeCompanyFilterType = "company"
	// TypeJobFilterType filter in job search
	TypeJobFilterType = "job"
	// TypeCandidateFilterType filter in candidate search
	TypeCandidateFilterType = "candidate"
	// TypeServiceFilterType filter in sevice search
	TypeServiceFilterType = "service"
	// TypeServiceRequestFilterType filter in sevice search
	TypeServiceRequestFilterType = "service_request"
)
