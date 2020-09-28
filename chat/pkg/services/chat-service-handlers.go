package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
)

func (s *ChatService) CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	span := s.tracer.MakeSpan(ctx, "CreateConversation")
	defer span.Finish()

	senderId := s.AuthenticateUser(ctx)
	// add sender as a participant
	res, err := s.createConversationFromParticipant(ctx, &model.Participant{Id: senderId, IsAdmin: true}, conversation)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return res, nil
}

func (s *ChatService) CreateConversationForCompany(ctx context.Context, companyId string, conversation *model.Conversation) (*model.Conversation, error) {
	span := s.tracer.MakeSpan(ctx, "CreateConversationForCompany")
	defer span.Finish()

	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	// add sender as a participant
	res, err := s.createConversationFromParticipant(
		ctx,
		&model.Participant{
			Id:        companyId,
			IsCompany: true,
			IsAdmin:   true,
		},
		conversation,
	)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return res, nil
}

func (s *ChatService) CreateUnverifiedConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	return s.createConversationFromParticipant(ctx, conversation.Participants[0], conversation)
}

func (s *ChatService) createConversationFromParticipant(ctx context.Context, participant *model.Participant, conversation *model.Conversation) (*model.Conversation, error) {
	span := s.tracer.MakeSpan(ctx, "createConversationFromParticipant")
	defer span.Finish()

	conversation.Participants = append(conversation.Participants, participant)
	// remove duplicate values
	conversation.Participants = clearDuplicates(conversation.Participants)

	if len(conversation.Participants) < 2 {
		err := errors.New("You should choose someone to create conversation")
		s.tracer.LogError(span, err)
		return nil, err
	}

	// don't recreate private conversation if it already exists
	if len(conversation.Participants) == 2 {
		// check if sender is not blocked
		if conversation.Participants[0].IsCompany == false {
			isBlocked, err := s.admin.IsBlockedByUser(ctx, conversation.Participants[0].Id)
			if err != nil {
				return nil, err
			}
			if isBlocked {
				return nil, errors.New("you_can_not_write_to_this_person")
			}
		} else {
			// TODO:
			// isBlocked, err := s.admin.IsBlockedCompanyByUser(ctx, conversation.Participants[0].Id)
			// if err != nil {
			// 	return nil, err
			// }
			// if isBlocked {
			// 	return nil, errors.New("you_can_not_write_to_this_person")
			// }
		}

		// check if conversation already exists
		conversations, err := s.repo.FindConversationsWithOnlyParticipants(ctx, conversation.Participants)
		if err != nil {
			s.tracer.LogError(span, err)
			return nil, err
		}

		if len(conversations) > 0 {
			return s.updateConversationFor(participant.Id, conversations[0]), nil
		}
	} else {
		conversation.IsGroup = true
	}
	// conversation doesn't exist, create new one
	conversation, err := s.repo.CreateConversation(ctx, conversation)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	return s.updateConversationFor(participant.Id, conversation), nil
}

func (s *ChatService) GetConversation(ctx context.Context, conversationId string) (*model.Conversation, error) {
	span := s.tracer.MakeSpan(ctx, "GetConversation")
	defer span.Finish()

	me, conversation := s.getParticipation(ctx, conversationId)

	var messagesUntil int64 = math.MaxInt64
	if me.HasLeft {
		messagesUntil = me.LeftTimestamp
	}
	lastMessage, err := s.repo.GetLastMessage(ctx, conversationId, me.Id, me.DeletedTimestamp, messagesUntil)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}
	conversation.LastMessage = lastMessage

	return s.updateConversationFor(me.Id, conversation), nil
}

func (s *ChatService) GetConversationForCompany(ctx context.Context, companyId, conversationId string) (*model.Conversation, error) {
	span := s.tracer.MakeSpan(ctx, "GetConversationForCompany")
	defer span.Finish()

	me, conversation := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.updateConversationFor(me.Id, conversation), nil
}

