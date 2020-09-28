package grpc_handlers

import (
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"golang.org/x/net/context"
)

type ChatService interface {
	CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)
	CreateConversationForCompany(ctx context.Context, companyId string, conversation *model.Conversation) (*model.Conversation, error)
	CreateUnverifiedConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)
	GetMyConversations(ctx context.Context, filter *model.ConversationFilter) ([]*model.Conversation, error)
	GetMyConversationsForCompany(ctx context.Context, companyId string, filter *model.ConversationFilter) ([]*model.Conversation, error)
	GetConversation(ctx context.Context, conversationId string) (*model.Conversation, error)
	GetConversationForCompany(ctx context.Context, companyId, conversationId string) (*model.Conversation, error)
	GetActiveConnections(ctx context.Context) ([]string, error)
	SendMessage(ctx context.Context, message *model.Message) error
	SendUnverifiedMessage(ctx context.Context, message *model.Message) error
	AddParticipantsToConversation(ctx context.Context, conversationId string, participants []string) (*model.Conversation, error)
	AddParticipantsToConversationForCompany(ctx context.Context, companyId, conversationId string, participants []string) (*model.Conversation, error)
	LeaveConversation(ctx context.Context, conversationId string) error
	LeaveConversationForCompany(ctx context.Context, companyId, conversationId string) error
	DeleteConversation(ctx context.Context, conversationId string) error
	DeleteConversationForCompany(ctx context.Context, companyId, conversationId string) error
	SetConversationUnreadFlag(ctx context.Context, conversationId string, value bool) error
	SetConversationUnreadFlagForCompany(ctx context.Context, companyId, conversationId string, value bool) error
	SearchInConversation(ctx context.Context, conversationId, query, file string) ([]*model.Message, error)
	SearchInConversationForCompany(ctx context.Context, companyId, conversationId, query, file string) ([]*model.Message, error)
	GetMessages(ctx context.Context, conversationId string) ([]*model.Message, error)
	GetMessagesForCompany(ctx context.Context, companyId, conversationId string) ([]*model.Message, error)
	MuteConversation(ctx context.Context, conversationId string, value bool) error
	MuteConversationForCompany(ctx context.Context, companyId, conversationId string, value bool) error
	BlockConversetionByParticipants(ctx context.Context, senderID string, targetID string, value bool) error
	BlockConversation(ctx context.Context, conversationId string, value bool) error
	BlockConversationForCompany(ctx context.Context, companyId, conversationId string, value bool) error
	ArchiveConversation(ctx context.Context, conversationId string, value bool) error
	ArchiveConversationForCompany(ctx context.Context, companyId, conversationId string, value bool) error
	RenameConversation(ctx context.Context, conversationId, name string) error
	RenameConversationForCompany(ctx context.Context, companyId, conversationId, name string) error
	ChangeConversationAvatar(ctx context.Context, conversationId, avatar string) error
	ChangeConversationAvatarForCompany(ctx context.Context, companyId, conversationId, avatar string) error

	CreateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error)
	CreateReplyForCompany(ctx context.Context, companyId string, reply *model.Reply) (*model.Reply, error)
	UpdateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error)
	UpdateReplyForCompany(ctx context.Context, companyId string, reply *model.Reply) (*model.Reply, error)
	DeleteReply(ctx context.Context, replyId string) error
	DeleteReplyForCompany(ctx context.Context, companyId, replyId string) error
	GetMyReplies(ctx context.Context, query string) ([]*model.Reply, error)
	GetMyRepliesForCompany(ctx context.Context, companyId, query string) ([]*model.Reply, error)

	CreateLabel(ctx context.Context, label *model.ConversationLabel) (*model.ConversationLabel, error)
	CreateLabelForCompany(ctx context.Context, companyId string, label *model.ConversationLabel) (*model.ConversationLabel, error)
	DeleteLabel(ctx context.Context, labelId string) error
	DeleteLabelForCompany(ctx context.Context, companyId, labelId string) error
	UpdateLabel(ctx context.Context, label *model.ConversationLabel) error
	UpdateLabelForCompany(ctx context.Context, companyId string, label *model.ConversationLabel) error
	GetAllLabel(ctx context.Context) ([]*model.ConversationLabel, error)
	GetAllLabelForCompany(ctx context.Context, companyId string) ([]*model.ConversationLabel, error)
	AddLabelToConversation(ctx context.Context, conversationId, labelId string) error
	AddLabelToConversationForCompany(ctx context.Context, companyId, conversationId, labelId string) error
	RemoveLabelFromConversation(ctx context.Context, conversationId, labelId string) error
	RemoveLabelFromConversationForCompany(ctx context.Context, companyId, conversationId, labelId string) error

	ReportConversation(ctx context.Context, report *model.Report) error
	ReportConversationForCompany(ctx context.Context, report *model.Report) error

	IsUserLive(ctx context.Context, userId string) (bool, error)
	SetParticipantOffline(ctx context.Context, status bool) error
	SetParticipantOfflineForCompany(ctx context.Context, companyId string, status bool) error
}

type GrpcHandlers struct {
	ns ChatService
}

func NewGrpcHandlers(networkService ChatService) *GrpcHandlers {
	return &GrpcHandlers{ns: networkService}
}
