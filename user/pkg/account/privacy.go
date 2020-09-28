package account

// Privacy ...
type Privacy struct {
	FindByEmail    PermissionType `bson:"find_by_email"`
	FindByPhone    PermissionType `bson:"find_by_phone"`
	ActiveStatus   PermissionType `bson:"active_status"`
	ShareEdits     PermissionType `bson:"sharing_edits"`
	ProfilePicture PermissionType `bson:"profile_pictures"`
	MyConnections  PermissionType `bson:"my_connections"`
}

// PrivacyItem ...
type PrivacyItem string

const (
	// PrivacyItemFindByEmail who can find by emeail
	PrivacyItemFindByEmail PrivacyItem = "find_by_email"
	// PrivacyItemFindByPhone who can find by phone
	PrivacyItemFindByPhone PrivacyItem = "find_by_phone"
	// PrivacyItemActiveStatus who can see active status
	PrivacyItemActiveStatus PrivacyItem = "active_status"
	// PrivacyItemShareEdits share profile edits
	PrivacyItemShareEdits PrivacyItem = "sharing_edits"
	// PrivacyItemProfilePicture who can see avatar
	PrivacyItemProfilePicture PrivacyItem = "profile_pictures"
	// PrivacyItemMyConnections who can connections
	PrivacyItemMyConnections PrivacyItem = "my_connections"
)
