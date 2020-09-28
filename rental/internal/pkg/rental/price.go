package rental

// Price ...
type Price struct {
	PriceType PriceType `bson:"price_type"`
	MinPrice  int32     `bson:"min_price"`
	MaxPrice  int32     `bson:"max_price"`
	FixPrice  int32     `bson:"fix_price"`
	Currency  string    `bson:"currency"`
}
