package rentaltype

<<<<<<< HEAD
import "gitlab.lan/Rightnao-site/microservices/house_rental/internal/pkg/rental"

// RuralFarm ...
type RuralFarm struct {
	rental.CommonRental  `bson:",inline"`
=======
import "gitlab.lan/Rightnao-site/microservices/rental/internal/pkg/rental"

// RuralFarm ...
type RuralFarm struct {
	RentalInfo           rental.CommonRental                 `bson:",inline"`
>>>>>>> 48f087f738961844b327305e1b55394bc07c0be7
	CommercialProperties []rental.CommericalProperty         `bson:"rural_farm_property"`
	CommericalLocation   []rental.CommericalPropertyLocation `bson:"rural_farm_location"`
	Status               rental.Status                       `bson:"status"`
	TotalArea            int32                               `bson:"total_area"`
	Details              []rental.Detail                     `bson:"detail"`
	Price                rental.Price                        `bson:"price"`
	Phones               []rental.Phone                      `bson:"phones"`
}
