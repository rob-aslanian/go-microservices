package account

// Notifications ...
type Notifications struct {
	ConnectionRequest    bool `bson:"connection_request"`
	AcceptInvitation     bool `bson:"accept_invitation"`
	NewChatMessage       bool `bson:"new_chat_message"`
	Endorsements         bool `bson:"endorsements"`
	EmailUpdates         bool `bson:"email_updates"`
	NewFollowers         bool `bson:"new_followers"`
	Birthdays            bool `bson:"birthdays"`
	JobChangesInNetwork  bool `bson:"job_changes_in_network"`
	ImportContactsJoined bool `bson:"import_contacts_joined"`
	JobRecommendations   bool `bson:"job_recommendations"`
}
