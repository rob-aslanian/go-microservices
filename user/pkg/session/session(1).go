package session

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/location"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
)

// Session represent session of user
type Session struct {
	ID         bson.ObjectId     `bson:"_id"`
	User       bson.ObjectId     `bson:"user"`
	Token      string            `bson:"token" json:"-"`
	IPAddress  string            `bson:"ip"`
	UserAgent  string            `bson:"user_agent"`
	DeviceInfo DeviceInfo        `bson:"device_info"`
	Location   location.Location `bson:"location"`
	Status     string            `bson:"status"`
	CreatedAt  time.Time         `bson:"created_at"`
}

// DeviceInfo represents information about user's device
type DeviceInfo struct {
	Browser Browser `bson:"browser"`
	OS      OS      `bson:"os"`
	Type    string  `bson:"type"`
}

// Browser represents information about user's browser
type Browser struct {
	Name    string  `bson:"name"`
	Version Version `bson:"version"`
}

// OS represents information about user's operation system
type OS struct {
	Name     string  `bson:"name"`
	Platform string  `bson:"platform"`
	Version  Version `bson:"version"`
}

// Version represents version of  OS or browser
type Version struct {
	Major int `bson:"major"`
	Minor int `bson:"minor"`
	Patch int `bson:"patch"`
}

// GetID returns id of session
func (s *Session) GetID() string {
	return s.ID.Hex()
}

// SetID applies id for session. If id has a wrong format returns usersErrors.ErrWrongID error.
func (s *Session) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		s.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for session
func (s *Session) GenerateID() string {
	id := bson.NewObjectId()
	s.ID = id
	return id.Hex()
}

// GetUserID returns id of user
func (s *Session) GetUserID() string {
	return s.ID.Hex()
}

// SetUserID applies id for user. If id has a wrong format returns usersErrors.ErrWrongID error.
func (s *Session) SetUserID(id string) error {
	if bson.IsObjectIdHex(id) {
		s.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}
