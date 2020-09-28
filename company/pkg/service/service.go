package service

import (
	"gitlab.lan/Rightnao-site/microservices/company/pkg/tracer"
)

// Service represents service itself
type Service struct {
	host       string
	repository Repository
	authRPC    AuthRPC
	networkRPC NetworkRPC
	jobsRPC    JobsRPC
	mailRPC    MailRPC
	infoRPC    InfoRPC
	userRPC    UserRPC
	chatRPC    ChatRPC
	stuffRPC   StuffRPC
	mq         MQ
	tracer     *tracing.Tracer
}

// Repository ...
type Repository struct {
	Company     CompanyRepository
	Reviews     ReviewsRepository
	Cache       CacheRepository
	arrangoRepo ArangoRepo
}

// Settings for service
type Settings struct {
	Host              string
	CompanyRepository CompanyRepository
	CacheRepository   CacheRepository
	ReviewsRepository ReviewsRepository
	AuthRPC           AuthRPC
	NetworkRPC        NetworkRPC
	MailRPC           MailRPC
	InfoRPC           InfoRPC
	JobsRPC           JobsRPC
	UserRPC           UserRPC
	ChatRPC           ChatRPC
	StuffRPC		  StuffRPC
	MQ                MQ
	ArangoRepo        ArangoRepo
	Tracer            *tracing.Tracer
}

// NewService creates new service instance
func NewService(settings Settings) (Service, error) {
	return Service{
		repository: Repository{
			Company:     settings.CompanyRepository,
			Reviews:     settings.ReviewsRepository,
			Cache:       settings.CacheRepository,
			arrangoRepo: settings.ArangoRepo,
		},
		host:       settings.Host,
		authRPC:    settings.AuthRPC,
		networkRPC: settings.NetworkRPC,
		mailRPC:    settings.MailRPC,
		infoRPC:    settings.InfoRPC,
		userRPC:    settings.UserRPC,
		jobsRPC:    settings.JobsRPC,
		chatRPC:    settings.ChatRPC,
		stuffRPC: 	settings.StuffRPC,
		mq:         settings.MQ,
		tracer:     settings.Tracer,
	}, nil
}
