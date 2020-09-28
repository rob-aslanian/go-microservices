package invitation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Invitation is each Invitation your recieved or gave out
type Invitation struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    *primitive.ObjectID `bson:"user_id"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	OfficeID  primitive.ObjectID  `bson:"office_id"`
	RequestID primitive.ObjectID  `bson:"request_id"`
	Title     string              `bson:"title"`
	CreatedAt time.Time           `bson:"created_at"`
}

// GetID returns id
func (p Invitation) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Invitation) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Invitation) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
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

	p.UserID = &objID
	return nil
}

// GetCompanyID returns user id
func (p Invitation) GetCompanyID() string {
	return p.CompanyID.Hex()
}

// SetCompanyID set user id
func (p *Invitation) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetRequestID returns id
func (p Invitation) GetRequestID() string {
	return p.ID.Hex()
}

// SetRequestID set id
func (p *Invitation) SetRequestID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GetOfficeID returns id
func (p Invitation) GetOfficeID() string {
	return p.ID.Hex()
}

// SetOfficeID set id
func (p *Invitation) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}