func (s *ChatService) GetMyConversations(ctx context.Context, filter *model.ConversationFilter) ([]*model.Conversation, error) {
	span := s.tracer.MakeSpan(ctx, "GetMyConversations")
	defer span.Finish()

	senderId := s.AuthenticateUser(ctx)

	res, err := s.getConversationsOfParticipant(ctx, senderId, filter)
	if err != nil {
		s.tracer.LogError(span, err)
		return nil, err
	}

	return res, nil
}

func (s *ChatService) GetMyConversationsForCompany(ctx context.Context, companyId string, filter *model.ConversationFilter) ([]*model.Conversation, error) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")

	return s.getConversationsOfParticipant(ctx, companyId, filter)
}

func (s *ChatService) getConversationsOfParticipant(ctx context.Context, participantId string, filter *model.ConversationFilter) ([]*model.Conversation, error) {
	conversations, err := s.repo.GetConversationsOf(ctx, participantId, filter)
	if err != nil {
		return nil, err
	}

	var filteredConversationsByCategory []*model.Conversation
	if filter.Category == "Active" {
		for _, c := range conversations {
			for _, p := range c.Participants {
				if p.Id != participantId && s.liveConnections.IsLive(p.Id) {
					filteredConversationsByCategory = append(filteredConversationsByCategory, c)
				}
			}
		}
		filteredConversationsByCategory = conversations
	} else {
		filteredConversationsByCategory = conversations
	}

	var filteredConversationsWithoutEmpty []*model.Conversation
	for _, conversation := range filteredConversationsByCategory {
		if conversation.LastMessage != nil {
			filteredConversationsWithoutEmpty = append(filteredConversationsWithoutEmpty, conversation)
		}
	}

	for _, conversation := range filteredConversationsWithoutEmpty {
		s.updateConversationFor(participantId, conversation)
	}

	return filteredConversationsWithoutEmpty, nil
}

func (s *ChatService) GetMessages(ctx context.Context, conversationId string) ([]*model.Message, error) {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.getMessagesOfParticipant(ctx, me, conversationId)
}
func (s *ChatService) GetMessagesForCompany(ctx context.Context, companyId, conversationId string) ([]*model.Message, error) {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)
	return s.getMessagesOfParticipant(ctx, me, conversationId)
}

func (s *ChatService) GetActiveConnections(ctx context.Context) ([]string, error) {
	ids, err := s.admin.GetAllFriendshipsID(ctx)
	if err != nil {
		return []string{}, nil
	}

	onlineIDs := make([]string, 0, len(ids))
	for i := range ids {
		if s.liveConnections.IsLive(ids[i]) {
			onlineIDs = append(onlineIDs, ids[i])
		}
	}

	return onlineIDs, nil
}

func (s *ChatService) getMessagesOfParticipant(ctx context.Context, participant *model.Participant, conversationId string) ([]*model.Message, error) {
	var messagesUntil int64 = math.MaxInt64
	if participant.HasLeft {
		messagesUntil = participant.LeftTimestamp
	}

	return s.repo.GetMessages(ctx, conversationId, participant.DeletedTimestamp, messagesUntil)
}

func (s *ChatService) AddParticipantsToConversation(ctx context.Context, conversationId string, participants []string) (*model.Conversation, error) {
	me, conversation := s.getParticipation(ctx, conversationId)

	return s.addParticipantsToConversationOfParticipant(ctx, me, conversation, participants)
}
func (s *ChatService) AddParticipantsToConversationForCompany(ctx context.Context, companyId, conversationId string, participants []string) (*model.Conversation, error) {
	me, conversation := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.addParticipantsToConversationOfParticipant(ctx, me, conversation, participants)
}

