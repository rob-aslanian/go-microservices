package stat

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRegistration ...
type UserRegistration struct {
	Statistic
	UserID         primitive.ObjectID `bson:"user_id"`
	Gender         string             `bson:"gender"`
	Location       Location           `bson:"location"`
	Birthday       time.Time          `bson:"birthday"`
	ActivationTime *time.Time         `bson:"created_at"`
}

// GetUserID ...
func (s UserRegistration) GetUserID() string {
	return s.UserID.Hex()
}

// SetUserID ...
func (s UserRegistration) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	s.UserID = objID

	return nil
}
