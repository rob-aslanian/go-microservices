package account

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/company-errors"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/location"
)

// Address ...
type Address struct {
	ID            bson.ObjectId     `bson:"id"`
	Name          string            `bson:"name"`
	Location      location.Location `bson:"location"`
	ZIPCode       string            `bson:"zip_code"`
	Apartment     string            `bson:"apartment"`
	Street        string            `bson:"street"`
	Phones        []*Phone          `bson:"phones"`
	GeoPos        GeoPos            `bson:"geopos"`
	IsPrimary     bool              `bson:"is_primary"`
	BusinessHours []*BusinessHour   `bson:"business_hours"`
	Websites      []*Website        `bson:"websites"`
}

// GeoPos ...
type GeoPos struct {
	Lantitude float64 `bson:"lantitude"`
	Longitude float64 `bson:"longitude"`
}

// BusinessHour ...
type BusinessHour struct {
	ID         bson.ObjectId `bson:"id"`
	Weekdays   []string      `bson:"weekdays"`
	StartHour  string        `bson:"start_hour"`
	FinishHour string        `bson:"finish_hour"`
}

// GetID returns id of company address
func (a BusinessHour) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company address. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *BusinessHour) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for company address
func (a *BusinessHour) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GetID returns id of company address
func (a Address) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of company address. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Address) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return companyErrors.ErrWrongID
}

// GenerateID creates new random id for company address
func (a *Address) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}
