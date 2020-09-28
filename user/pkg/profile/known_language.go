package profile

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// KnownLanguage ...
type KnownLanguage struct {
	ID       bson.ObjectId `bson:"id"`
	Language string        `bson:"language"`
	Rank     uint32        `bson:"rank"`
}

// GetID returns id of language
func (l KnownLanguage) GetID() string {
	return l.ID.Hex()
}

// SetID saves id of language. If id has a wrong format returns usersErrors.ErrWrongID error.
func (l *KnownLanguage) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		l.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for language
func (l *KnownLanguage) GenerateID() string {
	id := bson.NewObjectId()
	l.ID = id
	return id.Hex()
}
