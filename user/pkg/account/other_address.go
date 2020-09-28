package account

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// OtherAddress ...
type OtherAddress struct {
	ID        bson.ObjectId     `bson:"id,omitempty"`
	Name      string            `bson:"name"`
	Firstname string            `bson:"firstname"`
	Lastname  string            `bson:"lastname"`
	Apartment string            `bson:"apartment"`
	Street    string            `bson:"street"`
	ZIP       string            `bson:"zip"`
	Location  location.Location `bson:"location"`
}

// GetID returns id of address
func (a OtherAddress) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of address. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *OtherAddress) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates a new random id for address
func (a *OtherAddress) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}
