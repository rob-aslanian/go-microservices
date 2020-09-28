package post

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmojiType ...
type EmojiType string

const (
	// EmojiTypeLike ...
	EmojiTypeLike EmojiType = "like"
	// EmojiTypeHeart ...
	EmojiTypeHeart EmojiType = "heart"
	// EmojiTypeStop ...
	EmojiTypeStop EmojiType = "stop"
	// EmojiTypeHmm ...
	EmojiTypeHmm EmojiType = "hmm"
	// EmojiTypeClap ...
	EmojiTypeClap EmojiType = "clap"
	// EmojiTypeRocket ...
	EmojiTypeRocket EmojiType = "rocket"
	// EmojiTypeShit ...
	EmojiTypeShit EmojiType = "shit"
)

// LikesAmount ...
type LikesAmount struct {
	Like   uint32 `bson:"like"`
	Heart  uint32 `bson:"heart"`
	Stop   uint32 `bson:"stop"`
	Hmm    uint32 `bson:"hmm"`
	Clap   uint32 `bson:"clap"`
	Rocket uint32 `bson:"rocket"`
	Shit   uint32 `bson:"shit"`
}

// Like ...
type Like struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Type      string             `bson:"type" json:"type"`
	Emoji     EmojiType          `bson:"emoji" json:"emoji"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	// for subscription
	PostID    string  `bson:"-" json:"post_id"`
	CommentID *string `bson:"-" json:"comment_id"`
}

// GetID returns id of ad
func (l Like) GetID() string {
	return l.ID.Hex()
}

// SetID ...
func (l *Like) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	l.ID = objID
	return nil
}

// MarshalJSON ...
func (l *Like) MarshalJSON() ([]byte, error) {
	type LikeAlias Like

	c := &struct {
		ID string `json:"id"`

		*LikeAlias
	}{
		ID:        l.GetID(),
		LikeAlias: (*LikeAlias)(l),
	}

	m, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return m, nil
}
