package candidate

import "go.mongodb.org/mongo-driver/bson/primitive"

// ViewForCompany ...
type ViewForCompany struct {
	UserID          primitive.ObjectID `bson:"user_id"`
	CareerInterests CareerInterests    `bson:"career_interests"`
	IsSaved         *bool              `bson:"is_saved,omitempty"`
}

// GetUserID returns user id
func (p ViewForCompany) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *ViewForCompany) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}
