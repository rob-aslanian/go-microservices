package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

// CandidateSearch holds the field by which candidate will/can be searched
type CandidateSearch struct {
	CurrentCompany         []string
	Currency               string
	Period                 string
	Skill                  []string
	After                  string
	Industry               []string
	Language               []string
	Degree                 []string
	MinSalary              uint32
	IsPossibleToRelocate   bool
	First                  uint32
	JobType                []string
	IsMinSalaryNull        bool
	IsMaxSalaryNull        bool
	Keyword                []string
	PastCompany            []string
	School                 []string
	ExperienceLevel        ExperienceEnum
	IsWillingToTravel      bool
	SubIndustry            []string
	IsStudent              bool
	Country                []string
	CityID                 []string
	City                   []City `bson:"-"`
	FieldOfStudy           []string
	MaxSalary              uint32
	IsWillingToWorkRemotly bool
}

// CandidateSearchFilter holds the fields by which the CandidateSearchFilter will be saved
type CandidateSearchFilter struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID primitive.ObjectID  `bson:"company_id"`
	Type      FilterType          `bson:"type"`
	Name      string              `bson:"name"`
	CandidateSearch
}

// GetID returns id
func (p CandidateSearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *CandidateSearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *CandidateSearchFilter) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns id
func (p CandidateSearchFilter) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return p.UserID.Hex()
}

// SetUserID set id
func (p *CandidateSearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns id
func (p CandidateSearchFilter) GetCompanyID() string {
	return p.CompanyID.Hex()
}

// SetCompanyID set id
func (p *CandidateSearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = objID
	return nil
}
