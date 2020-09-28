package profile

// RecommendationRequest ...
type RecommendationRequest struct {
	ID        string
	Text      string
	CreatedAt int64
	Requestor *Profile
	Requested *Profile
	Title     string
	Relation  string
}
