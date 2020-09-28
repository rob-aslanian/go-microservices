package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Education ...
type Education struct {
	ID               bson.ObjectId     `bson:"id"`
	School           string            `bson:"school"`
	Degree           *string           `bson:"degree"`
	FieldStudy       string            `bson:"field_study"`
	Grade            *string           `bson:"grade"`
	CityID           *uint32           `bson:"city_id"`
	StartDate        time.Time         `bson:"start_date"`
	FinishDate       time.Time         `bson:"finish_date"`
	IsCurrentlyStudy bool              `bson:"is_currenlty_study"`
	Description      *string           `bson:"description"`
	Files            []*File           `bson:"files"`
	Links            []*Link           `bson:"links"`
	Location         location.Location `bson:"-"`

	Translations map[string]*EducationTranslation `bson:"translations"`
}

// EducationTranslation ...
type EducationTranslation struct {
	School      string  `bson:"school"`
	Degree      *string `bson:"degree"`
	FieldStudy  string  `bson:"field_study"`
	Grade       *string `bson:"grade"`
	Description *string `bson:"description"`
}

// GetID returns id of education
func (e Education) GetID() string {
	return e.ID.Hex()
}

// SetID saves id of education. If id has a wrong format returns usersErrors.ErrWrongID error.
func (e *Education) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		e.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for education
func (e *Education) GenerateID() string {
	id := bson.NewObjectId()
	e.ID = id
	return id.Hex()
}

// Translate ...
func (e *Education) Translate(lang string) {
	if e == nil || lang == "" {
		return
	}

	if tr, isExists := e.Translations[lang]; isExists {
		e.Degree = tr.Degree
		e.Description = tr.Description
		e.FieldStudy = tr.FieldStudy
		e.Grade = tr.Grade
		e.School = tr.School
	}
}
