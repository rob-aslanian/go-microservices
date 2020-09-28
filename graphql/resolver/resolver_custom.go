package resolver

import graphql "github.com/graph-gophers/graphql-go"

type NewsfeedCustom struct {
	Posts       []NewsfeedPostCustom
	Post_amount int32
}

type NewsfeedResolverCustom struct {
	R *NewsfeedCustom
}

func (r NewsfeedResolverCustom) Posts() []NewsfeedPostResolverCustom {
	items := []NewsfeedPostResolverCustom{}
	for i := range r.R.Posts {
		items = append(items, NewsfeedPostResolverCustom{&r.R.Posts[i]})
	}
	return items
}

func (r NewsfeedResolverCustom) Post_amount() int32 {
	return r.R.Post_amount
}

// -----

type NewsfeedPostCustom struct {
	Changed_at               string        `json:"-"`
	Hashtags                 []string      `json:"hashtags,omitempty"`
	ID                       string        `json:"id"`
	Text                     string        `json:"text"`
	Files                    []File        `json:"files,omitempty"`
	Tags                     []NewsfeedTag `json:"tags,omitempty"`
	Likes_amount             LikesAmount   `json:"-"`
	Shares_amount            int32         `json:"-"`
	User_profile             *Profile
	Shared_post_id           string `json:"shared_post_id,omitempty"`
	Company_profile          *CompanyProfile
	Created_at               string `json:"created_at"`
	Is_comments_disabled     bool
	Is_notification_disabled bool
	Liked                    string
	Comments_amount          int32

	// fieds below are requried for subscriptions
	NewsFeedUserID    string `json:"newsfeed_user_id,omitempty"`
	NewsFeedCompanyID string `json:"newsfeed_company_id,omitempty"`
}

type NewsfeedPostResolverCustom struct {
	R *NewsfeedPostCustom
}

func (r NewsfeedPostResolverCustom) User_profile() *ProfileResolver {
	if r.R.User_profile == nil {
		return nil
	}
	return &ProfileResolver{R: r.R.User_profile}
}

func (r NewsfeedPostResolverCustom) Company_profile() *CompanyProfileResolver {
	if r.R.Company_profile == nil {
		return nil
	}
	return &CompanyProfileResolver{r.R.Company_profile}
}

func (r NewsfeedPostResolverCustom) Files() *[]FileResolver {
	items := []FileResolver{}
	for i := range r.R.Files {
		items = append(items, FileResolver{&r.R.Files[i]})
	}
	return &items
}

func (r NewsfeedPostResolverCustom) Created_at() string {
	return r.R.Created_at
}

func (r NewsfeedPostResolverCustom) Is_notification_disabled() bool {
	return r.R.Is_notification_disabled
}

func (r NewsfeedPostResolverCustom) Likes_amount() LikesAmountResolver {
	return LikesAmountResolver{&r.R.Likes_amount}
}

func (r NewsfeedPostResolverCustom) Shares_amount() int32 {
	return r.R.Shares_amount
}

func (r NewsfeedPostResolverCustom) ID() graphql.ID {
	id := graphql.ID(r.R.ID)
	return id
}

func (r NewsfeedPostResolverCustom) Text() string {
	return r.R.Text
}

func (r NewsfeedPostResolverCustom) Changed_at() *string {
	return &r.R.Changed_at
}

func (r NewsfeedPostResolverCustom) Hashtags() *[]string {
	items := []string{}
	for _, itm := range r.R.Hashtags {
		items = append(items, itm)
	}
	return &items
}

func (r NewsfeedPostResolverCustom) Shared_post_id() *graphql.ID {
	if r.R.Shared_post_id == "" {
		return nil
	}

	id := graphql.ID(r.R.Shared_post_id)
	return &id
}

func (r NewsfeedPostResolverCustom) Is_comments_disabled() bool {
	return r.R.Is_comments_disabled
}

func (r NewsfeedPostResolverCustom) Comments_amount() int32 {
	return r.R.Comments_amount
}

func (r NewsfeedPostResolverCustom) Tags() *[]NewsfeedTagResolver {
	items := []NewsfeedTagResolver{}
	for i := range r.R.Tags {
		items = append(items, NewsfeedTagResolver{&r.R.Tags[i]})
	}
	return &items
}

