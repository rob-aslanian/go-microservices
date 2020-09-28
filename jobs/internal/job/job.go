package job

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/company"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Posting ...
type Posting struct {
	ID                  primitive.ObjectID `bson:"_id"`
	UserID              primitive.ObjectID `bson:"user_id"`
	CompanyID           primitive.ObjectID `bson:"company_id"`
	CompanyDetails      *company.Details   `bson:"company_details"`
	JobDetails          Details            `bson:"job_details"`
	JobMetadata         Meta               `bson:"job_metadata"`
	NormalizedSalaryMin float32            `bson:"normalized_salary_min"`
	NormalizedSalaryMax float32            `bson:"normalized_salary_max"`
	JobPriority         int                `bson:"priority"`

	CreatedAt      time.Time      `bson:"created_at"`
	ActivationDate time.Time      `bson:"activation_date"`
	ExpirationDate time.Time      `bson:"expiration_date"`
	LastPauseDate  time.Time      `bson:"last_pause_date"`
	PausedDays     int            `bson:"paused_days"`
	Status         Status         `bson:"status"` // TODO: should be enum
	Applications   []*Application `bson:"applications,omitempty"`
}

// GetID returns id
func (p Posting) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Posting) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Posting) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// GetUserID returns user id
func (p Posting) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Posting) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetCompanyID returns user id
func (p Posting) GetCompanyID() string {
	return p.CompanyID.Hex()
}

// SetCompanyID set user id
func (p *Posting) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = objID
	return nil
}
