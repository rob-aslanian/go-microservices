package location

// Location represents a location
type Location struct {
	CountryID   string  `bson:"country_id"`
	CityID      *string `bson:"city_id,omitempty"`
	City        *string `bson:"city,omitempty"`
	Subdivision *string `bson:"subdivision,omitempty"`
}
