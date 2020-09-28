package wallet


// WalletAmount ... 
type WalletAmount struct {
	GoldCoins     int32 	 `bson:"gold_coins"`
	SilverCoins   int32 	 `bson:"silver_coins"`
	PendingAmount int32  	 `bson:"pending_amount"`

}