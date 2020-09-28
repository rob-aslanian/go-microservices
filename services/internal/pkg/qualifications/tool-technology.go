package qualifications

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ToolTechnology is the list of tools and technologies in which the servicer is proficient in
//
type ToolTechnology struct {
	ID             primitive.ObjectID  `bson:"_id"`
	ToolTechnology string              `bson:"tool-technology"`
	Rank           *Level              `bson:"rank"`
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
	if p.ID.IsZero() {
		return ""
	}
	return p.ID.Hex()
}

// SetID set id
func (p *ToolTechnology) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *ToolTechnology) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}
