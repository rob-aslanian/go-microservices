package candidate

import (
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/company"
	suitable "gitlab.lan/Rightnao-site/microservices/jobs/internal/suitablefor"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CareerInterests represents careere interest
type CareerInterests struct {
	UserID           primitive.ObjectID     `bson:"user_id"`
	Jobs             []string               `bson:"jobs"`
	Industry         string                 `bson:"industry"`
	Subindustries    []string               `bson:"subindustries"`
	CompanySize      company.Size           `bson:"company_size"`
	JobTypes         []JobType              `bson:"job_types"`
	SalaryCurrency   string                 `bson:"salary_currency"`
	SalaryMin        int32                  `bson:"salary_min"`
	SalaryMax        int32                  `bson:"salary_max"`
	SalaryInterval   SalaryInterval         `bson:"salary_interval"`
	NormalizedSalary float32                `bson:"normalized_salary"`
	Relocate         bool                   `bson:"relocate"`
	Remote           bool                   `bson:"remote"`
	Travel           bool                   `bson:"travel"`
	Experience       ExperienceEnum         `bson:"experience"`
	Locations        []Location             `bson:"locations"`
	Suitable         []suitable.SuitableFor `bson:"suitable_for"`

	IsInvited bool `bson:"-"`
	IsSaved   bool `bson:"-"`
}

// GetUserID returns user id
func (p CareerInterests) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *CareerInterests) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// Location ...
type Location struct {
	CityID      int32  `bson:"city_id"`
	City        string `bson:"city"`
	Country     string `bson:"country"`
	Subdivision string `bson:"subdivision"`
}
