package account

import (
	"github.com/globalsign/mgo/bson"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Phone ...
type Phone struct {
	ID bson.ObjectId `bson:"id,omitempty"`
	// CountryAbbreviation string        `bson:"country_abbr"`
	CountryCode CountryCode `bson:"country_code"`
	Number      string      `bson:"number"`
	Activated   bool        `bson:"activated"`
	Primary     bool        `bson:"primary"`
	Permission  Permission  `bson:"permission"`
}

// CountryCode represents a country code of phone number.
//
// Example:
// ID:        88,
// Code:      "995",
// CountryID: "GE",
//
// Retrive id of country codes is possible by GraphQL query getListOfCountryCodes.
type CountryCode struct {
	ID        uint32 `bson:"id"`
	Code      string `bson:"code"`
	CountryID string `bson:"country_id"`
}

// GetID returns id of phone
func (p Phone) GetID() string {
	return p.ID.Hex()
}

// SetID saves id of phone. If id has a wrong format returns usersErrors.ErrWrongID error.
func (p *Phone) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		p.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for phone
func (p *Phone) GenerateID() string {
	id := bson.NewObjectId()
	p.ID = id
	return id.Hex()
}
