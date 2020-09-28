package rental

import "go.mongodb.org/mongo-driver/bson/primitive"

// Phone ...
type Phone struct {
	ID          primitive.ObjectID `bson:"_id"`
	CountryCode int32              `bson:"country_code"`
	Number      string             `bson:"number"`
}

// GetID returns id
func (p Phone) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Phone) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Phone) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}
