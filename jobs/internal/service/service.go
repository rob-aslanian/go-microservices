package service

import "gitlab.lan/Rightnao-site/microservices/jobs/internal/tracer"

// Service ...
type Service struct {
	authRPC    AuthRPC
	networkRPC NetworkRPC
	infoRPC    InfoRPC
	jobs       JobsRepository
	tracer     *tracing.Tracer
	mq         MQ
}

// Settings options for service
type Settings struct {
	AuthRPC        AuthRPC
	NetworkRPC     NetworkRPC
	InfoRPC        InfoRPC
	JobsRepository JobsRepository
	Tracer         *tracing.Tracer
	MQ             MQ
}

// NewService returns new instance of service
func NewService(settings *Settings) *Service {
	serv := Service{
		authRPC:    settings.AuthRPC,
		infoRPC:    settings.InfoRPC,
		networkRPC: settings.NetworkRPC,
		jobs:       settings.JobsRepository,
		tracer:     settings.Tracer,
		mq:         settings.MQ,
	}

	return &serv
}
