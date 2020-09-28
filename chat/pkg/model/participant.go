package model

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
)

type Participant struct {
	Id               string          `bson:"id"`
	IsCompany        bool            `bson:"is_company"`
	IsAdmin          bool            `bson:"is_admin"`
	Unread           bool            `bson:"unread"`
	Muted            bool            `bson:"muted"`
	Blocked          bool            `bson:"blocked"`
	Archived         bool            `bson:"archived"`
	HasLeft          bool            `bson:"has_left"`
	LeftTimestamp    int64           `bson:"left_timestamp"`
	DeletedTimestamp int64           `bson:"deleted_timestamp"`
	Labels           []bson.ObjectId `bson:"labels"`
	IsActive         bool            `bson:"-"`
}

func NewParticipant(p *chatRPC.Participant) *Participant {
	return &Participant{
		Id:        p.Id,
		IsCompany: p.IsCompany,
		IsAdmin:   p.IsAdmin,
		Unread:    p.Unread,
	}
}

func (p *Participant) ToRPC() *chatRPC.Participant {
	return &chatRPC.Participant{
		Id:        p.Id,
		IsCompany: p.IsCompany,
		IsAdmin:   p.IsAdmin,
		Unread:    p.Unread,
		HasLeft:   p.HasLeft,
		IsActive:  p.IsActive,
	}
}

type ConversationLabel struct {
	Id      bson.ObjectId `bson:"_id"`
	OwnerId string        `bson:"owner_id"`
	Name    string        `bson:"name"`
	Color   string        `bson:"color"`
	Count   int32         `bson:"count"`
}

func NewConversationLabel(label *chatRPC.ConversationLabel) *ConversationLabel {
	l := ConversationLabel{
		Name:  label.Name,
		Color: label.Color,
	}
	if label.Id != "" {
		l.Id = bson.ObjectIdHex(label.Id)
	}
	return &l
}

func (this *ConversationLabel) ToRPC() *chatRPC.ConversationLabel {
	return &chatRPC.ConversationLabel{
		Id:    this.Id.Hex(),
		Name:  this.Name,
		Color: this.Color,
		Count: this.Count,
	}
}

type ConversationLabelArr []*ConversationLabel

func (arr ConversationLabelArr) ToRPC() *chatRPC.ConversationLabelArr {
	res := make([]*chatRPC.ConversationLabel, len(arr))
	for i, curr := range arr {
		res[i] = curr.ToRPC()
	}
	return &chatRPC.ConversationLabelArr{List: res}
}
