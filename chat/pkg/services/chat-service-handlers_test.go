package services

//
// import (
// 	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
// 	"golang.org/x/net/context"
// 	"testing"
// )
//
// func TestChatService_CreateConversation(t *testing.T) {
// 	repoMock := ChatRepoMock{}
//
// 	authMock := AuthClientMock{users: map[string]string{
// 		"testtoken": "testuserid",
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, nil, AdminMock{"Admin"})
//
// 	ctx := contextFor("testtoken")
//
// 	t.Run("not authorized", func(t *testing.T) {
// 		defer func() {
// 			e := recover()
// 			if e == nil {
// 				t.Error("Auth error not fired")
// 			}
// 		}()
// 		service.CreateConversation(context.Background(), &model.Conversation{})
// 	})
// 	t.Run("chat with nobody", func(t *testing.T) {
// 		_, err := service.CreateConversation(ctx, &model.Conversation{})
//
// 		if err == nil {
// 			t.Error("Error was not returned for no participants")
// 		}
// 	})
//
// 	t.Run("chat with myself", func(t *testing.T) {
// 		_, err := service.CreateConversation(ctx, &model.Conversation{Participants: []*model.Participant{{Id: "testuserid"}}})
//
// 		if err == nil {
// 			t.Error("Error was not returned for only myself as participants")
// 		}
// 	})
//
// 	t.Run("pass scenario", func(t *testing.T) {
// 		c, err := service.CreateConversation(ctx, &model.Conversation{Participants: []*model.Participant{{Id: "testuserid2"}}})
//
// 		if err != nil {
// 			t.Error("Returned error ", err.Error())
// 		}
//
// 		if c == nil {
// 			t.Error("returned conversation is nil")
// 		}
//
// 		me := findParticipant("testuserid", c.Participants)
// 		if me == nil {
// 			t.Errorf("Conversation creator was not added to the participants")
// 		}
// 		if !me.IsAdmin {
// 			t.Errorf("conversation creator is not admin")
// 		}
// 	})
// }
//
// func TestChatService_GetConversation(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversations: map[string]*model.Conversation{
// 			"conv1": {
// 				Id:   "conv1",
// 				Name: "Conversation1",
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					},
// 				},
// 			},
// 			"conv2": {
// 				Id:   "conv2",
// 				Name: "Conversation2",
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user2",
// 						IsAdmin: true,
// 						Unread:  true,
// 					},
// 				},
// 			},
// 			"conv3": {
// 				Participants: []*model.Participant{
// 					{
// 						Id:     "user1",
// 						Unread: true,
// 					}, {
// 						Id:     "user2",
// 						Unread: false,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, nil, AdminMock{"Admin"})
// 	ctx := contextFor("token1")
//
// 	t.Run("get others conversations", func(t *testing.T) {
// 		defer func() {
// 			if e := recover(); e == nil {
// 				t.Error("Should panic because I am not participant")
// 			}
// 		}()
// 		_, err := service.GetConversation(contextFor("token2"), "conv1")
//
// 		if err == nil {
// 			t.Errorf("Error was not returned what requesting conversations where I am not participant")
// 		}
// 	})
//
// 	t.Run("get my conversation", func(t *testing.T) {
// 		conv, err := service.GetConversation(ctx, "conv1")
// 		if err != nil {
// 			t.Error("Error returned: ", err.Error())
// 		}
// 		if conv == nil {
// 			t.Error("Returned conversation is nil")
// 		}
// 		if conv.Name != "Conversation1" {
// 			t.Error("Returned different conversation")
// 		}
// 		if !conv.Unread {
// 			t.Error("Unread flag was not populated on conversation")
// 		}
// 	})
//
// 	t.Run("participant flags are correctly populated", func(t *testing.T) {
// 		conv, err := service.GetConversation(ctx, "conv3")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if !conv.Unread {
// 			t.Error("Unread flag should be true for the participant")
// 		}
//
// 		conv, err = service.GetConversation(contextFor("token2"), "conv3")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if conv.Unread {
// 			t.Error("Unread flag should be false for the participant")
// 		}
// 	})
//
// }
//
// func TestChatService_GetMyConversations(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversationsByuser: map[string][]*model.Conversation{
// 			"user1": {
// 				{
// 					Id:   "conv1",
// 					Name: "Conversation1",
// 					Participants: []*model.Participant{
// 						{
// 							Id:       "user1",
// 							IsAdmin:  true,
// 							Unread:   true,
// 							Archived: false,
// 						}, {
// 							Id: "user4",
// 						},
// 					},
// 				}, {
// 					Id:   "conv2",
// 					Name: "Conversation2",
// 					Participants: []*model.Participant{
// 						{
// 							Id:       "user1",
// 							IsAdmin:  true,
// 							Unread:   false,
// 							Archived: true,
// 						}, {
// 							Id: "user3",
// 						}, {
// 							Id: "user4",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 	}}
//
// 	liveMock := LiveConnectionsMock{lives: map[string]bool{
// 		"user3": true,
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 	t.Run("all", func(t *testing.T) {
// 		c, err := service.GetMyConversations(contextFor("token1"), "All", "")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(c) != 2 {
// 			t.Error("Num of conversations should be 2 but got ", len(c))
// 		}
// 	})
//
// 	t.Run("Archived", func(t *testing.T) {
// 		c, err := service.GetMyConversations(contextFor("token1"), "Archived", "")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(c) != 1 {
// 			t.Error("Num of conversations should be 1 but got ", len(c))
// 		}
// 		if c[0].Name != "Conversation2" {
// 			t.Error("Returned wrong conversations")
// 		}
// 	})
//
// 	t.Run("Unread", func(t *testing.T) {
// 		c, err := service.GetMyConversations(contextFor("token1"), "Unread", "")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(c) != 1 {
// 			t.Error("Num of conversations should be 1 but got ", len(c))
// 		}
// 		if c[0].Name != "Conversation1" {
// 			t.Error("Returned wrong conversations")
// 		}
// 	})
//
// 	t.Run("Active", func(t *testing.T) {
// 		c, err := service.GetMyConversations(contextFor("token1"), "Active", "")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(c) != 1 {
// 			t.Error("Num of conversations should be 1 but got ", len(c))
// 		}
// 		if c[0].Name != "Conversation2" {
// 			t.Error("Returned wrong conversations")
// 		}
// 	})
// }
//
// func TestChatService_RenameConversation(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversations: map[string]*model.Conversation{
// 			"conv1": {
// 				Id:   "conv1",
// 				Name: "Conversation1",
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					},
// 				},
// 			},
// 			"conv2": {
// 				Id:   "conv2",
// 				Name: "Conversation2",
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: false,
// 					}, {
// 						Id:      "user2",
// 						IsAdmin: true,
// 					},
// 				},
// 			},
// 			"conv3": {
// 				Participants: []*model.Participant{},
// 			},
// 		},
// 	}
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, nil, AdminMock{"Admin"})
//
// 	t.Run("no participant", func(t *testing.T) {
// 		defer func() {
// 			if e := recover(); e == nil {
// 				t.Error("Should panic because I am not participant")
// 			}
// 		}()
//
// 		err := service.RenameConversation(contextFor("token1"), "conv3", "name")
// 		if err == nil {
// 			t.Error("Was able to rename conversaiton where I am not participant")
// 		}
// 	})
//
// 	t.Run("no admin", func(t *testing.T) {
// 		err := service.RenameConversation(contextFor("token1"), "conv2", "name")
// 		if err == nil {
// 			t.Error("Was able to rename conversation where I am not admin")
// 		}
// 	})
//
// 	t.Run("admin", func(t *testing.T) {
// 		err := service.RenameConversation(contextFor("token2"), "conv2", "name")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 	})
// }
//
// func TestChatService_AddParticipantsToConversation(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversations: map[string]*model.Conversation{
// 			"private": {
// 				Id:      "private",
// 				Name:    "Conversation1",
// 				IsGroup: false,
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id: "user2",
// 					},
// 				},
// 			},
// 			"group": {
// 				Id:      "group",
// 				Name:    "Conversation2",
// 				IsGroup: true,
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id: "user2",
// 					}, {
// 						Id: "user3",
// 					},
// 				},
// 			},
// 			"groupWithLeft": {
// 				Id:      "groupWithLeft",
// 				Name:    "Conversation3",
// 				IsGroup: true,
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id:      "user2",
// 						HasLeft: true,
// 					}, {
// 						Id: "user3",
// 					},
// 				},
// 			},
// 			"noMoreInGroup": {
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						HasLeft: true,
// 					},
// 				},
// 			},
// 			"notInGroup": {
// 				Participants: []*model.Participant{},
// 			},
// 		},
// 	}
//
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	liveMock := LiveConnectionsMock{}
//
// 	service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 	ctx := contextFor("token1")
//
// 	t.Run("add to private conversation", func(t *testing.T) {
// 		conv, err := service.AddParticipantsToConversation(ctx, "private", []string{"user3"})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if !conv.IsGroup {
// 			t.Error("Returned conversation is not group")
// 		}
// 		if conv.Id == "private" {
// 			t.Error("Returned same conversation, should be created new for the group")
// 		}
// 		if len(conv.Participants) != 3 {
// 			t.Error("New participant was not added to the conversation")
// 		}
// 	})
//
// 	t.Run("add to group conversation", func(t *testing.T) {
// 		conv, err := service.AddParticipantsToConversation(ctx, "group", []string{"user4"})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if !conv.IsGroup {
// 			t.Error("Returned conversation is not group")
// 		}
// 		if conv.Id != "group" {
// 			t.Error("Created new conversations, participant should be added to the same group")
// 		}
// 		if len(conv.Participants) != 4 {
// 			t.Error("Participant was not added to the conversation")
// 		}
// 	})
//
// 	t.Run("add duplicate participants", func(t *testing.T) {
// 		conv, err := service.AddParticipantsToConversation(ctx, "group", []string{"user2", "user3", "user4"})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if conv.Id != "group" {
// 			t.Error("Created new conversations, participant should be added to the same group")
// 		}
// 		if len(conv.Participants) != 4 {
// 			t.Error("Participant was not added to the conversation")
// 		}
// 	})
//
// 	t.Run("add participant which has left", func(t *testing.T) {
// 		conv, err := service.AddParticipantsToConversation(ctx, "groupWithLeft", []string{"user2", "user4"})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(conv.Participants) != 4 {
// 			t.Error("Participant was not added to the conversation")
// 		}
// 		part := findParticipant("user2", conv.Participants)
// 		if part.HasLeft {
// 			t.Error("Participant was not readded to the conversations")
// 		}
// 	})
//
// 	t.Run("add to group which I have already left", func(t *testing.T) {
// 		_, err := service.AddParticipantsToConversation(ctx, "noMoreInGroup", []string{"user4"})
// 		if err == nil {
// 			t.Error("Could add participant to the group which I have already left")
// 		}
// 	})
//
// 	t.Run("add to group which I am not part of", func(t *testing.T) {
// 		defer func() {
// 			if e := recover(); e == nil {
// 				t.Error("Should panic because I am not participant")
// 			}
// 		}()
// 		_, err := service.AddParticipantsToConversation(ctx, "notInGroup", []string{"user4"})
// 		if err == nil {
// 			t.Error("Could add participant to the group which I am not part of")
// 		}
// 	})
// }
//
// func TestChatService_LeaveConversation(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversations: map[string]*model.Conversation{
// 			"conv": {
// 				Id:      "conv",
// 				Name:    "Conversation1",
// 				IsGroup: false,
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id: "user2",
// 					},
// 				},
// 			},
// 			"noMoreInGroup": {
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						HasLeft: true,
// 					},
// 				},
// 			},
// 			"notInGroup": {
// 				Participants: []*model.Participant{},
// 			},
// 		},
// 	}
//
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	liveMock := LiveConnectionsMock{}
//
// 	service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 	ctx := contextFor("token1")
//
// 	t.Run("not part of", func(t *testing.T) {
// 		defer func() {
// 			if e := recover(); e == nil {
// 				t.Error("Should panic because I am not participant")
// 			}
// 		}()
// 		err := service.LeaveConversation(ctx, "notInGroup")
// 		if err == nil {
// 			t.Error("I could leave conversation which I am not participant of")
// 		}
// 	})
//
// 	t.Run("have already left", func(t *testing.T) {
// 		err := service.LeaveConversation(ctx, "noMoreInGroup")
// 		if err == nil {
// 			t.Error("I could leave conversation which I have already left")
// 		}
// 	})
//
// 	t.Run("I am participant of", func(t *testing.T) {
// 		err := service.LeaveConversation(ctx, "conv")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 	})
// }
//
// func TestChatService_DeleteConversation(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversations: map[string]*model.Conversation{
// 			"conv": {
// 				Id:      "conv",
// 				Name:    "Conversation1",
// 				IsGroup: false,
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id: "user2",
// 					},
// 				},
// 			},
// 			"notInGroup": {
// 				Participants: []*model.Participant{},
// 			},
// 		},
// 	}
//
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	liveMock := LiveConnectionsMock{}
//
// 	service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 	ctx := contextFor("token1")
//
// 	t.Run("not part of", func(t *testing.T) {
// 		defer func() {
// 			if e := recover(); e == nil {
// 				t.Error("Should panic because I am not participant")
// 			}
// 		}()
// 		err := service.DeleteConversation(ctx, "notInGroup")
// 		if err == nil {
// 			t.Error("I could delete conversation which I am not participant of")
// 		}
// 	})
//
// 	t.Run("I am participant of", func(t *testing.T) {
// 		err := service.DeleteConversation(ctx, "conv")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		part := findParticipant("user1", repoMock.conversations["conv"].Participants)
// 		if part.DeletedTimestamp == 0 {
// 			t.Error("Conversation was not deleted")
// 		}
// 	})
// }
//
// func TestChatService_SendMessage(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		conversations: map[string]*model.Conversation{
// 			"conv": {
// 				Id:      "conv",
// 				Name:    "Conversation1",
// 				IsGroup: false,
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id: "user2",
// 					},
// 				},
// 			},
// 			"notInGroup": {
// 				Participants: []*model.Participant{},
// 			},
// 			"hasLeftGroup": {
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						HasLeft: true,
// 					},
// 				},
// 			},
// 			"groupWithLeftMember": {
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id:      "left",
// 						HasLeft: true,
// 					},
// 				},
// 			},
// 			"groupWithBlockedMember": {
// 				Participants: []*model.Participant{
// 					{
// 						Id:      "user1",
// 						IsAdmin: true,
// 						Unread:  true,
// 					}, {
// 						Id:      "blocked",
// 						Blocked: true,
// 					},
// 				},
// 			},
// 		},
// 	}
//
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	liveMock := LiveConnectionsMock{}
//
// 	service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 	ctx := contextFor("token1")
//
// 	t.Run("to group not part of", func(t *testing.T) {
// 		defer func() {
// 			if e := recover(); e == nil {
// 				t.Error("Should panic because I am not participant")
// 			}
// 		}()
// 		err := service.SendMessage(ctx, &model.Message{
// 			ConversationId: "notInGroup",
// 		})
// 		if err == nil {
// 			t.Error("Able to send message to conversation I am not part of")
// 		}
// 	})
//
// 	t.Run("to group I have already left", func(t *testing.T) {
// 		err := service.SendMessage(ctx, &model.Message{
// 			ConversationId: "hasLeftGroup",
// 		})
// 		if err == nil {
// 			t.Error("Able to send message to conversation I have already left")
// 		}
// 	})
//
// 	t.Run("message is sent to all participants", func(t *testing.T) {
// 		liveMock := LiveConnectionsMock{messages: make(map[string][]interface{})}
// 		service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 		err := service.SendMessage(ctx, &model.Message{
// 			ConversationId: "conv",
// 		})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(liveMock.messages["user1"]) == 0 || len(liveMock.messages["user2"]) == 0 {
// 			t.Error("Message was not sent to all participants")
// 		}
// 		message := liveMock.messages["user2"][0].(*model.Message)
// 		if message.SenderId != "user1" {
// 			t.Error("SenderId was not correctly set to message")
// 		}
// 	})
//
// 	t.Run("to participants who has already left", func(t *testing.T) {
// 		liveMock := LiveConnectionsMock{messages: make(map[string][]interface{})}
// 		service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 		err := service.SendMessage(ctx, &model.Message{
// 			ConversationId: "groupWithLeftMember",
// 		})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(liveMock.messages["left"]) > 0 {
// 			t.Error("Message sent to participant who already left the conversation")
// 		}
// 	})
//
// 	t.Run("to participants who has blocked the conversation", func(t *testing.T) {
// 		liveMock := LiveConnectionsMock{messages: make(map[string][]interface{})}
// 		service := NewChatService(repoMock, nil, authMock, liveMock, AdminMock{"Admin"})
//
// 		err := service.SendMessage(ctx, &model.Message{
// 			ConversationId: "groupWithBlockedMember",
// 		})
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 		if len(liveMock.messages["blocked"]) > 0 {
// 			t.Error("Message sent to participant who has blocked the conversation")
// 		}
// 	})
// }
//
// func TestChatService_CreateReply(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		repliesByUser: make(map[string][]*model.Reply),
// 	}
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, nil, AdminMock{"Admin"})
//
// 	reply, err := service.CreateReply(contextFor("token1"), &model.Reply{
// 		Title: "my reply",
// 		Text:  "my text",
// 	})
// 	if err != nil {
// 		t.Error("Returned error: ", err.Error())
// 	}
// 	if reply.Id.Hex() == "" {
// 		t.Error("Id was not populated")
// 	}
// 	if reply.OwnerId != "user1" {
// 		t.Error("OwnerId was not populated correctly")
// 	}
// }
//
// func TestChatService_GetMyReplies(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		repliesByUser: map[string][]*model.Reply{
// 			"user1": {
// 				{
// 					OwnerId: "user1",
// 				}, {
// 					OwnerId: "user1",
// 				},
// 			},
// 		},
// 	}
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, nil, AdminMock{"Admin"})
//
// 	tester := func(token string, expectedReplies int) func(t *testing.T) {
// 		return func(t *testing.T) {
// 			replies, err := service.GetMyReplies(contextFor(token))
// 			if err != nil {
// 				t.Error("Returned error: ", err.Error())
// 			}
// 			if len(replies) != expectedReplies {
// 				t.Errorf("Expected %d replies, returned %d", expectedReplies, len(replies))
// 			}
// 		}
// 	}
//
// 	t.Run("green path", tester("token1", 2))
//
// 	t.Run("no replies", tester("token2", 0))
// }
//
// func TestChatService_DeleteReply(t *testing.T) {
// 	repoMock := ChatRepoMock{
// 		replies: map[string]*model.Reply{
// 			"reply1": {
// 				Id:      "reply1",
// 				OwnerId: "user1",
// 			},
// 			"reply2": {
// 				Id:      "reply2",
// 				OwnerId: "user1",
// 			},
// 		},
// 	}
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 		"token2": "user2",
// 	}}
//
// 	service := NewChatService(repoMock, nil, authMock, nil, AdminMock{"Admin"})
//
// 	t.Run("green path", func(t *testing.T) {
// 		err := service.DeleteReply(contextFor("token1"), "reply1")
// 		if err != nil {
// 			t.Error("Returned error: ", err.Error())
// 		}
// 	})
//
// 	t.Run("other users reply", func(t *testing.T) {
// 		err := service.DeleteReply(contextFor("token2"), "reply1")
// 		if err == nil {
// 			t.Error("Was able to delete other users reply")
// 		}
// 	})
//
// 	t.Run("non existing reply", func(t *testing.T) {
// 		err := service.DeleteReply(contextFor("token1"), "reply3")
// 		if err == nil {
// 			t.Error("Was able to delete non existing reply")
// 		}
// 	})
// }
