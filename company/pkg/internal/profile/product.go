package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Product ...
type Product struct {
	ID        bson.ObjectId `bson:"id"`
	Name      string        `bson:"name"`
	Image     string        `bson:"image"`
	Website   string        `bson:"website"`
	CreatedAt time.Time     `bson:"created_at"`
}

// GetID returns id of product
func (a Product) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of product. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Product) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for product
func (a *Product) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}
