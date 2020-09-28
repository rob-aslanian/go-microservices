package resolver

import (
	"context"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
)

func (_ *Resolver) GetActiveConnections(ctx context.Context) ([]ProfileResolver, error) {
	res, err := chat.GetActiveConnections(ctx, &chatRPC.Empty{})
	if err != nil {
		return nil, err
	}

	profs := make([]ProfileResolver, 0, len(res.GetIDs()))

	if len(res.GetIDs()) > 0 {
		profiles, err := user.GetMapProfilesByID(ctx, &userRPC.UserIDs{
			ID: res.GetIDs(),
		})
		if err != nil {
			return nil, err
		}

		for _, id := range res.GetIDs() {
			p := ToProfile(ctx, profiles.GetProfiles()[id])
			profs = append(profs, ProfileResolver{
				R: &p,
			})
		}
	}

	return profs, nil
}

func (_ *Resolver) GetMyConversations(ctx context.Context, input GetMyConversationsRequest) ([]ConversationResolver, error) {
	filter := chatRPC.ConversationFilter{}
	if input.Category != nil {
		filter.Category = chatRPC.ConversationCategory(chatRPC.ConversationCategory_value[*input.Category])
	}
	if input.LabelId != nil {
		filter.LabelId = *input.LabelId
	}
	if input.ParticipantId != nil {
		filter.ParticipantId = *input.ParticipantId
	}
	if input.Text != nil {
		filter.Text = *input.Text
	}
	res, err := chat.GetMyConversations(ctx, &filter)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]ConversationResolver, len(res.Conversations))
	for i, o := range res.Conversations {
		conv, err := chat_conversationToResolver(ctx, o)
		if err != nil {
			return nil, err
		}
		list[i] = *conv
	}
	return list, nil
}

func (_ *Resolver) GetConversation(ctx context.Context, input GetConversationRequest) (*ConversationResolver, error) {
	res, err := chat.GetConversation(ctx, &chatRPC.Conversation{Id: input.ID})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	conv, err := chat_conversationToResolver(ctx, res)
	return conv, err
}

