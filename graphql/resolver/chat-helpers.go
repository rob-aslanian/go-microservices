package resolver

import (
	"context"
	"fmt"
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/utils"
)

func chat_conversationToResolver(ctx context.Context, c *chatRPC.Conversation) (*ConversationResolver, error) {
	userParticipantIds := userRPC.UserIDs{}
	companyParticipantIds := companyRPC.GetCompanyProfilesRequest{}
	for _, p := range c.Participants {
		if p.IsCompany {
			companyParticipantIds.Ids = append(companyParticipantIds.Ids, p.Id)
		} else {
			userParticipantIds.ID = append(userParticipantIds.ID, p.Id)
		}
	}
	var userParticipantProfiles *userRPC.ProfileList
	var companyParticipantProfiles *companyRPC.GetCompanyProfilesResponse
	var err error
	if len(userParticipantIds.ID) > 0 {
		userParticipantProfiles, err = user.GetProfilesByID(utils.ToOutContext(ctx), &userParticipantIds)
		if err != nil {
			return nil, err
		}
	}
	if len(companyParticipantIds.Ids) > 0 {
		companyParticipantProfiles, err = company.GetCompanyProfiles(utils.ToOutContext(ctx), &companyParticipantIds)
		if err != nil {
			return nil, err
		}
	}

	participants := make([]Participant, len(c.Participants))
	for i, p := range c.Participants {
		var fields = &chatParticipantFields{}
		if userParticipantProfiles != nil {
			for _, pr := range userParticipantProfiles.Profiles { // TODO needs optimization
				if pr != nil && pr.ID == p.Id {
					fields.id = p.Id
					fields.url = pr.URL
					fields.avatar = pr.Avatar
					fields.name = fmt.Sprint(pr.Firstname, " ", pr.Lastname)
					break
				}
			}
		}
		if companyParticipantProfiles != nil {
			for _, pr := range companyParticipantProfiles.Profiles { // TODO needs optimization
				if pr != nil && pr.Id == p.Id {
					fields.id = p.Id
					fields.url = pr.URL
					fields.avatar = pr.GetAvatar()
					fields.name = pr.Name
					break
				}
			}
		}
		participants[i] = *chat_participantToGql(p, fields)
	}
	var lastMessage *Message
	if c.LastMessage != nil {
		lastMessage = chat_messageToGql(c.LastMessage)
	} else {
		lastMessage = &Message{} // TODO should return nil, fix error
	}

	//conversationName := c.Name
	//if conversationName == "" {
	//	for _, p := range participants {
	//		conversationName = fmt.Sprint(conversationName, ", ", p.Name)
	//	}
	//	conversationName = conversationName[2:] // remove the starting ', '
	//}

	return &ConversationResolver{
		R: &Conversation{
			ID:           c.Id,
			Name:         c.Name,
			Avatar:       c.Avatar,
			Is_group:     c.IsGroup,
			Last_message: lastMessage,
			Participants: participants,
			Archived:     c.Archived,
			Blocked:      c.Blocked,
			Has_left:     c.HasLeft,
			Muted:        c.Muted,
			Unread:       c.Unread,
			Labels:       c.Labels,
		},
	}, nil
}

type chatParticipantFields struct {
	id, name, url, avatar string
}

func chat_participantToGql(p *chatRPC.Participant, fields *chatParticipantFields) *Participant {
	return &Participant{
		ID:         fields.id,
		Name:       fields.name,
		Avatar:     fields.avatar,
		Url:        fields.url,
		Is_admin:   p.IsAdmin,
		Is_company: p.IsCompany,
		Has_left:   p.HasLeft,
		Is_active:  p.IsActive,
	}
}

func chat_messageToResolver(m *chatRPC.Message) *MessageResolver {
	return &MessageResolver{
		R: chat_messageToGql(m),
	}
}

func chat_messageToGql(m *chatRPC.Message) *Message {
	mes := Message{
		ID:              m.Id,
		Type:            m.Type.String(),
		Conversation_id: m.ConversationId,
		Sender_id:       m.SenderId,
		Text:            m.Text,
		Files:           make([]FileMessage, 0, len(m.Files)),
		// File:            m.File,
		// File_name:       m.FileName,
		// File_size:       int32(m.FileSize),
		Received_by: m.ReceivedBy,
		Seen_by:     m.SeenBy,
		Timestamp:   time.Unix(0, m.Timestamp),
	}

	for _, me := range m.Files {
		mes.Files = append(mes.Files, FileMessage{
			File:      me.GetFile(),
			File_name: me.GetFileName(),
			File_size: int32(me.GetFileSize()),
		})
	}

	return &mes
}

func chat_chatReplyToResolver(r *chatRPC.Reply) *ChatReplyResolver {
	files := make([]ChatReplyFile, len(r.Files))
	for i, f := range r.Files {
		files[i] = ChatReplyFile{
			ID:   f.Id,
			Name: f.Name,
		}
	}
	return &ChatReplyResolver{
		R: &ChatReply{
			ID:    r.Id,
			Text:  r.Text,
			Title: r.Title,
			Files: files,
		},
	}
}

func chat_replyFileToRPC(f *ChatReplyFileInput) *chatRPC.ReplyFile {
	return &chatRPC.ReplyFile{
		Id:   f.ID,
		Name: f.Name,
	}
}
func chat_replyFileArrToRPC(f []ChatReplyFileInput) []*chatRPC.ReplyFile {
	files := make([]*chatRPC.ReplyFile, len(f))
	for i, c := range f {
		files[i] = chat_replyFileToRPC(&c)
	}
	return files
}

func chat_chatLabelToResolver(l *chatRPC.ConversationLabel) *ChatLabelResolver {
	return &ChatLabelResolver{
		R: &ChatLabel{
			ID:    l.Id,
			Name:  l.Name,
			Color: l.Color,
			Count: l.Count,
		},
	}
}
