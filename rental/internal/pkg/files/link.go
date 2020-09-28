package file

import "go.mongodb.org/mongo-driver/bson/primitive"

// Link ...
type Link struct {
	ID  primitive.ObjectID `bson:"id"`
	URL string             `bson:"url"`
}

// GetID returns id of file
func (f Link) GetID() string {
	return f.ID.Hex()
}

// SetID saves id of Link. If id has a wrong format returns usersErrors.ErrWrongID error.
func (f *Link) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f.ID = objID
	return nil
}

// GenerateID creates new random id for Link
func (f *Link) GenerateID() string {
	f.ID = primitive.NewObjectID()
	return f.ID.Hex()
}
