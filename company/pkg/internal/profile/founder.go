package profile

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Founder ...
type Founder struct {
	ID        bson.ObjectId  `bson:"id"`
	UserID    *bson.ObjectId `bson:"user_id,omitempty"`
	Name      string         `bson:"name"`
	Position  string         `bson:"position"`
	Avatar    string         `bson:"avatar"`
	CreatedAt time.Time      `bson:"created_at"`
	Approved  bool           `bson:"approved"`

	Translations map[string]*FounderTranslation `bson:"translations"`
}

// FounderTranslation ...
type FounderTranslation struct {
	Position string `bson:"position"`
	Name     string `bson:"name"`
}

// GetID returns id of founder
func (a Founder) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of  founder. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Founder) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for founder
func (a *Founder) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetUserID returns id of user
func (a Founder) GetUserID() string {
	if a.UserID != nil {
		return a.UserID.Hex()
	}
	return ""
}

// SetUserID saves id of user founder. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Founder) SetUserID(id string) error {
	if bson.IsObjectIdHex(id) {
		bid := bson.ObjectIdHex(id)
		a.UserID = &bid
		return nil
	}
	return companyErrors.ErrWrongID
}

// Translation ...
func (a *Founder) Translation(ctx context.Context, lang string) string {
	if a == nil || lang == "" {
		return "en"
	}

	if tr, isExists := a.Translations[lang]; isExists {
		a.Position = tr.Position
		a.Name = tr.Name
	} else {
		return "en"
	}

	return lang
}
