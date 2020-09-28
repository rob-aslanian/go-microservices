package invitation

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Invitation ...
type Invitation struct {
	ID        bson.ObjectId  `bson:"_id"`
	UserID    *bson.ObjectId `bson:"user_id"`
	CompanyID *bson.ObjectId `bson:"company_id"`
	Name      string         `bson:"name"`
	Email     string         `bson:"email"`
	CratedAt  time.Time      `bson:"created_at"`
}

// GetID returns id of user
func (a *Invitation) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of user. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Invitation) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for account
func (a *Invitation) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetUserID returns id of user
func (a *Invitation) GetUserID() string {
	if a.UserID != nil {
		return a.UserID.Hex()
	}
	return ""
}

// SetUserID saves id of user. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Invitation) SetUserID(id string) error {
	if bson.IsObjectIdHex(id) {
		id := bson.ObjectIdHex(id)
		a.UserID = &id
		return nil
	}
	return usersErrors.ErrWrongID
}

// GetCompanyID returns id of user
func (a *Invitation) GetCompanyID() string {
	if a.CompanyID != nil {
		return a.CompanyID.Hex()
	}
	return ""
}

// SetCompanyID saves id of user. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Invitation) SetCompanyID(id string) error {
	if bson.IsObjectIdHex(id) {
		id := bson.ObjectIdHex(id)
		a.CompanyID = &id
		return nil
	}
	return usersErrors.ErrWrongID
}
