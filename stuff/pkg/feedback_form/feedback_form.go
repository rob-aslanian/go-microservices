package feedback

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

// Form represents form of feedback
type Form struct {
	ID        bson.ObjectId  `bson:"_id"`
	UserID    *bson.ObjectId `bson:"user_id"`
	Name      string         `bson:"name"`
	Email     string         `bson:"email"`
	Message   string         `bson:"message"`
	CreatedAt time.Time      `bson:"created_at"`
}

// GenerateID ...
func (f *Form) GenerateID() {
	f.ID = bson.NewObjectId()
}

// GetID ...
func (f Form) GetID() string {
	return f.ID.Hex()
}

// Validate ...
func (f Form) Validate() error {
	if f.Name == "" {
		return errors.New("empty name")
	}

	if f.Message == "" {
		return errors.New("empty message")
	}

	if f.Email == "" {
		return errors.New("empty email")
	}

	// TODO:

	return nil
}

// SetUserID ...
func (f *Form) SetUserID(id string) {
	var userID bson.ObjectId

	if bson.IsObjectIdHex(id) {
		userID = bson.ObjectIdHex(id)
		f.UserID = new(bson.ObjectId)
		f.UserID = &userID
	}
}
