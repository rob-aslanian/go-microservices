package job

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToolTechnology struct {
	ID             primitive.ObjectID `bson:"_id"`
	ToolTechnology string             `bson:"tool"`
	Rank           *Level             `bson:"rank"`
}

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
