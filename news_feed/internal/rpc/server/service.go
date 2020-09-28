package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
)

// Service define functions inside Service
type Service interface {
	AddPost(ctx context.Context, p *post.Post) (string, error)
	ChangePost(ctx context.Context, p *post.Post) error
	RemovePost(ctx context.Context, postID string, companyID string) error
	GetNewsfeed(ctx context.Context, id, companyID string, pinned bool, first, after uint32) (*post.Newsfeed, error)
	AddComment(ctx context.Context, com *post.Comment) (string, error)
	ChangeComment(ctx context.Context, com *post.Comment) error
	RemoveComment(ctx context.Context, postID, commentID, companyID string) error
	GetComments(ctx context.Context, id string, companyID string, sort post.CommentSort, first, after uint32) (*post.Comments, error)
	GetCommentReplies(ctx context.Context, companyID, postID string, commentID string, first, after uint32) (*post.Comments, error)
	GetSharedPosts(ctx context.Context, companyID, id string, first, after uint32) (*post.Newsfeed, error)
	GetPostByID(ctx context.Context, postID string) (*post.Post, error)
	AddFile(ctx context.Context, userID, postID string, commentID, companyID string, f *file.File) (string, error)
	RemoveFile(ctx context.Context, postID string, commentID, companyID string, fileID string) error
	Like(ctx context.Context, postID string, commentID string, like *post.Like) error
	Unlike(ctx context.Context, postID string, commentID string, id string) error
	GetLikedList(ctx context.Context, postID, commentID string, emoji *post.EmojiType, first, after uint32) ([]*post.Like, error)
	SearchAmongPosts(ctx context.Context, companyID, newsfeedID string, keyword string, first, after uint32) (*post.Newsfeed, error)
}
