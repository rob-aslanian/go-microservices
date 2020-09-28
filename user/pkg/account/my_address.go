package account

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// MyAddress ...
type MyAddress struct {
	ID        bson.ObjectId     `bson:"id"`
	Name      string            `bson:"name"`
	Firstname string            `bson:"firstname"`
	Lastname  string            `bson:"lastname"`
	Apartment string            `bson:"apartment"`
	Street    string            `bson:"street"`
	ZIP       string            `bson:"zip"`
	Location  location.Location `bson:"location"`
	IsPrimary bool              `bson:"primary"`
}

// GetID returns id of address
func (a MyAddress) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of address. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *MyAddress) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates a new random id for address
func (a *MyAddress) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}
