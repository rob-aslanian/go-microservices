package account

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/status"
)

const (
	// MaximumAmountOfFollowersForDeactivation — maximum amount of followers which allow to deactivate a company
	MaximumAmountOfFollowersForDeactivation = 50
)

// Account represents company account
type Account struct {
	ID             bson.ObjectId        `bson:"_id"`
	OwnerID        bson.ObjectId        `bson:"owner_id"`
	Status         status.CompanyStatus `bson:"status"`
	Name           string               `bson:"name"`
	URL            string               `bson:"url"`
	Industry       Industry             `bson:"industry"`
	Type           Type                 `bson:"type"`
	Size           Size                 `bson:"size"`
	Parking        *Parking             `bson:"parking"`
	BusinessHours  []*BusinessHour      `bson:"business_hours"`
	Addresses      []*Address           `bson:"addresses"`
	FoundationDate time.Time            `bson:"foundation_date"`
	Emails         []*Email             `bson:"emails"`
	Phones         []*Phone             `bson:"phones"`
	CreatedAt      time.Time            `bson:"created_at"`
	Website        []*Website           `bson:"websites"`
	VAT            *string              `bson:"vat,omitempty"`
	InvitedBy      *bson.ObjectId 		`bson:"invited_by"`
	CompanyType    CompanyType          `bson:"company_type"`

	Email string `bson:"-"`
	Phone string `bson:"-"`
}

// GetID returns id of company account
func (a Account) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company account. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Account) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for company account
func (a *Account) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetInvitedByID ...
func (a *Account) GetInvitedByID() string {
	if a.InvitedBy == nil {
		return ""
	}
	return a.InvitedBy.Hex()
}

// SetInvitedByID ...
func (a *Account) SetInvitedByID(id string) error {
	if bson.IsObjectIdHex(id) {
		objID := bson.ObjectIdHex(id)
		a.InvitedBy = &objID
	}

	return companyErrors.ErrWrongID
}

// GetOwnerID ...
func (a Account) GetOwnerID() string {
	return a.OwnerID.Hex()
}

// SetOwnerID ...
func (a *Account) SetOwnerID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.OwnerID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}
