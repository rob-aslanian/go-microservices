package account

import (
	"math/rand"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	usersErrors "gitlab.lan/Rightnao-site/microservices/user/pkg/internal/users-errors"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/status"
)

// Account represents user account
type Account struct {
	ID             bson.ObjectId     `bson:"_id"`
	Status         status.UserStatus `bson:"status"`
	Username       string            `bson:"username"`
	FirstName      string            `bson:"first_name"`
	Lastname       string            `bson:"last_name"`
	Patronymic     *Patronymic       `bson:"patronymic,omitempty"`
	MiddleName     *MiddleName       `bson:"middlename,omitempty"`
	NativeName     *NativeName       `bson:"native_name,omitempty"`
	Nickname       *Nickname         `bson:"nickname,omitempty"`
	Birthday       *Birthday         `bson:"birthday"`
	Gender         Gender            `bson:"gender"`
	Emails         []*Email          `bson:"emails"`
	Phones         []*Phone          `bson:"phones"`
	MyAddresses    []*MyAddress      `bson:"my_addresses"`
	OtherAddresses []*OtherAddress   `bson:"other_addresses"`
	LanguageID     string            `bson:"language_id"`
	Location       *UserLocation     `bson:"location"`
	// Sessions       []Sessions      `bson:"-"`
	URL          string         `bson:"url"`
	Notification *Notifications `bson:"notification"`
	Privacy      *Privacy       `bson:"privacy"`
	CreatedAt    time.Time      `bson:"created_at"`

	BirthdayDate Date           `bson:"birthday_date"`
	InvitedBy    *bson.ObjectId `bson:"invited_by"`

	LastChangePassword time.Time `bson:"last_change_password"`
	AmountOfSessions   int32     `bson:"-"`
}

// User ...
type User struct {
	ID              bson.ObjectId     `bson:"id"`
	Status          status.UserStatus `bson:"status"`
	Avatar          string            `bson:"avatar"`
	FirstName       string            `bson:"first_name"`
	Lastname        string            `bson:"last_name"`
	Birthday        *Birthday         `bson:"birthday"`
	Gender          *string           `bson:"gender"`
	Email           *string           `bson:"email"`
	Phones          []*Phone          `bson:"phones"`
	Location        *UserLocation     `bson:"location"`
	URL             string            `bson:"url"`
	CreatedAt       time.Time         `bson:"created_at"`
	CompletePercent int8              `bson:"-"`
	BirthdayDate    Date              `bson:"birthday_date"`
}

// UserForAdvert ...
type UserForAdvert struct {
	OwnerID   string
	Gender    string
	AgeFrom   int32
	AgeTo     int32
	Locations []string
	Languages []string
}

// GetUserID ...
func (u *User) GetUserID() string {
	return u.ID.Hex()
}

// GetID returns id of user
func (a *Account) GetID() string {
	return a.ID.Hex()
}

// SetID saves id of user. If id has a wrong format returns usersErrors.ErrWrongID error.
func (a *Account) SetID(id string) error {
	if bson.IsObjectIdHex(id) {
		a.ID = bson.ObjectIdHex(id)
		return nil
	}
	return usersErrors.ErrWrongID
}

// GenerateID creates new random id for account
func (a *Account) GenerateID() string {
	id := bson.NewObjectId()
	a.ID = id
	return id.Hex()
}

// GenerateURL ...
func (a *Account) GenerateURL() string {

	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	url := strings.Join([]string{a.Lastname, a.FirstName, string(b)}, "_")
	a.URL = url
	return url
}

// GetInvitedByID ...
func (a *Account) GetInvitedByID() string {
	if a.InvitedBy == nil {
		return ""
	}
	return a.InvitedBy.Hex()
}

// SetInvitedByID ...
func (a *Account) SetInvitedByID(id string) error {
	if bson.IsObjectIdHex(id) {
		objID := bson.ObjectIdHex(id)
		a.InvitedBy = &objID
		return nil
	}
	return usersErrors.ErrWrongID
}
