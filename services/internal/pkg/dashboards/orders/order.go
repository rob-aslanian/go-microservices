package order

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order is each order your recieved or gave out
type Order struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    primitive.ObjectID  `bson:"user_id"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	OfficeID  primitive.ObjectID  `bson:"office_id"`
	ServiceID primitive.ObjectID  `bson:"service_id"`
	Status    Status              `bson:"status"`
	Note      string              `bson:"note"`
	Action    Action              `bson:"action"`
	CreatedAt time.Time           `bson:"created_at"`
}

// GetID returns id
func (p Order) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Order) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Order) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns user id
func (p Order) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Order) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetCompanyID returns user id
func (p Order) GetCompanyID() string {
	return p.CompanyID.Hex()
}

// SetCompanyID set user id
func (p *Order) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetServiceID returns user id
func (p Order) GetServiceID() string {
	return p.CompanyID.Hex()
}

// SetServiceID set user id
func (p *Order) SetServiceID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetOfficeID returns user id
func (p Order) GetOfficeID() string {
	return p.OfficeID.Hex()
}

// SetOfficeID set user id
func (p *Order) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}