func (s *ChatService) addParticipantsToConversationOfParticipant(ctx context.Context, participant *model.Participant, conversation *model.Conversation, participants []string) (*model.Conversation, error) {
	if participant.HasLeft {
		return nil, errors.New("You have already left the conversation")
	}

	newParticipants := make([]*model.Participant, len(participants))
	for i, p := range participants {
		newParticipants[i] = &model.Participant{Id: p}
	}

	var err error
	needsNewParticipantNotification := false
	if !conversation.IsGroup {
		participant.IsAdmin = true
		conversation = &model.Conversation{
			IsGroup:      true,
			Participants: clearDuplicates(append(conversation.Participants, newParticipants...)),
		}
		conversation, err = s.repo.CreateConversation(ctx, conversation)
	} else {
		for _, p := range newParticipants {
			existing := findParticipant(p.Id, conversation.Participants)
			if existing != nil {
				existing.HasLeft = false
			} else {
				conversation.Participants = append(conversation.Participants, p)
			}
		}
		err = s.repo.UpdateConversation(ctx, conversation)
		needsNewParticipantNotification = true
	}

	if err != nil {
		return nil, err
	}

	if needsNewParticipantNotification { // send notification only if new participant was added to existing conversation
		notificationMsg := &model.Message{
			Type:           chatRPC.MessageType_AddParticipant.String(),
			ConversationId: conversation.Id.Hex(),
			SenderId:       participant.Id,
			Text:           strings.Join(participants, ","),
		}

		go s.SendMessage(ctx, notificationMsg)
	}

	return conversation, nil
}

func (s *ChatService) LeaveConversation(ctx context.Context, conversationId string) error {
	me, conversation := s.getParticipation(ctx, conversationId)

	return s.leaveConversationOfParticipant(ctx, me, conversation)
}

func (s *ChatService) LeaveConversationForCompany(ctx context.Context, companyId, conversationId string) error {
	me, conversation := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.leaveConversationOfParticipant(ctx, me, conversation)
}

func (s *ChatService) leaveConversationOfParticipant(ctx context.Context, participant *model.Participant, conversation *model.Conversation) error {
	if participant.HasLeft {
		return errors.New("You have already left that conversation")
	}

	notificationMsg := &model.Message{
		Type:           chatRPC.MessageType_ParticipantLeft.String(),
		ConversationId: conversation.Id.Hex(),
		SenderId:       participant.Id,
		Text:           participant.Id,
	}
	s.SendMessage(ctx, notificationMsg)

	participant.HasLeft = true
	participant.LeftTimestamp = time.Now().UnixNano()

	err := s.repo.UpdateParticipant(ctx, conversation.Id.Hex(), participant)
	if err != nil {
		return err
	}

	return nil
}

func (s *ChatService) DeleteConversation(ctx context.Context, conversationId string) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.deleteConversationOfParticipant(ctx, me, conversationId)
}
func (s *ChatService) DeleteConversationForCompany(ctx context.Context, companyId, conversationId string) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.deleteConversationOfParticipant(ctx, me, conversationId)
}
func (s *ChatService) deleteConversationOfParticipant(ctx context.Context, participant *model.Participant, conversationId string) error {
	participant.DeletedTimestamp = time.Now().UnixNano()

	err := s.repo.UpdateParticipant(ctx, conversationId, participant)
	if err != nil {
		return err
	}

	return nil
}

