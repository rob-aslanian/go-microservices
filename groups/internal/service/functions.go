package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	companyadmin "gitlab.lan/Rightnao-site/microservices/groups/internal/company-admin"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/group"
	"gitlab.lan/Rightnao-site/microservices/groups/internal/location"
	"google.golang.org/grpc/metadata"
)

// GetGroupByURL ...
func (s Service) GetGroupByURL(ctx context.Context, url string) (*group.Group, error) {
	span := s.tracer.MakeSpan(ctx, "GetGroupByURL")
	defer span.Finish()

	gr, err := s.repository.GetGroupByURL(ctx, url)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	// TODO: return secret without invitation for not members?

	return gr, nil
}

// CreateGroup ...
func (s Service) CreateGroup(ctx context.Context, gr *group.Group) (string, error) {
	span := s.tracer.MakeSpan(ctx, "CreateGroup")
	defer span.Finish()

	err := gr.ValidateName()
	if err != nil {
		return "", err
	}

	id := gr.GenerateID()
	gr.CreatedAt = time.Now()
	gr.URL = id

	gr.Trim()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = gr.SetOwnerID(userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}
	err = gr.AddMember(userID, true)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	err = s.repository.SaveGroup(ctx, gr)
	if err != nil {
		s.tracer.LogError(span, err)
		return "", err
	}

	return id, nil
}

// ChangeTagline ...
func (s Service) ChangeTagline(ctx context.Context, groupID string, tagline string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeTagline")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	err = s.repository.ChangeTagline(ctx, groupID, tagline)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeGroupDescription ...
func (s Service) ChangeGroupDescription(ctx context.Context, groupID string, desc, rules string, loc *location.Location) error {
	span := s.tracer.MakeSpan(ctx, "ChangeGroupDescription")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	if loc != nil && loc.CityID != nil {
		id, _ := strconv.Atoi(*loc.CityID)
		_, _, countryID, err := s.infoRPC.GetCityInformationByID(ctx, int32(id), nil)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}

		loc.CountryID = countryID
	}

	// TODO: validate/trim desc and rules

	err = s.repository.ChangeGroupDescription(ctx, groupID, desc, rules, loc)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ChangeGroupName ...
func (s Service) ChangeGroupName(ctx context.Context, groupID string, name string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeGroupName")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	err = group.ValidateName(name)
	if err != nil {
		return err
	}

	err = s.repository.ChangeGroupName(ctx, groupID, name)
	if err != nil {
		return err
	}

	return nil
}

// ChangeGroupPrivacyType ...
func (s Service) ChangeGroupPrivacyType(ctx context.Context, groupID string, privacyType group.PrivacyType) error {
	span := s.tracer.MakeSpan(ctx, "ChangeGroupPrivacyType")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	err = s.repository.ChangeGroupPrivacyType(ctx, groupID, privacyType)
	if err != nil {
		return err
	}

	return nil
}

// IsURLBusy returns true if given url is already taken by another group.
func (s Service) IsURLBusy(ctx context.Context, url string) (bool, error) {
	span := s.tracer.MakeSpan(ctx, "ChangeGroupPrivacyType")
	defer span.Finish()

	isBusy, err := s.repository.IsURLBusy(ctx, url)
	if err != nil {
		return false, err
	}

	return isBusy, nil
}

// ChangeGroupURL ...
func (s Service) ChangeGroupURL(ctx context.Context, groupID string, url string) error {
	span := s.tracer.MakeSpan(ctx, "ChangeGroupPrivacyType")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	isBusy, err := s.IsURLBusy(ctx, url)
	if err != nil {
		return err
	}
	if isBusy {
		return errors.New("already_taken")
	}

	err = s.repository.ChangeGroupURL(ctx, groupID, url)
	if err != nil {
		return err
	}

	return nil
}

// AddAdmin ...
func (s Service) AddAdmin(ctx context.Context, groupID string, id string) error {
	span := s.tracer.MakeSpan(ctx, "AddAdmin")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	// check if member
	isMember, err := s.repository.IsMember(ctx, groupID, id)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("not_a_member")
	}

	err = s.repository.AddAdmin(ctx, groupID, id)
	if err != nil {
		return err
	}

	return nil
}

