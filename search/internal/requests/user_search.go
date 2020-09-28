package requests

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserSearch holds the field by which user will/can be searched
type UserSearch struct {
	First           uint32
	MyConnections   bool
	IsStudent       bool
	CurrentCompany  []string
	Firstname       []string
	Language        []string
	School          []string
	PastCompany     []string
	Position        []string
	IsMale          bool
	After           string
	ConnectionsOfID []string
	CountryID       []string
	FieldOfStudy    []string
	Industry        []string
	Lastname        []string
	IsFemale        bool
	Keyword         []string
	CityID          []string
	City            []City `bson:"-"`
	Degree          []string
	Nickname        []string
	Skill           []string
	Interest        []string
	MinAge          uint32
	MaxAge          uint32
	FullName        string
}

// Date ...
type Date struct {
	Day   uint32
	Month uint32
	Year  uint32
}

// UserSearchFilter holds the fields by which the UserSearchFilter will be saved
type UserSearchFilter struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Name      string              `bson:"name"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Type      FilterType          `bson:"type"`
	UserSearch
}

// GetID returns id
func (p UserSearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *UserSearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *UserSearchFilter) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns id
func (p UserSearchFilter) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return p.UserID.Hex()
}

// SetUserID set id
func (p *UserSearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns id
func (p UserSearchFilter) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}
	return p.CompanyID.Hex()
}

// SetCompanyID set id
func (p *UserSearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}
