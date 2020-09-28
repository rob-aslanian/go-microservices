package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Report ...
type Report struct {
	ID          bson.ObjectId `bson:"id"`
	Creator     bson.ObjectId `bson:"creator"`
	Reason      ReasonType    `bson:"reason"`
	Description string        `bson:"description"`
	CreatedAt   time.Time     `bson:"created_at"`
}

// ReasonType ...
type ReasonType string

const (
	// ReasonViolationTermOfUse ...
	ReasonViolationTermOfUse ReasonType = "violation_term_of_use"
	// ReasonNotRealOrganization ...
	ReasonNotRealOrganization ReasonType = "not_real_organization"
	// ReasonMayBeHacked ...
	ReasonMayBeHacked ReasonType = "may_be_hacked"
	// ReasonPictureNotLogo ...
	ReasonPictureNotLogo ReasonType = "picture_not_logo"
	// ReasonDuplicate ...
	ReasonDuplicate ReasonType = "duplicate"
	// ReasonOther ...
	ReasonOther ReasonType = "other"
)

// GetID returns id of report
func (a Report) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of report. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Report) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for company report
func (a *Report) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetCreatorID returns id of creator
func (a Report) GetCreatorID() string {
	return a.ID.Hex()
}

// SetCreatorID saves id of creator. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Report) SetCreatorID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}
