package grpc_handlers

import (
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"golang.org/x/net/context"
)

func (h *GrpcHandlers) CreateConversation(ctx context.Context, input *chatRPC.Conversation) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation := model.NewConversation(input)
	conversation, err = h.ns.CreateConversation(ctx, conversation)
	panicIf(err)

	return conversation.ToRPC(), nil
}

func (h *GrpcHandlers) CreateConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversation) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation := model.NewConversation(input.Conversation)
	conversation, err = h.ns.CreateConversationForCompany(ctx, input.CompanyId, conversation)
	panicIf(err)

	return conversation.ToRPC(), nil
}

func (this *GrpcHandlers) CreateUnverifiedConversation(ctx context.Context, data *chatRPC.Conversation) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation := model.NewConversation(data)
	conversation, err = this.ns.CreateUnverifiedConversation(ctx, conversation)
	panicIf(err)

	return conversation.ToRPC(), nil
}

func (h *GrpcHandlers) GetMyConversations(ctx context.Context, input *chatRPC.ConversationFilter) (response *chatRPC.ConversationArr, err error) {
	defer recoverHandler(&err)

	conversations, err := h.ns.GetMyConversations(ctx, &model.ConversationFilter{
		Category:      chatRPC.ConversationCategory_name[int32(input.Category)],
		LabelId:       input.LabelId,
		ParticipantId: input.ParticipantId,
		Text:          input.Text,
	})
	panicIf(err)

	return model.ConversationArr(conversations).ToRPC(), nil
}

func (h *GrpcHandlers) GetMyConversationsForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationFilter) (response *chatRPC.ConversationArr, err error) {
	defer recoverHandler(&err)

	conversations, err := h.ns.GetMyConversationsForCompany(ctx, input.CompanyId, &model.ConversationFilter{
		Category:      chatRPC.ConversationCategory_name[int32(input.Filter.Category)],
		LabelId:       input.Filter.LabelId,
		ParticipantId: input.Filter.ParticipantId,
		Text:          input.Filter.Text,
	})
	panicIf(err)

	return model.ConversationArr(conversations).ToRPC(), nil
}

func (h *GrpcHandlers) GetConversation(ctx context.Context, input *chatRPC.Conversation) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation, err := h.ns.GetConversation(ctx, input.Id)
	panicIf(err)

	return conversation.ToRPC(), nil
}

func (h *GrpcHandlers) GetConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationId) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation, err := h.ns.GetConversationForCompany(ctx, input.CompanyId, input.ConversationId)
	panicIf(err)

	return conversation.ToRPC(), nil
}

func (h *GrpcHandlers) GetActiveConnections(ctx context.Context, input *chatRPC.Empty) (response *chatRPC.IDs, err error) {

	ids, err := h.ns.GetActiveConnections(ctx)

	return &chatRPC.IDs{
		IDs: ids,
	}, nil
}

func (h *GrpcHandlers) GetMessages(ctx context.Context, input *chatRPC.Conversation) (response *chatRPC.MessageArr, err error) {
	defer recoverHandler(&err)

	messages, err := h.ns.GetMessages(ctx, input.Id)
	panicIf(err)

	return model.MessageArr(messages).ToRPC(), nil
}

func (h *GrpcHandlers) GetMessagesForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationId) (response *chatRPC.MessageArr, err error) {
	defer recoverHandler(&err)

	messages, err := h.ns.GetMessagesForCompany(ctx, input.CompanyId, input.ConversationId)
	panicIf(err)

	return model.MessageArr(messages).ToRPC(), nil
}

