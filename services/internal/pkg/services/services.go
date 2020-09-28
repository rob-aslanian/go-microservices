package service

import (
	additionaldetails "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/additional-details"
	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	file "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/files"
	"gitlab.lan/Rightnao-site/microservices/services/internal/pkg/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service is for v-office  ...
type Service struct {
	ID                primitive.ObjectID                   `bson:"_id"`
	UserID            primitive.ObjectID                   `bson:"user_id"`
	CompanyID         *primitive.ObjectID                  `bson:"company_id,omitempty"`
	OfficeID          primitive.ObjectID                   `bson:"office_id"`
	Tittle            string                               `bson:"tittle"`
	Description       string                               `bson:"description"`
	Category          string                               `bson:"category"`
	DeliveryTime      servicerequest.DeliveryTime          `bson:"delivery_time"`
	Price             servicerequest.Price                 `bson:"price"`
	Currency          string                               `bson:"currency"`
	FixedPriceAmmount int32                                `bson:"fixed_price_amount"`
	MinPriceAmmout    int32                                `bson:"min_price_amount"`
	MaxPriceAmmout    int32                                `bson:"max_price_amount"`
	AdditionalDetails *additionaldetails.AdditionalDetails `bson:"additional_details"`
	Location          location.Location                    `bson:"location"`
	IsDraft           bool                                 `bson:"is_draft"`
	IsRemote          bool                                 `bson:"is_remote"`
	IsPaused          bool                                 `bson:"is_paused"`
	Files             []*file.File                         `bson:"files"`
}

// GetID get service id
func (s Service) GetID() string {
	return s.ID.Hex()
}
