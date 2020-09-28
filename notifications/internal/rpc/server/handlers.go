package serverRPC

import (
	"context"
	"encoding/json"
	"strconv"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/notificationsRPC"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/notification"
)

// GetSettings ...
func (s *Server) GetSettings(ctx context.Context, data *notificationsRPC.Empty) (*notificationsRPC.Settings, error) {
	settings, err := s.service.GetSettings(ctx)
	if err != nil {
		return nil, err
	}

	return notificationSettingsToNotificationsRPCSettings(settings), nil
}

// ChangeSettings ...
func (s *Server) ChangeSettings(ctx context.Context, data *notificationsRPC.ChangeSettingsRequest) (*notificationsRPC.Empty, error) {
	err := s.service.ChangeSettings(
		ctx,
		propertOptionRPCToNotificationParameterSetting(data.GetProperty()),
		data.GetValue(),
	)
	if err != nil {
		return nil, err
	}

	//

	return &notificationsRPC.Empty{}, nil
}

// GetNotifications ...
func (s *Server) GetNotifications(ctx context.Context, data *notificationsRPC.Pagination) (*notificationsRPC.NotificationList, error) {
	var first, after uint32
	if value, err := strconv.Atoi(data.GetFirst()); err == nil {
		first = uint32(value)
	}
	if value, err := strconv.Atoi(data.GetAfter()); err == nil {
		after = uint32(value)
	}

	notsMap, amount, err := s.service.GetNotifications(
		ctx,
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	nots := make([]*notificationsRPC.Notification, 0, len(notsMap))

	for _, n := range notsMap {
		j, err := json.Marshal(n)
		if err != nil {
			return nil, err
		}

		nots = append(nots, &notificationsRPC.Notification{
			Notification: string(j),
		})
	}

	return &notificationsRPC.NotificationList{
		Notifications: nots,
		Amount:        amount,
	}, nil
}

// GetUnseenNotifications ...
func (s *Server) GetUnseenNotifications(ctx context.Context, data *notificationsRPC.Pagination) (*notificationsRPC.NotificationList, error) {
	var first, after uint32
	if value, err := strconv.Atoi(data.GetFirst()); err == nil {
		first = uint32(value)
	}
	if value, err := strconv.Atoi(data.GetAfter()); err == nil {
		after = uint32(value)
	}

	notsMap, amount, err := s.service.GetUnseenNotifications(
		ctx,
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	nots := make([]*notificationsRPC.Notification, 0, len(notsMap))

	for _, n := range notsMap {
		j, err := json.Marshal(n)
		if err != nil {
			return nil, err
		}

		nots = append(nots, &notificationsRPC.Notification{
			Notification: string(j),
		})
	}

	return &notificationsRPC.NotificationList{
		Notifications: nots,
		Amount:        amount,
	}, nil
}

// MarkAsSeen ...
func (s *Server) MarkAsSeen(ctx context.Context, data *notificationsRPC.IDs) (*notificationsRPC.Empty, error) {
	err := s.service.MarkAsSeen(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &notificationsRPC.Empty{}, nil
}

// RemoveNotification ...
func (s *Server) RemoveNotification(ctx context.Context, data *notificationsRPC.IDs) (*notificationsRPC.Empty, error) {
	err := s.service.RemoveNotification(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &notificationsRPC.Empty{}, nil
}

// GetCompanySettings ...
func (s *Server) GetCompanySettings(ctx context.Context, data *notificationsRPC.ID) (*notificationsRPC.CompanySettings, error) {
	settings, err := s.service.GetCompanySettings(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return notificationCompanySettingsToNotificationsCompanySettingsRPC(settings), nil
}

// ChangeCompanySettings ...
func (s *Server) ChangeCompanySettings(ctx context.Context, data *notificationsRPC.ChangeCompanySettingsRequest) (*notificationsRPC.Empty, error) {
	err := s.service.ChangeCompanySettings(
		ctx,
		data.GetCompanyID(),
		propertCompanyOptionRPCToNotificationParameterSetting(data.GetProperty()),
		data.GetValue(),
	)
	if err != nil {
		return nil, err
	}

	return &notificationsRPC.Empty{}, nil
}

// GetCompanyNotifications ...
func (s *Server) GetCompanyNotifications(ctx context.Context, data *notificationsRPC.PaginationWithID) (*notificationsRPC.NotificationList, error) {
	var first, after uint32
	if value, err := strconv.Atoi(data.GetFirst()); err == nil {
		first = uint32(value)
	}
	if value, err := strconv.Atoi(data.GetAfter()); err == nil {
		after = uint32(value)
	}

	notsMap, amount, err := s.service.GetCompanyNotifications(
		ctx,
		data.GetID(),
		first,
		after,
	)
	if err != nil {
		return nil, err
	}

	nots := make([]*notificationsRPC.Notification, 0, len(notsMap))

	for _, n := range notsMap {
		j, err := json.Marshal(n)
		if err != nil {
			return nil, err
		}

		nots = append(nots, &notificationsRPC.Notification{
			Notification: string(j),
		})
	}

	return &notificationsRPC.NotificationList{
		Notifications: nots,
		Amount:        amount,
	}, nil
}

// MarkAsSeenForCompany ...
func (s *Server) MarkAsSeenForCompany(ctx context.Context, data *notificationsRPC.IDWithIDs) (*notificationsRPC.Empty, error) {
	err := s.service.MarkAsSeenForCompany(
		ctx,
		data.GetID(),
		data.GetIDs(),
	)
	if err != nil {
		return nil, err
	}

	return &notificationsRPC.Empty{}, nil
}

// RemoveNotificationForCompany ...
func (s *Server) RemoveNotificationForCompany(ctx context.Context, data *notificationsRPC.IDWithIDs) (*notificationsRPC.Empty, error) {
	err := s.service.RemoveNotificationForCompany(ctx, data.GetID(), data.GetIDs())
	if err != nil {
		return nil, err
	}

	return &notificationsRPC.Empty{}, nil
}

// To RPC

func notificationSettingsToNotificationsRPCSettings(data *notification.Settings) *notificationsRPC.Settings {
	if data == nil {
		return nil
	}

	set := notificationsRPC.Settings{
		// ApprovedConnection:    data.ApprovedConnection,
		// NewConnection:         data.NewConnection,
		// NewEndorsement:        data.NewEndorsement,
		// NewFollow:             data.NewFollow,
		// NewRecommendation:     data.NewRecommendation,
		// RecommendationRequest: data.RecommendationRequest,
		// NewJobInvitation:      data.NewJobInvitation,
	}

	if data.ApprovedConnection == nil {
		set.ApprovedConnection = true
	} else {
		set.ApprovedConnection = *data.ApprovedConnection
	}
	if data.NewConnection == nil {
		set.NewConnection = true
	} else {
		set.NewConnection = *data.NewConnection
	}
	if data.NewEndorsement == nil {
		set.NewEndorsement = true
	} else {
		set.NewEndorsement = *data.NewEndorsement
	}
	if data.NewFollow == nil {
		set.NewFollow = true
	} else {
		set.NewFollow = *data.NewFollow
	}
	if data.NewRecommendation == nil {
		set.NewRecommendation = true
	} else {
		set.NewRecommendation = *data.NewRecommendation
	}
	if data.RecommendationRequest == nil {
		set.RecommendationRequest = true
	} else {
		set.RecommendationRequest = *data.RecommendationRequest
	}
	if data.NewJobInvitation == nil {
		set.NewJobInvitation = true
	} else {
		set.NewJobInvitation = *data.NewJobInvitation
	}

	return &set
}

func notificationCompanySettingsToNotificationsCompanySettingsRPC(data *notification.CompanySettings) *notificationsRPC.CompanySettings {
	if data == nil {
		return nil
	}

	set := notificationsRPC.CompanySettings{}
	if data.NewFollow == nil {
		set.NewFollow = true
	} else {
		set.NewFollow = *data.NewFollow
	}
	if data.NewReview == nil {
		set.NewReview = true
	} else {
		set.NewReview = *data.NewReview
	}
	if data.NewApplicant == nil {
		set.NewApplicant = true
	} else {
		set.NewApplicant = *data.NewApplicant
	}

	return &set
}

// from RPC

func propertOptionRPCToNotificationParameterSetting(data notificationsRPC.ChangeSettingsRequest_PropertyOption) notification.ParameterSetting {
	switch data {
	case notificationsRPC.ChangeSettingsRequest_ApprovedConnection:
		return notification.ApprovedConnection
	case notificationsRPC.ChangeSettingsRequest_NewConnection:
		return notification.NewConnection
	case notificationsRPC.ChangeSettingsRequest_NewEndorsement:
		return notification.NewEndorsement
	case notificationsRPC.ChangeSettingsRequest_NewFollow:
		return notification.NewFollow
	case notificationsRPC.ChangeSettingsRequest_NewRecommendation:
		return notification.NewRecommendation
	case notificationsRPC.ChangeSettingsRequest_RecommendationRequest:
		return notification.RecommendationRequest
	case notificationsRPC.ChangeSettingsRequest_NewJobInvitation:
		return notification.NewJobInvitation
	}

	return notification.ParameterUnknonwn
}

func propertCompanyOptionRPCToNotificationParameterSetting(data notificationsRPC.ChangeCompanySettingsRequest_PropertyOption) notification.ParameterSetting {
	switch data {
	case notificationsRPC.ChangeCompanySettingsRequest_NewFollow:
		return notification.NewFollow
	case notificationsRPC.ChangeCompanySettingsRequest_NewReview:
		return notification.NewReview
	case notificationsRPC.ChangeCompanySettingsRequest_NewApplicant:
		return notification.NewApplicant
	}

	return notification.ParameterUnknonwn
}

// ---------

// func stringDateToTime(s string) time.Time {
// 	if date, err := time.Parse("2-1-2006", s); err == nil {
// 		return date
// 	}
// 	return time.Time{}
// }
//
// func stringDayMonthAndYearToTime(s string) time.Time {
// 	if date, err := time.Parse("1-2006", s); err == nil {
// 		return date
// 	}
// 	return time.Time{}
// }
//
// func timeToStringMonthAndYear(t time.Time) string {
// 	if t == (time.Time{}) {
// 		return ""
// 	}
//
// 	y, m, _ := t.UTC().Date()
// 	return strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
// }
//
// func timeToStringDayMonthAndYear(t time.Time) string {
// 	if t == (time.Time{}) {
// 		return ""
// 	}
//
// 	y, m, d := t.UTC().Date()
// 	return strconv.Itoa(d) + "-" + strconv.Itoa(int(m)) + "-" + strconv.Itoa(y)
// }
