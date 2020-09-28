package servicesrepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository represents what repositories we have in db
type Repository struct {
	db               *mongo.Database
	housesCollection *mongo.Collection
}

// Settings parameters of services repository
type Settings struct {
	DBAddresses []string
	User        string
	Password    string

	DBName               string
	HousesCollectionName string
}

// NewRepository reates a new instance of services repository
func NewRepository(settings *Settings) (*Repository, error) {
	repo := connect(settings)

	repo.housesCollection = repo.db.Collection(settings.HousesCollectionName)

	return repo, nil
}

func connect(settings *Settings) *Repository {
	repo := Repository{}

	opts := options.Client()
	opts.SetHosts(settings.DBAddresses)
	opts.SetAuth(options.Credential{
		Username: settings.User,
		Password: settings.Password,
	})
	err := opts.Validate()
	if err != nil {
		panic(err)
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		panic(err)
	}

	// connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	ctx.Done()

	// ping
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	ctx.Done()

	repo.db = client.Database(settings.DBName)

	return &repo
}

// Close ...
func (r *Repository) Close() {
	ctx := context.Background()
	err := r.db.Client().Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}
