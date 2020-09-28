package review

import (
	"time"

	servicerequest "gitlab.lan/Rightnao-site/microservices/services/internal/pkg/dashboards/service-request"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Review ...
type Review struct {
	ID             primitive.ObjectID     `bson:"_id"`
	OwnerID        primitive.ObjectID     `bson:"owner_id"`
	IsOwnerCompnay bool                   `bson:"is_owner_company"`
	ServiceID      primitive.ObjectID     `bson:"service_id,omitempty"`
	RequestID      primitive.ObjectID     `bson:"request_id,omitempty"`
	OfficeID       *primitive.ObjectID    `bson:"office_id,omitempty"`
	ReviewDetail   ReviewDetail           `bson:"review_detail"`
	ReviewAVG      float64                `bson:"review_avg,omitempty"`
	Service        servicerequest.Service `bson:"-"`
	Request        servicerequest.Request `bson:"-"`
}

// GetReview ...
type GetReview struct {
	ReviewAmount     int32    `bson:"review_amount"`
	ClartityAVG      float64  `bson:"clartity_avg"`
	CommunicationAVG float64  `bson:"communication_avg"`
	PaymentAVG       float64  `bson:"payment_avg"`
	Reviews          []Review `bson:"reviews"`
}

// ReviewDetail ...
type ReviewDetail struct {
	ProfileID     primitive.ObjectID `bson:"profile_id"`
	IsCompany     bool               `bson:"is_company"`
	Clarity       uint32             `bson:"clarity"`
	Communication uint32             `bson:"communication"`
	Payment       uint32             `bson:"payment"`
	Hire          Hire               `bson:"hire"`
	Description   string             `bson:"description"`
	CreatedAt     time.Time          `bson:"created_at"`
}

// Hire will the user hire this servicer again in the future
type Hire string

const (
	// WillHire ...
	WillHire Hire = "will_hire"
	// NotHire ...
	NotHire Hire = "not_hire"
	// NotAnswer ...
	NotAnswer Hire = "not_answer"
)

// GetID returns id
func (p Review) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Review) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Review) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetOwnerID returns owner id
func (p Review) GetOwnerID() string {
	return p.OwnerID.Hex()
}

// SetOwnerID set owner id
func (p *Review) SetOwnerID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OwnerID = objID
	return nil
}

// GetRequestID returns request id
func (p Review) GetRequestID() string {
	return p.RequestID.Hex()
}

// SetRequestID set requst id
func (p *Review) SetRequestID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.RequestID = objID
	return nil
}

// GetServiceID returns service id
func (p Review) GetServiceID() string {
	return p.ServiceID.Hex()
}

// SetServiceID set service id
func (p *Review) SetServiceID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ServiceID = objID
	return nil
}

// GetOfficeID returns user id
func (p Review) GetOfficeID() string {
	if p.OfficeID == nil {
		return ""
	}
	return p.OfficeID.Hex()
}

// SetOfficeID set user id
func (p *Review) SetOfficeID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.OfficeID = &objID
	return nil
}

// GetProfileID returns profile id
func (p ReviewDetail) GetProfileID() string {
	return p.ProfileID.Hex()
}

// SetProfileID set profile id
func (p *ReviewDetail) SetProfileID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ProfileID = objID
	return nil
}
