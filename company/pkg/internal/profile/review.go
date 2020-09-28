package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Review ...
type Review struct {
	ID          bson.ObjectId `bson:"_id"`
	AuthorID    bson.ObjectId `bson:"author_id"`
	CompanyID   bson.ObjectId `bson:"company_id"`
	Rate        uint8         `bson:"rate"`
	Headline    string        `bson:"headline"`
	Description string        `bson:"description"`
	CreatedAt   time.Time     `bson:"created_at"`

	Company Profile `bson:"-"`
}

// GetID returns id of review
func (a Review) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company review. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Review) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for review
func (a *Review) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetAuthorID returns id of review's author
func (a Review) GetAuthorID() string {
	return a.AuthorID.Hex()
}

// SetAuthorID saves id of review's author. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Review) SetAuthorID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.AuthorID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GetCompanyID returns id of review's company
func (a Review) GetCompanyID() string {
	return a.CompanyID.Hex()
}

// SetCompanyID saves id of review's company. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Review) SetCompanyID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.CompanyID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}
