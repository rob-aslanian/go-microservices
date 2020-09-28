package post

// Newsfeed ...
type Newsfeed struct {
	Posts         []*Post `bson:"posts"`
	AmountOfPosts uint32  `bson:"amount"`
}
