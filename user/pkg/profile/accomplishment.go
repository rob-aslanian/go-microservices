package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Accomplishment ...
type Accomplishment struct {
	ID            bson.ObjectId      `bson:"id"`
	Name          string             `bson:"name"`
	Type          AccomplishmentType `bson:"type"`
	Issuer        *string            `bson:"issuer"`
	LicenseNumber *string            `bson:"license_number"`
	IsExpire      *bool              `bson:"is_expire"`
	URL           *string            `bson:"url"`
	StartDate     *time.Time         `bson:"start_date"`
	FinishDate    *time.Time         `bson:"finish_date"`
	Description   *string            `bson:"description"`
	Score         *float32           `bson:"score"`
	Links         []*Link            `bson:"links"`
	Files         []*File            `bson:"files"`

	Translations map[string]*AccomplishmentTranslation `bson:"translations"`
}

// AccomplishmentTranslation ...
type AccomplishmentTranslation struct {
	Name        string  `bson:"name"`
	Issuer      *string `bson:"issuer"`
	Description *string `bson:"description"`
}

// GetID returns id of accomplishment
func (a Accomplishment) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of accomplishment. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Accomplishment) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for accomplishment
func (a *Accomplishment) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// Translate ...
func (a *Accomplishment) Translate(lang string) {
	if a == nil || lang == "" {
		return
	}

	if tr, isExists := a.Translations[lang]; isExists {
		a.Name = tr.Name
		a.Issuer = tr.Issuer
		a.Description = tr.Description
	}
}
