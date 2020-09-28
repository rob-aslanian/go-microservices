package account

// Permission represent a permission of some field or privacy
type Permission struct {
	Type PermissionType `bson:"type"`
}

// PermissionType ...
type PermissionType string

const (
	// PermissionTypeNone ...
	PermissionTypeNone PermissionType = ""

	// PermissionTypeMe means only me can see
	PermissionTypeMe PermissionType = "ME"

	// PermissionTypeMyConnections means can see me and my connenctions
	PermissionTypeMyConnections PermissionType = "MY_CONNECTIONS"

	// PermissionTypeMembers means can see everybody who registered
	PermissionTypeMembers PermissionType = "MEMBERS"

	// PermissionTypeEveryone means can see anybody
	PermissionTypeEveryone PermissionType = "EVERYONE"
)