func (h *GrpcHandlers) SendMessage(ctx context.Context, input *chatRPC.Message) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	message := model.NewMessage(input)
	err = h.ns.SendMessage(ctx, message)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) SendUnverifiedMessage(ctx context.Context, data *chatRPC.Message) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	message := model.NewMessage(data)
	err = this.ns.SendUnverifiedMessage(ctx, message)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) AddParticipants(ctx context.Context, input *chatRPC.AddParticipantsRequest) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation, err := h.ns.AddParticipantsToConversation(ctx, input.ConversationId, input.Participants)
	panicIf(err)

	return conversation.ToRPC(), nil
}
func (h *GrpcHandlers) AddParticipantsForCompany(ctx context.Context, input *chatRPC.AddParticipantsForCompanyRequest) (response *chatRPC.Conversation, err error) {
	defer recoverHandler(&err)

	conversation, err := h.ns.AddParticipantsToConversationForCompany(ctx, input.CompanyId, input.ConversationId, input.Participants)
	panicIf(err)

	return conversation.ToRPC(), nil
}

func (h *GrpcHandlers) LeaveConversation(ctx context.Context, input *chatRPC.Conversation) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.LeaveConversation(ctx, input.Id)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) LeaveConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.LeaveConversationForCompany(ctx, input.CompanyId, input.ConversationId)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) DeleteConversation(ctx context.Context, input *chatRPC.Conversation) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.DeleteConversation(ctx, input.Id)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) DeleteConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.DeleteConversationForCompany(ctx, input.CompanyId, input.ConversationId)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) SetConversationUnreadFlag(ctx context.Context, input *chatRPC.SetConversationUnreadFlagRequest) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.SetConversationUnreadFlag(ctx, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) SetConversationUnreadFlagForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationIdAndBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.SetConversationUnreadFlagForCompany(ctx, input.CompanyId, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) SearchInConversation(ctx context.Context, input *chatRPC.SearchInConversationRequest) (response *chatRPC.MessageArr, err error) {
	defer recoverHandler(&err)

	messages, err := h.ns.SearchInConversation(ctx, input.ConversationId, input.Query, input.File)
	panicIf(err)

	return model.MessageArr(messages).ToRPC(), nil
}
func (h *GrpcHandlers) SearchInConversationForCompany(ctx context.Context, input *chatRPC.SearchInConversationForCompanuRequest) (response *chatRPC.MessageArr, err error) {
	defer recoverHandler(&err)

	messages, err := h.ns.SearchInConversationForCompany(ctx, input.CompanyId, input.ConversationId, input.Query, input.File)
	panicIf(err)

	return model.MessageArr(messages).ToRPC(), nil
}

