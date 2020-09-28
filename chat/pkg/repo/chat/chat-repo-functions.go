package chat

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"golang.org/x/net/context"
)

func (r *ChatRepo) CreateConversation(ctx context.Context, conversation *model.Conversation) (*model.Conversation, error) {
	conversation.Id = bson.NewObjectId()
	err := r.conversations.Insert(conversation)
	if err != nil {
		return nil, err
	}
	return r.GetConversation(ctx, conversation.Id.Hex())
}

func (r *ChatRepo) GetConversation(ctx context.Context, id string) (*model.Conversation, error) {
	var conversation model.Conversation
	err := r.conversations.FindId(bson.ObjectIdHex(id)).One(&conversation)
	return &conversation, err
}

func (r *ChatRepo) GetLastMessage(ctx context.Context, conversationId, senderId string, from, to int64) (*model.Message, error) {
	var message model.Message
	cond := bson.M{
		"type":            chatRPC.MessageType_UserMessage.String(),
		"conversation_id": conversationId,
		"sender_id":       senderId,
		"timestamp": bson.M{
			"$gt": from,
			"$lt": to,
		},
	}
	err := r.messages.Find(cond).Sort("-timestamp").Limit(1).One(&message)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *ChatRepo) FindConversationsByParticipants(ctx context.Context, participants []*model.Participant) ([]*model.Conversation, error) {
	participantIds := make([]string, len(participants))
	for i, p := range participants {
		participantIds[i] = p.Id
	}
	var conversations []*model.Conversation
	err := r.conversations.Find(bson.M{
		"participants.id": bson.M{"$all": participantIds},
	}).All(&conversations)
	return conversations, err
}

func (r *ChatRepo) FindConversationsWithOnlyParticipants(ctx context.Context, participants []*model.Participant) ([]*model.Conversation, error) {
	conversations, err := r.FindConversationsByParticipants(ctx, participants)
	if err != nil {
		return nil, err
	}

	filtered := make([]*model.Conversation, 0)
	for _, conversation := range conversations {
		if len(conversation.Participants) == len(participants) {
			filtered = append(filtered, conversation)
		}
	}
	return filtered, nil
}

func (r *ChatRepo) GetConversationsOf(ctx context.Context, id string, filter *model.ConversationFilter) ([]*model.Conversation, error) {
	participantFilter := bson.M{
		"id": id,
	}
	matchQuery := bson.M{"participants": bson.M{"$elemMatch": participantFilter}}
	if filter != nil {
		if filter.Category == "All" {
			participantFilter["archived"] = false
			participantFilter["blocked"] = false
		}
		if filter.Category == "Unread" {
			participantFilter["unread"] = true
		}
		if filter.Category == "Archived" {
			participantFilter["archived"] = true
		}
		if filter.LabelId != "" {
			participantFilter["labels"] = bson.ObjectIdHex(filter.LabelId)
		}
		if filter.ParticipantId != "" {
			matchQuery["participants.id"] = filter.ParticipantId
		}
		if filter.Text != "" {
			var ids []map[string]bson.ObjectId
			err := r.conversations.Find(bson.M{
				"participants.id": id,
			}).Select(bson.M{"_id": 1}).All(&ids)
			if err == nil {
				convIds := make([]string, len(ids))
				for i, id := range ids {
					convIds[i] = id["_id"].Hex()
				}
				ids := make([]string, 0)
				err = r.messages.Find(bson.M{
					"conversation_id": bson.M{"$in": convIds},
					"text": bson.M{
						"$regex":   ".*" + filter.Text + ".*",
						"$options": "$i",
					},
				}).Select(bson.M{"conversation_id": 1}).Distinct("conversation_id", &ids)
				if err == nil {
					convObjectIds := make([]bson.ObjectId, len(ids))
					for i, id := range ids {
						convObjectIds[i] = bson.ObjectIdHex(id)
					}
					matchQuery["_id"] = bson.M{"$in": convObjectIds}
				}
			}
		}
	}

	var conversations []*model.Conversation
	err := r.conversations.Pipe([]bson.M{
		{
			"$match": matchQuery,
		}, {
			"$addFields": bson.M{
				"id_hex": bson.M{"$toString": "$_id"},
				"me": bson.M{
					"$arrayElemAt": []interface{}{"$participants", bson.M{
						"$indexOfArray": []interface{}{"$participants.id", id},
					}},
				},
			},
		}, {
			"$lookup": bson.M{
				"from":         "messages",
				"localField":   "id_hex",
				"foreignField": "conversation_id",
				"as":           "last_message",
			},
		}, {
			"$addFields": bson.M{
				"last_message": bson.M{
					"$arrayElemAt": []interface{}{bson.M{
						"$filter": bson.M{
							"input": "$last_message",
							"as":    "msg",
							"cond": bson.M{
								"$and": []interface{}{
									bson.M{
										"$cond": bson.M{
											"if":   "$me.has_left",
											"then": bson.M{"$lt": []interface{}{"$$msg.timestamp", "$me.left_timestamp"}},
											"else": true,
										},
									},
									bson.M{
										"$gt": []interface{}{"$$msg.timestamp", "$me.deleted_timestamp"},
									},
								},
							},
						},
					}, -1},
				},
			},
		}, {
			"$sort": bson.M{
				"last_message.timestamp": -1,
			},
		},
	}).All(&conversations)

	return conversations, err
}

