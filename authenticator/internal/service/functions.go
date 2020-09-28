package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/avct/uasurfer"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/session"
	"google.golang.org/grpc/metadata"
)

// InsertSession is a function to insert the session, get it's info and save it in db
func (s *Service) InsertSession(ctx context.Context, userID string, ses session.Session) error {
	span := s.tracer.MakeSpan(ctx, "InsertSession")
	defer span.Finish()

	// retrive ip and user-agent from context
	mt, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// TODO: handle error
		return errors.New("can't extract metadata from context")
	}
	var ip, ua string
	if len(mt["ip"]) > 0 {
		ip = mt["ip"][0]
	} else {
		return errors.New("IP is empty") // TODO: handle error
	}
	if len(mt["http_user_agent"]) > 0 {
		ua = mt["http_user_agent"][0]
	} else {
		return errors.New("User-Agent is empty") // TODO: handle error
	}

	ses.GenerateID()
	ses.SetUser(userID)

	// identify location by ip address
	loc, err := s.ipDB.GetCityByIP(ip)
	if err != nil {
		s.tracer.LogError(span, err)
	}
	ses.Location = loc

	ses.Status = session.StatusActive
	ses.CreatedAt = time.Now()
	ses.IPAddress = ip
	ses.UserAgent = ua

	// get device info by user-agent
	uaParsed := uasurfer.Parse(ua)
	ses.DeviceInfo.Browser.Name = uaParsed.Browser.Name.String()[7:]
	ses.DeviceInfo.Browser.Version.Major = uaParsed.Browser.Version.Major
	ses.DeviceInfo.Browser.Version.Minor = uaParsed.Browser.Version.Minor
	ses.DeviceInfo.Browser.Version.Patch = uaParsed.Browser.Version.Patch
	ses.DeviceInfo.OS.Name = uaParsed.OS.Name.String()[2:]
	ses.DeviceInfo.OS.Platform = uaParsed.OS.Platform.String()[8:]
	ses.DeviceInfo.OS.Version.Major = uaParsed.OS.Version.Major
	ses.DeviceInfo.OS.Version.Minor = uaParsed.OS.Version.Minor
	ses.DeviceInfo.OS.Version.Patch = uaParsed.OS.Version.Patch
	ses.DeviceInfo.Type = uaParsed.DeviceType.String()[6:]

	// save session info in db
	err = s.auth.InsertSession(ctx, ses)
	if err != nil {
		return err
	}

	return nil
}

// LoginUser generates token for user, saves it in session and then saves it in db
// It also uses insertSession function and saves token in it.
func (s *Service) LoginUser(ctx context.Context, userID string) (string, error) {
	token := randString(32)
	err := s.tokens.SaveToken(token, userID)
	if err != nil {
		return "", err
	}

	ses := session.Session{}
	ses.Token = token

	err = s.InsertSession(ctx, userID, ses)
	if err != nil {
		// TODO: what behavior should be in this case?
		return "", err
	}

	return token, nil
}

// GetUser gets the user with it's token, sets the token to live for another 24 hours
// It also updates activity time. we call this function every time this user does something
func (s *Service) GetUser(ctx context.Context, token string) (string, error) {
	id, err := s.tokens.GetUserByToken(token)
	if err != nil {
		return "", err
	}

	s.tokens.ExpireToken(token)

	s.auth.UpdateActivityTime(ctx, token)

	return id, nil
}

// LogoutUser deletes the token from redis when user logs out and also deactivates the
// session with the passed token
func (s *Service) LogoutUser(ctx context.Context, token string) error {
	err := s.tokens.DeleteToken(token)
	if err != nil {
		return err
	}

	err = s.auth.DeactivateSessionByToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

// LogoutOtherSession logs user out to any session by his choosing  NEW
func (s *Service) LogoutOtherSession(ctx context.Context, sessionID string) error {

	token := retriveToken(ctx)
	if token == "" {
		return errors.New("token_is_empty")
	}

	userID, err := s.GetUser(ctx, token)
	if err != nil {
		return err
	}

	tokenOld, err := s.auth.DeactivateSpecificSessionByToken(ctx, userID, sessionID)
	if err != nil {
		return err
	}

	err = s.tokens.DeleteToken(tokenOld)
	if err != nil {
		return err
	}

	return nil
}

// SignOutFromAll ...
func (s *Service) SignOutFromAll(ctx context.Context, token string) error {
	userID, tokens, err := s.auth.DeactivateSessions(ctx, token)
	if err != nil {
		return err
	}

	log.Println(userID, "current token:", token, "other tokens:", tokens)

	if len(tokens) > 0 {
		err = s.tokens.DeleteTokens(tokens)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetListOfSessions brings us back all the user's sessions with passed ID. when did he log in,
// from where and what device etc
func (s *Service) GetListOfSessions(ctx context.Context, first, after int32) ([]session.Session, error) {
	token := retriveToken(ctx)
	if token == "" {
		return nil, errors.New("token_is_empty")
	}

	userID, err := s.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}

	ses, err := s.auth.GetListOfSessions(ctx, userID, first, after)
	if err != nil {
		return nil, err
	}

	for i := range ses {
		if ses[i].Token == token {
			ses[i].CurrentSession = true
			break
		}
	}

	return ses, nil
}

// GetAmountOfSessions ...
func (s *Service) GetAmountOfSessions(ctx context.Context) (int32, error) {
	token := retriveToken(ctx)
	if token == "" {
		return 0, errors.New("token_is_empty")
	}

	userID, err := s.GetUser(ctx, token)
	if err != nil {
		return 0, err
	}

	amount, err := s.auth.GetAmountOfSessions(ctx, userID)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

// GetTimeOfLastActivity gives us the last activity time of the user with passed ID
func (s *Service) GetTimeOfLastActivity(ctx context.Context, id string) (time.Time, error) {
	act, err := s.auth.GetTimeOfLastActivity(ctx, id)
	if err != nil {
		return time.Time{}, err
	}

	return act, nil
}

func retriveToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}
