package services

import (
	"context"
	"errors"
	"fmt"

	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/utils"
	"gitlab.lan/Rightnao-site/microservices/shared/hc-errors"
)

func (s *ChatService) AuthenticateUser(ctx context.Context) string {
	md := utils.ExtractMetadata(ctx, "token", "user-id")
	userId, ok := md["user-id"]
	if ok {
		return userId
	}

	token, ok := md["token"]
	if !ok {
		panic(hc_errors.NOT_AUTHENTICATED_ERROR)
	}

	senderId, err := s.auth.GetUserId(ctx, token)
	if err != nil {
		panic(err)
	}
	utils.AddToIncomingMetadata(ctx, "user-id", senderId)

	return senderId
}

func (s *ChatService) RequireAdminLevelForCompany(ctx context.Context, companyKey string, levels ...string) string {
	userId := s.AuthenticateUser(ctx)

	level, err := s.admin.GetAdminLevelFor(ctx, companyKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(level)

	for _, l := range levels {
		if l == level {
			return userId
		}
	}

	panic("You don't have required admin level for the company")
}

func (s *ChatService) getParticipation(ctx context.Context, conversationId string) (*model.Participant, *model.Conversation) {
	userId := s.AuthenticateUser(ctx)
	conversation, err := s.repo.GetConversation(ctx, conversationId)
	if err != nil {
		panic(err)
	}

	// check if sender is part of the conversation
	me := findParticipant(userId, conversation.Participants)
	if me == nil {
		panic(errors.New("You are not participant of the conversations"))
	}

	return me, conversation
}

func (s *ChatService) getParticipationForCompany(ctx context.Context, companyId, conversationId string) (*model.Participant, *model.Conversation) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	conversation, err := s.repo.GetConversation(ctx, conversationId)
	if err != nil {
		panic(err)
	}

	// check if sender is part of the conversation
	me := findParticipant(companyId, conversation.Participants)
	if me == nil {
		panic(errors.New("You are not participant of the conversations"))
	}

	return me, conversation
}

func clearDuplicates(list []*model.Participant) []*model.Participant {
	mapped := make(map[string]*model.Participant)
	for _, p := range list {
		mapped[p.Id] = p
	}

	uniques := make([]*model.Participant, 0, len(mapped))
	for _, value := range mapped {
		uniques = append(uniques, value)
	}
	return uniques
}

func findParticipant(id string, list []*model.Participant) *model.Participant {
	for _, p := range list {
		if p.Id == id {
			return p
		}
	}
	return nil
}

func (s *ChatService) updateConversationFor(participantId string, conversation *model.Conversation) *model.Conversation {
	participant := findParticipant(participantId, conversation.Participants)
	if participant != nil {
		conversation.Unread = participant.Unread
		conversation.Muted = participant.Muted
		conversation.Blocked = participant.Blocked
		conversation.Archived = participant.Archived
		conversation.HasLeft = participant.HasLeft
		conversation.Labels = participant.Labels
	}

	for _, p := range conversation.Participants {
		p.IsActive = s.liveConnections.IsLive(p.Id)
	}

	return conversation
}

func (s *ChatService) SendUnreadConversationsStatus(ctx context.Context, conversation *model.Conversation) {
	for _, participant := range conversation.Participants {
		if excludeBlockedAndMuted(participant) {
			count, err := s.repo.GetUnreadConversationsStatus(participant.Id)
			if err == nil {
				s.liveConnections.SendTo(ctx, participant.Id, model.NewTotalUnreadCountMessage(count))
			}
		}
	}
}
