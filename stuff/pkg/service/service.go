package service

// Service represents service itself
type Service struct {
	repository Repository
	authRPC    AuthRPC
	userRPC    UserRPC 
}

// Settings for service
type Settings struct {
	Repository Repository
	AuthRPC    AuthRPC
	UserRPC    UserRPC
}

// NewService creates new service instance
func NewService(settings Settings) (Service, error) {
	return Service{
		repository: settings.Repository,
		authRPC:    settings.AuthRPC,
		userRPC: 	settings.UserRPC,
	}, nil
}
