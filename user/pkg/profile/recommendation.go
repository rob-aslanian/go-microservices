package profile

// Recommendation ...
type Recommendation struct {
	ID            string
	Text          string
	IsHidden      *bool
	CreatedAt     int64
	Receiver      *Profile
	Recommendator *Profile
	Title         string
	Relation      string
}
