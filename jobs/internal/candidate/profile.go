package candidate

import "go.mongodb.org/mongo-driver/bson/primitive"

// Profile ..
type Profile struct {
	UserID          primitive.ObjectID   `bson:"user_id"`
	IsOpen          bool                 `bson:"is_open"`
	CareerInterests *CareerInterests     `bson:"career_interests"`
	SavedJobs       []primitive.ObjectID `bson:"saved_jobs"`
	SkippedJobs     []primitive.ObjectID `bson:"skipped_jobs"`
}

// GetUserID returns user id
func (p Profile) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Profile) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}
