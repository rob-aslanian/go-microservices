package offer

import (
	"time"

	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"
	serviceorder "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/service-order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Proposal is each Offer your recieved or gave out
type Proposal struct {
	ID             primitive.ObjectID      `bson:"_id"`
	OwnerID        primitive.ObjectID      `bson:"owner_id"`
	IsOwnerCompany bool                    `bson:"is_owner_company"`
	RequestID      primitive.ObjectID      `bson:"request_id"`
	ProposalDetail ProposalDetail          `bson:"proposal_detail"`
	Service        servicerequest.Service  `bson:"-"`
	Request        *servicerequest.Request `bson:"-"`
}

// ProposalDetail ...
type ProposalDetail struct {
	OfficeID       primitive.ObjectID          `bson:"office_id"`
	ServiceID      primitive.ObjectID          `bson:"service_id"`
	ProfileID      primitive.ObjectID          `bson:"profile_id"`
	IsCompany      bool                        `bson:"is_company"`
	Status         serviceorder.OrderStatus    `bson:"status"`
	Message        string                      `bson:"message"`
	PriceType      servicerequest.Price        `bson:"price_type"`
	PriceAmount    int32                       `bson:"price_amount"`
	MinPriceAmount int32                       `bson:"min_price"`
	MaxPriceAmount int32                       `bson:"max_price"`
	Currency       string                      `bson:"currency"`
	DeliveryTime   servicerequest.DeliveryTime `bson:"delivery_time"`
	CustomDate     string                      `bson:"custom_date"`
	ExperationTime int32                       `bson:"experation_time"`
	CreatedAt      time.Time                   `bson:"created_at"`
}

// GetProposal ...
type GetProposal struct {
	ProposalAmount int32      `bson:"proposal_amount"`
	Proposals      []Proposal `bson:"proposals"`
}

// GetID returns id
func (p Proposal) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Proposal) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Proposal) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetOwnerID returns owner id
func (p Proposal) GetOwnerID() string {
	return p.OwnerID.Hex()
}

// SetOwnerID set owner id
func (p *Proposal) SetOwnerID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OwnerID = objID
	return nil
}

// GetRequestID returns user id
func (p Proposal) GetRequestID() string {
	return p.RequestID.Hex()
}

// SetRequestID set user id
func (p *Proposal) SetRequestID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.RequestID = objID
	return nil
}

// GetProfileID returns profile id
func (p ProposalDetail) GetProfileID() string {
	return p.ProfileID.Hex()
}

// SetProfileID set profile id
func (p *ProposalDetail) SetProfileID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ProfileID = objID
	return nil
}

// GetServiceID returns service id
func (p ProposalDetail) GetServiceID() string {
	return p.ServiceID.Hex()
}

// SetServiceID set service id
func (p *ProposalDetail) SetServiceID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ServiceID = objID
	return nil
}

// GetOfficeID returns service id
func (p ProposalDetail) GetOfficeID() string {
	return p.OfficeID.Hex()
}

// SetOfficeID set service id
func (p *ProposalDetail) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OfficeID = objID
	return nil
}