func (r *ChatRepo) GetMessages(ctx context.Context, conversationId string, from, to int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.messages.Find(bson.M{
		"conversation_id": conversationId,
		"timestamp": bson.M{
			"$gt": from,
			"$lt": to,
		},
	}).All(&messages)
	return messages, err
}

func (r *ChatRepo) GetMessage(ctx context.Context, messageId string) (*model.Message, error) {
	var message model.Message
	err := r.messages.FindId(bson.ObjectIdHex(messageId)).One(&message)
	return &message, err
}

func (r *ChatRepo) SetMessageStatus(ctx context.Context, messageId, senderId, status string) error {
	return r.messages.UpdateId(bson.ObjectIdHex(messageId), bson.M{
		"$addToSet": bson.M{
			status + "_by": senderId,
		},
	})
}

func (r *ChatRepo) UpdateConversation(ctx context.Context, conversation *model.Conversation) error {
	return r.conversations.UpdateId(conversation.Id, conversation)
}

func (r *ChatRepo) UpdateParticipant(ctx context.Context, conversationId string, participant *model.Participant) error {
	return r.conversations.Update(bson.M{
		"_id":             bson.ObjectIdHex(conversationId),
		"participants.id": participant.Id,
	}, bson.M{
		"$set": bson.M{
			"participants.$": participant,
		},
	})
}

func (r *ChatRepo) SendMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	message.Id = bson.NewObjectId()
	err := r.messages.Insert(message)

	return message, err
}

func (r *ChatRepo) SetUnreadFor(ctx context.Context, conversationId, participantId string, value bool) error {
	return r.conversations.Update(bson.M{
		"_id":             bson.ObjectIdHex(conversationId),
		"participants.id": participantId,
	}, bson.M{
		"$set": bson.M{
			"participants.$.unread": value,
		},
	})
}

func (r *ChatRepo) SetUnreadForAll(ctx context.Context, conversationId string, value bool) error {
	return r.conversations.UpdateId(bson.ObjectIdHex(conversationId), bson.M{
		"$set": bson.M{
			"participants.$[].unread": value,
		},
	})
}

func (r *ChatRepo) GetUnreadConversationsStatus(participantId string) (int, error) {
	return r.conversations.Find(bson.M{
		"participants": bson.M{
			"$elemMatch": bson.M{
				"id":       participantId,
				"unread":   true,
				"has_left": false,
			},
		},
	}).Count()
}

func (r *ChatRepo) SearchInConversation(ctx context.Context, conversationId, query, file string, from, to int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.messages.Find(bson.M{
		"conversation_id": conversationId,
		"text": bson.M{
			"$regex":   query,
			"$options": "$i",
		},
		"files.file_name": bson.M{
			"$regex": file,
		},
		"timestamp": bson.M{
			"$gt": from,
			"$lt": to,
		},
	}).All(&messages)
	return messages, err
}

func (r *ChatRepo) MuteConversation(ctx context.Context, participantId, conversationId string, value bool) error {
	return r.conversations.Update(bson.M{
		"_id":             bson.ObjectIdHex(conversationId),
		"participants.id": participantId,
	}, bson.M{
		"$set": bson.M{
			"participants.$.muted": value,
		},
	})
}

