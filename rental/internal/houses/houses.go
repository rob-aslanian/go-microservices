package houses

import (
	tracing "gitlab.lan/Rightnao-site/microservices/rental/internal/tracer"
)

// HouseRental ...
type HouseRental struct {
	host       string
	authRPC    AuthRPC
	repository Repository
	networkRPC NetworkRPC
	infoRPC    InfoRPC
	userRPC    UserRPC
	chatRPC    ChatRPC
	tracer     *tracing.Tracer
	mq         MQ
}

// Repository ...
type Repository struct {
	Houses HousesRepository
	Cache  CacheRepository
}

// Settings for service
type Settings struct {
	Host             string
	CacheRepository  CacheRepository
	HousesRepository HousesRepository
	AuthRPC          AuthRPC
	NetworkRPC       NetworkRPC
	InfoRPC          InfoRPC
	UserRPC          UserRPC
	ChatRPC          ChatRPC
	MQ               MQ
	Tracer           *tracing.Tracer
}

// NewHouseRental creates new service instance
func NewHouseRental(settings Settings) *HouseRental {
	return &HouseRental{
		repository: Repository{
			Houses: settings.HousesRepository,
			Cache:  settings.CacheRepository,
		},
		host:       settings.Host,
		authRPC:    settings.AuthRPC,
		networkRPC: settings.NetworkRPC,
		infoRPC:    settings.InfoRPC,
		userRPC:    settings.UserRPC,
		chatRPC:    settings.ChatRPC,
		mq:         settings.MQ,
		tracer:     settings.Tracer,
	}
}
