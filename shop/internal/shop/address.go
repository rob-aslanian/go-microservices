package shop

// Address ...
type Address struct {
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	CityID    uint   `bson:"city_id"`
	ZIPCode   string `bson:"zip_code"`
	Phone     string `bson:"phone"`
	Address   string `bson:"address"`
	Apartment string `bson:"apartment"`
	Comments  string `bson:"comments"`
}
