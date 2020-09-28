package comingsoon

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// ComingSoon ... 
type ComingSoon struct {
	UserID       bson.ObjectId      `bson:"user_id"`
	Email		 string				`bson:"email"`
	Type		 string				`bson:"type"`
	CreatedAt    time.Time          `bson:"created_at"`
}


// SetUserID ...
func (f *ComingSoon) SetUserID(id string) {
	var userID bson.ObjectId

	if bson.IsObjectIdHex(id) {
		userID = bson.ObjectIdHex(id)
		f.UserID = userID
	}
}