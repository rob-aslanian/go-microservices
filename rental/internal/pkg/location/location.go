package location

// Location represnts general location
type Location struct {
	Country *Country `bson:"country"`
	City    *City    `bson:"city"`
	Street  *string  `bson:"street,omitempty"`
	Address *string  `bson:"address,omitempty"`
}

// Country ...
type Country struct {
	ID   string `bson:"id"`
	Name string `bson:"name,omitempty"`
}

// City ...
type City struct {
	ID          string `bson:"id"`
	Name        string `bson:"name,omitempty"`
	Subdivision string `bson:"subdivision,omitempty"`
}
