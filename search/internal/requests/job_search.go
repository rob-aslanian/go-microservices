package requests

import (
	"gitlab.lan/Rightnao-site/microservices/search/internal/company"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JobSearch holds the field by which job will/can be searched
type JobSearch struct {
	IsMaxSalaryNull    bool
	First              uint32
	JobType            []string
	Industry           []string
	CompanySize        company.Size
	WithSalary         bool
	After              string
	Subindustry        []string
	CompanyName        []string
	ExperienceLevel    ExperienceEnum
	Country            []string
	CityID             []string
	City               []City `bson:"-"`
	Language           []string
	Skill              []string
	MaxSalary          uint32
	IsFollowing        bool
	WithoutCoverLetter bool
	Keyword            []string
	Degree             []string
	Currency           string
	Period             string
	MinSalary          uint32
	IsMinSalaryNull    bool
	DatePosted         DatePostedEnum
	CompanyIDs         []string
}

// City ...
type City struct {
	ID          string
	City        string
	Subdivision string
	Country     string
}

// JobSearchFilter holds the fields by which the JobSearchFilter will be saved
type JobSearchFilter struct {
	ID        primitive.ObjectID  `bson:"_id"`
	UserID    *primitive.ObjectID `bson:"user_id,omitempty"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty"`
	Type      FilterType          `bson:"type"`
	Name      string              `bson:"name"`
	JobSearch 
}

// GetID returns id
func (p JobSearchFilter) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *JobSearchFilter) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *JobSearchFilter) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns id
func (p JobSearchFilter) GetUserID() string {
	if p.UserID == nil {
		return ""
	}
	return p.UserID.Hex()
}

// SetUserID set id
func (p *JobSearchFilter) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = &objID
	return nil
}

// GetCompanyID returns id
func (p JobSearchFilter) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}
	return p.CompanyID.Hex()
}

// SetCompanyID set id
func (p *JobSearchFilter) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

type DatePostedEnum string

const (
	//DateEnumAnytime
	DateEnumAnytime DatePostedEnum = "anytime"
	//DateEnumPast24Hours...
	DateEnumPast24Hours DatePostedEnum = "past_24_hours"
	//DateEnumPastWeek ...
	DateEnumPastWeek DatePostedEnum = "past_week"
	//DateEnumPastWeek
	DateEnumPastMonth DatePostedEnum = "past_month"
)

// ExperienceEnum ...
type ExperienceEnum string

const (
	//ExperienceEnumUnknownExperience
	ExperienceEnumExpericenUnknown ExperienceEnum = "experience_unknown"
	//ExperienceEnumWithoutExperience ...
	ExperienceEnumWithoutExperience ExperienceEnum = "without_experience"
	// ExperienceEnumLessThenOneYear ...
	ExperienceEnumLessThenOneYear ExperienceEnum = "less_then_one_year"
	// ExperienceEnumOneTwoYears ...
	ExperienceEnumOneTwoYears ExperienceEnum = "one_two_years"
	// ExperienceEnumTwoThreeYears ...
	ExperienceEnumTwoThreeYears ExperienceEnum = "two_three_years"
	// ExperienceEnumThreeFiveYears ...
	ExperienceEnumThreeFiveYears ExperienceEnum = "three_five_years"
	// ExperienceEnumFiveSevenYears ...
	ExperienceEnumFiveSevenYears ExperienceEnum = "five_seven_years"
	// ExperienceEnumSevenTenYears ...
	ExperienceEnumSevenTenYears ExperienceEnum = "seven_ten_years"
	// ExperienceEnumTenYearsAndMore ...
	ExperienceEnumTenYearsAndMore ExperienceEnum = "ten_years_and_more"
)