func (r *ChatRepo) BlockConversation(ctx context.Context, participantId, conversationId string, value bool) error {
	return r.conversations.Update(bson.M{
		"_id":             bson.ObjectIdHex(conversationId),
		"participants.id": participantId,
	}, bson.M{
		"$set": bson.M{
			"participants.$.blocked": value,
		},
	})
}

func (r *ChatRepo) ArchiveConversation(ctx context.Context, participantId, conversationId string, value bool) error {
	return r.conversations.Update(bson.M{
		"_id":             bson.ObjectIdHex(conversationId),
		"participants.id": participantId,
	}, bson.M{
		"$set": bson.M{
			"participants.$.archived": value,
		},
	})
}

func (r *ChatRepo) RenameConversation(ctx context.Context, conversationId, name string) error {
	return r.conversations.UpdateId(bson.ObjectIdHex(conversationId), bson.M{
		"$set": bson.M{
			"name": name,
		},
	})
}
func (r *ChatRepo) ChangeConversationAvatar(ctx context.Context, conversationId, avatar string) error {
	return r.conversations.UpdateId(bson.ObjectIdHex(conversationId), bson.M{
		"$set": bson.M{
			"avatar": avatar,
		},
	})
}

func (r *ChatRepo) GetReply(ctx context.Context, id string) (*model.Reply, error) {
	var reply model.Reply
	err := r.replies.FindId(bson.ObjectIdHex(id)).One(&reply)
	return &reply, err
}

func (r *ChatRepo) CreateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	reply.Id = bson.NewObjectId()
	err := r.replies.Insert(reply)
	if err != nil {
		return nil, err
	}
	return r.GetReply(ctx, reply.Id.Hex())
}

func (r *ChatRepo) UpdateReply(ctx context.Context, reply *model.Reply) (*model.Reply, error) {
	err := r.replies.UpdateId(reply.Id, reply)
	if err != nil {
		return nil, err
	}
	return r.GetReply(ctx, reply.Id.Hex())
}

func (r *ChatRepo) GetUserReplies(ctx context.Context, ownerId, query string) ([]*model.Reply, error) {
	var replies []*model.Reply

	condition := bson.M{
		"owner_id": ownerId,
	}
	if query != "" {
		condition["$or"] = []interface{}{
			bson.M{"title": bson.M{"$regex": ".*" + query + ".*"}},
			bson.M{"text": bson.M{"$regex": ".*" + query + ".*"}},
		}
	}
	err := r.replies.Find(condition).All(&replies)

	return replies, err
}

func (r *ChatRepo) DeleteReply(ctx context.Context, replyId string) error {
	return r.replies.RemoveId(bson.ObjectIdHex(replyId))
}

func (r *ChatRepo) CreateLabel(ctx context.Context, label *model.ConversationLabel) (*model.ConversationLabel, error) {
	label.Id = bson.NewObjectId()
	err := r.labels.Insert(label)
	return label, err
}

func (r *ChatRepo) DeleteLabel(ctx context.Context, ownerId, labelId string) error {
	labelObjectId := bson.ObjectIdHex(labelId)
	err := r.labels.Remove(bson.M{"owner_id": ownerId, "_id": labelObjectId})
	if err != nil {
		return err
	}
	_, err = r.conversations.UpdateAll(bson.M{
		"participants.id":     ownerId,
		"participants.labels": labelObjectId,
	}, bson.M{
		"$pull": bson.M{
			"participants.$.labels": labelObjectId,
		},
	})
	return err
}

func (r *ChatRepo) UpdateLabel(ctx context.Context, label *model.ConversationLabel) error {
	return r.labels.UpdateId(label.Id, label)
}

func (r *ChatRepo) GetAllLabel(ctx context.Context, ownerId string) ([]*model.ConversationLabel, error) {
	var labels []*model.ConversationLabel
	err := r.labels.Pipe([]bson.M{
		{
			"$match": bson.M{
				"owner_id": ownerId,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "conversations",
				"localField":   "_id",
				"foreignField": "participants.labels",
				"as":           "conversations",
			},
		}, {
			"$addFields": bson.M{
				"count": bson.M{"$size": "$conversations"},
			},
		},
	}).All(&labels)

	return labels, err
}

func (r *ChatRepo) ReportConversation(ctx context.Context, report *model.Report) error {
	return r.reports.Insert(report)
}
