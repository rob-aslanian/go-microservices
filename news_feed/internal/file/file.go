package file

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File ...
type File struct {
	ID       primitive.ObjectID `bson:"id" json:"id"`
	Name     string             `bson:"name" json:"name"`
	MimeType string             `bson:"mime" json:"mime"`
	URL      string             `bson:"url" json:"url"`
	Position uint32             `bson:"position" json:"position"`
}

// GetID ...
func (f File) GetID() string {
	return f.ID.Hex()
}

// SetID ...
func (f *File) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f.ID = objID
	return nil
}

// GenerateID creates new id
func (f *File) GenerateID() string {
	f.ID = primitive.NewObjectID()
	return f.ID.Hex()
}

// MarshalJSON ...
func (f *File) MarshalJSON() ([]byte, error) {
	type FileAlias File

	m, err := json.Marshal(&struct {
		ID string `json:"id"`

		*FileAlias
	}{
		ID:        f.GetID(),
		FileAlias: (*FileAlias)(f),
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}
