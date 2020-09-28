package service

import "context"

// Repository ...
type Repository interface {
	SaveUserEvent(ctx context.Context, event map[string]interface{}) error
	SaveCompanyEvent(ctx context.Context, event map[string]interface{}) error
}
