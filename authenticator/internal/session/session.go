package session

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/device_info"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/location"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Status is a type to then easily use the below variables
type Status string

const (
	// StatusActive holds the string of Active
	StatusActive Status = "ACTIVE"
	// StatusSignedout holds the string of SIGNED_OUT
	StatusSignedout Status = "SIGNED_OUT"
)

// Session is a struct which holds all the information about each session.
type Session struct {
	ID               primitive.ObjectID    `bson:"_id"`
	User             *primitive.ObjectID   `bson:"user"`
	Token            string                `bson:"token" json:"-"`
	IPAddress        string                `bson:"ip"`
	UserAgent        string                `bson:"user_agent"`
	DeviceInfo       deviceinfo.DeviceInfo `bson:"device_info"`
	Location         location.Location     `bson:"location"`
	Status           Status                `bson:"status"`
	LastActivityTime time.Time             `bson:"last_activity_time"`
	CreatedAt        time.Time             `bson:"created_at"`
	CurrentSession   bool                  `bson:"-"`
}

// GenerateID creates new id
func (s *Session) GenerateID() string {
	s.ID = primitive.NewObjectID()
	return s.ID.Hex()
}

// GetID is used to transform primitive type into a string
func (s *Session) GetID() string {
	return s.ID.Hex()
}

// SetID is used to set the recieved string as the session ID which is in primitive type
func (s *Session) SetID(id string) error {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	s.ID = obj
	return nil
}

// GetUser is used to transform primitive type into a string
func (s Session) GetUser() string {
	if s.User != nil {
		return s.User.Hex()
	}
	return ""
}

// SetUser is used to set the recieved string as the session user which is in primitive type
func (s *Session) SetUser(user string) error {
	obj, err := primitive.ObjectIDFromHex(user)
	if err != nil {
		return err
	}

	s.User = &obj
	return nil
}
