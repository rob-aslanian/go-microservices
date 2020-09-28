package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

// ServiceRequest ...
type ServiceRequest struct {
	First         uint32 `bson:"-"`
	After         string `bson:"-"`
	Keyword       []string
	Category      string
	CityID        []string
	CountryID     []string
	LocationType  LocationType
	ProjectType   []ProjectType
	DeliveryTime  DeliveryTime
	PriceType     Price
	FixedPrice    int32
	MinPrice      int32
	MaxPrice      int32
	CurrencyPrice string
	Skills        []string
	Languages     []string
	Tools         []string
	ServiceOwner  ServiceOwner
}

// ServiceRequestSearchFilter holds the fields by which the JobSearchFilter will be saved
type ServiceRequestSearchFilter struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Type      FilterType          `bson:"type"`
	Name      string              `bson:"name"`
	ServiceRequest
}

// GetID returns id
func (p ServiceRequestSearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *ServiceRequestSearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *ServiceRequestSearchFilter) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns id
func (p ServiceRequestSearchFilter) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return p.UserID.Hex()
}

// SetUserID set id
func (p *ServiceRequestSearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns id
func (p ServiceRequestSearchFilter) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}
	return p.CompanyID.Hex()
}

// SetCompanyID set id
func (p *ServiceRequestSearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// ProjectType ...
type ProjectType string

const (
	// ProjectTypeOneTime ...
	ProjectTypeOneTime ProjectType = "OneTime"

	// ProjectTypeOnGoing ...
	ProjectTypeOnGoing ProjectType = "OnGoing"

	// ProjectTypeNotSure ...
	ProjectTypeNotSure ProjectType = "NotSure"
)
