package profile

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Gallery ...
type Gallery struct {
	ID       bson.ObjectId `bson:"id"`
	Files    []*File       `bson:"files"`
	Position uint32        `bson:"position"`
}

// GalleryTranslation ...
type GalleryTranslation struct {
	Files []*File `bson:"files"`
}

// GetID returns id of Gallery
func (s Gallery) GetID() string {
	return s.ID.Hex()
}

// SetID saves id of Gallery. If id has a wrong format returns companyErrors.ErrWrongID error.
func (s *Gallery) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		s.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for gallery
func (s *Gallery) GenerateID() string {
	id := bson.NewObjectId()
	s.ID = id
	return id.Hex()
}
