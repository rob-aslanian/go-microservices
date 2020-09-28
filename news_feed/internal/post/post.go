package post

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	space   = regexp.MustCompile(` +`)
	tab     = regexp.MustCompile(`\t+`)
	newLine = regexp.MustCompile(`\n+`)
)

// Post ...
type Post struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id"`
	UserID    primitive.ObjectID  `bson:"user_id" json:"user_id"`
	CompanyID *primitive.ObjectID `bson:"company_id,omitempty" json:"company_id,omitempty"`
	// NewsFeedID id of user, this post belongs to.
	NewsFeedUserID *primitive.ObjectID `bson:"newsfeed_user_id,omitempty" json:"newsfeed_user_id,omitempty"`
	// NewsFeedID id of company, this post belongs to.
	NewsFeedCompanyID      *primitive.ObjectID `bson:"newsfeed_company_id,omitempty" json:"newsfeed_company_id,omitempty"`
	Text                   string              `bson:"text" json:"text"`
	Files                  []*file.File        `bson:"files,omitempty" json:"files,omitempty"`
	CreatedAt              time.Time           `bson:"created_at" json:"created_at"`
	ChangedAt              *time.Time          `bson:"changed_at" json:"-"`
	SharedPost             *primitive.ObjectID `bson:"shared_post_id,omitempty" json:"shared_post_id,omitempty"`
	Hashtags               []string            `bson:"hashtags,omitempty" json:"hashtags,omitempty"`
	Tags                   []Tag               `bson:"tags,omitempty" json:"tags,omitempty"`
	IsPinned               bool                `bson:"is_pinned" json:"is_pinned"`
	IsCommentDisabled      bool                `bson:"comment_disabled" json:"comment_disabled"`
	IsNotificationDisabled bool                `bson:"-" json:"-"`

	Liked          *Like        `bson:"liked,omitempty" json:"-"`
	LikesAmount    *LikesAmount `bson:"likes_amount,omitempty" json:"-"`
	CommentsAmount uint32       `bson:"comments_amount,omitempty" json:"-"`
	SharesAmount   uint32       `bson:"-" json:"-"`
}

// GetID returns id of ad
func (p Post) GetID() string {
	return p.ID.Hex()
}

// SetID ...
func (p *Post) SetID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.ID = objID
	return nil
}

// GenerateID creates new id
func (p *Post) GenerateID() string {
	p.ID = primitive.NewObjectID()
	return p.ID.Hex()
}

// SetUserID ...
func (p *Post) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}

// GetUserID ...
func (p Post) GetUserID() string {
	return p.UserID.Hex()
}

// SetCompanyID ...
func (p *Post) SetCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.CompanyID = &objID
	return nil
}

// GetCompanyID ...
func (p Post) GetCompanyID() string {
	if p.CompanyID == nil {
		return ""
	}
	return p.CompanyID.Hex()
}

// SetNewsFeedCompanyID ...
func (p *Post) SetNewsFeedCompanyID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.NewsFeedCompanyID = &objID
	return nil
}

// GetNewsFeedCompanyID ...
func (p Post) GetNewsFeedCompanyID() string {
	if p.NewsFeedCompanyID == nil {
		return ""
	}
	return p.NewsFeedCompanyID.Hex()
}

// SetNewsFeedUserID ...
func (p *Post) SetNewsFeedUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.NewsFeedUserID = &objID
	return nil
}

// GetNewsFeedUserID ...
func (p Post) GetNewsFeedUserID() string {
	if p.NewsFeedUserID == nil {
		return ""
	}
	return p.NewsFeedUserID.Hex()
}

// SetSharedPost ...
func (p *Post) SetSharedPost(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.SharedPost = &objID
	return nil
}

// GetSharedPost ...
func (p Post) GetSharedPost() string {
	if p.SharedPost == nil {
		return ""
	}
	return p.SharedPost.Hex()
}

// Validate ...
func (p Post) Validate() error {
	if p.NewsFeedUserID == nil && p.NewsFeedCompanyID == nil {
		return errors.New("newsfeed_id_empty")
	}

	err := checkLengthText(p.Text, 1000)
	if err != nil {
		return errors.New("text_too_long")
	}

	return nil
}

// ValidateText ...
func (p Post) ValidateText() error {
	err := checkLengthText(p.Text, 1000)
	if err != nil {
		return errors.New("text_too_long")
	}

	return nil
}

func checkLengthText(s string, length int) error {
	if utf8.RuneCountInString(s) > length {
		return errors.New("text_too_long")
	}
	return nil
}

// Trim ...
func (p *Post) Trim() {
	p.Text = strings.TrimSpace(p.Text)
	p.Text = space.ReplaceAllString(p.Text, " ")
	p.Text = tab.ReplaceAllString(p.Text, " ")
	p.Text = newLine.ReplaceAllString(p.Text, "\n")
}

// FindHashtags find all hashtags in text and saves them.
func (p *Post) FindHashtags() {
	if p.Hashtags == nil {
		p.Hashtags = make([]string, 0)
	}

	words := strings.Fields(p.Text)

	for _, word := range words {
		if utf8.RuneCountInString(word) < 2 {
			continue
		}

		if !strings.HasPrefix(word, "#") {
			continue
		}

		if !isHashtag(word) {
			continue
		}

		p.Hashtags = append(p.Hashtags, strings.ToLower(word))
	}
}

func isHashtag(s string) bool {
	for _, letter := range s {
		if unicode.IsLetter(letter) ||
			unicode.IsDigit(letter) ||
			letter == rune('_') {
			return true
		}
	}

	return false
}

// MarshalJSON ...
func (p *Post) MarshalJSON() ([]byte, error) {
	type PostAlias Post

	c := &struct {
		ID                string  `json:"id"`
		UserID            string  `json:"user_id"`
		CompanyID         *string `json:"company_id,omitempty"`
		NewsFeedUserID    *string `json:"newsfeed_user_id,omitempty"`
		NewsFeedCompanyID *string `json:"newsfeed_company_id,omitempty"`
		SharedPost        *string `json:"shared_post_id,omitempty"`

		*PostAlias
	}{
		ID:        p.GetID(),
		UserID:    p.GetUserID(),
		PostAlias: (*PostAlias)(p),
	}

	if p.CompanyID != nil {
		companyID := p.GetCompanyID()
		c.CompanyID = &companyID
	}

	if p.NewsFeedUserID != nil {
		id := p.GetNewsFeedUserID()
		c.NewsFeedUserID = &id
	}

	if p.NewsFeedCompanyID != nil {
		id := p.GetNewsFeedCompanyID()
		c.NewsFeedCompanyID = &id
	}

	if p.SharedPost != nil {
		id := p.GetSharedPost()
		c.SharedPost = &id
	}

	m, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return m, nil
}
