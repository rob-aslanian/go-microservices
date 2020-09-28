package price

// Price ...
type Price struct {
	Amount   uint32 `bson:"amount"`
	Currency string `bson:"currency"`
}
