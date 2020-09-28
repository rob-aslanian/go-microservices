package post

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
)

// CommentSort ...
type CommentSort uint8

const (
	// CommentSortRecent sort by creation time
	CommentSortRecent CommentSort = iota
	// CommentSortTop sort by amount of likes
	CommentSortTop
)

// Comments ...
type Comments struct {
	Comments []Comment `bson:"comments"`
	Amount   uint32    `bson:"amount"`
}

// Comment ...
type Comment struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id"`
	PostID    primitive.ObjectID  `bson:"-" json:"post_id"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"user_id"`
	CompanyID *primitive.ObjectID `bson:"company_id" json:"company_id"`
	Text      string              `bson:"text" json:"text"`
	Tags      []Tag               `bson:"tags,omitempty" json:"tags,omitempty"`
	ParentID  *primitive.ObjectID `bson:"parent_id" json:"parent_id"`
	CreatedAt time.Time           `bson:"created_at" json:"created_at"`
	ChangedAt *time.Time          `bson:"changed_at" json:"-"`
	Files     []*file.File        `bson:"files,omitempty" json:"-"`

	Liked           *Like        `bson:"liked,omitempty" json:"-"`
	LikesAmount     *LikesAmount `bson:"likes_amount,omitempty" json:"-"`
	AmountOfLikes   uint32       `bson:"-" json:"-"`
	AmountOfReplies uint32       `bson:"amount_replies,omitempty" json:"-"`
}

// GetID returns id of ad
func (c Comment) GetID() string {
	return c.ID.Hex()
}

// SetID ...
func (c *Comment) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	c.ID = objID
	return nil
}

// GetPostID returns id of ad
func (c Comment) GetPostID() string {
	return c.PostID.Hex()
}

// SetPostID ...
func (c *Comment) SetPostID(id string) error {
	objPostID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	c.PostID = objPostID
	return nil
}

// SetUserID ...
func (c *Comment) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	c.UserID = objID
	return nil
}

// GetUserID ...
func (c Comment) GetUserID() string {
	return c.UserID.Hex()
}

// GenerateID creates new id
func (c *Comment) GenerateID() string {
	c.ID = primitive.NewObjectID()
	return c.ID.Hex()
}

// GetCompanyID returns id of ad
func (c Comment) GetCompanyID() string {
	if c.CompanyID == nil {
		return ""
	}
	return c.CompanyID.Hex()
}

// SetCompanyID ...
func (c *Comment) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	c.CompanyID = &objID
	return nil
}

// GetParentID returns id of ad
func (c Comment) GetParentID() string {
	if c.ParentID == nil {
		return ""
	}
	return c.ParentID.Hex()
}

// SetParentID ...
func (c *Comment) SetParentID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	c.ParentID = &objID
	return nil
}

// Validate ...
func (c Comment) Validate() error {
	if utf8.RuneCountInString(c.Text) > 300 {
		return errors.New("too_long")
	}
	return nil
}

// Trim ...
func (c *Comment) Trim() {
	c.Text = strings.TrimSpace(c.Text)
	c.Text = space.ReplaceAllString(c.Text, " ")
	c.Text = tab.ReplaceAllString(c.Text, " ")
	c.Text = newLine.ReplaceAllString(c.Text, "\n")
}

// MarshalJSON ...
func (c *Comment) MarshalJSON() ([]byte, error) {
	type CommentAlias Comment

	com := &struct {
		ID        string  `json:"id"`
		PostID    string  `bson:"-" json:"post_id"`
		UserID    string  `json:"user_id"`
		CompanyID *string `json:"company_id,omitempty"`
		ParentID  *string `json:"parent_id"`
		// NewsFeedUserID    *string `json:"newsfeed_user_id,omitempty"`
		// NewsFeedCompanyID *string `json:"newsfeed_company_id,omitempty"`
		// SharedPost        *string `json:"shared_post_id,omitempty"`

		*CommentAlias
	}{
		ID:           c.GetID(),
		PostID:       c.GetPostID(),
		UserID:       c.GetUserID(),
		CommentAlias: (*CommentAlias)(c),
	}

	if c.CompanyID != nil {
		companyID := c.GetCompanyID()
		com.CompanyID = &companyID
	}

	if c.ParentID != nil {
		parentID := c.GetParentID()
		com.ParentID = &parentID
	}

	m, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}

	return m, nil
}
