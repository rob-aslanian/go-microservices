package feedback

import (
	"github.com/globalsign/mgo/bson"
)

// File ...
type File struct {
	ID       bson.ObjectId 		`bson:"id"`
	Name     string             `bson:"name"`
	MimeType string             `bson:"mime"`
	URL      string             `bson:"url"`
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
	return nil
}

// GenerateID creates new random id for file
func (f *File) GenerateID() string {
	id := bson.NewObjectId()
	f.ID = id
	return id.Hex()
}
