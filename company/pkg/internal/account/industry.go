package account

// Industry ...
type Industry struct {
	Main string   `bson:"main"`
	Sub  []string `bson:"sub"`
}
