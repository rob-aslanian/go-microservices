package serverRPC

import (
	"context"
	"strconv"
	"time"

	"gitlab.lan/Rightnao-site/microservices/groups/internal/group"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/location"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/groupsRPC"
)

// RegisterGroup ...
func (s Server) RegisterGroup(ctx context.Context, data *groupsRPC.RegisterGroupRequest) (*groupsRPC.ID, error) {
	id, err := s.service.CreateGroup(ctx, &group.Group{
		Name:    data.GetName(),
		Type:    data.GetType(),
		Privacy: groupsRPCPrivacyTypeToGroupPrivacyType(data.GetPrivacyType()),
	})
	if err != nil {
		return nil, err
	}

	return &groupsRPC.ID{
		ID: id,
	}, nil
}

// ChangeTagline ...
func (s Server) ChangeTagline(ctx context.Context, data *groupsRPC.ChangeTaglineRequest) (*groupsRPC.Empty, error) {
	err := s.service.ChangeTagline(ctx, data.GetID(), data.GetTagline())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// ChangeGroupDescription ...
func (s Server) ChangeGroupDescription(ctx context.Context, data *groupsRPC.ChangeGroupDescriptionRequest) (*groupsRPC.Empty, error) {
	err := s.service.ChangeGroupDescription(
		ctx,
		data.GetID(),
		data.GetDescription(),
		data.GetRules(),
		groupsRPCLocationToLocation(data.GetLocation()),
	)
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// ChangeGroupName ...
func (s Server) ChangeGroupName(ctx context.Context, data *groupsRPC.ChangeGroupNameRequest) (*groupsRPC.Empty, error) {
	err := s.service.ChangeGroupName(
		ctx,
		data.GetID(),
		data.GetName(),
	)
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// ChangeGroupPrivacyType ...
func (s Server) ChangeGroupPrivacyType(ctx context.Context, data *groupsRPC.ChangeGroupPrivacyTypeRequest) (*groupsRPC.Empty, error) {
	err := s.service.ChangeGroupPrivacyType(
		ctx,
		data.GetID(),
		groupsRPCPrivacyTypeToGroupPrivacyType(data.GetType()),
	)
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// IsGroupURLBusy ...
func (s Server) IsGroupURLBusy(ctx context.Context, data *groupsRPC.URL) (*groupsRPC.BooleanValue, error) {
	isBusy, err := s.service.IsURLBusy(ctx, data.GetURL())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.BooleanValue{
		Value: isBusy,
	}, nil
}

// ChangeGroupURL ...
func (s Server) ChangeGroupURL(ctx context.Context, data *groupsRPC.ChangeGroupURLRequest) (*groupsRPC.Empty, error) {
	err := s.service.ChangeGroupURL(
		ctx,
		data.GetID(),
		data.GetURL(),
	)
	if err != nil {
		return nil, err
	}
	return &groupsRPC.Empty{}, nil
}

// AddAdmin ...
func (s Server) AddAdmin(ctx context.Context, data *groupsRPC.AddAdminRequest) (*groupsRPC.Empty, error) {
	err := s.service.AddAdmin(ctx, data.GetID(), data.GetUserID())
	if err != nil {
		return nil, err
	}
	return &groupsRPC.Empty{}, nil
}

// JoinGroup ...
func (s Server) JoinGroup(ctx context.Context, data *groupsRPC.ID) (*groupsRPC.Empty, error) {
	err := s.service.JoinGroup(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// LeaveGroup ...
func (s Server) LeaveGroup(ctx context.Context, data *groupsRPC.ID) (*groupsRPC.Empty, error) {
	err := s.service.LeaveGroup(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// RemoveMemberFromGroup ...
func (s Server) RemoveMemberFromGroup(ctx context.Context, data *groupsRPC.RemoveMemberFromGroupRequest) (*groupsRPC.Empty, error) {
	err := s.service.RemoveFromGroup(ctx, data.GetID(), data.GetUserID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// SentInvitations ...
func (s Server) SentInvitations(ctx context.Context, data *groupsRPC.SentInvitationsRequest) (*groupsRPC.Empty, error) {
	err := s.service.SentInvitations(ctx, data.GetID(), data.GetUserID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// AcceptInvitation ...
func (s Server) AcceptInvitation(ctx context.Context, data *groupsRPC.ID) (*groupsRPC.Empty, error) {
	err := s.service.AcceptInvitation(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// DeclineInvitation ...
func (s Server) DeclineInvitation(ctx context.Context, data *groupsRPC.ID) (*groupsRPC.Empty, error) {
	err := s.service.DeclineInvitation(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// SentJoinRequest ...
func (s Server) SentJoinRequest(ctx context.Context, data *groupsRPC.ID) (*groupsRPC.Empty, error) {
	err := s.service.SentJoinRequest(ctx, data.GetID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// ApproveJoinRequest ...
func (s Server) ApproveJoinRequest(ctx context.Context, data *groupsRPC.ApproveJoinRequestRequest) (*groupsRPC.Empty, error) {
	err := s.service.ApproveJoinRequest(ctx, data.GetID(), data.GetUserID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// DeclineJoinRequest ...
func (s Server) DeclineJoinRequest(ctx context.Context, data *groupsRPC.DeclineJoinRequestRequest) (*groupsRPC.Empty, error) {
	err := s.service.DeclineJoinRequest(ctx, data.GetID(), data.GetUserID())
	if err != nil {
		return nil, err
	}

	return &groupsRPC.Empty{}, nil
}

// GetGroupByURL ...
func (s Server) GetGroupByURL(ctx context.Context, data *groupsRPC.URL) (*groupsRPC.Group, error) {
	gr, err := s.service.GetGroupByURL(ctx, data.GetURL())
	if err != nil {
		return nil, err
	}

	return groupGroupToGroupsRPCGroup(gr), nil
}

// GetMembers ...
func (s Server) GetMembers(ctx context.Context, data *groupsRPC.GetMembersRequest) (*groupsRPC.Members, error) {
	var /*first,*/ after uint32

	if data.GetPagination() != nil {
		// f, err := strconv.Atoi(data.GetPagination().GetFirst())
		// if err == nil {
		// 	first = uint32(f)
		// }

		a, err := strconv.Atoi(data.GetPagination().GetAfter())
		if err == nil {
			after = uint32(a)
		}
	}

	members, err := s.service.GetMembers(
		ctx,
		data.GetID(),
		data.GetPagination().GetFirst(),
		after,
	)
	if err != nil {
		return nil, err
	}

	m := make([]*groupsRPC.Member, 0, len(members))

	for i := range members {
		m = append(m, groupMemberToMemberRPC(members[i]))
	}

	return &groupsRPC.Members{
		Members: m,
	}, nil
}

// -----

func groupsRPCPrivacyTypeToGroupPrivacyType(data groupsRPC.GroupPrivacyType) group.PrivacyType {
	switch data {
	case groupsRPC.GroupPrivacyType_Closed:
		return group.PrivacyTypeClosed
	case groupsRPC.GroupPrivacyType_Secret:
		return group.PrivacyTypeSecret

	}
	return group.PrivacyTypePublic
}

func groupsRPCLocationToLocation(data *groupsRPC.Location) *location.Location {
	if data == nil {
		return nil
	}

	loc := location.Location{
		CountryID: data.GetCountryID(),
	}

	if data.GetCityID() != "" {
		id := data.GetCityID()
		loc.CityID = &id
	}

	return &loc
}

// -----

func groupGroupToGroupsRPCGroup(data *group.Group) *groupsRPC.Group {
	if data == nil {
		return nil
	}

	gr := groupsRPC.Group{
		ID:              data.GetID(),
		Cover:           data.Cover,
		OriginCover:     data.CoverOriginal,
		Name:            data.Name,
		Description:     data.Description,
		CreatedAt:       timeToString(data.CreatedAt),
		OwnerID:         data.GetOwnerID(),
		Rules:           data.Rules,
		TagLine:         data.Tagline,
		Type:            data.Type,
		PrivacyType:     groupPrivacyTypeToGroupsRPCPrivacyType(data.Privacy),
		Location:        locationToGroupsRPCLocation(data.Location),
		AmountOfMembers: data.AmountOfMembers,
		URL:             data.URL,
	}

	return &gr
}

func groupPrivacyTypeToGroupsRPCPrivacyType(data group.PrivacyType) groupsRPC.GroupPrivacyType {
	switch data {
	case group.PrivacyTypeClosed:
		return groupsRPC.GroupPrivacyType_Closed
	case group.PrivacyTypeSecret:
		return groupsRPC.GroupPrivacyType_Secret
	}

	return groupsRPC.GroupPrivacyType_Public
}

func locationToGroupsRPCLocation(data *location.Location) *groupsRPC.Location {
	if data == nil {
		return nil
	}

	loc := groupsRPC.Location{
		CountryID:       data.CountryID,
		CityName:        data.CityName,
		CitySubdivision: data.CitySubdivision,
	}

	if data.CityID != nil {
		loc.CityID = *data.CityID
	}

	return &loc
}

func groupMemberToMemberRPC(data *group.Member) *groupsRPC.Member {
	if data == nil {
		return nil
	}

	m := groupsRPC.Member{
		UserID:    data.GetID(),
		CreatedAt: timeToString(data.CreatedAt),
	}

	if data.IsAdmin != nil {
		m.IsAdmin = *data.IsAdmin
	}

	return &m
}

// -----

func timeToString(s time.Time) string {
	return s.Format(time.RFC3339)
}
