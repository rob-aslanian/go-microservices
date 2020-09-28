package advert

// Place define where advert will be located
type Place string

const (
	// PlaceUser on the user profile
	PlaceUser Place = "user"
	// PlaceCompany on the company profile
	PlaceCompany Place = "company"
	// PlaceLocalBusiness on ...
	PlaceLocalBusiness Place = "local_business"
	// PlaceBrands on ...
	PlaceBrands Place = "brands"
	// PlaceGroups on ...
	PlaceGroups Place = "groups"
)
