package file

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File ...
type File struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    primitive.ObjectID  `bson:"user_id"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Name      string              `bson:"name"`
	MimeType  string              `bson:"mime"`
	URL       string              `bson:"url"`
	Position  uint32              `bson:"position"`
}

// GetID returns id of file
func (f File) GetID() string {
	return f.ID.Hex()
}

// SetID ...
func (f File) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f.ID = objID
	return nil
}

// GenerateID creates new file
func (f *File) GenerateID() string {
	f.ID = primitive.NewObjectID()
	return f.ID.Hex()
}

// GetUserID returns id of user
func (f File) GetUserID() string {
	return f.UserID.Hex()
}

// SetUserID ...
func (f *File) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f.UserID = objID
	return nil
}

// GetCompanyID returns id of company
func (f File) GetCompanyID() string {
	if f.CompanyID == nil {
		return ""
	}

	return f.CompanyID.Hex()
}

// SetCompanyID ...
func (f *File) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	f.CompanyID = &objID
	return nil
}
