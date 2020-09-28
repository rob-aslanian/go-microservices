package jobsrepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Repository represents jobs repository
type Repository struct {
	db                         *mongo.Database
	jobsCollection             *mongo.Collection
	profileCollection          *mongo.Collection
	companiesCollection        *mongo.Collection
	jobReportsCollection       *mongo.Collection
	candidateReportsCollection *mongo.Collection
	jobFiltersCollection       *mongo.Collection
	candidateFiltersCollection *mongo.Collection
	pricesCollection           *mongo.Collection
}

// Settings parameters of jobs repository
type Settings struct {
	DBAddresses []string
	User        string
	Password    string

	DBName                         string
	JobsCollectionName             string
	ProfileCollectionName          string
	CompaniesCollectionName        string
	JobReportsCollectionName       string
	CandidateReportsCollectionName string
	JobFiltersCollectionName       string
	CandidateFiltersCollectionName string
	PricesCollectionCollectionName string
}

// NewRepository reates a new instance of jobs repository
func NewRepository(settings *Settings) (*Repository, error) {
	repo := connect(settings)

	repo.jobsCollection = repo.db.Collection(settings.JobsCollectionName)
	repo.profileCollection = repo.db.Collection(settings.ProfileCollectionName)
	repo.companiesCollection = repo.db.Collection(settings.CompaniesCollectionName)
	repo.jobReportsCollection = repo.db.Collection(settings.JobReportsCollectionName)
	repo.candidateReportsCollection = repo.db.Collection(settings.CandidateReportsCollectionName)
	repo.jobFiltersCollection = repo.db.Collection(settings.JobFiltersCollectionName)
	repo.candidateFiltersCollection = repo.db.Collection(settings.CandidateFiltersCollectionName)
	repo.pricesCollection = repo.db.Collection(settings.PricesCollectionCollectionName)

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
	opts.ReadPreference = readpref.SecondaryPreferred()
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
