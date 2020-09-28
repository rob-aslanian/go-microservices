package requests

import (
	"gitlab.lan/Rightnao-site/microservices/search/internal/company"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CompanySearch holds the field by which company will/can be searched
type CompanySearch struct {
	Keyword                []string
	CityID                 []string
	City                   []City `bson:"-"`
	Industry               []string
	Type                   company.Type
	IsCompany              bool
	SubIndustry            []string
	Size                   company.Size
	IsJobOffers            bool
	Rating                 []string
	Name                   []string
	IsOrganization         bool
	Country                []string
	FoundersName           []string
	FoundersID             []string
	BusinessHours          []string
	IsCareerCenterOpenened bool
	First                  uint32
	After                  string
}

// CompanySearchFilter holds the fields by which the CompanySearchFilter will be saved
type CompanySearchFilter struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Type      FilterType          `bson:"type"`
	Name      string              `bson:"name"`
	CompanySearch
}

// GetID returns id
func (p CompanySearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *CompanySearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *CompanySearchFilter) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns id
func (p CompanySearchFilter) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return p.UserID.Hex()
}

// SetUserID set id
func (p *CompanySearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns id
func (p CompanySearchFilter) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}
	return p.CompanyID.Hex()
}

// SetCompanyID set id
func (p *CompanySearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}
