package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/internal/account"
	notmes "gitlab.lan/Rightnao-site/microservices/company/pkg/notification_messages"
)

// NetworkRPC represents a Network gRPC client
type NetworkRPC interface {
	MakeCompanyOwner(ctx context.Context, userID string, companyID string) error
	MakeCompanyAdmin(ctx context.Context, userID string, companyID string, adminLevel account.AdminLevel) error
	GetCompanyFollowersNumber(ctx context.Context, companyID string) (uint32, error)
	GetAdminLevel(ctx context.Context, companyID string) (account.AdminLevel, error)
	GetCompanyCountings(ctx context.Context, companyID string) (followings, followers, employees int32, err error)
	AddCompanyAdmin(ctx context.Context, companyID string, userID string, level account.AdminLevel) error
	DeleteCompanyAdmin(ctx context.Context, companyID string, userID string) error
	IsFollow(ctx context.Context, companyID string) (bool, error)
	IsFavourite(ctx context.Context, companyID string) (bool, error)
	// IsFavouriteCompany(ctx context.Context, userID, companyID string) (bool, error)
	IsBlockedCompany(ctx context.Context, companyID string) (bool, error)
	IsFollowForCompany(ctx context.Context, userID, companyID string) (bool, error)
	IsBlockedCompanyForCompany(ctx context.Context, userID, companyID string) (bool, error)
	IsBlockedCompanyByUser(ctx context.Context, companyID string) (bool, error)
	IsBlockedCompanyByCompany(ctx context.Context, userID, companyID string) (bool, error)
}

// AuthRPC represents a Auth gRPC client
type AuthRPC interface {
	GetUserID(ctx context.Context) (string, error)
}

// MailRPC represents a Mail gRPC client
type MailRPC interface {
	SendEmail(ctx context.Context, email string, message string) error
}

// InfoRPC represents a Info gRPC client
type InfoRPC interface {
	GetCountryIDAndCountryCode(ctx context.Context, countryCodeID int32) (countryCode string, countryID string, err error)
	GetCityInformationByID(ctx context.Context, cityID int32, lang *string) (cityName, subdivision, countryID string, err error)
}

// JobsRPC ...
type JobsRPC interface {
	GetAmountOfActiveJobsOfCompany(ctx context.Context, companyID string) (int32, error)
}

// StuffRPC ... 
type StuffRPC interface {
	AddGoldCoinsToWallet(ctx context.Context , userID string , coins int32) error 
}

// UserRPC represents a Info gRPC client
type UserRPC interface {
	CheckPassword(ctx context.Context, password string) error
}

// ChatRPC ...
type ChatRPC interface {
	IsLive(ctx context.Context, id string) (bool, error)
}

// MQ ...
type MQ interface {
	SendNewCompanyReview(companyID string, message *notmes.NewCompanyReview) error
	SendNewFounderRequest(companyID string, message *notmes.NewFounderRequest) error

	SendEmail(address string, message string) error
}
