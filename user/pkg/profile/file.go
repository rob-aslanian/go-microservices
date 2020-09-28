package profile

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// File ...
type File struct {
	ID       bson.ObjectId `bson:"id"`
	Name     string        `bson:"name"`
	MimeType string        `bson:"mime"`
	URL      string        `bson:"url"`
	Position uint32        `bson:"position"`
}

// GetID returns id of file
func (f File) GetID() string {
	return f.ID.Hex()
}

// SetID saves id of file. If id has a wrong format returns usersErrors.ErrWrongID error.
func (f *File) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		f.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for file
func (f *File) GenerateID() string {
	id := bson.NewObjectId()
	f.ID = id
	return id.Hex()
}
