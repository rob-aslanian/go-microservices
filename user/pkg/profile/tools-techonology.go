package profile

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

// ToolTechnology is the list of tools and technologies in which the servicer is proficient in
type ToolTechnology struct {
	ID             bson.ObjectId `bson:"id"`
	ToolTechnology string        `bson:"tool-technology"`
	Rank           Level         `bson:"rank"`
	CreatedAt      time.Time     `bson:"created_at"`

	Translations map[string]*ToolTechnologyTranslation `bson:"translations"`
}

// ToolTechnologyTranslation ...
type ToolTechnologyTranslation struct {
	ToolTechnology string `bson:"tool-technology"`
	Rank           Level  `bson:"rank"`
}

// Level presents the level of poficiency of each tool 7 technologie
type Level string

const (
	// LevelBeginner to choose the tools knowledge level
	LevelBeginner Level = "beginner"

	// LevelIntermediate to choose the tools knowledge level
	LevelIntermediate Level = "intermediate"

	// LevelAdvanced to choose the tools knowledge level
	LevelAdvanced Level = "advanced"

	// LevelMaster to choose the tools knowledge level
	LevelMaster Level = "master"
)

// GetID returns id
func (p ToolTechnology) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *ToolTechnology) SetID(id string) error {
	if ok := bson.IsObjectIdHex(id); !ok {
		return errors.New(`wrong_id`)
	}

	objID := bson.ObjectIdHex(id)

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *ToolTechnology) GenerateID() string {
	p.ID = bson.NewObjectId()
	return p.ID.Hex()
}
