package service

import (
	"context"

	companyadmin "gitlab.lan/Rightnao-site/microservices/news_feed/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
)

// AuthRPC represents a Auth gRPC client
type AuthRPC interface {
	GetUser(ctx context.Context, token string) (string, error)
}

// NetworkRPC ...
type NetworkRPC interface {
	GetAdminLevel(ctx context.Context, companyID string) (companyadmin.AdminLevel, error)
	GetFollowersIDs(ctx context.Context, id string, isCompany bool) ([]string, error)
}

// MQ ...
type MQ interface {
	SendNewPostEvent(p *post.Post) error
	SendNewCommentEvent(p *post.Comment) error
	SendNewLikeEvent(c *post.Like) error
}
