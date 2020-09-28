package profile

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Interest ...
type Interest struct {
	ID          bson.ObjectId `bson:"id"`
	Interest    string        `bson:"interest"`
	Image       *File         `bson:"image"`
	Description *string       `bson:"description"`

	Translations map[string]*InterestTranslation `bson:"translations"`
}

// InterestTranslation ...
type InterestTranslation struct {
	Interest    string  `bson:"interest"`
	Description *string `bson:"description"`
}

// GetID returns id of interest
func (i Interest) GetID() string {
	return i.ID.Hex()
}

// SetID saves id of interest. If id has a wrong format returns usersErrors.ErrWrongID error.
func (i *Interest) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		i.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for interest
func (i *Interest) GenerateID() string {
	id := bson.NewObjectId()
	i.ID = id
	return id.Hex()
}

// Translate ...
func (i *Interest) Translate(lang string) {
	if i == nil || lang == "" {
		return
	}

	if tr, isExists := i.Translations[lang]; isExists {
		i.Interest = tr.Interest
		i.Description = tr.Description
	}
}
