package httpserver

import (
	"context"
)

// Service represent service
type Service interface {
	SaveUserEvent(ctx context.Context, request map[string]interface{}, typeEvent string) error
	SaveCompanyEvent(ctx context.Context, request map[string]interface{}, typeEvent string) error
}