func (s *ChatService) SendMessage(ctx context.Context, message *model.Message) error {
	span := s.tracer.MakeSpan(ctx, "SendMessage")
	defer span.Finish()

	me, conversation := s.getParticipation(ctx, message.ConversationId)
	message.SenderId = me.Id
	message.Timestamp = time.Now().UnixNano()
	if message.Type == "" {
		message.Type = chatRPC.MessageType_UserMessage.String()
	}

	if me.HasLeft {
		return errors.New("You have left the conversation")
	}

	if message.Type != chatRPC.MessageType_MessageStatus.String() {

		for i := range message.Files {
			file, err := s.repo.ReadFile(ctx, message.Files[i].File)
			if err != nil {
				s.tracer.LogError(span, err)
				return err
			}
			if message.Files[i].FileName == "" {
				message.Files[i].FileName = file.Name()
			}
			message.Files[i].FileSize = file.Size()
			file.Close()
		}

		// save message in db
		message, err := s.repo.SendMessage(ctx, message)
		if err != nil {
			s.tracer.LogError(span, err)
			return err
		}
		err = s.repo.SetUnreadForAll(ctx, message.ConversationId, true)
		if err != nil {
			s.tracer.LogError(span, err)
		}
		err = s.repo.SetUnreadFor(ctx, message.ConversationId, message.SenderId, false)
		if err != nil {
			s.tracer.LogError(span, err)
		}
		s.SendUnreadConversationsStatus(ctx, conversation)
	} else {
		if message.Text != "seen" && message.Text != "received" {
			return errors.New("Message status should be 'received' or 'seen'")
		} else {
			s.repo.SetMessageStatus(ctx, message.Id.Hex(), message.SenderId, message.Text)
			if message.Text == "seen" {
				s.repo.SetUnreadFor(ctx, message.ConversationId, message.SenderId, false)
			}
			newMsg, err := s.repo.GetMessage(ctx, message.Id.Hex())
			if err == nil {
				message.ReceivedBy = newMsg.ReceivedBy
				message.SeenBy = newMsg.SeenBy
			}
		}
	}

	// send message to participants
	// fix timestamp
	message.Timestamp = message.Timestamp / 1000000
	s.sendToParticipants(ctx, conversation, message, excludeBlockedAndMuted)

	return nil
}

func (s *ChatService) SendUnverifiedMessage(ctx context.Context, message *model.Message) error {
	conversation, err := s.repo.GetConversation(ctx, message.ConversationId)
	me := findParticipant(message.SenderId, conversation.Participants)
	if me == nil {
		panic(errors.New("You are not participant of the conversations"))
	}

	message.Timestamp = time.Now().UnixNano()

	if me.HasLeft {
		return errors.New("You have left the conversation")
	}

	// save message in db
	message, err = s.repo.SendMessage(ctx, message)
	if err != nil {
		return err
	}
	s.repo.SetUnreadForAll(ctx, message.ConversationId, true) // it's ok if we have error here, message is still sent

	// send message to participants
	for _, p := range conversation.Participants {
		if !p.Blocked && !p.HasLeft {
			s.liveConnections.SendTo(ctx, p.Id, message)
		}
	}

	return nil
}

func (s *ChatService) SearchInConversation(ctx context.Context, conversationId, query, file string) ([]*model.Message, error) {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.searchInConversationOfParticipant(ctx, me, conversationId, query, file)
}
func (s *ChatService) SearchInConversationForCompany(ctx context.Context, companyId, conversationId, query, file string) ([]*model.Message, error) {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.searchInConversationOfParticipant(ctx, me, conversationId, query, file)
}
func (s *ChatService) searchInConversationOfParticipant(ctx context.Context, participant *model.Participant, conversationId, query, file string) ([]*model.Message, error) {
	var messagesUntil int64 = math.MaxInt64
	if participant.HasLeft {
		messagesUntil = participant.LeftTimestamp
	}

	return s.repo.SearchInConversation(ctx, conversationId, query, file, participant.DeletedTimestamp, messagesUntil)
}

func (s *ChatService) SetConversationUnreadFlag(ctx context.Context, conversationId string, value bool) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.repo.SetUnreadFor(ctx, conversationId, me.Id, value)
}
func (s *ChatService) SetConversationUnreadFlagForCompany(ctx context.Context, companyId, conversationId string, value bool) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.repo.SetUnreadFor(ctx, conversationId, me.Id, value)
}

func (s *ChatService) MuteConversation(ctx context.Context, conversationId string, value bool) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.repo.MuteConversation(ctx, me.Id, conversationId, value)
}
func (s *ChatService) MuteConversationForCompany(ctx context.Context, companyId, conversationId string, value bool) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.repo.MuteConversation(ctx, me.Id, conversationId, value)
}

