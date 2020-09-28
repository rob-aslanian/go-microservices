package serverRPC

import (
	"context"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
)

// AddPost ...
func (s Server) AddPost(ctx context.Context, data *newsfeedRPC.Post) (*newsfeedRPC.ID, error) {
	id, err := s.service.AddPost(ctx, newsfeedRPCPostTopostPost(data))
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.ID{
		ID: id,
	}, nil
}

// ChangePost ...
func (s Server) ChangePost(ctx context.Context, data *newsfeedRPC.Post) (*newsfeedRPC.Empty, error) {
	err := s.service.ChangePost(ctx, newsfeedRPCPostTopostPost(data))
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// RemovePost ...
func (s Server) RemovePost(ctx context.Context, data *newsfeedRPC.RemovePostRequest) (*newsfeedRPC.Empty, error) {
	err := s.service.RemovePost(
		ctx,
		data.GetPostID(),
		data.GetCompanyID(),
	)
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// GetNewsfeed ...
func (s Server) GetNewsfeed(ctx context.Context, data *newsfeedRPC.GetNewsfeedRequest) (*newsfeedRPC.Newsfeed, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, err := strconv.Atoi(data.GetPagination().GetFirst())
		if err == nil {
			first = uint32(f)
		}

		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	feed, err := s.service.GetNewsfeed(
		ctx,
		data.GetID(),
		data.GetCompanyID(),
		data.GetPinned(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	return postNewsfeedToNewsfeedRPCNewsfeed(feed), nil
}

// AddComment ...
func (s Server) AddComment(ctx context.Context, data *newsfeedRPC.Comment) (*newsfeedRPC.ID, error) {
	id, err := s.service.AddComment(ctx, newsfeedRPCCommentTopostComment(data))
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.ID{
		ID: id,
	}, nil
}

// ChangeComment ...
func (s Server) ChangeComment(ctx context.Context, data *newsfeedRPC.Comment) (*newsfeedRPC.Empty, error) {
	err := s.service.ChangeComment(ctx, newsfeedRPCCommentTopostComment(data))
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// RemoveComment ...
func (s Server) RemoveComment(ctx context.Context, data *newsfeedRPC.RemoveCommentRequest) (*newsfeedRPC.Empty, error) {
	err := s.service.RemoveComment(
		ctx,
		data.GetPostID(),
		data.GetCommentID(),
		data.GetCompanyID(),
	)
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// GetComments ...
func (s Server) GetComments(ctx context.Context, data *newsfeedRPC.GetCommentsRequest) (*newsfeedRPC.Comments, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, err := strconv.Atoi(data.GetPagination().GetFirst())
		if err == nil {
			first = uint32(f)
		}

		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	com, err := s.service.GetComments(
		ctx,
		data.GetID(),
		data.GetCompanyID(),
		sortOptionRPCToPostCommentSort(data.GetSort()),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	return commentsNewsfeedToNewsfeedRPCComments(com), nil
}

// GetCommentReplies ...
func (s Server) GetCommentReplies(ctx context.Context, data *newsfeedRPC.GetCommentRepliesRequest) (*newsfeedRPC.Comments, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, err := strconv.Atoi(data.GetPagination().GetFirst())
		if err == nil {
			first = uint32(f)
		}
		// GetCommentReplies(ctx context.Context, postID string, commentID string, first, after uint32) (*post.Comments, error)
		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	com, err := s.service.GetCommentReplies(
		ctx,
		data.GetCompanyID(),
		data.GetPostID(),
		data.GetCommentID(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	return commentsNewsfeedToNewsfeedRPCComments(com), nil
}

// GetSharedPosts ...
func (s Server) GetSharedPosts(ctx context.Context, data *newsfeedRPC.GetSharedPostsRequest) (*newsfeedRPC.Newsfeed, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, err := strconv.Atoi(data.GetPagination().GetFirst())
		if err == nil {
			first = uint32(f)
		}

		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	feed, err := s.service.GetSharedPosts(
		ctx,
		data.GetCompanyID(),
		data.GetID(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	return postNewsfeedToNewsfeedRPCNewsfeed(feed), nil
}

// GetPostByID ...
func (s Server) GetPostByID(ctx context.Context, data *newsfeedRPC.ID) (*newsfeedRPC.Post, error) {
	p, err := s.service.GetPostByID(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return postPostToNewsfeedRPCPost(p), nil
}

// AddFileInPost ...
func (s Server) AddFileInPost(ctx context.Context, data *newsfeedRPC.File) (*newsfeedRPC.ID, error) {
	id, err := s.service.AddFile(
		ctx,
		data.GetUserID(),
		data.GetTargetID(),
		data.GetItemID(),
		data.GetCompanyID(),
		newsfeedRPCFileToFile(data),
	)
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.ID{
		ID: id,
	}, nil
}

// RemoveFileInPost ...
func (s Server) RemoveFileInPost(ctx context.Context, data *newsfeedRPC.RemoveFileInPostRequest) (*newsfeedRPC.Empty, error) {
	err := s.service.RemoveFile(
		ctx,
		data.GetPostID(),
		data.GetCommentID(),
		data.GetCompanyID(),
		data.GetFileID(),
	)
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// Like ...
func (s Server) Like(ctx context.Context, data *newsfeedRPC.LikeRequest) (*newsfeedRPC.Empty, error) {
	err := s.service.Like(
		ctx,
		data.GetPostID(),
		data.GetCommentID(),
		newsfeedRPCLikeToPostLike(data.GetLike()),
	)
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// Unlike ...
func (s Server) Unlike(ctx context.Context, data *newsfeedRPC.UnlikeRequest) (*newsfeedRPC.Empty, error) {
	err := s.service.Unlike(
		ctx,
		data.GetPostID(),
		data.GetCommentID(),
		data.GetID(),
	)
	if err != nil {
		return nil, err
	}

	return &newsfeedRPC.Empty{}, nil
}

// GetLikedList ...
func (s Server) GetLikedList(ctx context.Context, data *newsfeedRPC.GetLikedListRequest) (*newsfeedRPC.GetLikedListResponse, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, err := strconv.Atoi(data.GetPagination().GetFirst())
		if err == nil {
			first = uint32(f)
		}

		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	var emoji *post.EmojiType
	if data.GetEmoji() != "" {
		e := stringToPostEmojiType(data.GetEmoji())
		emoji = &e
	}

	likes, err := s.service.GetLikedList(
		ctx,
		data.GetPostID(),
		data.GetCommentID(),
		emoji,
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	likesRPC := make([]*newsfeedRPC.Like, 0, len(likes))

	for _, l := range likes {
		if l != nil {
			likesRPC = append(likesRPC, postLikeToNewsfeedRPCLike(l))
		}
	}

	return &newsfeedRPC.GetLikedListResponse{
		Likes: likesRPC,
	}, nil
}

// SearchAmongPosts ...
func (s Server) SearchAmongPosts(ctx context.Context, data *newsfeedRPC.SearchAmongPostsRequest) (*newsfeedRPC.Newsfeed, error) {
	var first, after uint32

	if data.GetPagination() != nil {
		f, err := strconv.Atoi(data.GetPagination().GetFirst())
		if err == nil {
			first = uint32(f)
		}

		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	feed, err := s.service.SearchAmongPosts(
		ctx,
		data.GetCompanyID(),
		data.GetNewsfeedID(),
		data.GetKeyword(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	return postNewsfeedToNewsfeedRPCNewsfeed(feed), nil
}

func newsfeedRPCPostTopostPost(data *newsfeedRPC.Post) *post.Post {
	if data == nil {
		return nil
	}

	p := post.Post{
		Text:                   data.GetText(),
		IsCommentDisabled:      data.GetIsCommentedDisabled(),
		IsNotificationDisabled: data.GetIsNotificationDisabled(),
		IsPinned:               data.GetIsPinned(),
		Tags:                   make([]post.Tag, 0, len(data.GetTags())),
	}

	if data.GetID() != "" {
		_ = p.SetID(data.GetID())
	}

	if data.GetUserID() != "" {
		_ = p.SetUserID(data.GetUserID())
	}

	if data.GetCompanyID() != "" {
		_ = p.SetCompanyID(data.GetCompanyID())
	}

	if data.GetNewsFeedUserID() != "" {
		_ = p.SetNewsFeedUserID(data.GetNewsFeedUserID())
	}

	if data.GetNewsFeedCompanyID() != "" {
		_ = p.SetNewsFeedCompanyID(data.GetNewsFeedCompanyID())
	}

	if data.GetSharedPostID() != "" {
		_ = p.SetSharedPost(data.GetSharedPostID())
	}

	for _, t := range data.GetTags() {
		if tag := newsfeedRPCTagToPostTag(t); tag != nil {
			p.Tags = append(p.Tags, *tag)
		}
	}

	return &p
}

func newsfeedRPCTagToPostTag(data *newsfeedRPC.Tag) *post.Tag {
	if data == nil {
		return nil
	}

	t := post.Tag{
		Type: newsfeedRPCTagEntityTypeToString(data.GetEntity()),
	}
	_ = t.SetID(data.GetID())

	return &t
}

func newsfeedRPCTagEntityTypeToString(data newsfeedRPC.EntityType) string {
	switch data {
	case newsfeedRPC.EntityType_Group:
		return "group"
	case newsfeedRPC.EntityType_Company:
		return "company"
	case newsfeedRPC.EntityType_Community:
		return "community"

	}
	return "user"
}

func newsfeedRPCCommentTopostComment(data *newsfeedRPC.Comment) *post.Comment {
	if data == nil {
		return nil
	}

	com := post.Comment{
		Text: data.GetText(),
		Tags: make([]post.Tag, 0, len(data.GetTags())),
	}

	_ = com.SetID(data.GetID())
	_ = com.SetUserID(data.GetUserID())
	_ = com.SetCompanyID(data.GetCompanyID())
	_ = com.SetParentID(data.GetParentID())
	_ = com.SetPostID(data.GetPostID())

	for _, t := range data.GetTags() {
		if tag := newsfeedRPCTagToPostTag(t); tag != nil {
			com.Tags = append(com.Tags, *tag)
		}
	}

	return &com
}

func sortOptionRPCToPostCommentSort(data newsfeedRPC.GetCommentsRequest_SortOption) post.CommentSort {
	switch data {
	case newsfeedRPC.GetCommentsRequest_ByTopLiked:
		return post.CommentSortTop
	}
	return post.CommentSortRecent
}

func newsfeedRPCFileToFile(data *newsfeedRPC.File) *file.File {
	if data == nil {
		return nil
	}

	f := file.File{
		MimeType: data.GetMimeType(),
		Name:     data.GetName(),
		// Position: data.GetPosition(),
		URL: data.GetURL(),
	}

	return &f
}

func newsfeedRPCLikeToPostLike(data *newsfeedRPC.Like) *post.Like {
	if data == nil {
		return nil
	}

	like := post.Like{
		Type:  newsfeedRPCTagEntityTypeToString(data.GetEntity()),
		Emoji: stringToPostEmojiType(data.GetEmoji()),
	}
	like.SetID(data.GetID())

	return &like
}

func stringToPostEmojiType(s string) post.EmojiType {
	switch s {
	case "clap":
		return post.EmojiTypeClap
	case "hmm":
		return post.EmojiTypeHmm
	case "shit":
		return post.EmojiTypeShit
	case "stop":
		return post.EmojiTypeStop
	case "heart":
		return post.EmojiTypeHeart
	case "rocket":
		return post.EmojiTypeRocket
	}

	return post.EmojiTypeLike
}

// -----------------------

func postNewsfeedToNewsfeedRPCNewsfeed(data *post.Newsfeed) *newsfeedRPC.Newsfeed {
	if data == nil {
		return nil
	}

	feed := newsfeedRPC.Newsfeed{
		Amount: data.AmountOfPosts,
		Posts:  make([]*newsfeedRPC.Post, 0, len(data.Posts)),
	}

	for _, p := range data.Posts {
		feed.Posts = append(feed.Posts, postPostToNewsfeedRPCPost(p))
	}

	return &feed
}

func postPostToNewsfeedRPCPost(data *post.Post) *newsfeedRPC.Post {
	if data == nil {
		return nil
	}

	p := newsfeedRPC.Post{
		ID:                     data.GetID(),
		UserID:                 data.GetUserID(),
		CompanyID:              data.GetCompanyID(),
		NewsFeedUserID:         data.GetNewsFeedUserID(),
		NewsFeedCompanyID:      data.GetNewsFeedCompanyID(),
		Text:                   data.Text,
		Hashtags:               data.Hashtags,
		Tags:                   make([]*newsfeedRPC.Tag, 0, len(data.Tags)),
		SharedPostID:           data.GetSharedPost(),
		IsCommentedDisabled:    data.IsCommentDisabled,
		IsNotificationDisabled: data.IsNotificationDisabled,
		CommentsAmount:         data.CommentsAmount,
		SharesAmount:           data.SharesAmount,
		CreatedAt:              timeToString(data.CreatedAt),
		Files:                  make([]*newsfeedRPC.File, 0, len(data.Files)),
	}

	if data.Liked != nil {
		p.Liked = postEmojiTypeToString(data.Liked.Emoji)
	}

	if data.LikesAmount != nil {
		p.LikesAmount = postLikesAmountToNewsfeedLikesAmount(data.LikesAmount)
	}

	if data.ChangedAt != nil {
		p.ChangedAt = timeToString(*data.ChangedAt)
	}

	for _, f := range data.Files {
		p.Files = append(p.Files, fileFileToNewsfeedRPCFile(f))
	}

	for _, t := range data.Tags {
		p.Tags = append(p.Tags, postTagToNewsfeedRPCTag(&t))
	}

	return &p
}

func postLikesAmountToNewsfeedLikesAmount(data *post.LikesAmount) *newsfeedRPC.LikesAmount {
	if data == nil {
		return nil
	}

	likes := newsfeedRPC.LikesAmount{
		Clap:   data.Clap,
		Heart:  data.Heart,
		Hmm:    data.Hmm,
		Like:   data.Like,
		Rocket: data.Rocket,
		Shit:   data.Shit,
		Stop:   data.Stop,
	}

	return &likes
}

func postTagToNewsfeedRPCTag(data *post.Tag) *newsfeedRPC.Tag {
	if data == nil {
		return nil
	}

	t := newsfeedRPC.Tag{
		ID:     data.GetID(),
		Entity: stringToNewsfeedRPCTagEntityType(data.Type),
	}

	return &t
}

func stringToNewsfeedRPCTagEntityType(data string) newsfeedRPC.EntityType {
	switch data {
	case "group":
		return newsfeedRPC.EntityType_Group
	case "company":
		return newsfeedRPC.EntityType_Company
	case "community":
		return newsfeedRPC.EntityType_Community
	}

	return newsfeedRPC.EntityType_User
}

func fileFileToNewsfeedRPCFile(data *file.File) *newsfeedRPC.File {
	if data == nil {
		return nil
	}

	f := newsfeedRPC.File{
		ID:       data.GetID(),
		MimeType: data.MimeType,
		Name:     data.Name,
		URL:      data.URL,
	}

	return &f
}

func commentsNewsfeedToNewsfeedRPCComments(data *post.Comments) *newsfeedRPC.Comments {
	if data == nil {
		return nil
	}

	com := newsfeedRPC.Comments{
		Amount:   data.Amount,
		Comments: make([]*newsfeedRPC.Comment, 0, len(data.Comments)),
	}

	for _, p := range data.Comments {
		com.Comments = append(com.Comments, commentToNewsfeedRPCComment(&p))
	}

	return &com
}

func commentToNewsfeedRPCComment(data *post.Comment) *newsfeedRPC.Comment {
	if data == nil {
		return nil
	}

	p := newsfeedRPC.Comment{
		ID:            data.GetID(),
		UserID:        data.GetUserID(),
		CompanyID:     data.GetCompanyID(),
		Text:          data.Text,
		Tags:          make([]*newsfeedRPC.Tag, 0, len(data.Tags)),
		ParentID:      data.GetParentID(),
		PostID:        data.GetPostID(),
		RepliesAmount: data.AmountOfReplies,
		CreatedAt:     timeToString(data.CreatedAt),
		Files:         make([]*newsfeedRPC.File, 0, len(data.Files)),
	}

	if data.Liked != nil {
		p.Liked = postEmojiTypeToString(data.Liked.Emoji)
	}

	if data.LikesAmount != nil {
		p.LikesAmount = postLikesAmountToNewsfeedLikesAmount(data.LikesAmount)
	}

	if data.ChangedAt != nil {
		p.ChangedAt = timeToString(*data.ChangedAt)
	}

	for _, f := range data.Files {
		p.Files = append(p.Files, fileFileToNewsfeedRPCFile(f))
	}

	for _, t := range data.Tags {
		p.Tags = append(p.Tags, postTagToNewsfeedRPCTag(&t))
	}

	return &p
}

func postEmojiTypeToString(s post.EmojiType) string {
	switch s {
	case post.EmojiTypeClap:
		return "clap"
	case post.EmojiTypeHmm:
		return "hmm"
	case post.EmojiTypeShit:
		return "shit"
	case post.EmojiTypeStop:
		return "stop"
	case post.EmojiTypeHeart:
		return "heart"
	case post.EmojiTypeRocket:
		return "rocket"
	}

	return "like"
}

func postLikeToNewsfeedRPCLike(data *post.Like) *newsfeedRPC.Like {
	if data == nil {
		return nil
	}

	like := newsfeedRPC.Like{
		ID:     data.GetID(),
		Emoji:  postEmojiTypeToString(data.Emoji),
		Entity: stringToNewsfeedRPCTagEntityType(data.Type),
	}

	return &like
}

func timeToString(s time.Time) string {
	return s.Format(time.RFC3339)
}
