package serverHTTP

import "context"

// Service represent service
type Service interface {
	ActivateCompany(ctx context.Context, companyID string, code string) error
	ActivateEmail(ctx context.Context, companyID string, code string) error
}
