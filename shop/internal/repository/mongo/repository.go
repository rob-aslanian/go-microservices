package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository represents storage for adverts
type Repository struct {
	db                 *mongo.Database
	shopsCollection    *mongo.Collection
	productsCollection *mongo.Collection
	ordersCollection   *mongo.Collection
	wishlistCollection *mongo.Collection
}

// Settings for repository
type Settings struct {
	Addresses []string
	User      string
	Password  string

	Database           string
	ShopsCollection    string
	ProductsCollection string
	OrdersCollection   string
	WishlistCollection string
}

// NewRepository creates new repository
func NewRepository(settings Settings) (Repository, error) {
	repo := connect(settings)
	if repo == nil {
		panic("can't connect to db")
	}

	repo.shopsCollection = repo.db.Collection(settings.ShopsCollection)
	repo.productsCollection = repo.db.Collection(settings.ProductsCollection)
	repo.ordersCollection = repo.db.Collection(settings.OrdersCollection)
	repo.wishlistCollection = repo.db.Collection(settings.WishlistCollection)

	repo.createIndexes()

	return *repo, nil
}
