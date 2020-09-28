package service

import "gitlab.lan/Rightnao-site/microservices/shop/internal/tracer"

// Service represents service itself
type Service struct {
	repository Repository
	authRPC    AuthRPC
	networkRPC NetworkRPC
	infoRPC    InfoRPC
	tracer     *tracing.Tracer
	mq         MQ
}

// Settings for service
type Settings struct {
	Tracer     *tracing.Tracer
	Repository Repository
	AuthRPC    AuthRPC
	NetworkRPC NetworkRPC
	InfoRPC    InfoRPC
	MQ         MQ
}

// NewService creates new service instance
func NewService(settings Settings) (Service, error) {
	return Service{
		repository: settings.Repository,
		authRPC:    settings.AuthRPC,
		networkRPC: settings.NetworkRPC,
		infoRPC:    settings.InfoRPC,
		tracer:     settings.Tracer,
		mq:         settings.MQ,
	}, nil
}
