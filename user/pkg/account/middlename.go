package account

// MiddleName ...
type MiddleName struct {
	Middlename string      `bson:"name"`
	Permission *Permission `bson:"permission"`
}
