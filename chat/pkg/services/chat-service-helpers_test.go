package services

//
// import (
// 	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
// 	"google.golang.org/grpc/metadata"
// 	"testing"
// )
//
// func TestChatService_authenticateUser(t *testing.T) {
// 	authMock := AuthClientMock{users: map[string]string{
// 		"token1": "user1",
// 	}}
//
// 	service := NewChatService(nil, nil, authMock, nil, AdminMock{"Admin"})
//
// 	t.Run("invalid token", func(t *testing.T) {
// 		defer func() {
// 			e := recover()
// 			if e == nil {
// 				t.Error("Was able to authenticate invalid token")
// 			}
// 		}()
// 		service.AuthenticateUser(contextFor("invalidtoken"))
// 	})
//
// 	t.Run("valid token", func(t *testing.T) {
// 		defer func() {
// 			e := recover()
// 			if e != nil {
// 				t.Error("Panic: ", e)
// 			}
// 		}()
// 		ctx := contextFor("token1")
// 		userId := service.AuthenticateUser(ctx)
// 		if userId != "user1" {
// 			t.Error("Returned wrong userId")
// 		}
// 		if md, ok := metadata.FromIncomingContext(ctx); ok {
// 			if len(md.Get("user-id")) == 0 {
// 				t.Error("user-id was not populated on metadata during authentication")
// 			}
// 			if md.Get("user-id")[0] != "user1" {
// 				t.Error("Wrong user id was set into metadata")
// 			}
// 		} else {
// 			t.Error("metadata was not set to incomming context")
// 		}
// 	})
//
// 	t.Run("with preauthenticated user-id", func(t *testing.T) {
// 		defer func() {
// 			e := recover()
// 			if e != nil {
// 				t.Error("Panic: ", e)
// 			}
// 		}()
//
// 		ctx := contextFor("token1")
// 		md, _ := metadata.FromIncomingContext(ctx)
// 		md.Set("user-id", "userid")
//
// 		userId := service.AuthenticateUser(ctx)
// 		if userId != "userid" {
// 			t.Error("User id was not taken from metadata")
// 		}
// 	})
//
// }
//
// func TestFindParticipant(t *testing.T) {
// 	t.Run("empty participants", func(t *testing.T) {
// 		participant := findParticipant("id", []*model.Participant{})
// 		if participant != nil {
// 			t.Error("Found something in empty list", participant)
// 		}
// 	})
//
// 	t.Run("nil participants", func(t *testing.T) {
// 		participant := findParticipant("id", nil)
// 		if participant != nil {
// 			t.Error("Found something in nil", participant)
// 		}
// 	})
//
// 	t.Run("green path", func(t *testing.T) {
// 		participant := findParticipant("id", []*model.Participant{
// 			{
// 				Id: "other1",
// 			}, {
// 				Id: "other2",
// 			}, {
// 				Id: "id",
// 			}, {
// 				Id: "other3",
// 			},
// 		})
// 		if participant == nil {
// 			t.Error("Participant was not found")
// 		}
// 		if participant.Id != "id" {
// 			t.Error("Returned wrong participant")
// 		}
// 	})
// }
//
// func TestUpdateConversation(t *testing.T) {
// 	t.Run("for not participant", func(t *testing.T) {
// 		conv := updateConversationFor("other", &model.Conversation{
// 			Participants: []*model.Participant{
// 				{
// 					Id:      "user1",
// 					Blocked: true,
// 				}, {
// 					Id:      "user2",
// 					Blocked: true,
// 				},
// 			},
// 		})
// 		if conv == nil {
// 			t.Error("Returned nil")
// 		}
// 		if conv.Blocked {
// 			t.Error("Updated conversation for wrong participant")
// 		}
// 	})
//
// 	t.Run("for all fields", func(t *testing.T) {
// 		conv := updateConversationFor("user1", &model.Conversation{
// 			Participants: []*model.Participant{
// 				{
// 					Id:       "user1",
// 					Blocked:  true,
// 					HasLeft:  true,
// 					Archived: true,
// 					Unread:   true,
// 					Muted:    true,
// 				}, {
// 					Id:      "user2",
// 					Blocked: false,
// 				},
// 			},
// 		})
// 		if conv == nil {
// 			t.Error("Returned nil")
// 		}
// 		if !conv.Blocked || !conv.HasLeft || !conv.Archived || !conv.Unread || !conv.Muted {
// 			t.Error("Updated conversation for wrong participant")
// 		}
// 	})
// }
