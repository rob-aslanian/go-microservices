package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/notifications/internal/notification"
)

// NotificationsRepository contains functions which have to be in users repository
type NotificationsRepository interface {
	GetNotificationsSettings(ctx context.Context, userID string) (*notification.Settings, error)
	ChangeNotificationsSettings(ctx context.Context, userID string, parameter notification.ParameterSetting, value bool) error
	GetNotifications(ctx context.Context, userID string, first uint32, after uint32) ([]map[string]interface{}, error)
	GetUnseenNotifications(ctx context.Context, userID string, first uint32, after uint32) ([]map[string]interface{}, error)
	MarkAsSeen(ctx context.Context, userID string, ids []string) error
	RemoveNotification(ctx context.Context, userID string, ids []string) error
	GetAmountOfNotSeenNotifications(ctx context.Context, userID string) (int32, error)

	GetCompanyNotificationsSettings(ctx context.Context, companyID string) (*notification.CompanySettings, error)
	ChangeCompanyNotificationsSettings(ctx context.Context, companyID string, parameter notification.ParameterSetting, value bool) error
	GetCompanyNotifications(ctx context.Context, companyID string, first uint32, after uint32) ([]map[string]interface{}, error)
	GetAmountOfNotSeenCompanyNotifications(ctx context.Context, companyID string) (int32, error)
	MarkAsSeenForCompany(ctx context.Context, companyID string, ids []string) error
	RemoveNotificationForCompany(ctx context.Context, companyID string, ids []string) error
}
