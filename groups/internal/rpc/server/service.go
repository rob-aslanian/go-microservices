package serverRPC

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/groups/internal/group"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/location"
)

// Service define functions inside Service
type Service interface {
	CreateGroup(ctx context.Context, gr *group.Group) (string, error)
	ChangeTagline(ctx context.Context, groupID string, tagline string) error
	ChangeGroupDescription(ctx context.Context, groupID string, desc, rules string, loc *location.Location) error
	ChangeGroupName(ctx context.Context, groupID string, name string) error
	ChangeGroupPrivacyType(ctx context.Context, groupID string, privacyType group.PrivacyType) error
	IsURLBusy(ctx context.Context, url string) (bool, error)
	ChangeGroupURL(ctx context.Context, groupID string, url string) error
	AddAdmin(ctx context.Context, groupID string, id string) error
	JoinGroup(ctx context.Context, groupID string) error
	LeaveGroup(ctx context.Context, groupID string) error
	RemoveFromGroup(ctx context.Context, groupID string, userID string) error
	SentInvitations(ctx context.Context, groupID string, userIDs []string) error
	AcceptInvitation(ctx context.Context, groupID string) error
	DeclineInvitation(ctx context.Context, groupID string) error
	SentJoinRequest(ctx context.Context, groupID string) error
	ApproveJoinRequest(ctx context.Context, groupID, userID string) error
	DeclineJoinRequest(ctx context.Context, groupID, userID string) error
	GetGroupByURL(ctx context.Context, url string) (*group.Group, error)
	GetMembers(ctx context.Context, groupID string, first, after uint32) ([]*group.Member, error)
}
