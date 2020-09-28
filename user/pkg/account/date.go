package account

// Date ...
type Date struct {
	Day   int `bson:"day"`
	Month int `bson:"month"`
	Year  int `bson:"year"`
}
