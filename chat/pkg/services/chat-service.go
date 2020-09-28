package services

import (
	"io"

	"github.com/globalsign/mgo"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/tracer"
	"golang.org/x/net/context"
)

type AuthClient interface {
	GetUserId(context.Context, string) (string, error)
}

type ChatValidator interface {
}

type LiveConnectionsInterface interface {
	SendTo(ctx context.Context, userId string, data interface{}) bool
	IsLive(userId string) bool
	SetParticipantOffline(participantId string, status bool)
}

type ChatRepo interface {
	CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error)
	UpdateConversation(ctx context.Context, conversation *model.Conversation) error
	UpdateParticipant(ctx context.Context, conversationId string, participant *model.Participant) error
	GetConversation(ctx context.Context, id string) (*model.Conversation, error)
	GetLastMessage(ctx context.Context, conversationId, senderId string, from, to int64) (*model.Message, error)
	GetConversationsOf(ctx context.Context, id string, filter *model.ConversationFilter) ([]*model.Conversation, error)
	SendMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	FindConversationsWithOnlyParticipants(ctx context.Context, participants []*model.Participant) ([]*model.Conversation, error)
	SetUnreadFor(ctx context.Context, conversationId, participantId string, value bool) error
	SetUnreadForAll(ctx context.Context, conversationId string, value bool) error
	GetUnreadConversationsStatus(participantId string) (int, error)
	SearchInConversation(ctx context.Context, conversationId, query, file string, from, to int64) ([]*model.Message, error)
	GetMessages(ctx context.Context, conversationId string, from, to int64) ([]*model.Message, error)
	GetMessage(ctx context.Context, messageId string) (*model.Message, error)
	SetMessageStatus(ctx context.Context, messageId, senderId, status string) error
	MuteConversation(ctx context.Context, participantId, conversationId string, value bool) error
	BlockConversation(ctx context.Context, participantId, conversationId string, value bool) error
	ArchiveConversation(ctx context.Context, participantId, conversationId string, value bool) error
	RenameConversation(ctx context.Context, conversationId, name string) error
	ChangeConversationAvatar(ctx context.Context, conversationId, avatar string) error

	GetReply(ctx context.Context, id string) (*model.Reply, error)
	CreateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error)
	UpdateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error)
	GetUserReplies(ctx context.Context, userId, query string) ([]*model.Reply, error)
	DeleteReply(ctx context.Context, replyId string) error

	CreateLabel(ctx context.Context, label *model.ConversationLabel) (*model.ConversationLabel, error)
	DeleteLabel(ctx context.Context, userId, labelId string) error
	UpdateLabel(ctx context.Context, label *model.ConversationLabel) error
	GetAllLabel(ctx context.Context, userId string) ([]*model.ConversationLabel, error)

	ReportConversation(ctx context.Context, report *model.Report) error

	SaveFile(ctx context.Context, input io.Reader, name string, md interface{}) (*mgo.GridFile, error)
	ReadFile(ctx context.Context, fileId string) (*mgo.GridFile, error)
}

type AdminManager interface {
	GetAdminLevelFor(ctx context.Context, companyId string) (string, error)
	GetAllFriendshipsID(ctx context.Context) ([]string, error)
	GetBlockedIDs(ctx context.Context) ([]string, error)
	IsBlockedByUser(ctx context.Context, userID string) (bool, error)
	IsBlockedCompanyByUser(ctx context.Context, companyID string) (bool, error)
}

type OfflineRepo interface {
	SetOffline(id string, isOffline bool) error
	IsOffline(id string) (bool, error)
}

type ChatService struct {
	repo            ChatRepo
	auth            AuthClient
	liveConnections LiveConnectionsInterface
	admin           AdminManager
	validator       ChatValidator
	tracer          *tracing.Tracer
	offlineRepo     OfflineRepo
}

func NewChatService(repo ChatRepo, validator ChatValidator, auth AuthClient, lc LiveConnectionsInterface, am AdminManager, tracer *tracing.Tracer) *ChatService {
	return &ChatService{
		repo:            repo,
		validator:       validator,
		auth:            auth,
		liveConnections: lc,
		admin:           am,
		tracer:          tracer,
	}
}
