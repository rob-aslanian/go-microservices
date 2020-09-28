package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// ReviewReport ...
type ReviewReport struct {
	ID          bson.ObjectId      `bson:"id"`
	CompanyID   bson.ObjectId      `bson:"company_id"`
	ReviewID    bson.ObjectId      `bson:"review_id"`
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

// GetID returns id of review
func (a ReviewReport) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company review. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for review
func (a *ReviewReport) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetReviewID returns id of review's author
func (a ReviewReport) GetReviewID() string {
	return a.ReviewID.Hex()
}

// SetReviewID saves id of review's author. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetReviewID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ReviewID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GetCompanyID returns id of review's company
func (a ReviewReport) GetCompanyID() string {
	return a.CompanyID.Hex()
}

// SetCompanyID saves id of review's company. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetCompanyID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.CompanyID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GetAuthorID returns id of review's company
func (a ReviewReport) GetAuthorID() string {
	return a.CompanyID.Hex()
}

// SetAuthorID saves id of review's company. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *ReviewReport) SetAuthorID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.CompanyID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}
