package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/notifications/internal/notification"
)

// Service define functions inside Service
type Service interface {
	GetSettings(ctx context.Context) (*notification.Settings, error)
	ChangeSettings(ctx context.Context, par notification.ParameterSetting, value bool) error
	GetNotifications(ctx context.Context, first uint32, after uint32) ([]map[string]interface{}, int32, error)
	GetUnseenNotifications(ctx context.Context, first uint32, after uint32) ([]map[string]interface{}, int32, error)
	MarkAsSeen(ctx context.Context, ids []string) error
	RemoveNotification(ctx context.Context, ids []string) error

	GetCompanySettings(ctx context.Context, companyID string) (*notification.CompanySettings, error)
	ChangeCompanySettings(ctx context.Context, companyID string, par notification.ParameterSetting, value bool) error
	GetCompanyNotifications(ctx context.Context, companyID string, first uint32, after uint32) ([]map[string]interface{}, int32, error)
	MarkAsSeenForCompany(ctx context.Context, companyID string, ids []string) error
	RemoveNotificationForCompany(ctx context.Context, companyID string, ids []string) error
}
