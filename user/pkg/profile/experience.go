package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Experience represent user's experience
type Experience struct {
	ID            bson.ObjectId     `bson:"id"`
	Position      string            `bson:"position"`
	Company       string            `bson:"company"`
	CityID        *uint32           `bson:"city_id"`
	StartDate     time.Time         `bson:"start_date"`
	FinishDate    *time.Time        `bson:"finish_date"`
	CurrentlyWork bool              `bson:"currently_work"`
	Description   *string           `bson:"description"`
	Links         []*Link           `bson:"links"`
	Files         []*File           `bson:"files"`
	Location      location.Location `bson:"-"`

	Translations map[string]*ExperienceTranslation `bson:"translations"`
}

// ExperienceTranslation ...
type ExperienceTranslation struct {
	Position    string  `bson:"position"`
	Company     string  `bson:"company"`
	Description *string `bson:"description"`
}

// GetID returns id of experience
func (e Experience) GetID() string {
	return e.ID.Hex()
}

// SetID saves id of experience. If id has a wrong format returns usersErrors.ErrWrongID error.
func (e *Experience) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		e.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for experience
func (e *Experience) GenerateID() string {
	id := bson.NewObjectId()
	e.ID = id
	return id.Hex()
}

// Translate ...
func (e *Experience) Translate(lang string) {
	if e == nil || lang == "" {
		return
	}

	if tr, isExists := e.Translations[lang]; isExists {
		e.Company = tr.Company
		e.Description = tr.Description
		e.Position = tr.Position
	}
}
