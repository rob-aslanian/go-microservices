package profile

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Link ...
type Link struct {
	ID  bson.ObjectId `bson:"id"`
	URL string        `bson:"url"`
}

// GetID returns id of link
func (l Link) GetID() string {
	return l.ID.Hex()
}

// SetID saves id of link. If id has a wrong format returns usersErrors.ErrWrongID error.
func (l *Link) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		l.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for link
func (l *Link) GenerateID() string {
	id := bson.NewObjectId()
	l.ID = id
	return id.Hex()
}
