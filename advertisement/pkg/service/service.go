package service

import (
	tracing "gitlab.lan/Rightnao-site/microservices/advertisement/pkg/tracer"
)

// Service represents service itself
type Service struct {
	repository Repository
	authRPC    AuthRPC
	userRPC    UserRPC
	tracer     *tracing.Tracer
}

// Settings for service
type Settings struct {
	Tracer     *tracing.Tracer
	Repository Repository
	AuthRPC    AuthRPC
	UserRPC    UserRPC
}

// NewService creates new service instance
func NewService(settings Settings) (Service, error) {
	return Service{
		repository: settings.Repository,
		authRPC:    settings.AuthRPC,
		userRPC:    settings.UserRPC,
		tracer:     settings.Tracer,
	}, nil
}
