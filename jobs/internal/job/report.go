package job

import "go.mongodb.org/mongo-driver/bson/primitive"

// Report ...
type Report struct {
	ID         primitive.ObjectID `bson:"_id"`
	ReporterID primitive.ObjectID `bson:"reporter_id"`
	JobID      primitive.ObjectID `bson:"job_id"`
	Type       ReportType         `bson:"type"`
	Text       string             `bson:"text"`
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

// GetReporterID returns reporter id
func (p Report) GetReporterID() string {
	return p.ReporterID.Hex()
}

// SetReporterID set reporter id
func (p *Report) SetReporterID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ReporterID = objID
	return nil
}

// GetJobID returns job id
func (p Report) GetJobID() string {
	return p.JobID.Hex()
}

// SetJobID set job id
func (p *Report) SetJobID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.JobID = objID
	return nil
}

// ReportType ...
type ReportType string

const (
	// ReportTypeOther ...
	ReportTypeOther ReportType = "other"
	// ReportTypeScam ...
	ReportTypeScam ReportType = "scam"
	// ReportTypeOffensive ...
	ReportTypeOffensive ReportType = "offensive"
	// ReportTypeIncorrect ...
	ReportTypeIncorrect ReportType = "incorrect"
	// ReportTypeExpired ...
	ReportTypeExpired ReportType = "expired"
)