func (_ *Resolver) GetMessages(ctx context.Context, input GetMessagesRequest) ([]MessageResolver, error) {
	res, err := chat.GetMessages(ctx, &chatRPC.Conversation{Id: input.ConversationId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]MessageResolver, len(res.Messages))
	for i, o := range res.Messages {
		list[i] = *chat_messageToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) SearchInConversation(ctx context.Context, input SearchInConversationRequest) ([]MessageResolver, error) {
	res, err := chat.SearchInConversation(ctx, &chatRPC.SearchInConversationRequest{
		ConversationId: input.ConversationId,
		Query:          input.Query,
		File:           input.File,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]MessageResolver, len(res.Messages))
	for i, o := range res.Messages {
		list[i] = *chat_messageToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) GetMyReplies(ctx context.Context, input GetMyRepliesRequest) ([]ChatReplyResolver, error) {
	query := ""
	if input.Query != nil {
		query = *input.Query
	}
	res, err := chat.GetMyReplies(ctx, &chatRPC.StringValue{Value: query})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]ChatReplyResolver, len(res.Replies))
	for i, o := range res.Replies {
		list[i] = *chat_chatReplyToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) GetAllLabels(ctx context.Context) ([]ChatLabelResolver, error) {
	res, err := chat.GetAllLabel(ctx, &chatRPC.Empty{})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]ChatLabelResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *chat_chatLabelToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) CreateConversation(ctx context.Context, input CreateConversationRequest) (*ConversationResolver, error) {

	participants := make([]*chatRPC.Participant, 0, len(input.Participants))

	for _, p := range input.Participants {
		participants = append(participants, &chatRPC.Participant{
			Id:        p.ID,
			IsAdmin:   p.Is_admin,
			IsCompany: p.Is_company,
		})
	}
	res, err := chat.CreateConversation(ctx, &chatRPC.Conversation{
		Name:         input.Name,
		Avatar:       input.Avatar,
		Participants: participants,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	conv, err := chat_conversationToResolver(ctx, res)
	return conv, err
}

func (_ *Resolver) AddParticipants(ctx context.Context, input AddParticipantsRequest) (*ConversationResolver, error) {
	participants := make([]string, len(input.Participants))
	for i, p := range input.Participants {
		participants[i] = p.ID
	}
	conversation, err := chat.AddParticipants(ctx, &chatRPC.AddParticipantsRequest{
		ConversationId: input.ConversationId,
		Participants:   participants,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return chat_conversationToResolver(ctx, conversation)
}

func (_ *Resolver) LeaveConversation(ctx context.Context, input LeaveConversationRequest) (*bool, error) {
	_, err := chat.LeaveConversation(ctx, &chatRPC.Conversation{Id: input.ConversationId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) DeleteConversation(ctx context.Context, input DeleteConversationRequest) (*bool, error) {
	_, err := chat.DeleteConversation(ctx, &chatRPC.Conversation{Id: input.ConversationId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetConversationUnreadFlag(ctx context.Context, input SetConversationUnreadFlagRequest) (*bool, error) {
	_, err := chat.SetConversationUnreadFlag(ctx, &chatRPC.SetConversationUnreadFlagRequest{
		ConversationId: input.ConversationId,
		Value:          input.Flag,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) MuteConversation(ctx context.Context, input MuteConversationRequest) (*bool, error) {
	_, err := chat.MuteConversation(ctx, &chatRPC.ConversationIdWithBool{
		ConversationId: input.ConversationId,
		Value:          input.Mute,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) BlockConversation(ctx context.Context, input BlockConversationRequest) (*bool, error) {
	_, err := chat.BlockConversation(ctx, &chatRPC.ConversationIdWithBool{
		ConversationId: input.ConversationId,
		Value:          input.Block,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ArchiveConversation(ctx context.Context, input ArchiveConversationRequest) (*bool, error) {
	_, err := chat.ArchiveConversation(ctx, &chatRPC.ConversationIdWithBool{
		ConversationId: input.ConversationId,
		Value:          input.Archive,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RenameConversation(ctx context.Context, input RenameConversationRequest) (*bool, error) {
	_, err := chat.RenameConversation(ctx, &chatRPC.ConversationIdWithString{
		ConversationId: input.ConversationId,
		Value:          input.Name,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}
func (_ *Resolver) ChangeConversationAvatar(ctx context.Context, input ChangeConversationAvatarRequest) (*bool, error) {
	_, err := chat.ChangeConversationAvatar(ctx, &chatRPC.ConversationIdWithString{
		ConversationId: input.ConversationId,
		Value:          input.Avatar,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) CreateReply(ctx context.Context, input CreateReplyRequest) (*ChatReplyResolver, error) {
	req := chatRPC.Reply{
		Title: input.Title,
		Text:  input.Text,
	}
	if input.Files != nil {
		req.Files = chat_replyFileArrToRPC(*input.Files)
	}
	reply, err := chat.CreateReply(ctx, &req)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return chat_chatReplyToResolver(reply), nil
}
func (_ *Resolver) UpdateReply(ctx context.Context, input UpdateReplyRequest) (*ChatReplyResolver, error) {
	req := chatRPC.Reply{
		Id:    input.ID,
		Title: input.Title,
		Text:  input.Text,
	}
	if input.Files != nil {
		req.Files = chat_replyFileArrToRPC(*input.Files)
	}
	reply, err := chat.UpdateReply(ctx, &req)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return chat_chatReplyToResolver(reply), nil
}

func (_ *Resolver) DeleteReply(ctx context.Context, input DeleteReplyRequest) (*bool, error) {
	_, err := chat.DeleteReply(ctx, &chatRPC.Reply{
		Id: input.ID,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) CreateLabel(ctx context.Context, input CreateLabelRequest) (ChatLabelResolver, error) {
	label, err := chat.CreateLabel(ctx, &chatRPC.ConversationLabel{
		Name:  input.Name,
		Color: input.Color,
	})
	if e, isErr := handleError(err); isErr {
		return ChatLabelResolver{R: &ChatLabel{}}, e
	}
	return *chat_chatLabelToResolver(label), nil
}

func (_ *Resolver) UpdateLabel(ctx context.Context, input UpdateLabelRequest) (*bool, error) {
	_, err := chat.UpdateLabel(ctx, &chatRPC.ConversationLabel{
		Id:    input.ID,
		Name:  input.Name,
		Color: input.Color,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) DeleteLabel(ctx context.Context, input DeleteLabelRequest) (*bool, error) {
	_, err := chat.DeleteLabel(ctx, &chatRPC.ConversationLabel{
		Id: input.ID,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) AddLabelToConversation(ctx context.Context, input AddLabelToConversationRequest) (*bool, error) {
	_, err := chat.AddLabelToConversation(ctx, &chatRPC.ConversationIdWithLabelId{
		ConversationId: input.ConversationId,
		LabelId:        input.LabelId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveLabelFromConversation(ctx context.Context, input RemoveLabelFromConversationRequest) (*bool, error) {
	_, err := chat.RemoveLabelFromConversation(ctx, &chatRPC.ConversationIdWithLabelId{
		ConversationId: input.ConversationId,
		LabelId:        input.LabelId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ReportConversation(ctx context.Context, input ReportConversationRequest) (*bool, error) {
	_, err := chat.ReportConversation(ctx, &chatRPC.ConversationIdWithString{
		ConversationId: input.ConversationId,
		Value:          input.Text,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetOffline(ctx context.Context, input SetOfflineRequest) (*bool, error) {
	_, err := chat.SetParticipantOffline(ctx, &chatRPC.BoolValue{Value: input.Offline})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetOfflineForCompany(ctx context.Context, input SetOfflineForCompanyRequest) (*bool, error) {
	_, err := chat.SetParticipantOfflineForCompany(ctx, &chatRPC.CompanyIdWithBool{CompanyId: input.CompanyId, Value: input.Offline})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) GetMyConversationsForCompany(ctx context.Context, input GetMyConversationsForCompanyRequest) ([]ConversationResolver, error) {
	filter := chatRPC.ConversationFilter{}
	if input.Category != nil {
		filter.Category = chatRPC.ConversationCategory(chatRPC.ConversationCategory_value[*input.Category])
	}
	if input.LabelId != nil {
		filter.LabelId = *input.LabelId
	}
	if input.ParticipantId != nil {
		filter.ParticipantId = *input.ParticipantId
	}
	if input.Text != nil {
		filter.Text = *input.Text
	}
	res, err := chat.GetMyConversationsForCompany(ctx, &chatRPC.CompanyIdWithConversationFilter{CompanyId: input.CompanyId, Filter: &filter})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]ConversationResolver, len(res.Conversations))
	for i, o := range res.Conversations {
		conv, err := chat_conversationToResolver(ctx, o)
		if err != nil {
			return nil, err
		}
		list[i] = *conv
	}
	return list, nil
}

func (_ *Resolver) GetConversationForCompany(ctx context.Context, input GetConversationForCompanyRequest) (*ConversationResolver, error) {
	res, err := chat.GetConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationId{CompanyId: input.CompanyId, ConversationId: input.ID})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	conv, err := chat_conversationToResolver(ctx, res)
	return conv, err
}

func (_ *Resolver) GetMessagesForCompany(ctx context.Context, input GetMessagesForCompanyRequest) ([]MessageResolver, error) {
	res, err := chat.GetMessagesForCompany(ctx, &chatRPC.CompanyIdWithConversationId{CompanyId: input.CompanyId, ConversationId: input.ConversationId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]MessageResolver, len(res.Messages))
	for i, o := range res.Messages {
		list[i] = *chat_messageToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) SearchInConversationForCompany(ctx context.Context, input SearchInConversationForCompanyRequest) ([]MessageResolver, error) {
	res, err := chat.SearchInConversationForCompany(ctx, &chatRPC.SearchInConversationForCompanuRequest{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Query:          input.Query,
		File:           input.File,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]MessageResolver, len(res.Messages))
	for i, o := range res.Messages {
		list[i] = *chat_messageToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) GetMyRepliesForCompany(ctx context.Context, input GetMyRepliesForCompanyRequest) ([]ChatReplyResolver, error) {
	query := ""
	if input.Query != nil {
		query = *input.Query
	}
	res, err := chat.GetMyRepliesForCompany(ctx, &chatRPC.CompanyIdWithString{CompanyId: input.CompanyId, Value: query})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]ChatReplyResolver, len(res.Replies))
	for i, o := range res.Replies {
		list[i] = *chat_chatReplyToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) GetAllLabelsForCompany(ctx context.Context, input GetAllLabelsForCompanyRequest) ([]ChatLabelResolver, error) {
	res, err := chat.GetAllLabelForCompany(ctx, &chatRPC.CompanyId{CompanyId: input.CompanyId})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	list := make([]ChatLabelResolver, len(res.List))
	for i, o := range res.List {
		list[i] = *chat_chatLabelToResolver(o)
	}
	return list, nil
}

func (_ *Resolver) CreateConversationForCompany(ctx context.Context, input CreateConversationForCompanyRequest) (*ConversationResolver, error) {
	participants := make([]*chatRPC.Participant, len(input.Participants))
	for i, p := range input.Participants {
		participants[i] = &chatRPC.Participant{
			Id:        p.ID,
			IsCompany: p.Is_company,
			IsAdmin:   p.Is_admin,
		}
	}
	res, err := chat.CreateConversationForCompany(ctx, &chatRPC.CompanyIdWithConversation{
		CompanyId: input.CompanyId,
		Conversation: &chatRPC.Conversation{
			Name:         input.Name,
			Avatar:       input.Avatar,
			Participants: participants,
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	conv, err := chat_conversationToResolver(ctx, res)
	return conv, err
}

func (_ *Resolver) AddParticipantsForCompany(ctx context.Context, input AddParticipantsForCompanyRequest) (*ConversationResolver, error) {
	participants := make([]string, len(input.Participants))
	for i, p := range input.Participants {
		participants[i] = p.ID
	}
	res, err := chat.AddParticipantsForCompany(ctx, &chatRPC.AddParticipantsForCompanyRequest{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Participants:   participants,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return chat_conversationToResolver(ctx, res)
}

func (_ *Resolver) LeaveConversationForCompany(ctx context.Context, input LeaveConversationForCompanyRequest) (*bool, error) {
	_, err := chat.LeaveConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationId{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) DeleteConversationForCompany(ctx context.Context, input DeleteConversationForCompanyRequest) (*bool, error) {
	_, err := chat.DeleteConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationId{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) SetConversationUnreadFlagForCompany(ctx context.Context, input SetConversationUnreadFlagForCompanyRequest) (*bool, error) {
	_, err := chat.SetConversationUnreadFlagForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndBool{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Value:          input.Flag,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) MuteConversationForCompany(ctx context.Context, input MuteConversationForCompanyRequest) (*bool, error) {
	_, err := chat.MuteConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndBool{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Value:          input.Mute,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) BlockConversationForCompany(ctx context.Context, input BlockConversationForCompanyRequest) (*bool, error) {
	_, err := chat.BlockConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndBool{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Value:          input.Block,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ArchiveConversationForCompany(ctx context.Context, input ArchiveConversationForCompanyRequest) (*bool, error) {
	_, err := chat.ArchiveConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndBool{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Value:          input.Archive,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RenameConversationForCompany(ctx context.Context, input RenameConversationForCompanyRequest) (*bool, error) {
	_, err := chat.RenameConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndString{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Value:          input.Name,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}
func (_ *Resolver) ChangeConversationAvatarForCompany(ctx context.Context, input ChangeConversationAvatarForCompanyRequest) (*bool, error) {
	_, err := chat.ChangeConversationAvatarForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndString{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Value:          input.Avatar,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) CreateReplyForCompany(ctx context.Context, input CreateReplyForCompanyRequest) (*ChatReplyResolver, error) {
	req := chatRPC.CreateReplyForCompanyRequest{
		CompanyId: input.CompanyId,
		Reply: &chatRPC.Reply{
			Title: input.Title,
			Text:  input.Text,
		},
	}
	if input.Files != nil {
		req.Reply.Files = chat_replyFileArrToRPC(*input.Files)
	}
	reply, err := chat.CreateReplyForCompany(ctx, &req)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return chat_chatReplyToResolver(reply), nil
}
func (_ *Resolver) UpdateReplyForCompany(ctx context.Context, input UpdateReplyForCompanyRequest) (*ChatReplyResolver, error) {
	req := chatRPC.CreateReplyForCompanyRequest{
		CompanyId: input.CompanyId,
		Reply: &chatRPC.Reply{
			Id:    input.ID,
			Title: input.Title,
			Text:  input.Text,
		},
	}
	if input.Files != nil {
		req.Reply.Files = chat_replyFileArrToRPC(*input.Files)
	}
	reply, err := chat.UpdateReplyForCompany(ctx, &req)
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return chat_chatReplyToResolver(reply), nil
}

func (_ *Resolver) DeleteReplyForCompany(ctx context.Context, input DeleteReplyForCompanyRequest) (*bool, error) {
	_, err := chat.DeleteReplyForCompany(ctx, &chatRPC.CompanyIdWithId{
		CompanyId: input.CompanyId,
		Id:        input.ID,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) CreateLabelForCompany(ctx context.Context, input CreateLabelForCompanyRequest) (ChatLabelResolver, error) {
	label, err := chat.CreateLabelForCompany(ctx, &chatRPC.CompanyIdWithConversationLabel{
		CompanyId: input.CompanyId,
		Label: &chatRPC.ConversationLabel{
			Name:  input.Name,
			Color: input.Color,
		},
	})
	if e, isErr := handleError(err); isErr {
		return ChatLabelResolver{R: &ChatLabel{}}, e
	}
	return *chat_chatLabelToResolver(label), nil
}

func (_ *Resolver) UpdateLabelForCompany(ctx context.Context, input UpdateLabelForCompanyRequest) (*bool, error) {
	_, err := chat.UpdateLabelForCompany(ctx, &chatRPC.CompanyIdWithConversationLabel{
		CompanyId: input.CompanyId,
		Label: &chatRPC.ConversationLabel{
			Id:    input.ID,
			Name:  input.Name,
			Color: input.Color,
		},
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) DeleteLabelForCompany(ctx context.Context, input DeleteLabelForCompanyRequest) (*bool, error) {
	_, err := chat.DeleteLabelForCompany(ctx, &chatRPC.CompanyIdWithId{
		CompanyId: input.CompanyId,
		Id:        input.ID,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) AddLabelToConversationForCompany(ctx context.Context, input AddLabelToConversationForCompanyRequest) (*bool, error) {
	_, err := chat.AddLabelToConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndLabelId{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		LabelId:        input.LabelId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) RemoveLabelFromConversationForCompany(ctx context.Context, input RemoveLabelFromConversationForCompanyRequest) (*bool, error) {
	_, err := chat.RemoveLabelFromConversationForCompany(ctx, &chatRPC.CompanyIdWithConversationIdAndLabelId{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		LabelId:        input.LabelId,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}

func (_ *Resolver) ReportConversationForCompany(ctx context.Context, input ReportConversationForCompanyRequest) (*bool, error) {
	_, err := chat.ReportConversationForCompany(ctx, &chatRPC.ReportConversationForCompanyRequest{
		CompanyId:      input.CompanyId,
		ConversationId: input.ConversationId,
		Text:           input.Text,
	})
	if e, isErr := handleError(err); isErr {
		return nil, e
	}
	return nil, nil
}
