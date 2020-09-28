package profile

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
)

// Service ...
type Service struct {
	ID        bson.ObjectId `bson:"id"`
	Image     string        `bson:"image"`
	Name      string        `bson:"name"`
	Website   string        `bson:"website"`
	CreatedAt time.Time     `bson:"created_at"`
}

// GetID returns id of service
func (a Service) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of service. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Service) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for service
func (a *Service) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}
