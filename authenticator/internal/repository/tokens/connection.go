package authtoken

import "github.com/go-redis/redis"

// UserSessionsRepo connecting redis through this struct
type UserSessionsRepo struct {
	Client *redis.Client
}

// Settings ...
type Settings struct {
	Address  string
	Password string
	DB       int
}

// NewRepo ...
func NewRepo(settings Settings) (*UserSessionsRepo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     settings.Address,
		Password: settings.Password,
		DB:       settings.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &UserSessionsRepo{
		Client: client,
	}, nil
}
