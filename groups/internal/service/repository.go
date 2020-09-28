package service

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/groups/internal/group"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/location"
)

// Repository ...
type Repository interface {
	SaveGroup(ctx context.Context, gr *group.Group) error
	ChangeTagline(ctx context.Context, groupID string, tagline string) error
	ChangeGroupDescription(ctx context.Context, groupID string, desc, rules string, loc *location.Location) error
	ChangeGroupName(ctx context.Context, groupID string, name string) error
	ChangeGroupPrivacyType(ctx context.Context, groupID string, privacyType group.PrivacyType) error
	IsURLBusy(ctx context.Context, url string) (bool, error)
	ChangeGroupURL(ctx context.Context, groupID, url string) error
	AddAdmin(ctx context.Context, groupID, userID string) error
	AddToMembers(ctx context.Context, groupID string, m group.Member) error
	LeaveGroup(ctx context.Context, groupID, userID string) error
	GetGroupByID(ctx context.Context, groupID string) (*group.Group, error)
	GetGroupByURL(ctx context.Context, groupID string) (*group.Group, error)
	IsMember(ctx context.Context, groupID, userID string) (bool, error)
	AddInvitations(ctx context.Context, groupID string, users []group.InvitedMember) error
	IsInvited(ctx context.Context, groupID, userID string) (bool, error)
	RemoveInvitations(ctx context.Context, groupID string, userIDs []string) error
	AddJoinRequest(ctx context.Context, groupID string, users group.Member) error
	IsRequestSend(ctx context.Context, groupID, userID string) (bool, error)
	RemoveInvitationRequest(ctx context.Context, groupID, userID string) error
	GetMembers(ctx context.Context, groupID string, first, after uint32) ([]*group.Member, error)
}
