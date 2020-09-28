package statisticsrepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository represents jobs repository
type Repository struct {
	db                          *mongo.Database
	usersStatisticsCollection   *mongo.Collection
	companyStatisticsCollection *mongo.Collection
}

// Settings parameters of jobs repository
type Settings struct {
	DBAddresses []string
	User        string
	Password    string

	DBName                      string
	UsersStatisticsCollection   string
	CompanyStatisticsCollection string
}

// NewRepository reates a new instance of jobs repository
func NewRepository(settings *Settings) (*Repository, error) {
	repo := connect(settings)

	repo.usersStatisticsCollection = repo.db.Collection(settings.UsersStatisticsCollection)
	repo.companyStatisticsCollection = repo.db.Collection(settings.CompanyStatisticsCollection)

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
