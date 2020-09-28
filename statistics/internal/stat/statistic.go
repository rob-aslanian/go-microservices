package stat

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Statistic ...
type Statistic struct {
	ID        primitive.ObjectID `bson:"_id"`
	IP        string             `bson:"ip"`
	UserAgent string             `bson:"user_agent"`
	CreatedAt time.Time          `bson:"created_at"`
}

// GenerateID ...
func (s Statistic) GenerateID() {
	s.ID = primitive.NewObjectID()
}

// GetID ...
func (s Statistic) GetID() string {
	return s.ID.Hex()
}

// SetID ...
func (s Statistic) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	s.ID = objID

	return nil
}
