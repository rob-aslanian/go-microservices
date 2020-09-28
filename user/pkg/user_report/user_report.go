package userReport

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Report ...
type Report struct {
	ID          bson.ObjectId `bson:"_id"`
	Type        Type          `bson:"type"`
	UserID      bson.ObjectId `bson:"user_id"`
	CreatorID   bson.ObjectId `bson:"creator_id"`
	Description string        `bson:"description"`
	CreatedAt   time.Time     `bson:"created_at"`
}

// GetID returns id of report
func (r Report) GetID() string {
	return r.ID.Hex()
}

// SetID saves id of report. If id has a wrong format returns usersErrors.ErrWrongID error.
func (r *Report) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		r.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for report
func (r *Report) GenerateID() string {
	id := bson.NewObjectId()
	r.ID = id
	return id.Hex()
}

// --------------------------

// GetUserID returns id of user
func (r Report) GetUserID() string {
	return r.UserID.Hex()
}

// SetUserID saves id of user. If id has a wrong format returns usersErrors.ErrWrongID error.
func (r *Report) SetUserID(id string) error {
	if bson.IsObjectIdHex(id) {
		r.UserID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// --------------------------

// GetCreatorID returns id of creator
func (r Report) GetCreatorID() string {
	return r.CreatorID.Hex()
}

// SetCreatorID saves id of creator. If id has a wrong format returns usersErrors.ErrWrongID error.
func (r *Report) SetCreatorID(id string) error {
	if bson.IsObjectIdHex(id) {
		r.CreatorID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// --------------------------

// Type ...
type Type string

const (
	// ReportUserRequestOther ...
	ReportUserRequestOther Type = "Other"
	// ReportUserRequestVolatationTermsOfUse ...
	ReportUserRequestVolatationTermsOfUse Type = "VolatationTermsOfUse"
	// ReportUserRequestNotRealIndividual ...
	ReportUserRequestNotRealIndividual Type = "NotRealIndividual"
	// ReportUserRequestPretendingToBeSomone ...
	ReportUserRequestPretendingToBeSomone Type = "PretendingToBeSomone"
	// ReportUserRequestMayBeHacked ...
	ReportUserRequestMayBeHacked Type = "MayBeHacked"
	// ReportUserRequestPictureIsNotPerson ...
	ReportUserRequestPictureIsNotPerson Type = "PictureIsNotPerson"
	// ReportUserRequestPictureIsOffensive ...
	ReportUserRequestPictureIsOffensive Type = "PictureIsOffensive"
)
