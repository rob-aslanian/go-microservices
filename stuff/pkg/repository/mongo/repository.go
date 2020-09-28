package repository

import (
	"log"

	"github.com/globalsign/mgo"
)

// Repository represents storage for adverts
type Repository struct {
	collections map[string]*mgo.Collection
}

// Settings for repository
type Settings struct {
	Addresses   []string
	User        string
	Password    string
	Collections []string
	Database    string
}

// NewRepository creates new repository
func NewRepository(settings Settings) (Repository, error) {
	r := Repository{}
	err := r.connect(settings)
	if err != nil {
		log.Fatalln("error while creating repository:", err)
	}

	return r, nil
}
