package profile

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Milestone ...
type Milestone struct {
	ID          bson.ObjectId `bson:"id"`
	Image       string        `bson:"image"`
	Date        time.Time     `bson:"date"`
	Title       string        `bson:"title"`
	Description string        `bson:"description"`
	CreatedAt   time.Time     `bson:"created_at"`

	Translations map[string]*MilestoneTranslation `bson:"translations"`
}

// MilestoneTranslation ...
type MilestoneTranslation struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
}

// GetID returns id of milestone
func (m Milestone) GetID() string {
	return m.ID.Hex()
}

// SetID saves id of company milestone. If id has a wrong format returns usersErrors.ErrWrongID error.
func (m *Milestone) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		m.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for milestone
func (m *Milestone) GenerateID() string {
	id := bson.NewObjectId()
	m.ID = id
	return id.Hex()
}

// Translate ...
func (m *Milestone) Translate(ctx context.Context, lang string) string {
	if m == nil || lang == "" {
		return "en"
	}

	if tr, isExists := m.Translations[lang]; isExists {
		m.Title = tr.Title
		m.Description = tr.Description
	} else {
		return "en"
	}

	return lang
}
