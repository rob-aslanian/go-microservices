package rental

import "go.mongodb.org/mongo-driver/bson/primitive"

// Detail ...
type Detail struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	HouseRules  string             `bson:"house_rules,omitempty"`
	Description string             `bson:"description,omitempty"`
}

// GetID returns id
func (p Detail) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Detail) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Detail) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}
