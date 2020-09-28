package service

import (
	"time"

	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"

	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderType ...
type OrderType string

const (
	// OrderTypeSeller ...
	OrderTypeSeller OrderType = "seller"
	// OrderTypeBuyer ...
	OrderTypeBuyer OrderType = "buyer"
)

// Order ...
type Order struct {
	ID             primitive.ObjectID `bson:"_id"`
	OwnerID        primitive.ObjectID `bson:"owner_id"`
	IsOwnerCompany bool               `bson:"is_owner_company"`
	ServiceID      primitive.ObjectID `bson:"service_id"`
	RequestID      primitive.ObjectID `bson:"request_id,omitempty"`
	OfficeID       primitive.ObjectID `bson:"office_id,omitempty"`
	ReferalID      primitive.ObjectID `bson:"referal_id,omitempty"`

	OrderDetail OrderDetail            `bson:"order_detail"`
	OrderType   OrderType              `bson:"order_type"`
	Service     servicerequest.Service `bson:"-"`
	Request     servicerequest.Request `bson:"-"`
}

// OrderDetail ...
type OrderDetail struct {
	ProfileID      primitive.ObjectID   `bson:"profile_id"`
	IsCompany      bool                 `bson:"is_company"`
	Status         OrderStatus          `bson:"status"`
	Description    string               `bson:"description"`
	Files          []*file.File         `bson:"files,omitempty"`
	PriceType      servicerequest.Price `bson:"price_type"`
	PriceAmount    int32                `bson:"price_amount"`
	MinPriceAmount int32                `bson:"min_price"`
	MaxPriceAmount int32                `bson:"max_price"`
	Currency       string               `bson:"currency"`
	CustomDate     string               `bson:"custom_date"`

	DeliveryTime servicerequest.DeliveryTime `bson:"delivery_time"`
	Note         *string                     `bson:"note"`
	CreatedAt    time.Time                   `bson:"created_at"`
}

// GetOrder ...
type GetOrder struct {
	OrderAmount int32   `bson:"order_amount"`
	Orders      []Order `bson:"orders"`
}

// SavedServices ...
type SavedServices struct {
	ID       primitive.ObjectID   `bson:"_id"`
	Services []primitive.ObjectID `bson:"saved_services"`
}

// GetID returns  id
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

// GetOwnerID returns owner id
func (p Order) GetOwnerID() string {
	return p.OwnerID.Hex()
}

// SetOwnerID set owner id
func (p *Order) SetOwnerID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OwnerID = objID
	return nil
}

// GetOfficeID returns service id
func (p Order) GetOfficeID() string {
	return p.OfficeID.Hex()
}

// SetOfficeID set service id
func (p *Order) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OfficeID = objID
	return nil
}

// GetServiceID returns service id
func (p Order) GetServiceID() string {
	return p.ServiceID.Hex()
}

// SetServiceID set service id
func (p *Order) SetServiceID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ServiceID = objID
	return nil
}

// GetRequestID returns service id
func (p Order) GetRequestID() string {
	return p.RequestID.Hex()
}

// SetRequestID set service id
func (p *Order) SetRequestID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.RequestID = objID
	return nil
}

// GetProfileID returns profle id
func (p OrderDetail) GetProfileID() string {
	return p.ProfileID.Hex()
}

// SetProfileID set profile id
func (p *OrderDetail) SetProfileID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ProfileID = objID
	return nil
}

// GetReferalID returns referal id
func (p Order) GetReferalID() string {
	return p.ReferalID.Hex()
}

// SetReferalID set referal id
func (p *Order) SetReferalID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ReferalID = objID
	return nil
}
