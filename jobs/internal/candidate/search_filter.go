package candidate

import "go.mongodb.org/mongo-driver/bson/primitive"

// NamedSearchFilter ...
type NamedSearchFilter struct {
	UserID    primitive.ObjectID `bson:"user_id"`
	ID        primitive.ObjectID `bson:"_id"`
	CompanyID primitive.ObjectID `bson:"company_id"`
	Name      string             // `bson:"name"`
	Filter    *SearchFilter
}

// GetID returns id
func (p NamedSearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *NamedSearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *NamedSearchFilter) GenerateID() {
	p.ID = primitive.NewObjectID()
}

// GetUserID returns user id
func (p NamedSearchFilter) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *NamedSearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetCompanyID returns user id
func (p NamedSearchFilter) GetCompanyID() string {
	return p.CompanyID.Hex()
}

// SetCompanyID set user id
func (p *NamedSearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = objID
	return nil
}

// SearchFilter ...
type SearchFilter struct {
	Keyword []string

	Country []string
	City    []string

	CurrentCompany []string
	PastCompany    []string

	Industry    []string
	SubIndustry []string

	ExperienceLevel ExperienceEnum
	JobType         []string

	Skill    []string
	Language []string

	School       []string
	Degree       []string
	FieldOfStudy []string
	IsStudent    bool

	Currency  string
	Period    string
	MinSalary uint32
	MaxSalary uint32

	IsWillingToTravel      bool
	IsWillingToWorkRemotly bool
	IsPossibleToRelocate   bool
}
