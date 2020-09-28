package services

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func contextFor(token string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("token", token))
}

type ChatRepoMock struct {
	ChatRepo
	conversations       map[string]*model.Conversation
	conversationsByuser map[string][]*model.Conversation
	replies             map[string]*model.Reply
	repliesByUser       map[string][]*model.Reply
}

func (mock ChatRepoMock) CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	conversation.Id = bson.NewObjectId()
	return conversation, nil
}

func (mock ChatRepoMock) GetConversation(ctx context.Context, id string) (*model.Conversation, error) {
	if c, ok := mock.conversations[id]; ok {
		return c, nil
	}
	return nil, errors.New("Not found")
}

func (mock ChatRepoMock) GetConversationsOf(ctx context.Context, id string) ([]*model.Conversation, error) {
	if c, ok := mock.conversationsByuser[id]; ok {
		return c, nil
	}
	return nil, errors.New("Not found")
}

func (mock ChatRepoMock) RenameConversation(ctx context.Context, conversationId, name string) error {
	return nil
}

func (mock ChatRepoMock) UpdateConversation(ctx context.Context, conversation *model.Conversation) error {
	return nil
}

func (mock ChatRepoMock) UpdateParticipant(ctx context.Context, conversationId string, participant *model.Participant) error {
	return nil
}

func (mock ChatRepoMock) SendMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	return message, nil
}

func (mock ChatRepoMock) SetUnreadForAll(ctx context.Context, conversationId string, value bool) error {
	return nil
}

func (mock ChatRepoMock) FindConversationsWithOnlyParticipants(ctx context.Context, participants []*model.Participant) ([]*model.Conversation, error) {
	return nil, nil
}

func (mock ChatRepoMock) CreateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	reply.Id = bson.NewObjectId()
	mock.repliesByUser[reply.OwnerId] = append(mock.repliesByUser[reply.OwnerId], reply)
	return reply, nil
}

func (mock ChatRepoMock) GetReply(ctx context.Context, replyId string) (*model.Reply, error) {
	if r, ok := mock.replies[replyId]; ok {
		return r, nil
	}
	return nil, mgo.ErrNotFound
}

func (mock ChatRepoMock) GetUserReplies(ctx context.Context, userId string) ([]*model.Reply, error) {
	return mock.repliesByUser[userId], nil
}

func (mock ChatRepoMock) DeleteReply(ctx context.Context, replyId string) error {
	return nil
}

type AuthClientMock struct {
	AuthClient
	users map[string]string
}

func (mock AuthClientMock) GetUserId(ctx context.Context, token string) (string, error) {
	if user, ok := mock.users[token]; ok {
		return user, nil
	}
	return "", errors.New("Not authenticated")
}

type LiveConnectionsMock struct {
	LiveConnectionsInterface
	lives    map[string]bool
	messages map[string][]interface{}
}

func (mock LiveConnectionsMock) IsLive(userId string) bool {
	return mock.lives[userId]
}

func (mock LiveConnectionsMock) SendTo(userId string, data interface{}) bool {
	mock.messages[userId] = append(mock.messages[userId], data)
	return true
}

type AdminMock struct {
	level string
}

func (mock AdminMock) GetAdminLevelFor(ctx context.Context, companyId string) (string, error) {
	return mock.level, nil
}