func (s *ChatService) BlockConversetionByParticipants(ctx context.Context, senderID string, targetID string, value bool) error {
	span := s.tracer.MakeSpan(ctx, "BlockConversetionByParticipants")
	defer span.Finish()

	participants := make([]*model.Participant, 2)

	participants[0] = &model.Participant{
		Id: senderID,
	}
	participants[1] = &model.Participant{
		Id: targetID,
	}

	conversations, err := s.repo.FindConversationsWithOnlyParticipants(ctx, participants)
	if err != nil {
		s.tracer.LogError(span, err)
	}

	log.Println("BlockConversetionByParticipants", senderID, targetID)

	if len(conversations) > 0 {
		err = s.BlockConversation(ctx, conversations[0].Id.Hex(), value)
		if err != nil {
			log.Println("BlockConversetionByParticipants ID", conversations[0].Id)
			s.tracer.LogError(span, err)
		}
	}

	return nil
}

func (s *ChatService) BlockConversation(ctx context.Context, conversationId string, value bool) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.repo.BlockConversation(ctx, me.Id, conversationId, value)
}
func (s *ChatService) BlockConversationForCompany(ctx context.Context, companyId, conversationId string, value bool) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.repo.BlockConversation(ctx, me.Id, conversationId, value)
}

func (s *ChatService) ArchiveConversation(ctx context.Context, conversationId string, value bool) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.repo.ArchiveConversation(ctx, me.Id, conversationId, value)
}
func (s *ChatService) ArchiveConversationForCompany(ctx context.Context, companyId, conversationId string, value bool) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.repo.ArchiveConversation(ctx, me.Id, conversationId, value)
}

func (s *ChatService) RenameConversation(ctx context.Context, conversationId, name string) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.renameConversationOfParticipant(ctx, me, conversationId, name)
}
func (s *ChatService) RenameConversationForCompany(ctx context.Context, companyId, conversationId, name string) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.renameConversationOfParticipant(ctx, me, conversationId, name)
}
func (s *ChatService) renameConversationOfParticipant(ctx context.Context, participant *model.Participant, conversationId, name string) error {
	// if !participant.IsAdmin {
	// 	return errors.New("You are not admin of the conversation")
	// }

	return s.repo.RenameConversation(ctx, conversationId, name)
}
func (s *ChatService) ChangeConversationAvatar(ctx context.Context, conversationId, avatar string) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.changeConversationAvatarOfParticipant(ctx, me, conversationId, avatar)
}
func (s *ChatService) ChangeConversationAvatarForCompany(ctx context.Context, companyId, conversationId, avatar string) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.changeConversationAvatarOfParticipant(ctx, me, conversationId, avatar)
}
func (s *ChatService) changeConversationAvatarOfParticipant(ctx context.Context, participant *model.Participant, conversationId, avatar string) error {
	// if !participant.IsAdmin {
	// 	return errors.New("You are not admin of the conversation")
	// }

	return s.repo.ChangeConversationAvatar(ctx, conversationId, avatar)
}

func (s *ChatService) CreateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	senderId := s.AuthenticateUser(ctx)

	reply.OwnerId = senderId

	return s.repo.CreateReply(ctx, reply)
}
func (s *ChatService) CreateReplyForCompany(ctx context.Context, companyId string, reply *model.Reply) (*model.Reply, error) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")

	reply.OwnerId = companyId

	return s.repo.CreateReply(ctx, reply)
}
func (s *ChatService) UpdateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	senderId := s.AuthenticateUser(ctx)

	reply.OwnerId = senderId

	existing, err := s.repo.GetReply(ctx, reply.Id.Hex())
	if err != nil {
		return nil, err
	}
	if existing.OwnerId != reply.OwnerId {
		return nil, errors.New("You are not owner of the reply")
	}

	return s.repo.UpdateReply(ctx, reply)
}
func (s *ChatService) UpdateReplyForCompany(ctx context.Context, companyId string, reply *model.Reply) (*model.Reply, error) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")

	reply.OwnerId = companyId

	existing, err := s.repo.GetReply(ctx, reply.Id.Hex())
	if err != nil {
		return nil, err
	}
	if existing.OwnerId != reply.OwnerId {
		return nil, errors.New("You are not owner of the reply")
	}

	return s.repo.UpdateReply(ctx, reply)
}

