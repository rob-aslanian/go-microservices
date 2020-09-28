package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/file"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
)

// Repository ...
type Repository interface {
	SavePost(ctx context.Context, p *post.Post) error
	GetPostByID(ctx context.Context, id string) (*post.Post, error)
	ChangePost(ctx context.Context, p *post.Post) error
	RemovePost(ctx context.Context, postID string) error
	GetNewsfeed(ctx context.Context, requestorID, id string, pinned bool, first, after uint32) (*post.Newsfeed, error)
	GetNewsfeedOfFollowings(ctx context.Context, requestorID string, ids []string, first, after uint32) (*post.Newsfeed, error)
	AddComment(ctx context.Context, com *post.Comment) error
	ChangeComment(ctx context.Context, com *post.Comment) error
	RemoveComment(ctx context.Context, postID, commentID string) error
	GetComments(ctx context.Context, myID, id string, sort post.CommentSort, first, after uint32) (*post.Comments, error)
	GetCommentReplies(ctx context.Context, myID, postID string, commentID string, first, after uint32) (*post.Comments, error)
	GetSharedPosts(ctx context.Context, requestorID, id string, first, after uint32) (*post.Newsfeed, error)
	GetAmountOfSharedPosts(ctx context.Context, postIDs []string) (map[string]uint32, error)
	GetAmountOfReplies(ctx context.Context, postID string, commentsIDs []string) (map[string]uint32, error)
	AddFile(ctx context.Context, postID string, f *file.File) error
	AddFileInComment(ctx context.Context, postID, commentID string, f *file.File) error
	RemoveFile(ctx context.Context, postID string, fileID string) error
	RemoveFileInComment(ctx context.Context, postID, commentID string, fileID string) error
	GetCommentByID(ctx context.Context, postID string, commentID string) (*post.Comment, error)
	LikePost(ctx context.Context, postID string, like *post.Like) error
	LikeComment(ctx context.Context, postID, commentID string, like *post.Like) error
	UnlikePost(ctx context.Context, postID, id string) error
	UnlikeComment(ctx context.Context, postID, commentID string, id string) error
	GetLikeInPostByID(ctx context.Context, postID, id string) (*post.Like, error)
	GetLikeInCommentByID(ctx context.Context, postID, commentID, id string) (*post.Like, error)
	GetLikedListInPost(ctx context.Context, postID string, emoji *post.EmojiType, first, after uint32) ([]*post.Like, error)
	GetLikedListInComment(ctx context.Context, postID, commentID string, emoji *post.EmojiType, first, after uint32) ([]*post.Like, error)
	Search(ctx context.Context, myID, newsfeedID string, keywords string, hashtags []string, first, after uint32) (*post.Newsfeed, error)
}
