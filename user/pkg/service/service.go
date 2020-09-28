package service

import (
	"gitlab.lan/Rightnao-site/microservices/user/pkg/tracer"
)

// Service represents service itself
type Service struct {
	Host       string
	repository Repository
	authRPC    AuthRPC
	// mailRPC    MailRPC
	infoRPC    InfoRPC
	networkRPC NetworkRPC
	companyRPC CompanyRPC
	stuffRPC   StuffRPC
	chatRPC    ChatRPC
	tracer     *tracing.Tracer
	mq         MQ
	tpl        Template
}

// Repository ...
type Repository struct {
	Users       UsersRepository
	Cache       CacheRepository
	GeoIP       GeoipRepository
	arrangoRepo ArangoRepo
}

// Settings for service
type Settings struct {
	Host            string
	UsersRepository UsersRepository
	CacheRepository CacheRepository
	GeoIPRepository GeoipRepository
	ArrangoRepo     ArangoRepo
	AuthRPC         AuthRPC
	// MailRPC         MailRPC
	InfoRPC    InfoRPC
	NetworkRPC NetworkRPC
	CompanyRPC CompanyRPC
	ChatRPC    ChatRPC
	StuffRPC   StuffRPC
	Tracer     *tracing.Tracer
	MQ         MQ
}

// NewService creates new service instance
func NewService(settings Settings) (Service, error) {
	t := Template{}
	err := t.Parse()
	if err != nil {
		return Service{}, err
	}

	return Service{
		repository: Repository{
			Users:       settings.UsersRepository,
			Cache:       settings.CacheRepository,
			GeoIP:       settings.GeoIPRepository,
			arrangoRepo: settings.ArrangoRepo,
		},
		Host:    settings.Host,
		authRPC: settings.AuthRPC,
		// mailRPC:    settings.MailRPC,
		networkRPC: settings.NetworkRPC,
		infoRPC:    settings.InfoRPC,
		stuffRPC:   settings.StuffRPC,
		companyRPC: settings.CompanyRPC,
		chatRPC:    settings.ChatRPC,
		tracer:     settings.Tracer,
		mq:         settings.MQ,
		tpl:        t,
	}, nil
}
