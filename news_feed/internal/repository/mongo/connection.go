package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connect(settings Settings) *Repository {
	repo := Repository{}

	opts := options.Client()
	opts.SetHosts(settings.Addresses)
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

	repo.db = client.Database(settings.Database)

	return &repo
}

func (r *Repository) createIndexes() {
	indexes := r.postsCollection.Indexes()

	_, err := indexes.CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"user_id": 1,
			},
		},
		{
			Keys: bson.M{
				"newsfeed_user_id": 1,
			},
		},
		{
			Keys: bson.M{
				"company_id": 1,
			},
		},
		{
			Keys: bson.M{
				"newsfeed_company_id": 1,
			},
		},
		{
			Keys: bson.M{
				"hashtags": 1,
			},
		},
		{
			Keys: bson.M{
				"text": "text",
			},
			Options: options.Index().SetTextVersion(3),
		},
	})
	if err != nil {
		log.Println(err)
	}
}

// Close ...
func (r *Repository) Close() {
	ctx := context.Background()
	err := r.db.Client().Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}