func (h *GrpcHandlers) MuteConversation(ctx context.Context, input *chatRPC.ConversationIdWithBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.MuteConversation(ctx, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) MuteConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationIdAndBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.MuteConversationForCompany(ctx, input.CompanyId, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) BlockConversation(ctx context.Context, input *chatRPC.ConversationIdWithBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.BlockConversation(ctx, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) BlockConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationIdAndBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.BlockConversationForCompany(ctx, input.CompanyId, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) ArchiveConversation(ctx context.Context, input *chatRPC.ConversationIdWithBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.ArchiveConversation(ctx, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) ArchiveConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationIdAndBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.ArchiveConversationForCompany(ctx, input.CompanyId, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) RenameConversation(ctx context.Context, input *chatRPC.ConversationIdWithString) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.RenameConversation(ctx, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) RenameConversationForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationIdAndString) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.RenameConversationForCompany(ctx, input.CompanyId, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) ChangeConversationAvatar(ctx context.Context, input *chatRPC.ConversationIdWithString) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.ChangeConversationAvatar(ctx, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) ChangeConversationAvatarForCompany(ctx context.Context, input *chatRPC.CompanyIdWithConversationIdAndString) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.ChangeConversationAvatarForCompany(ctx, input.CompanyId, input.ConversationId, input.Value)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) CreateReply(ctx context.Context, input *chatRPC.Reply) (response *chatRPC.Reply, err error) {
	defer recoverHandler(&err)

	reply, err := h.ns.CreateReply(ctx, model.NewReply(input))
	panicIf(err)

	return reply.ToRPC(), nil
}
func (h *GrpcHandlers) CreateReplyForCompany(ctx context.Context, input *chatRPC.CreateReplyForCompanyRequest) (response *chatRPC.Reply, err error) {
	defer recoverHandler(&err)

	reply, err := h.ns.CreateReplyForCompany(ctx, input.CompanyId, model.NewReply(input.Reply))
	panicIf(err)

	return reply.ToRPC(), nil
}

func (h *GrpcHandlers) UpdateReply(ctx context.Context, input *chatRPC.Reply) (response *chatRPC.Reply, err error) {
	defer recoverHandler(&err)

	reply, err := h.ns.UpdateReply(ctx, model.NewReply(input))
	panicIf(err)

	return reply.ToRPC(), nil
}
func (h *GrpcHandlers) UpdateReplyForCompany(ctx context.Context, input *chatRPC.CreateReplyForCompanyRequest) (response *chatRPC.Reply, err error) {
	defer recoverHandler(&err)

	reply, err := h.ns.UpdateReplyForCompany(ctx, input.CompanyId, model.NewReply(input.Reply))
	panicIf(err)

	return reply.ToRPC(), nil
}

func (h *GrpcHandlers) DeleteReply(ctx context.Context, input *chatRPC.Reply) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.DeleteReply(ctx, input.Id)
	panicIf(err)

	return EMPTY, nil
}
func (h *GrpcHandlers) DeleteReplyForCompany(ctx context.Context, input *chatRPC.CompanyIdWithId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = h.ns.DeleteReplyForCompany(ctx, input.CompanyId, input.Id)
	panicIf(err)

	return EMPTY, nil
}

func (h *GrpcHandlers) GetMyReplies(ctx context.Context, input *chatRPC.StringValue) (response *chatRPC.ReplyArr, err error) {
	defer recoverHandler(&err)

	replies, err := h.ns.GetMyReplies(ctx, input.Value)
	panicIf(err)

	return model.ReplyArr(replies).ToRPC(), nil
}
func (h *GrpcHandlers) GetMyRepliesForCompany(ctx context.Context, input *chatRPC.CompanyIdWithString) (response *chatRPC.ReplyArr, err error) {
	defer recoverHandler(&err)

	replies, err := h.ns.GetMyRepliesForCompany(ctx, input.CompanyId, input.Value)
	panicIf(err)

	return model.ReplyArr(replies).ToRPC(), nil
}

func (this *GrpcHandlers) CreateLabel(ctx context.Context, data *chatRPC.ConversationLabel) (response *chatRPC.ConversationLabel, err error) {
	defer recoverHandler(&err)

	label, err := this.ns.CreateLabel(ctx, model.NewConversationLabel(data))
	panicIf(err)

	return label.ToRPC(), nil
}
func (this *GrpcHandlers) CreateLabelForCompany(ctx context.Context, data *chatRPC.CompanyIdWithConversationLabel) (response *chatRPC.ConversationLabel, err error) {
	defer recoverHandler(&err)

	label, err := this.ns.CreateLabelForCompany(ctx, data.CompanyId, model.NewConversationLabel(data.Label))
	panicIf(err)

	return label.ToRPC(), nil
}

