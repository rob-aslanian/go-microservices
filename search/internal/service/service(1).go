package service

import "gitlab.lan/Rightnao-site/microservices/search/internal/tracer"

// Config ...
type Config struct {
	Tracer           *tracing.Tracer
	Repository       Repository
	AuthRPC          AuthRPC
	CompanyRPC       CompanyRPC
	NetworkRPC       NetworkRPC
	UserRPC          UserRPC
	JobsRPC          JobsRPC
	InfoRPC          InfoRPC
	FilterRepository FilterRepository
}

// Service ...
type Service struct {
	tracer           *tracing.Tracer
	search           Repository
	authRPC          AuthRPC
	networkRPC       NetworkRPC
	userRPC          UserRPC
	jobsRPC          JobsRPC
	infoRPC          InfoRPC
	companyRPC       CompanyRPC
	filterRepository FilterRepository
}

// Create ...
func Create(config *Config) *Service {
	service := Service{
		tracer:           config.Tracer,
		search:           config.Repository,
		authRPC:          config.AuthRPC,
		networkRPC:       config.NetworkRPC,
		userRPC:          config.UserRPC,
		companyRPC:       config.CompanyRPC,
		jobsRPC:          config.JobsRPC,
		infoRPC:          config.InfoRPC,
		filterRepository: config.FilterRepository,
	}

	return &service
}
