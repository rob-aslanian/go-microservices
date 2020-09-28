package repository

import (
	"log"

	"github.com/globalsign/mgo"
)

func (r *Repository) connect(settings Settings) error {
	session, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:    settings.Addresses,
			Username: settings.User,
			Password: settings.Password,
		},
	)
	if err != nil {
		log.Panic("error while connection with db:", err)
	}

	db := session.DB(settings.Database)

	collection := make(map[string]*mgo.Collection, len(settings.Collections))

	for i := range settings.Collections {
		collection[settings.Collections[i]] = db.C(settings.Collections[i])
	}

	r.collections = collection

	return nil
}
