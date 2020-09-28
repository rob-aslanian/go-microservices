package feedback

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type FeedBack struct {
	ID           bson.ObjectId      `bson:"_id"`
	UserID       *bson.ObjectId     `bson:"user_id"`
	CompanyId    *bson.ObjectId     `bson:"company_id"`
	Reactions    FeedBackReaction   `bson:"reaction"`
	Compliment   FeedBackCompliment `bson:"compliment"`
	Complaint    FeedBackComplaint  `bson:"complaint"`
	Bugs         FeedBackBugs       `bson:"bugs"`
	CouldNotFind string             `bson:"could_not_find"`
	Suggestion   FeedBackSuggestion `bson:"suggestion"`
	Other        FeedBackOther      `bson:"other"`
	CreatedAt    time.Time          `bson:"created_at"`
}

// GenerateID ...
func (f *FeedBack) GenerateID() {
	f.ID = bson.NewObjectId()
}

// GetID ...
func (f FeedBack) GetID() string {
	return f.ID.Hex()
}

func (f *FeedBack) SetUserID(id string) {
	var userID bson.ObjectId

	if bson.IsObjectIdHex(id) {
		userID = bson.ObjectIdHex(id)
		f.UserID = new(bson.ObjectId)
		f.UserID = &userID
	}
}

func (f *FeedBack) SetCompanyID(id string) {
	var companyID bson.ObjectId

	if bson.IsObjectIdHex(id) {
		companyID = bson.ObjectIdHex(id)
		f.CompanyId = new(bson.ObjectId)
		f.CompanyId = &companyID
	}
}
