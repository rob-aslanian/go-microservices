package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository represents storage for adverts
type Repository struct {
	db               *mongo.Database
	groupsCollection *mongo.Collection
}

// Settings for repository
type Settings struct {
	Addresses []string
	User      string
	Password  string

	Database         string
	GroupsCollection string
}

// NewRepository creates new repository
func NewRepository(settings Settings) (Repository, error) {
	repo := connect(settings)
	if repo == nil {
		panic("can't connect to db")
	}

	repo.groupsCollection = repo.db.Collection(settings.GroupsCollection)

	repo.createIndexes()

	return *repo, nil
}
