package companiesRepository

import (
	"log"

	"github.com/globalsign/mgo"
	// "github.com/mongodb/mongo-go-driver/mongo"
)

// Repository represents storage for adverts
type Repository struct {
	collections map[string]*mgo.Collection // mgo driver
	// collections map[string]*mongo.Collection // mongo driver
}

// Settings for repository
type Settings struct {
	Addresses []string
	User      string
	Password  string
	Database  string
}

// NewRepository creates new users repository
func NewRepository(settings Settings) (Repository, error) {
	r := Repository{}
	err := r.connect(settings)
	if err != nil {
		log.Fatalln("error while creating repository:", err)
	}

	return r, nil
}
