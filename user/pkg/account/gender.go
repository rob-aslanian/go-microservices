package account

// Gender ...
type Gender struct {
	Gender     string     `bson:"gender"`
	Permission Permission `bson:"permission"`
}
