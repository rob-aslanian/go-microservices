package service

import (
	"context"
	"time"

	notmes "gitlab.lan/Rightnao-site/microservices/user/pkg/notification_messages"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/profile"
)

// AuthRPC represents Auth gRPC client
type AuthRPC interface {
	GetUserID(ctx context.Context, token string) (string, error)
	LoginUser(ctx context.Context, userID string) (token string, err error)
	SignOut(ctx context.Context, token string) error
	SignOutSession(ctx context.Context, sessionID string) error
	SignOutFromAll(ctx context.Context, token string) error
	GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error)
	GetAmountOfSessions(ctx context.Context) (int32, error)
}

// // MailRPC represents Mail gRPC client
// type MailRPC interface {
// 	SendEmail(ctx context.Context, email string, message string) error
// }

// InfoRPC represents Info gRPC client
type InfoRPC interface {
	GetCountryIDAndCountryCode(ctx context.Context, countryCodeID int32) (countryCode string, countryID string, err error)
	GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error)
}

// NetworkRPC represents Network gRPC client
type NetworkRPC interface {
	IsFriend(ctx context.Context, userID string) (bool, error)
	IsBlocked(ctx context.Context, userID string) (bool, error)
	IsBlockedByUser(ctx context.Context, userID string) (bool, error)
	IsFavourite(ctx context.Context, userID string) (bool, error)
	IsFollowing(ctx context.Context, userID string) (bool, error)
	IsFriendRequestSend(ctx context.Context, userID string) (bool, error)
	IsFriendRequestRecieved(ctx context.Context, userID string) (bool, string, error)
	GetUserCompanies(ctx context.Context, userID string) ([]string, error)
	GetReceivedRecommendationByID(ctx context.Context, userID string, first int32, after int32) ([]*profile.Recommendation, error)
	GetGivenRecommendationByID(ctx context.Context, userID string, first int32, after int32) ([]*profile.Recommendation, error)
	GetReceivedRecommendationRequests(ctx context.Context, userID string, first int32, after int32) ([]*profile.RecommendationRequest, error)
	GetRequestedRecommendationRequests(ctx context.Context, userID string, first int32, after int32) ([]*profile.RecommendationRequest, error)
	GetHiddenRecommendationByID(ctx context.Context, userID string, first int32, after int32) ([]*profile.Recommendation, error)
	GetFriendshipID(ctx context.Context, targetUserID string) (string, error)
	IsBlockedForCompany(ctx context.Context, userID string, companyID string) (bool, error)
	IsBlockedByCompany(ctx context.Context, userID string, companyID string) (bool, error)
	IsFollowingForCompany(ctx context.Context, userID string, companyID string) (bool, error)
	GetAmountOfMutualConnections(ctx context.Context, targetUserID string) (int32, error)
}

// CompanyRPC represents Company gRPC client
type CompanyRPC interface {
	GetCompanies(ctx context.Context, ids []string) (interface{}, error)
}

// StuffRPC ...
type StuffRPC interface {
	ContactInvitationForWallet(ctx context.Context , name string , email string , message string , coins int32) error
	AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32) error 
	CreateWalletAccount(ctx context.Context , userID string) error
}

// ChatRPC represents Chat gRPC client
type ChatRPC interface {
	IsOnline(ctx context.Context, userID string) (bool, error)
}

// MQ represents a message queue
type MQ interface {
	SendNewEndorsement(userID string, message *notmes.NewEndorsement) error
	SendEmail(address string, subject string, message string) error
}
