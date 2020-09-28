package job

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Application ...
type Application struct {
	UserID      primitive.ObjectID `bson:"user_id"`
	Email       string
	Phone       string
	CoverLetter string
	Documents   []*File
	CreatedAt   time.Time `bson:"created_at"`

	Metadata ApplicationMeta
}

// GetUserID returns user id
func (p Application) GetUserID() string {
	return p.UserID.Hex()
}

// SetUserID set user id
func (p *Application) SetUserID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	p.UserID = objID
	return nil
}
