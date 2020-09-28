package candidate

import "go.mongodb.org/mongo-driver/bson/primitive"

// Report ...
type Report struct {
	ID                primitive.ObjectID `bson:"_id"`
	ReporterUserID    primitive.ObjectID `bson:"reporter_user_id"`
	ReporterCompanyID primitive.ObjectID `bson:"reporter_company_id"`
	CandidateID       primitive.ObjectID `bson:"candidate_id"`
	Text              string             `bson:"text"`
}

// GetID returns id
func (p Report) GetID() string {
	return p.ID.Hex()
}

// SetID set id
func (p *Report) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Report) GenerateID() {
	p.ID = primitive.NewObjectID()
}

// GetReporterUserID returns reporter id
func (p Report) GetReporterUserID() string {
	return p.ReporterUserID.Hex()
}

// SetReporterUserID set reporter id
func (p *Report) SetReporterUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ReporterUserID = objID
	return nil
}

// GetReporterCompanyID returns reporter id
func (p Report) GetReporterCompanyID() string {
	return p.ReporterCompanyID.Hex()
}

// SetReporterCompanyID set reporter id
func (p *Report) SetReporterCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ReporterCompanyID = objID
	return nil
}

// GetCandidateID returns job id
func (p Report) GetCandidateID() string {
	return p.CandidateID.Hex()
}

// SetCandidateID set job id
func (p *Report) SetCandidateID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CandidateID = objID
	return nil
}
