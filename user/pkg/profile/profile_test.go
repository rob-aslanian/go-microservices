package profile

import (
	"testing"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/account"
)

func TestIsAllowed(t *testing.T) {
	tables := []struct {
		isMe       bool
		isFriend   bool
		isMember   bool
		Permission account.PermissionType
		Result     bool
	}{

		{
			Permission: account.PermissionTypeMe,
			isMe:       true,
			isMember:   false,
			isFriend:   false,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeMe,
			isMe:       false,
			isMember:   true,
			isFriend:   false,
			Result:     false,
		},
		{
			Permission: account.PermissionTypeMe,
			isMe:       false,
			isMember:   false,
			isFriend:   true,
			Result:     false,
		},
		{
			Permission: account.PermissionTypeMembers,
			isMe:       true,
			isMember:   false,
			isFriend:   false,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeMembers,
			isMe:       false,
			isMember:   true,
			isFriend:   false,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeMembers,
			isMe:       false,
			isMember:   false,
			isFriend:   true,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeMembers,
			isMe:       false,
			isMember:   false,
			isFriend:   false,
			Result:     false,
		},
		{
			Permission: account.PermissionTypeMyConnections,
			isMe:       true,
			isMember:   false,
			isFriend:   false,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeMyConnections,
			isMe:       false,
			isMember:   false,
			isFriend:   true,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeMyConnections,
			isMe:       false,
			isMember:   true,
			isFriend:   false,
			Result:     false,
		},
		{
			Permission: account.PermissionTypeEveryone,
			isMe:       false,
			isMember:   true,
			isFriend:   false,
			Result:     true,
		},
		{
			Permission: account.PermissionTypeEveryone,
			isMe:       true,
			isMember:   true,
			isFriend:   true,
			Result:     true,
		},
	}

	for _, ta := range tables {
		res := isAllowed(ta.isMe, ta.isFriend, ta.isMember, ta.Permission)
		if res != ta.Result {
			t.Errorf("Error: isAllowed(%v, %v, %v, %q) got %v, expected: %v", ta.isMe, ta.isFriend, ta.isMember, ta.Permission, res, ta.Result)
		}
	}
}
