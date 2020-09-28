package qualifications

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Language ...
type Language struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Language  string              `bson:"language"`
	Rank      *Level              `bson:"rank"`
}

// GetID returns id
func (p Language) GetID() string {
	if p.ID.IsZero() {
		return ""
	}
	return p.ID.Hex()
}

// SetID set id
func (p *Language) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Language) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}



