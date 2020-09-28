package servicerequest

import (
	"time"

	additionaldetails "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/additional-details"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/category"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Request is a struct that includes the whole service request
type Request struct {
	ID                primitive.ObjectID                   `bson:"_id"`
	UserID            primitive.ObjectID                   `bson:"user_id,omitempty"`
	CompanyID         *primitive.ObjectID                  `bson:"company_id,omitempty"`
	Tittle            string                               `bson:"tittle"`
	Description       string                               `bson:"description"`
	avatar            string                               `bson:"avatar"`
	Status            Status                               `bson:"status"`
	Category          category.Category                    `bson:"category"`
	Currency          string                               `bson:"currency"`
	Price             Price                                `bson:"price"`
	DeliveryTime      DeliveryTime                         `bson:"delivery_time"`
	FixedPriceAmmount int32                                `bson:"fixed_price_amount"`
	MinPriceAmmout    int32                                `bson:"min_price_amount"`
	MaxPriceAmmout    int32                                `bson:"max_price_amount"`
	Files             []*file.File                         `bson:"files,omitempty"`
	AdditionalDetails *additionaldetails.AdditionalDetails `bson:"additional_details"`
	IsRemote          bool                                 `bson:"is_remote"`
	Location          *location.Location                   `bson:"location"`
	LocationType      location.LocationType                `bson:"location_type"`
	CustomDate        string                               `bson:"custom_date"`
	ProjectType       ProjectType                          `bson:"project_type"`
	CreatedAt         time.Time                            `bson:"created_at"`
	IsDraft           bool                                 `bson:"is_draft"`
	IsClosed          bool                                 `bson:"is_closed"`
	IsPaused          bool                                 `bson:"is_paused"`
	ProposalAmount    int32                                `bson:"-"`
	HasLiked          bool                                 `bson:"-"`
}

// GetServicesRequest ...
type GetServicesRequest struct {
	ServiceAmount int32      `bson:"service_amount"`
	Services      []*Request `bson:"services"`
}

// ServiceReqestStatus ...
type ServiceReqestStatus struct {
	IsDraft  bool `bson:"is_draft"`
	IsClosed bool `bson:"is_closed"`
	IsPaused bool `bson:"is_paused"`
}

// GetID returns id
func (p Request) GetID() string {
	return p.ID.Hex()
}

// GetStatus ...
func (p *Request) GetStatus() Status {
	if p.Status == "" {
		return StatusActive
	}

	return p.Status
}

// SetID set id
func (p *Request) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Request) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns user id
func (p Request) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Request) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetCompanyID returns user id
func (p Request) GetCompanyID() string {
	if p.CompanyID != nil {
		return p.CompanyID.Hex()
	}

	return ""

}

// SetCompanyID set user id
func (p *Request) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}
