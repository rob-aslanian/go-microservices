package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository represents storage for adverts
type Repository struct {
	db                *mongo.Database
	advertCollection  *mongo.Collection
	galleryCollection *mongo.Collection
}

// Settings for repository
type Settings struct {
	Addresses []string
	User      string
	Password  string

	Database                    string
	AdvertismentCollection      string
	GalleryCollectionCollection string
}

// NewRepository creates new repository
func NewRepository(settings Settings) (Repository, error) {
	repo := connect(settings)
	if repo == nil {
		panic("can't connect to db")
	}

	repo.advertCollection = repo.db.Collection(settings.AdvertismentCollection)
	repo.galleryCollection = repo.db.Collection(settings.GalleryCollectionCollection)

	return *repo, nil
}