func (s *ChatService) DeleteReply(ctx context.Context, replyId string) error {
	senderId := s.AuthenticateUser(ctx)
	return s.deleteReplyOfOwner(ctx, senderId, replyId)
}
func (s *ChatService) DeleteReplyForCompany(ctx context.Context, companyId, replyId string) error {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	return s.deleteReplyOfOwner(ctx, companyId, replyId)
}
func (s *ChatService) deleteReplyOfOwner(ctx context.Context, ownerId, replyId string) error {
	reply, err := s.repo.GetReply(ctx, replyId)
	if err != nil {
		return err
	}

	if reply.OwnerId != ownerId {
		return errors.New("This reply doesn't belong to you")
	}

	return s.repo.DeleteReply(ctx, replyId)
}

func (s *ChatService) GetMyReplies(ctx context.Context, query string) ([]*model.Reply, error) {
	senderId := s.AuthenticateUser(ctx)

	return s.repo.GetUserReplies(ctx, senderId, query)
}
func (s *ChatService) GetMyRepliesForCompany(ctx context.Context, companyId, query string) ([]*model.Reply, error) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")

	return s.repo.GetUserReplies(ctx, companyId, query)
}

func (s *ChatService) CreateLabel(ctx context.Context, label *model.ConversationLabel) (*model.ConversationLabel, error) {
	senderId := s.AuthenticateUser(ctx)
	label.OwnerId = senderId
	return s.repo.CreateLabel(ctx, label)
}
func (s *ChatService) CreateLabelForCompany(ctx context.Context, companyId string, label *model.ConversationLabel) (*model.ConversationLabel, error) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")

	label.OwnerId = companyId
	return s.repo.CreateLabel(ctx, label)
}

func (s *ChatService) DeleteLabel(ctx context.Context, labelId string) error {
	senderId := s.AuthenticateUser(ctx)
	return s.repo.DeleteLabel(ctx, senderId, labelId)
}
func (s *ChatService) DeleteLabelForCompany(ctx context.Context, companyId, labelId string) error {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	return s.repo.DeleteLabel(ctx, companyId, labelId)
}

func (s *ChatService) UpdateLabel(ctx context.Context, label *model.ConversationLabel) error {
	senderId := s.AuthenticateUser(ctx)
	label.OwnerId = senderId
	return s.repo.UpdateLabel(ctx, label)
}
func (s *ChatService) UpdateLabelForCompany(ctx context.Context, companyId string, label *model.ConversationLabel) error {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	label.OwnerId = companyId
	return s.repo.UpdateLabel(ctx, label)
}

func (s *ChatService) GetAllLabel(ctx context.Context) ([]*model.ConversationLabel, error) {
	senderId := s.AuthenticateUser(ctx)
	return s.repo.GetAllLabel(ctx, senderId)
}
func (s *ChatService) GetAllLabelForCompany(ctx context.Context, companyId string) ([]*model.ConversationLabel, error) {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	return s.repo.GetAllLabel(ctx, companyId)
}

func (s *ChatService) AddLabelToConversation(ctx context.Context, conversationId, labelId string) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.addLabelToConversationOfParticipant(ctx, me, conversationId, labelId)
}
func (s *ChatService) AddLabelToConversationForCompany(ctx context.Context, companyId, conversationId, labelId string) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.addLabelToConversationOfParticipant(ctx, me, conversationId, labelId)
}
func (s *ChatService) addLabelToConversationOfParticipant(ctx context.Context, participant *model.Participant, conversationId, labelId string) error {
	for _, label := range participant.Labels {
		if label.Hex() == labelId {
			return errors.New("This label already exists on this conversation")
		}
	}

	participant.Labels = append(participant.Labels, bson.ObjectIdHex(labelId))
	return s.repo.UpdateParticipant(ctx, conversationId, participant)
}

