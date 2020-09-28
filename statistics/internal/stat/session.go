package stat

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRegistration ...
type Session struct {
	Statistic
	UserID   primitive.ObjectID `bson:"user_id"`
	Location Location           `bson:"location"`
	Device   Device             `bson:"device_type"`
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
