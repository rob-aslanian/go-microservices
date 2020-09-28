package service

import (
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/tracer"
)

// Service represents service itself
type Service struct {
	repository Repository
	authRPC    AuthRPC
	networkRPC NetworkRPC
	tracer     *tracing.Tracer
}

// Repository ...
type Repository struct {
	Notifications NotificationsRepository
}

// Settings for service
type Settings struct {
	AuthRPC       AuthRPC
	NetworkRPC    NetworkRPC
	Notifications NotificationsRepository
	Tracer        *tracing.Tracer
}

// NewService creates new service instance
func NewService(settings Settings) (Service, error) {
	return Service{
		repository: Repository{
			Notifications: settings.Notifications,
		},
		authRPC:    settings.AuthRPC,
		networkRPC: settings.NetworkRPC,
		tracer:     settings.Tracer,
	}, nil
}