func (s *ChatService) RemoveLabelFromConversation(ctx context.Context, conversationId, labelId string) error {
	me, _ := s.getParticipation(ctx, conversationId)

	return s.removeLabelFromConversationOfParticipant(ctx, me, conversationId, labelId)
}
func (s *ChatService) RemoveLabelFromConversationForCompany(ctx context.Context, companyId, conversationId, labelId string) error {
	me, _ := s.getParticipationForCompany(ctx, companyId, conversationId)

	return s.removeLabelFromConversationOfParticipant(ctx, me, conversationId, labelId)
}
func (s *ChatService) removeLabelFromConversationOfParticipant(ctx context.Context, participant *model.Participant, conversationId, labelId string) error {
	for i, label := range participant.Labels {
		if label.Hex() == labelId {
			participant.Labels = append(participant.Labels[:i], participant.Labels[i+1:]...)
			break
		}
	}
	return s.repo.UpdateParticipant(ctx, conversationId, participant)
}

func (s *ChatService) ReportConversation(ctx context.Context, report *model.Report) error {
	me, _ := s.getParticipation(ctx, report.ConversationId)

	if me.HasLeft {
		return errors.New("You have already left the conversation")
	}
	report.UserId = me.Id

	return s.repo.ReportConversation(ctx, report)
}

func (s *ChatService) ReportConversationForCompany(ctx context.Context, report *model.Report) error {
	senderId := s.AuthenticateUser(ctx)
	me, _ := s.getParticipationForCompany(ctx, report.CompanyId, report.ConversationId)

	if me.HasLeft {
		return errors.New("You have already left the conversation")
	}
	report.UserId = senderId

	return s.repo.ReportConversation(ctx, report)
}

func (s *ChatService) IsUserLive(ctx context.Context, userId string) (bool, error) {
	return s.liveConnections.IsLive(userId), nil
}

func (s *ChatService) SetParticipantOffline(ctx context.Context, status bool) error {
	senderId := s.AuthenticateUser(ctx)
	s.liveConnections.SetParticipantOffline(senderId, status)
	return nil
}

func (s *ChatService) SetParticipantOfflineForCompany(ctx context.Context, companyId string, status bool) error {
	s.RequireAdminLevelForCompany(ctx, companyId, "Admin")
	s.liveConnections.SetParticipantOffline(companyId, status)
	return nil
}

func (s *ChatService) UploadFile(ctx context.Context, header *multipart.FileHeader) (string, error) {
	userId := s.AuthenticateUser(ctx)

	file, err := header.Open()
	if err != nil {
		return "", err
	}
	filename := header.Filename
	fmt.Println("save file ", filename, " with size ", header.Size)
	gridFile, err := s.repo.SaveFile(ctx, file, filename, map[string]interface{}{
		"sender_id": userId,
	})
	if err != nil {
		return "", err
	}
	fmt.Println("grid file id: ", gridFile.Id())
	fileId := gridFile.Id().(bson.ObjectId)

	//message := &model.Message{
	//	ConversationId: conversationId,
	//	File:           fileId.Hex(),
	//	FileName:       filename,
	//	FileSize:       header.Size,
	//}
	//
	//err = s.SendMessage(ctx, message)
	//if err != nil {
	//	return "", err
	//}

	return fileId.Hex(), nil
}

func (s *ChatService) ReadFile(ctx context.Context, fileId string) (*mgo.GridFile, error) {
	s.AuthenticateUser(ctx)
	return s.repo.ReadFile(ctx, fileId)
}

func excludeBlockedAndMuted(participant *model.Participant) bool {
	return !(participant.HasLeft || participant.Blocked)
}

func (s *ChatService) sendToParticipants(ctx context.Context, conv *model.Conversation, msg interface{}, predicate func(participant *model.Participant) bool) {
	span := s.tracer.MakeSpan(ctx, "sendToParticipants")
	defer span.Finish()

	if predicate == nil {
		predicate = excludeBlockedAndMuted
	}

	for _, p := range conv.Participants {
		if predicate(p) {
			s.liveConnections.SendTo(ctx, p.Id, msg)
		}
	}
}
