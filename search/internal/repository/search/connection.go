package searchrepo

import (
	"log"
	"os"

	"github.com/olivere/elastic"
)

// Config ...
type Config struct {
	Addresses []string
}

// Repository ...
type Repository struct {
	client *elastic.Client
}

// Connect ...
func Connect(config *Config) *Repository {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Addresses...),
		elastic.SetInfoLog(log.New(os.Stdout, "info:", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "trace:", log.LstdFlags)),
	)
	if err != nil {
		panic(err)
	}

	repo := Repository{
		client: client,
	}

	return &repo
}
