package account

// Patronymic represents a patronymic
type Patronymic struct {
	Patronymic string      `bson:"name"`
	Permission *Permission `bson:"permission"`
}
