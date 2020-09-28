package account

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Website ...
type Website struct {
	ID   bson.ObjectId `bson:"id"`
	Site string        `bson:"website"`
}

// GetID returns id of website
func (w Website) GetID() string {
	return w.ID.Hex()
}

// SetID saves id of website. If id has a wrong format returns usersErrors.ErrWrongID error.
func (w *Website) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		w.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for website
func (w *Website) GenerateID() string {
	id := bson.NewObjectId()
	w.ID = id
	return id.Hex()
}
