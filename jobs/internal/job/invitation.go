package job

import "go.mongodb.org/mongo-driver/bson/primitive"

// Invitation ...
type Invitation struct {
	UserID primitive.ObjectID `bson:"user_id"`
	Text   string             `bson:"text"`
}

// GetUserID returns user id
func (p Invitation) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Invitation) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}
