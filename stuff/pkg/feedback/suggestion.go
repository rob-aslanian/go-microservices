package feedback

type FeedBackSuggestion struct {
	Idea     string `bson:"idea"`
	Proposal string `bson:"proposal"`
	Files    []File `bson:"files,omitempty"`
}
