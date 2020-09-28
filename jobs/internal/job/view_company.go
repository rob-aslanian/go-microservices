package job

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ViewForCompany ...
type ViewForCompany struct {
	ID                   primitive.ObjectID `bson:"_id"`
	UserID               primitive.ObjectID `bson:"user_id"`
	JobDetails           Details            `bson:"job_details"`
	Metadata             Meta               `bson:"job_metadata"`
	CreatedAt            time.Time          `bson:"created_at"`
	PublishedAt          time.Time          `bson:"activation_date"`
	ExpiredAt            time.Time          `bson:"expiration_date"`
	Status               string             `bson:"status"`
	NumberOfApplications int32              `bson:"num_of_applications"`
	NumberOfViews        int32              `bson:"num_of_views"`
}

// GetID returns id
func (p ViewForCompany) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *ViewForCompany) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *ViewForCompany) GenerateID() {
	p.ID = primitive.NewObjectID()
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
