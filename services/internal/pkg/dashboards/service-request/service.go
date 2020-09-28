package servicerequest

import (
	"time"

	additionaldetails "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/additional-details"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/category"
	order "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/orders"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service ...
type Service struct {
	ID                primitive.ObjectID                   `bson:"_id"`
	UserID            primitive.ObjectID                   `bson:"user_id,omitempty"`
	CompanyID         *primitive.ObjectID                  `bson:"company_id,omitempty"`
	OfficeID          primitive.ObjectID                   `bson:"office_id"`
	Title             string                               `bson:"title"`
	Description       string                               `bson:"description"`
	Status            Status                               `bson:"status"`
	Category          *category.Category                   `bson:"category"`
	Currency          string                               `bson:"currency"`
	Price             Price                                `bson:"price"`
	DeliveryTime      DeliveryTime                         `bson:"delivery_time"`
	FixedPriceAmmount int32                                `bson:"fixed_price_amount"`
	MinPriceAmmout    int32                                `bson:"min_price_amount"`
	MaxPriceAmmout    int32                                `bson:"max_price_amount"`
	Orders            []order.Order                        `bson:"order"`
	Image             string                               `bson:"image"`
	AdditionalDetails *additionaldetails.AdditionalDetails `bson:"additional_details"`
	LocationType      location.LocationType                `bson:"location_type"`
	Location          *location.Location                   `bson:"location"`
	Files             []*file.File                         `bson:"files,omitempty"`
	Cancellations     int32                                `bson:"cancellations"`
	CreatedAt         time.Time                            `bson:"created_at"`
	IsSaved           bool                                 `bson:"is_saved"`
	IsDraft           bool                                 `bson:"is_draft"`
	IsRemote          bool                                 `bson:"is_remote"`
	IsPaused          bool                                 `bson:"is_paused"`
	Action            Action                               `bson:"action"`
	HasLiked          bool                                 `bson:"-"`
	Clicks            int32                                `bson:"-"`
	Views             int32                                `bson:"-"`
	WorkingHours      *WorkingHours                        `bson:"working_hours"`
}

// GetServices ...
type GetServices struct {
	ServiceAmount int32      `bson:"service_amount"`
	Services      []*Service `bson:"services"`
}

// ServiceStatus ...
type ServiceStatus struct {
	IsDraft  bool `bson:"is_draft"`
	IsPaused bool `bson:"is_paused"`
}

// GetID returns id
func (p Service) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Service) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Service) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns user id
func (p Service) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Service) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetCompanyID returns user id
func (p Service) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}

	return p.CompanyID.Hex()
}

// SetCompanyID set user id
func (p *Service) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetOfficeID returns user id
func (p Service) GetOfficeID() string {
	return p.OfficeID.Hex()
}

// SetOfficeID set user id
func (p *Service) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OfficeID = objID
	return nil
}
