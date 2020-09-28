package serverRPC

import (
	"context"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/session"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
)

// LoginUser uses service functions. logs in user and adds it's token to the session of gRPC
func (c *Server) LoginUser(ctx context.Context, data *authRPC.User) (*authRPC.Session, error) {

	token, err := c.service.LoginUser(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &authRPC.Session{
		Token: token,
	}, nil
}

// GetUser uses service functions. ...
func (c *Server) GetUser(ctx context.Context, session *authRPC.Session) (*authRPC.User, error) {
	cmd, err := c.service.GetUser(ctx, session.Token)
	if err != nil {
		return nil, err
	}
	return &authRPC.User{
		Id: cmd,
	}, nil
}

// LogoutSession uses service functions. ...
func (c *Server) LogoutSession(ctx context.Context, session *authRPC.Session) (*authRPC.Empty, error) {
	err := c.service.LogoutUser(ctx, session.Token)
	if err != nil {
		return nil, err
	}
	return &authRPC.Empty{}, nil

}

// LogoutOtherSession ... NEW
func (c *Server) LogoutOtherSession(ctx context.Context, sessionID *authRPC.Session) (*authRPC.Empty, error) {
	err := c.service.LogoutOtherSession(ctx, sessionID.GetID())
	if err != nil {
		return nil, err
	}

	return &authRPC.Empty{}, nil
}

// SignOutFromAll sight out from all sessions except current one
func (c *Server) SignOutFromAll(ctx context.Context, session *authRPC.Session) (*authRPC.Empty, error) {
	err := c.service.SignOutFromAll(ctx, session.Token)
	if err != nil {
		return nil, err
	}
	return &authRPC.Empty{}, nil

}

// GetListOfSessions uses service functions. ...
func (c *Server) GetListOfSessions(ctx context.Context, data *authRPC.ListOfSessionsQuery) (*authRPC.ListOfSessions, error) {
	sessions, err := c.service.GetListOfSessions(
		ctx,
		data.GetFirst(),
		data.GetAfter(),
	)
	if err != nil {
		return nil, err
	}

	ses := make([]*authRPC.Sessions, 0, len(sessions))

	for i := range sessions {
		ses = append(ses, authSessionToAuthSessionRPC(&(sessions[i])))

	}
	return &authRPC.ListOfSessions{
		Sessions: ses,
	}, nil
}

// GetAmountOfSessions ...
func (c *Server) GetAmountOfSessions(ctx context.Context, data *authRPC.Empty) (*authRPC.Amount, error) {
	amount, err := c.service.GetAmountOfSessions(ctx)
	if err != nil {
		return nil, err
	}

	return &authRPC.Amount{
		Amount: amount,
	}, nil
}

// GetTimeOfLastActivity uses service functions. ...
func (c *Server) GetTimeOfLastActivity(ctx context.Context, data *authRPC.User) (*authRPC.Time, error) {
	tm, err := c.service.GetTimeOfLastActivity(ctx, data.GetId())
	if err != nil {
		return nil, err
	}

	return &authRPC.Time{
		Time: tm.UnixNano(),
	}, nil
}

func authSessionToAuthSessionRPC(data *session.Session) *authRPC.Sessions {
	if data == nil {
		return nil
	}

	ses := authRPC.Sessions{
		ID:               data.GetID(),
		OS:               data.DeviceInfo.OS.Name,
		OSVersion:        strconv.Itoa(data.DeviceInfo.OS.Version.Major),
		DeviceType:       data.DeviceInfo.Type,
		Browser:          data.DeviceInfo.Browser.Name,
		BrowserVersion:   strconv.Itoa(data.DeviceInfo.Browser.Version.Major),
		LastActivityTime: timeToString(data.LastActivityTime),
		City:             nullUintToUint32(data.Location.City),
		CountryID:        nullStringToString(data.Location.Country),
		CurrentSession:   data.CurrentSession,
	}
	return &ses
}

func nullStringToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func nullUintToUint32(s *uint) uint32 {
	if s == nil {
		return 0
	}
	return uint32(*s)
}

func timeToString(t time.Time) string {
	var ti string

	ti = t.Format(time.RFC3339)

	return ti
}
