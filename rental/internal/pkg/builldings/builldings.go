package rentaltype

import "gitlab.lan/Rightnao-site/microservices/house_rental/internal/pkg/rental"

// BuilldingAndStorateRooms ...
type BuilldingAndStorateRooms struct {
	rental.CommonRental `bson:",inline"`
	Status              *rental.Status  `bson:"status,omitempty"`
	TotalArea           int32           `bson:"total_area"`
	Details             []rental.Detail `bson:"detail"`
	Price               rental.Price    `bson:"price"`
	Phones              []rental.Phone  `bson:"phones"`
}
