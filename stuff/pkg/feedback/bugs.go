package feedback

type FeedBackBugs struct {
	Description string `bson:"description"`
	Files       []File `bson:"files,omitempty"`
}