func (this *GrpcHandlers) UpdateLabel(ctx context.Context, data *chatRPC.ConversationLabel) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.UpdateLabel(ctx, model.NewConversationLabel(data))
	panicIf(err)

	return EMPTY, nil
}
func (this *GrpcHandlers) UpdateLabelForCompany(ctx context.Context, data *chatRPC.CompanyIdWithConversationLabel) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.UpdateLabelForCompany(ctx, data.CompanyId, model.NewConversationLabel(data.Label))
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) DeleteLabel(ctx context.Context, data *chatRPC.ConversationLabel) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.DeleteLabel(ctx, data.Id)
	panicIf(err)

	return EMPTY, nil
}
func (this *GrpcHandlers) DeleteLabelForCompany(ctx context.Context, data *chatRPC.CompanyIdWithId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.DeleteLabelForCompany(ctx, data.CompanyId, data.Id)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) GetAllLabel(ctx context.Context, data *chatRPC.Empty) (response *chatRPC.ConversationLabelArr, err error) {
	defer recoverHandler(&err)

	labels, err := this.ns.GetAllLabel(ctx)
	panicIf(err)

	return model.ConversationLabelArr(labels).ToRPC(), nil
}
func (this *GrpcHandlers) GetAllLabelForCompany(ctx context.Context, data *chatRPC.CompanyId) (response *chatRPC.ConversationLabelArr, err error) {
	defer recoverHandler(&err)

	labels, err := this.ns.GetAllLabelForCompany(ctx, data.CompanyId)
	panicIf(err)

	return model.ConversationLabelArr(labels).ToRPC(), nil
}

func (this *GrpcHandlers) AddLabelToConversation(ctx context.Context, data *chatRPC.ConversationIdWithLabelId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.AddLabelToConversation(ctx, data.ConversationId, data.LabelId)
	panicIf(err)

	return EMPTY, nil
}
func (this *GrpcHandlers) AddLabelToConversationForCompany(ctx context.Context, data *chatRPC.CompanyIdWithConversationIdAndLabelId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.AddLabelToConversationForCompany(ctx, data.CompanyId, data.ConversationId, data.LabelId)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) RemoveLabelFromConversation(ctx context.Context, data *chatRPC.ConversationIdWithLabelId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.RemoveLabelFromConversation(ctx, data.ConversationId, data.LabelId)
	panicIf(err)

	return EMPTY, nil
}
func (this *GrpcHandlers) RemoveLabelFromConversationForCompany(ctx context.Context, data *chatRPC.CompanyIdWithConversationIdAndLabelId) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.RemoveLabelFromConversationForCompany(ctx, data.CompanyId, data.ConversationId, data.LabelId)
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) ReportConversation(ctx context.Context, data *chatRPC.ConversationIdWithString) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.ReportConversation(ctx, &model.Report{ConversationId: data.ConversationId, Text: data.Value})
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) ReportConversationForCompany(ctx context.Context, data *chatRPC.ReportConversationForCompanyRequest) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.ReportConversationForCompany(ctx, &model.Report{
		CompanyId:      data.CompanyId,
		ConversationId: data.ConversationId,
		Text:           data.Text,
	})
	panicIf(err)

	return EMPTY, nil
}

func (this *GrpcHandlers) IsUserLive(ctx context.Context, data *chatRPC.IsUserLiveRequest) (response *chatRPC.BoolValue, err error) {
	defer recoverHandler(&err)

	isLive, err := this.ns.IsUserLive(ctx, data.UserId)
	panicIf(err)

	return &chatRPC.BoolValue{Value: isLive}, nil
}
func (this *GrpcHandlers) SetParticipantOffline(ctx context.Context, data *chatRPC.BoolValue) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.SetParticipantOffline(ctx, data.Value)
	panicIf(err)

	return &chatRPC.Empty{}, nil
}
func (this *GrpcHandlers) SetParticipantOfflineForCompany(ctx context.Context, data *chatRPC.CompanyIdWithBool) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.SetParticipantOfflineForCompany(ctx, data.CompanyId, data.Value)
	panicIf(err)

	return &chatRPC.Empty{}, nil
}

func (this *GrpcHandlers) BlockConversetionByParticipants(ctx context.Context, data *chatRPC.BlockRequest) (response *chatRPC.Empty, err error) {
	defer recoverHandler(&err)

	err = this.ns.BlockConversetionByParticipants(ctx, data.SenderID, data.TargetID, data.Value)
	panicIf(err)

	return &chatRPC.Empty{}, nil
}
