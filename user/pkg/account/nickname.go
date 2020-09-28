package account

// Nickname ...
type Nickname struct {
	Nickname   string      `bson:"name"`
	Permission *Permission `bson:"permission"`
}
