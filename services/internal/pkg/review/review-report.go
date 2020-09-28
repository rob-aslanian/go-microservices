package review

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/services/internal/services-errors"
)

// ReviewReport ...
type ReviewReport struct {
	ID          bson.ObjectId      `bson:"id"`
	OfficeID    bson.ObjectId      `bson:"office_id"`
	ReviewID    bson.ObjectId      `bson:"review"`
	AuthorID    bson.ObjectId      `bson:"author_id"`
	Reason      ReviewReportReason `bson:"reason"`
	Explanation string             `bson:"explanation"`
	CreatedAt   time.Time          `bson:"created_at"`
}

// ReviewReportReason ...
type ReviewReportReason string

const (
	// ReviewReportReasonUnknown ...
	ReviewReportReasonUnknown ReviewReportReason = "unknown"
	// ReviewReportReasonSpam ...
	ReviewReportReasonSpam ReviewReportReason = "spam"
	// ReviewReportReasonScam ...
	ReviewReportReasonScam ReviewReportReason = "scam"
	// ReviewReportReasonOffensive ...
	ReviewReportReasonOffensive ReviewReportReason = "offensive"
	// ReviewReportReasonFake ...
	ReviewReportReasonFake ReviewReportReason = "fake"
	// ReviewReportReasonOffTopic ...
	ReviewReportReasonOffTopic ReviewReportReason = "off_topic"
	// ReviewReportReasonOther ...
	ReviewReportReasonOther ReviewReportReason = "other"
)

// GetID returns id of Review
func (a *ReviewReport) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company Review. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return serviceErrors.ErrWrongID
}

// GenerateID creates new random id for Review
func (a *ReviewReport) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetReviewID returns id of Review's author
func (a ReviewReport) GetReviewID() string {
	return a.ReviewID.Hex()
}

// SetReviewID saves id of Review's author. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetReviewID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ReviewID = bson.ObjectIdHex(id)
		return nil
	}
	return serviceErrors.ErrWrongID
}

// GetOfficeID returns id of Review's office
func (a ReviewReport) GetOfficeID() string {
	return a.OfficeID.Hex()
}

// SetOfficeID saves id of Review's office. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetOfficeID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.OfficeID = bson.ObjectIdHex(id)
		return nil
	}
	return serviceErrors.ErrWrongID
}

// GetAuthorID returns id of Review's office
func (a ReviewReport) GetAuthorID() string {
	return a.OfficeID.Hex()
}

// SetAuthorID saves id of Review's office. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetAuthorID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.OfficeID = bson.ObjectIdHex(id)
		return nil
	}
	return serviceErrors.ErrWrongID
}
