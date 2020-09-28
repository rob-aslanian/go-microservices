package rentaltype

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/house_rental/internal/pkg/rental"
)

// Appartament ...
type Appartament struct {
	rental.CommonRental `bson:",inline"`
	TypeOfProperty      []rental.TypeOfProperty  `bson:"type_of_property"`
	Status              rental.Status            `bson:"status"`
	BadRooms            int32                    `bson:"badrooms"`
	BathRooms           int32                    `bson:"bathrooms"`
	TotalArea           int32                    `bson:"total_area"`
	Floor               int32                    `bson:"floor"`
	Floors              int32                    `bson:"floors"`
	CarSpaces           int32                    `bson:"car_spaces"`
	OutdoorFeatures     []rental.OutdoorFeatures `bson:"outdoor_features"`
	IndoorFeatures      []rental.IndoorFeatures  `bson:"indoor_features"`
	ClimatControl       []rental.ClimatControl   `bson:"climate_control"`
	AvailibatiFrom      time.Time                `bson:"avilable_from"`
	AvailibatiTo        time.Time                `bson:"avilable_to"`
	Details             []rental.Detail          `bson:"detail"`
	Price               rental.Price             `bson:"price"`
	Phones              []rental.Phone           `bson:"phones"`
	IsAgent             bool                     `bson:"is_agent"`
	HasRepossesed       bool                     `bson:"has_repossesed"`
}
