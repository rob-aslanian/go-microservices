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
	db                      *mongo.Database
	officeCollection        *mongo.Collection
	servicesCollection      *mongo.Collection
	requestCollection       *mongo.Collection
	savedServicesCollection *mongo.Collection
	invitationCollection    *mongo.Collection
	offerCollection         *mongo.Collection
	orderCollection         *mongo.Collection
	reviewCollections       *mongo.Collection
}

// Settings parameters of services repository
type Settings struct {
	DBAddresses []string
	User        string
	Password    string

	DBName                      string
	ServicesCollectionName      string
	RequestCollectionName       string
	OfficeCollectionName        string
	InvitationCollectionName    string
	OfferCollectionName         string
	OrderCollectionName         string
	ReviewCollectionName        string
	SavedServicesCollectionName string
}

// NewRepository reates a new instance of services repository
func NewRepository(settings *Settings) (*Repository, error) {
	repo := connect(settings)

	repo.officeCollection = repo.db.Collection(settings.OfficeCollectionName)
	repo.servicesCollection = repo.db.Collection(settings.ServicesCollectionName)
	repo.requestCollection = repo.db.Collection(settings.RequestCollectionName)
	repo.invitationCollection = repo.db.Collection(settings.InvitationCollectionName)
	repo.offerCollection = repo.db.Collection(settings.OfferCollectionName)
	repo.orderCollection = repo.db.Collection(settings.OrderCollectionName)
	repo.reviewCollections = repo.db.Collection(settings.ReviewCollectionName)
	repo.savedServicesCollection = repo.db.Collection(settings.SavedServicesCollectionName)

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
