package service

import (
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/tracer"
)

// Service ...
type Service struct {
	auth   AuthentificationRepository
	tokens TokensRepository
	ipDB   IPDBRepository
	tracer *tracing.Tracer
}

// Settings options for service
type Settings struct {
	Auth   AuthentificationRepository
	IPDB   IPDBRepository
	Tracer *tracing.Tracer
	Tokens TokensRepository
}

// NewService ...
func NewService(settings *Settings) *Service {
	serv := Service{
		ipDB:   settings.IPDB,
		auth:   settings.Auth,
		tracer: settings.Tracer,
		tokens: settings.Tokens,
	}

	return &serv
}
