package job

import (
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/company"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NamedSearchFilter ...
type NamedSearchFilter struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Name   string             // `bson:"name"`
	Filter *SearchFilter      // `bson:"filter"`
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

// SearchFilter ...
type SearchFilter struct {
	Keyword    []string
	DatePosted DatePostedEnum

	ExperienceLevel ExperienceEnum

	Degrees []string

	Countries []string
	Cities    []string

	JobTypes []string // TODO:

	Languages []string

	Industries    []string
	Subindustries []string
	CompanyNames  []string
	CompanySizes  company.Size // TODO:

	Currency  string
	MinSalary uint32
	MaxSalary uint32
	Period    string // TODO:

	Skills []string

	FollowingCompanies bool
	WithoutCoverLetter bool
	WithSalary         bool
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
