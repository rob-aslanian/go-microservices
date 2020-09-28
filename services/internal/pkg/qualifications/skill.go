package qualifications

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Skill ...
type Skill struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Skill     string              `bson:"skill"`

}

// SkillTranslation ...
type SkillTranslation struct {
	Skill string `bson:"skill"`
}

// GetID returns id
func (p Skill) GetID() string {
	if p.ID.IsZero() {
		return ""
	}
	return p.ID.Hex()
}

// SetID set id
func (p *Skill) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Skill) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// Translate ...
// func (p *Skill) Translate(lang string) {
// 	if p == nil || lang == "" {
// 		return
// 	}

// 	if tr, isExists := p.Translations[lang]; isExists {
// 		p.Skill = tr.Skill
// 	}
// }
