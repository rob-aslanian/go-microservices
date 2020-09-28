package cacheRepository

import (
	"github.com/go-redis/redis"
)

// Repository ...
type Repository struct {
	client *redis.Client
}

// Settings for repository
type Settings struct {
	Address  string
	Password string
	Database int
}

// NewRepository creates new users repository
func NewRepository(settings Settings) (*Repository, error) {
	client, err := connect(settings)
	if err != nil {
		return nil, err
	}

	return &Repository{
		client: client,
	}, nil
}

func connect(settings Settings) (*redis.Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:     settings.Address,
			DB:       settings.Database,
			Password: settings.Password,
		},
	)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return client, nil
}
