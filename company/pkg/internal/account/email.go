package account

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Email ...
type Email struct {
	ID        bson.ObjectId `bson:"id,omitempty"`
	Email     string        `bson:"email"`
	Activated bool          `bson:"activated"`
	Primary   bool          `bson:"primary"`
}

// GetID returns id of email
func (e Email) GetID() string {
	return e.ID.Hex()
}

// SetID saves id of email. If id has a wrong format returns usersErrors.ErrWrongID error.
func (e *Email) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		e.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for email
func (e *Email) GenerateID() string {
	id := bson.NewObjectId()
	e.ID = id
	return id.Hex()
}
