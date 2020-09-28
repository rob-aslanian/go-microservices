package file

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File ...
type File struct {
	ID       primitive.ObjectID `bson:"id"`
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f.ID = objID
	return nil
}

// GenerateID creates new random id for file
func (f *File) GenerateID() string {
	f.ID = primitive.NewObjectID()
	return f.ID.Hex()
}
