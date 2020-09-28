package profile

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Award ...
type Award struct {
	ID        bson.ObjectId `bson:"id"`
	Title     string        `bson:"title"`
	Issuer    string        `bson:"issuer"`
	Date      time.Time     `bson:"date"`
	Files     []*File       `bson:"files"`
	Links     []*Link       `bson:"links"`
	CreatedAt time.Time     `bson:"created_at"`

	Translations map[string]*AwardTranslation `bson:"translations"`
}

// AwardTranslation ...
type AwardTranslation struct {
	Issuer string `bson:"issuer"`
	Title  string `bson:"title"`
}

// GetID returns id of award
func (a Award) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company award. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Award) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for award
func (a *Award) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// Translate ...
func (a *Award) Translate(ctx context.Context, lang string) string {
	if a == nil || lang == "" {
		return "en"
	}

	if tr, isExists := a.Translations[lang]; isExists {
		a.Issuer = tr.Issuer
		a.Title = tr.Title
	} else {
		return "en"
	}

	return lang
}
