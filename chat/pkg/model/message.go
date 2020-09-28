package model

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
)

type Message struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	Type           string        `bson:"type" json:"type"`
	ConversationId string        `bson:"conversation_id" json:"conversation_id"`
	SenderId       string        `bson:"sender_id" json:"sender_id"`
	Text           string        `bson:"text" json:"text"`
	// File           string        `bson:"file" json:"file"`
	// FileName       string        `bson:"file_name" json:"file_name"`
	// FileSize       int64         `bson:"file_size" json:"file_size"`
	Files      []File   `bson:"files" json:"files"`
	ReceivedBy []string `bson:"received_by" json:"received_by"`
	SeenBy     []string `bson:"seen_by" json:"seen_by"`
	Timestamp  int64    `bson:"timestamp" json:"timestamp"`
}

func NewMessage(m *chatRPC.Message) *Message {
	message := &Message{
		Type:           m.Type.String(),
		ConversationId: m.ConversationId,
		SenderId:       m.SenderId,
		Text:           m.Text,
		Files:          make([]File, 0, len(m.Files)),
		// File:           m.File,
		// FileName:       m.FileName,
		Timestamp: m.Timestamp,
	}
	if m.Id != "" {
		message.Id = bson.ObjectIdHex(m.Id)
	}

	for _, f := range message.Files {
		message.Files = append(message.Files, File{
			File:     f.File,
			FileName: f.FileName,
		})
	}

	return message
}

func (m *Message) ToRPC() *chatRPC.Message {
	mes := chatRPC.Message{
		Id:             m.Id.Hex(),
		Type:           chatRPC.MessageType(chatRPC.MessageType_value[m.Type]),
		ConversationId: m.ConversationId,
		SenderId:       m.SenderId,
		Text:           m.Text,
		Files:          make([]*chatRPC.File, 0, len(m.Files)),
		// File:           m.File,
		// FileName:       m.FileName,
		// FileSize:       m.FileSize,
		ReceivedBy: m.ReceivedBy,
		SeenBy:     m.SeenBy,
		Timestamp:  m.Timestamp,
	}

	for _, f := range m.Files {
		mes.Files = append(mes.Files, f.ToRPC())
	}

	return &mes
}

type MessageArr []*Message

func (messages MessageArr) ToRPC() *chatRPC.MessageArr {
	rpcMessages := make([]*chatRPC.Message, len(messages))
	for i, m := range messages {
		rpcMessages[i] = m.ToRPC()
	}

	return &chatRPC.MessageArr{Messages: rpcMessages}
}

type TotalUnreadCountMessage struct {
	Type             string `json:"type"`
	TotalUnreadCount int    `json:"total_unread_count"`
}

func NewTotalUnreadCountMessage(count int) *TotalUnreadCountMessage {
	return &TotalUnreadCountMessage{Type: chatRPC.MessageType_TotalUnreadCount.String(), TotalUnreadCount: count}
}

type File struct {
	File     string `bson:"file" json:"file"`
	FileName string `bson:"file_name" json:"file_name"`
	FileSize int64  `bson:"file_size" json:"file_size"`
}

func (f File) ToRPC() *chatRPC.File {
	return &chatRPC.File{
		File:     f.File,
		FileName: f.FileName,
		FileSize: f.FileSize,
	}
}
