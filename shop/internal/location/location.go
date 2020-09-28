package location

// Location ...
type Location struct {
	CityID          *string `bson:"city_id,omitempty"`
	CityName        string  `bson:"-"`
	CitySubdivision string  `bson:"-"`
	CountryID       string  `bson:"country_id"`
}