// GetMembers ...
func (s Service) GetMembers(ctx context.Context, groupID string, first, after uint32) ([]*group.Member, error) {
	span := s.tracer.MakeSpan(ctx, "GetMembers")
	defer span.Finish()

	// if group is secret, requestor should be a member of group
	// check if secret group
	gr, err := s.getGroupByID(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	if gr.Privacy == group.PrivacyTypeSecret {
		token := s.retriveToken(ctx)
		userID, err := s.authRPC.GetUser(ctx, token)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}
		// check if not a member
		isMember, err := s.repository.IsMember(ctx, groupID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, nil
		}
	}

	members, err := s.repository.GetMembers(ctx, groupID, first, after)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// JoinGroup join group. Only for public groups.
func (s Service) JoinGroup(ctx context.Context, groupID string) error {
	span := s.tracer.MakeSpan(ctx, "JoinGroup")
	defer span.Finish()

	// check if public
	gr, err := s.getGroupByID(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if gr.Privacy != group.PrivacyTypePublic {
		return errors.New("not_allowed")
	}

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if not a member
	isMember, err := s.repository.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("already_joined")
	}

	m := group.Member{
		CreatedAt: time.Now(),
	}
	m.SetID(userID)

	err = s.repository.AddToMembers(ctx, groupID, m)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// LeaveGroup ...
func (s Service) LeaveGroup(ctx context.Context, groupID string) error {
	span := s.tracer.MakeSpan(ctx, "LeaveGroup")
	defer span.Finish()

	// TODO: what to do with admin and owner???

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		return err
	}

	err = s.repository.LeaveGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFromGroup ...
func (s Service) RemoveFromGroup(ctx context.Context, groupID string, userID string) error {
	span := s.tracer.MakeSpan(ctx, "RemoveFromGroup")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	err = s.repository.LeaveGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}

	return nil
}

// ----------------------

// SentInvitations ...
func (s Service) SentInvitations(ctx context.Context, groupID string, userIDs []string) error {
	span := s.tracer.MakeSpan(ctx, "SentInvitation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	isMember, err := s.repository.IsMember(ctx, groupID, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isMember {
		return errors.New("not_allowed")
	}

	users := make([]group.InvitedMember, 0, len(userIDs))

	for _, id := range userIDs {
		// check if is not a member
		isMember, err := s.repository.IsMember(ctx, groupID, id)
		if err != nil {
			s.tracer.LogError(span, err)
			continue
		}
		if isMember {
			continue
		} else {
			m := group.InvitedMember{
				Member: group.Member{
					CreatedAt: time.Now(),
				},
			}

			m.SetID(id)
			m.SetInvitedByID(userID)
			users = append(users, m)
		}
	}

	// save in DB
	err = s.repository.AddInvitations(ctx, groupID, users)
	if err != nil {
		return err
	}

	// TODO: sent notification

	return nil
}

// AcceptInvitation ...
func (s Service) AcceptInvitation(ctx context.Context, groupID string) error {
	span := s.tracer.MakeSpan(ctx, "AcceptInvitation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check invitation
	isInvited, err := s.repository.IsInvited(ctx, groupID, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isInvited {
		return errors.New("not_invited")
	}

	// add to members
	m := group.Member{
		CreatedAt: time.Now(),
	}
	m.SetID(userID)

	err = s.repository.AddToMembers(ctx, groupID, m)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// remove invitation
	err = s.repository.RemoveInvitations(ctx, groupID, []string{userID})
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// DeclineInvitation ...
func (s Service) DeclineInvitation(ctx context.Context, groupID string) error {
	span := s.tracer.MakeSpan(ctx, "DeclineInvitation")
	defer span.Finish()

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// remove invitation
	err = s.repository.RemoveInvitations(ctx, groupID, []string{userID})
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ----------------------

// SentJoinRequest ...
func (s Service) SentJoinRequest(ctx context.Context, groupID string) error {
	span := s.tracer.MakeSpan(ctx, "SentJoinRequest")
	defer span.Finish()

	// check if closed group
	gr, err := s.getGroupByID(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if gr.Privacy != group.PrivacyTypeClosed {
		return errors.New("not_allowed")
	}

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// check if not a member
	isMember, err := s.repository.IsMember(ctx, groupID, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isMember {
		return errors.New("not_allowed")
	}

	// check if request wasn't send
	isSend, err := s.repository.IsRequestSend(ctx, groupID, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if isSend {
		return errors.New("not_allowed")
	}

	m := group.Member{
		CreatedAt: time.Now(),
	}
	m.SetID(userID)

	// save in DB
	err = s.repository.AddJoinRequest(ctx, groupID, m)
	if err != nil {
		return err
	}

	// TODO: sent notification

	return nil
}

// ApproveJoinRequest ...
func (s Service) ApproveJoinRequest(ctx context.Context, groupID, targetUserID string) error {
	span := s.tracer.MakeSpan(ctx, "ApproveJoinRequest")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	// check if request exists
	isSend, err := s.repository.IsRequestSend(ctx, groupID, targetUserID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isSend {
		return errors.New("not_allowed")
	}

	// delete request
	err = s.repository.RemoveInvitationRequest(ctx, groupID, targetUserID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	// add to members
	m := group.Member{
		CreatedAt: time.Now(),
	}
	m.SetID(targetUserID)

	err = s.repository.AddToMembers(ctx, groupID, m)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// DeclineJoinRequest ...
func (s Service) DeclineJoinRequest(ctx context.Context, groupID, userID string) error {
	span := s.tracer.MakeSpan(ctx, "DeclineJoinRequest")
	defer span.Finish()

	// check if admin
	isAdmin, err := s.isAdmin(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}
	if !isAdmin {
		return errors.New("not_allowed")
	}

	// remove request
	err = s.repository.RemoveInvitationRequest(ctx, groupID, userID)
	if err != nil {
		s.tracer.LogError(span, err)
		return err
	}

	return nil
}

// ----------------------

func (s Service) getGroupByID(ctx context.Context, groupID string) (*group.Group, error) {
	span := s.tracer.MakeSpan(ctx, "getGroupByID")
	defer span.Finish()

	return s.repository.GetGroupByID(ctx, groupID)
}

func (s Service) isAdmin(ctx context.Context, groupID string) (bool, error) {
	span := s.tracer.MakeSpan(ctx, "isAdmin")
	defer span.Finish()

	gr, err := s.getGroupByID(ctx, groupID)
	if err != nil {
		s.tracer.LogError(span, err)
		return false, err
	}

	if gr == nil {
		return false, nil
	}

	token := s.retriveToken(ctx)
	userID, err := s.authRPC.GetUser(ctx, token)
	if err != nil {
		s.tracer.LogError(span, err)
		return false, err
	}

	return gr.IsAdmin(userID)
}

func (s Service) retriveToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		arr := md.Get("token")
		if len(arr) > 0 {
			return arr[0]
		}
	}
	return ""
}

// checkAdminLevel return false if level doesn't much
func (s Service) checkAdminLevel(ctx context.Context, companyID string, requiredLevels ...companyadmin.AdminLevel) bool {
	span := s.tracer.MakeSpan(ctx, "checkAdminLevel")
	defer span.Finish()

	actualLevel, err := s.networkRPC.GetAdminLevel(ctx, companyID)
	if err != nil {
		s.tracer.LogError(span, err)
		log.Println("Error: checkAdminLevel:", err)
		return false
	}

	for _, lvl := range requiredLevels {
		if lvl == actualLevel {
			return true
		}
	}

	return false
}
