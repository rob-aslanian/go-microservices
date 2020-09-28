package account

// NativeName ...
type NativeName struct {
	Name       string      `bson:"name"`
	LanguageID string      `bson:"language"` // should it be just string?
	Permission *Permission `bson:"permission"`
}
