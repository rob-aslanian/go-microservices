package service

import (
	"context"
	"errors"

	"gitlab.lan/Rightnao-site/microservices/notifications/internal/notification"
	"google.golang.org/grpc/metadata"
)

// user

// GetSettings ...
func (s Service) GetSettings(ctx context.Context) (*notification.Settings, error) {
	span := s.tracer.MakeSpan(ctx, "GetSettings")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	set, err := s.repository.Notifications.GetNotificationsSettings(
		ctx,
		userID,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return set, nil
}

// ChangeSettings ...
func (s Service) ChangeSettings(ctx context.Context, par notification.ParameterSetting, value bool) error {
	span := s.tracer.MakeSpan(ctx, "ChangeSettings")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	if par == notification.ParameterUnknonwn {
		return errors.New("unknown_parameter")
	}

	err = s.repository.Notifications.ChangeNotificationsSettings(ctx, userID, par, value)
	if err != nil {
		return err
	}

	return nil
}

// GetNotifications ...
func (s Service) GetNotifications(ctx context.Context, first uint32, after uint32) ([]map[string]interface{}, int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetNotifications")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, 0, errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	nots, err := s.repository.Notifications.GetNotifications(
		ctx,
		userID,
		first,
		after,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	amount, err := s.repository.Notifications.GetAmountOfNotSeenNotifications(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return nots, amount, nil
}

// GetUnseenNotifications ...
func (s Service) GetUnseenNotifications(ctx context.Context, first uint32, after uint32) ([]map[string]interface{}, int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetNotifications")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return nil, 0, errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	nots, err := s.repository.Notifications.GetUnseenNotifications(
		ctx,
		userID,
		first,
		after,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	amount, err := s.repository.Notifications.GetAmountOfNotSeenNotifications(ctx, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return nots, amount, nil
}

// MarkAsSeen ...
func (s Service) MarkAsSeen(ctx context.Context, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "MarkAsSeen")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Notifications.MarkAsSeen(ctx, userID, ids)
	if err != nil {
		return err
	}

	return nil
}

// RemoveNotification ...
func (s Service) RemoveNotification(ctx context.Context, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "MarkAsSeen")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	userID, err := s.authRPC.GetUserID(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	err = s.repository.Notifications.RemoveNotification(ctx, userID, ids)
	if err != nil {
		return err
	}

	return nil
}

// company

// GetCompanySettings ...
func (s Service) GetCompanySettings(ctx context.Context, companyID string) (*notification.CompanySettings, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanySettings")
	defer span.Finish()

	isAdmin, err := s.networkRPC.IsAdmin(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	if !isAdmin {
		return nil, errors.New("not_allowed")
	}

	settings, err := s.repository.Notifications.GetCompanyNotificationsSettings(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return settings, nil
}

// ChangeCompanySettings ...
func (s Service) ChangeCompanySettings(ctx context.Context, companyID string, par notification.ParameterSetting, value bool) error {
	span := s.tracer.MakeSpan(ctx, "ChangeCompanySettings")
	defer span.Finish()

	isAdmin, err := s.networkRPC.IsAdmin(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	if par == notification.ParameterUnknonwn {
		return errors.New("unknown_parameter")
	}

	err = s.repository.Notifications.ChangeCompanyNotificationsSettings(ctx, companyID, par, value)
	if err != nil {
		return err
	}

	return nil
}

// GetCompanyNotifications ...
func (s Service) GetCompanyNotifications(ctx context.Context, companyID string, first uint32, after uint32) ([]map[string]interface{}, int32, error) {
	span := s.tracer.MakeSpan(ctx, "GetCompanyNotifications")
	defer span.Finish()

	// token := s.retriveToken(ctx)
	// if token == "" {
	// 	return nil, 0, errors.New("token_is_empty")
	// }
	//
	// userID, err := s.authRPC.GetUserID(ctx, token)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return nil, 0, err
	// }

	isAdmin, err := s.networkRPC.IsAdmin(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}
	if !isAdmin {
		return nil, 0, errors.New("not_allowed")
	}

	nots, err := s.repository.Notifications.GetCompanyNotifications(
		ctx,
		companyID,
		first,
		after,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	amount, err := s.repository.Notifications.GetAmountOfNotSeenCompanyNotifications(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, 0, err
	}

	return nots, amount, nil
}

// MarkAsSeenForCompany ...
func (s Service) MarkAsSeenForCompany(ctx context.Context, companyID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "MarkAsSeenForCompany")
	defer span.Finish()

	// token := s.retriveToken(ctx)
	// if token == "" {
	// 	return errors.New("token_is_empty")
	// }
	//
	// userID, err := s.authRPC.GetUserID(ctx, token)
	// if err != nil {
	// 	s.tracer.LogError(span, err)
	// 	return err
	// }

	isAdmin, err := s.networkRPC.IsAdmin(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	err = s.repository.Notifications.MarkAsSeenForCompany(ctx, companyID, ids)
	if err != nil {
		return err
	}

	return nil
}

// RemoveNotificationForCompany ...
func (s Service) RemoveNotificationForCompany(ctx context.Context, companyID string, ids []string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveNotificationForCompany")
	defer span.Finish()

	token := s.retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	// TODO: check if admin

	err := s.repository.Notifications.RemoveNotificationForCompany(ctx, companyID, ids)
	if err != nil {
		return err
	}

	return nil
}

// --------------------

func (s Service) retriveToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}
