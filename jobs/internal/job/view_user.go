package job

import (
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/company"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ViewForUser ...
type ViewForUser struct {
	ID             primitive.ObjectID `bson:"_id"`
	CompanyDetails company.Details    `bson:"company_details"`
	JobDetails     Details            `bson:"job_details"`
	Metadata       Meta               `bson:"job_metadata"`
	Application    Application        `bson:"application"`
	IsSaved        bool               `bson:"is_saved,omitempty"`
	IsApplied      bool               `bson:"-"`
	InvitationText string             `bson:"invitation_text,omitempty"`
}

// GetID returns id
func (p ViewForUser) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *ViewForUser) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *ViewForUser) GenerateID() {
	p.ID = primitive.NewObjectID()
}
