package model

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
)

type Reply struct {
	Id      bson.ObjectId `bson:"_id"`
	OwnerId string        `bson:"owner_id"`
	Title   string        `bson:"title"`
	Text    string        `bson:"text"`
	Files   []*ReplyFile  `bson:"files"`
}

func NewReply(r *chatRPC.Reply) *Reply {
	reply := &Reply{
		Title: r.Title,
		Text:  r.Text,
	}
	if r.Id != "" {
		reply.Id = bson.ObjectIdHex(r.Id)
	}
	reply.Files = make([]*ReplyFile, len(r.Files))
	for i, f := range r.Files {
		reply.Files[i] = NewReplyFile(f)
	}

	return reply
}

func (r *Reply) ToRPC() *chatRPC.Reply {
	reply := chatRPC.Reply{
		Id:    r.Id.Hex(),
		Title: r.Title,
		Text:  r.Text,
	}
	reply.Files = make([]*chatRPC.ReplyFile, len(r.Files))
	for i, f := range r.Files {
		reply.Files[i] = f.ToRPC()
	}

	return &reply
}

type ReplyArr []*Reply

func (replies ReplyArr) ToRPC() *chatRPC.ReplyArr {
	arr := make([]*chatRPC.Reply, len(replies))
	for i, r := range replies {
		arr[i] = r.ToRPC()
	}
	return &chatRPC.ReplyArr{Replies: arr}
}

type ReplyFile struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

func NewReplyFile(file *chatRPC.ReplyFile) *ReplyFile {
	return &ReplyFile{
		Id:   file.Id,
		Name: file.Name,
	}
}

func (this *ReplyFile) ToRPC() *chatRPC.ReplyFile {
	return &chatRPC.ReplyFile{
		Id:   this.Id,
		Name: this.Name,
	}
}
