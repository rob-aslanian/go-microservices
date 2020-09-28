package model

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
)

type Conversation struct {
	Id            bson.ObjectId   `bson:"_id"`
	Name          string          `bson:"name"`
	Participants  []*Participant  `bson:"participants"`
	Avatar        string          `bson:"avatar"`
	IsGroup       bool            `bson:"is_group"`
	LastMessage   *Message        `bson:"last_message"`
	Unread        bool            `bson:"-"`
	Archived      bool            `bson:"-"`
	Muted         bool            `bson:"-"`
	Blocked       bool            `bson:"-"`
	HasLeft       bool            `bson:"-"`
	LeftTimestamp int64           `bson:"-"`
	DeletedTime   int64           `bson:"-"`
	Labels        []bson.ObjectId `bson:"-"`
}

func NewConversation(c *chatRPC.Conversation) *Conversation {
	participants := make([]*Participant, len(c.Participants))
	for i, p := range c.Participants {
		participants[i] = NewParticipant(p)
	}

	conv := Conversation{
		Name:         c.Name,
		Avatar:       c.Avatar,
		Participants: participants,
		IsGroup:      c.IsGroup,
	}
	if c.Id != "" {
		conv.Id = bson.ObjectIdHex(c.Id)
	}

	return &conv
}

func (c *Conversation) ToRPC() *chatRPC.Conversation {
	participants := make([]*chatRPC.Participant, len(c.Participants))
	for i, p := range c.Participants {
		participants[i] = p.ToRPC()
	}

	labels := make([]string, len(c.Labels))
	for i, l := range c.Labels {
		labels[i] = l.Hex()
	}

	var lastMessage *chatRPC.Message
	if c.LastMessage != nil {
		lastMessage = c.LastMessage.ToRPC()
	}

	return &chatRPC.Conversation{
		Id:           c.Id.Hex(),
		Name:         c.Name,
		Avatar:       c.Avatar,
		Participants: participants,
		IsGroup:      c.IsGroup,
		LastMessage:  lastMessage,

		Unread:   c.Unread,
		Muted:    c.Muted,
		Blocked:  c.Blocked,
		Archived: c.Archived,
		HasLeft:  c.HasLeft,
		Labels:   labels,
	}
}

type ConversationArr []*Conversation

func (conversations ConversationArr) ToRPC() *chatRPC.ConversationArr {
	rpcConversations := make([]*chatRPC.Conversation, len(conversations))
	for i, c := range conversations {
		rpcConversations[i] = c.ToRPC()
	}
	return &chatRPC.ConversationArr{Conversations: rpcConversations}
}

type ConversationFilter struct {
	Category      string
	LabelId       string
	ParticipantId string
	Text          string
}
