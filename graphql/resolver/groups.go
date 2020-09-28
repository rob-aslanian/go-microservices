package resolver

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/groupsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

// RegisterGroup ...
func (r Resolver) RegisterGroup(ctx context.Context, input RegisterGroupRequest) (*SuccessResolver, error) {
	resp, err := groups.RegisterGroup(ctx, &groupsRPC.RegisterGroupRequest{
		Name:        input.Input.Name,
		Type:        input.Input.Type,
		PrivacyType: stringToGroupsRPCGroupPrivacyType(input.Input.Privacy_type),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			ID:      resp.GetID(),
			Success: true,
		},
	}, nil
}

// ChangeGroupTagline ...
func (r Resolver) ChangeGroupTagline(ctx context.Context, input ChangeGroupTaglineRequest) (*SuccessResolver, error) {
	_, err := groups.ChangeTagline(ctx, &groupsRPC.ChangeTaglineRequest{
		ID:      input.Group_id,
		Tagline: input.Tagline,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ChangeGroupDescription ...
func (r Resolver) ChangeGroupDescription(ctx context.Context, input ChangeGroupDescriptionRequest) (*SuccessResolver, error) {
	_, err := groups.ChangeGroupDescription(ctx, &groupsRPC.ChangeGroupDescriptionRequest{
		ID:          input.Group_id,
		Description: input.Description.Description,
		Rules:       input.Description.Rules,
		Location:    locationInputToGroupsRPCLocation(input.Description.Location),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ChangeGroupName ...
func (r Resolver) ChangeGroupName(ctx context.Context, input ChangeGroupNameRequest) (*SuccessResolver, error) {
	_, err := groups.ChangeGroupName(ctx, &groupsRPC.ChangeGroupNameRequest{
		ID:   input.Group_id,
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ChangeGroupPrivacyType ...
func (r Resolver) ChangeGroupPrivacyType(ctx context.Context, input ChangeGroupPrivacyTypeRequest) (*SuccessResolver, error) {
	_, err := groups.ChangeGroupPrivacyType(ctx, &groupsRPC.ChangeGroupPrivacyTypeRequest{
		ID:   input.Group_id,
		Type: stringToGroupsRPCGroupPrivacyType(input.Type),
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ChangeGroupURL ...
func (r Resolver) ChangeGroupURL(ctx context.Context, input ChangeGroupURLRequest) (*SuccessResolver, error) {
	_, err := groups.ChangeGroupURL(ctx, &groupsRPC.ChangeGroupURLRequest{
		ID:  input.Group_id,
		URL: input.Url,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// AddAdminInGroup ...
func (r Resolver) AddAdminInGroup(ctx context.Context, input AddAdminInGroupRequest) (*SuccessResolver, error) {
	_, err := groups.AddAdmin(ctx, &groupsRPC.AddAdminRequest{
		ID:     input.Group_id,
		UserID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// IsGroupURLBusy ...
func (r Resolver) IsGroupURLBusy(ctx context.Context, input IsGroupURLBusyRequest) (bool, error) {
	resp, err := groups.IsGroupURLBusy(ctx, &groupsRPC.URL{
		URL: input.Url,
	})
	if err != nil {
		return false, err
	}

	return resp.GetValue(), nil
}

// JoinGroup ...
func (r Resolver) JoinGroup(ctx context.Context, input JoinGroupRequest) (*SuccessResolver, error) {
	_, err := groups.JoinGroup(ctx, &groupsRPC.ID{
		ID: input.Group_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// LeaveGroup ...
func (r Resolver) LeaveGroup(ctx context.Context, input LeaveGroupRequest) (*SuccessResolver, error) {
	_, err := groups.LeaveGroup(ctx, &groupsRPC.ID{
		ID: input.Group_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// RemoveMemberFromGroup ...
func (r Resolver) RemoveMemberFromGroup(ctx context.Context, input RemoveMemberFromGroupRequest) (*SuccessResolver, error) {
	_, err := groups.RemoveMemberFromGroup(ctx, &groupsRPC.RemoveMemberFromGroupRequest{
		ID:     input.Group_id,
		UserID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// SentInvitations ...
func (r Resolver) SentInvitations(ctx context.Context, input SentInvitationsRequest) (*SuccessResolver, error) {
	_, err := groups.SentInvitations(ctx, &groupsRPC.SentInvitationsRequest{
		ID:     input.Group_id,
		UserID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// AcceptInvitation ...
func (r Resolver) AcceptInvitation(ctx context.Context, input AcceptInvitationRequest) (*SuccessResolver, error) {
	_, err := groups.AcceptInvitation(ctx, &groupsRPC.ID{
		ID: input.Group_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// DeclineInvitation ...
func (r Resolver) DeclineInvitation(ctx context.Context, input DeclineInvitationRequest) (*SuccessResolver, error) {
	_, err := groups.DeclineInvitation(ctx, &groupsRPC.ID{
		ID: input.Group_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// SentJoinRequest ...
func (r Resolver) SentJoinRequest(ctx context.Context, input SentJoinRequestRequest) (*SuccessResolver, error) {
	_, err := groups.SentJoinRequest(ctx, &groupsRPC.ID{
		ID: input.Group_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// ApproveJoinRequest ...
func (r Resolver) ApproveJoinRequest(ctx context.Context, input ApproveJoinRequestRequest) (*SuccessResolver, error) {
	_, err := groups.ApproveJoinRequest(ctx, &groupsRPC.ApproveJoinRequestRequest{
		ID:     input.Group_id,
		UserID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// DeclineJoinRequest ...
func (r Resolver) DeclineJoinRequest(ctx context.Context, input DeclineJoinRequestRequest) (*SuccessResolver, error) {
	_, err := groups.DeclineJoinRequest(ctx, &groupsRPC.DeclineJoinRequestRequest{
		ID:     input.Group_id,
		UserID: input.User_id,
	})
	if err != nil {
		return nil, err
	}

	return &SuccessResolver{
		R: &Success{
			Success: true,
		},
	}, nil
}

// GetGroupByURL ...
func (r Resolver) GetGroupByURL(ctx context.Context, input GetGroupByURLRequest) (*GroupResolver, error) {
	resp, err := groups.GetGroupByURL(ctx, &groupsRPC.URL{
		URL: input.Url,
	})
	if err != nil {
		return nil, err
	}

	gr := groupsRPCGroupToGroup(resp)
	if gr == nil {
		return nil, nil
	}

	return &GroupResolver{
		R: gr,
	}, nil
}

// GetMembersOfGroup ...
func (r Resolver) GetMembersOfGroup(ctx context.Context, input GetMembersOfGroupRequest) (*[]ProfileResolver, error) {
	members, err := groups.GetMembers(ctx, &groupsRPC.GetMembersRequest{
		ID: input.Group_id,
		Pagination: &groupsRPC.Pagination{
			After: NullToString(input.Pagination.After),
			First: uint32(NullToInt32(input.Pagination.First)),
		},
	})
	if err != nil {
		return nil, err
	}

	membersIDs := make([]string, 0, len(members.GetMembers()))
	for i := range members.GetMembers() {
		membersIDs = append(membersIDs, members.GetMembers()[i].GetUserID())
	}

	resp, err := user.GetProfilesByID(ctx, &userRPC.UserIDs{
		ID: membersIDs,
	})
	if err != nil {
		return nil, err
	}

	profiles := make([]ProfileResolver, 0, len(resp.GetProfiles()))

	for i := range resp.GetProfiles() {
		pr := ToProfile(ctx, resp.GetProfiles()[i])
		profiles = append(profiles, ProfileResolver{
			R: &pr,
		})
	}

	return &profiles, nil
}

// ---------

func stringToGroupsRPCGroupPrivacyType(data string) groupsRPC.GroupPrivacyType {
	switch data {
	case "closed":
		return groupsRPC.GroupPrivacyType_Closed
	case "secret":
		return groupsRPC.GroupPrivacyType_Secret
	}

	return groupsRPC.GroupPrivacyType_Public
}

func locationInputToGroupsRPCLocation(data *LocationInput) *groupsRPC.Location {
	if data == nil {
		return nil
	}

	loc := groupsRPC.Location{
		CityID:          NullToString(data.City.ID),
		CityName:        NullToString(data.City.City),
		CitySubdivision: NullToString(data.City.Subdivision),
		CountryID:       data.Country_id,
	}

	return &loc
}

// ---------

func groupsRPCGroupToGroup(data *groupsRPC.Group) *Group {
	if data == nil {
		return nil
	}

	gr := Group{
		ID:                data.GetID(),
		Url:               data.GetURL(),
		Name:              data.GetName(),
		Type:              data.GetType(),
		Privacy_type:      groupsRPCGroupPrivacyTypeToString(data.GetPrivacyType()),
		Amount_of_members: int32(data.GetAmountOfMembers()),
		Tagline:           data.GetTagLine(),
		Description:       data.GetDescription(),
		Rules:             data.GetRules(),
		Location:          groupsRPCLocationToLocationInput(data.GetLocation()),
		// Owner: data.GetOwnerID()
		Owner:        &Profile{},
		Created_at:   data.GetCreatedAt(),
		Cover:        data.GetCover(),
		Cover_origin: data.GetOriginCover(),
	}

	return &gr
}

func groupsRPCGroupPrivacyTypeToString(data groupsRPC.GroupPrivacyType) string {
	switch data {
	case groupsRPC.GroupPrivacyType_Closed:
		return "closed"
	case groupsRPC.GroupPrivacyType_Secret:
		return "secret"
	}

	return "public"
}

func groupsRPCLocationToLocationInput(data *groupsRPC.Location) *City {
	if data == nil {
		return &City{}
		// return nil
	}

	loc := City{
		City:        data.GetCityName(),
		Country:     data.GetCountryID(),
		ID:          data.GetCityID(),
		Subdivision: data.GetCitySubdivision(),
	}

	return &loc
}
