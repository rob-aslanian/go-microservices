package service

import (
	tracing "gitlab.lan/Rightnao-site/microservices/services/internal/tracer"
)

// Service ...
type Service struct {
	host       string
	authRPC    AuthRPC
	repository Repository
	networkRPC NetworkRPC
	jobsRPC    JobsRPC
	mailRPC    MailRPC
	infoRPC    InfoRPC
	userRPC    UserRPC
	chatRPC    ChatRPC
	tracer     *tracing.Tracer
	mq         MQ
}

// Repository ...
type Repository struct {
	Services ServicesRepository
	Reviews  ReviewsRepository
	Cache    CacheRepository
	Request  RequestRepository
}

// Settings for service
type Settings struct {
	Host               string
	CacheRepository    CacheRepository
	ReviewsRepository  ReviewsRepository
	ServicesRepository ServicesRepository
	AuthRPC            AuthRPC
	NetworkRPC         NetworkRPC
	MailRPC            MailRPC
	InfoRPC            InfoRPC
	JobsRPC            JobsRPC
	UserRPC            UserRPC
	ChatRPC            ChatRPC
	MQ                 MQ
	Tracer             *tracing.Tracer
}

// NewService creates new service instance
func NewService(settings Settings) *Service {
	return &Service{
		repository: Repository{
			Services: settings.ServicesRepository,
			Reviews:  settings.ReviewsRepository,
			Cache:    settings.CacheRepository,
		},
		host:       settings.Host,
		authRPC:    settings.AuthRPC,
		networkRPC: settings.NetworkRPC,
		mailRPC:    settings.MailRPC,
		infoRPC:    settings.InfoRPC,
		userRPC:    settings.UserRPC,
		jobsRPC:    settings.JobsRPC,
		chatRPC:    settings.ChatRPC,
		mq:         settings.MQ,
		tracer:     settings.Tracer,
	}
}
