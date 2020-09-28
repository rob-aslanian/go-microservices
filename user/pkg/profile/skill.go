package profile

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Skill ...
type Skill struct {
	ID       bson.ObjectId `bson:"id"`
	Skill    string        `bson:"skill"`
	Position uint32        `bson:"position"`

	Translations map[string]*SkillTranslation `bson:"translations"`
}

// SkillTranslation ...
type SkillTranslation struct {
	Skill string `bson:"skill"`
}

// GetID returns id of skill
func (s Skill) GetID() string {
	return s.ID.Hex()
}

// SetID saves id of skill. If id has a wrong format returns usersErrors.ErrWrongID error.
func (s *Skill) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		s.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for skill
func (s *Skill) GenerateID() string {
	id := bson.NewObjectId()
	s.ID = id
	return id.Hex()
}

// Translate ...
func (s *Skill) Translate(lang string) {
	if s == nil || lang == "" {
		return
	}

	if tr, isExists := s.Translations[lang]; isExists {
		s.Skill = tr.Skill
	}
}