func (r NewsfeedPostResolverCustom) Liked() *string {
	if r.R.Liked == "" {
		return nil
	}

	return &r.R.Liked
}

// -----

type NewsfeedPostCommentsCustom struct {
	Comments []NewsfeedPostCommentCustom
	Amount   int32
}

type NewsfeedPostCommentsResolverCustom struct {
	R *NewsfeedPostCommentsCustom
}

func (r NewsfeedPostCommentsResolverCustom) Comments() *[]NewsfeedPostCommentResolverCustom {
	items := []NewsfeedPostCommentResolverCustom{}
	for i := range r.R.Comments {
		items = append(items, NewsfeedPostCommentResolverCustom{&r.R.Comments[i]})
	}
	return &items
}

func (r NewsfeedPostCommentsResolverCustom) Amount() int32 {
	return r.R.Amount
}

// -----

type NewsfeedPostCommentCustom struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	CompanyID       string          `json:"company_id"`
	Created_at      string          `json:"created_at"`
	Changed_at      string          `json:"-"`
	Liked           string          `json:"-"`
	Likes_amount    LikesAmount     `json:"-"`
	Replies_amount  int32           `json:"-"`
	User_profile    *Profile        `json:"-"`
	Company_profile *CompanyProfile `json:"-"`
	Text            string          `json:"text"`
	Tags            []NewsfeedTag   `json:"tags,omitempty"`
	Files           []File          `json:"-"`
	ParentID        *string         `json:"parent_id"`
	// subscription
	PostID string `json:"post_id"`
}

type NewsfeedPostCommentResolverCustom struct {
	R *NewsfeedPostCommentCustom
}

func (r NewsfeedPostCommentResolverCustom) User_profile() *ProfileResolver {
	if r.R.User_profile == nil {
		return nil
	}
	return &ProfileResolver{R: r.R.User_profile}
}

func (r NewsfeedPostCommentResolverCustom) Text() string {
	return r.R.Text
}

func (r NewsfeedPostCommentResolverCustom) Likes_amount() LikesAmountResolver {
	return LikesAmountResolver{&r.R.Likes_amount}
}

func (r NewsfeedPostCommentResolverCustom) Replies_amount() int32 {
	return r.R.Replies_amount
}

func (r NewsfeedPostCommentResolverCustom) ID() graphql.ID {
	id := graphql.ID(r.R.ID)
	return id
}

func (r NewsfeedPostCommentResolverCustom) Company_profile() *CompanyProfileResolver {
	if r.R.Company_profile == nil {
		return nil
	}

	return &CompanyProfileResolver{r.R.Company_profile}
}

func (r NewsfeedPostCommentResolverCustom) Files() *[]FileResolver {
	items := []FileResolver{}
	for i := range r.R.Files {
		items = append(items, FileResolver{&r.R.Files[i]})
	}
	return &items
}

func (r NewsfeedPostCommentResolverCustom) Created_at() string {
	return r.R.Created_at
}

func (r NewsfeedPostCommentResolverCustom) Changed_at() *string {
	return &r.R.Changed_at
}

func (r NewsfeedPostCommentResolverCustom) Parent_id() *graphql.ID {
	if r.R.ParentID == nil {
		return nil
	}
	id := graphql.ID(*r.R.ParentID)
	return &id
}

func (r NewsfeedPostCommentResolverCustom) Tags() *[]NewsfeedTagResolver {
	items := []NewsfeedTagResolver{}
	for i := range r.R.Tags {
		items = append(items, NewsfeedTagResolver{&r.R.Tags[i]})
	}
	return &items
}

func (r NewsfeedPostCommentResolverCustom) Liked() *string {
	if r.R.Liked == "" {
		return nil
	}

	return &r.R.Liked
}

// -----

type LikeCustom struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Emoji string `json:"emoji"`
	// for Subscription
	PostID    string  `json:"post_id"`
	CommentID *string `json:"comment_id"`
}

type LikeResolverCustom struct {
	R *LikeCustom
}

func (r LikeResolverCustom) ID() graphql.ID {
	id := graphql.ID(r.R.ID)
	return id
}

func (r LikeResolverCustom) Type() string {
	return r.R.Type
}

func (r LikeResolverCustom) Emoji() string {
	return r.R.Emoji
}
