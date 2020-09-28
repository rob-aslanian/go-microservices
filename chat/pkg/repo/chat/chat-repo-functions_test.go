package chat

import (
	"fmt"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/model"
	"math"
	"os"
	"testing"
	"time"
)

var repo *ChatRepo
var conv1, conv2 *model.Conversation

func TestMain(m *testing.M) {
	repo = createDB()
	prepareDB(repo)

	retCode := m.Run()

	dropDB(repo)
	os.Exit(retCode)
}
func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return def
}
func createDB() *ChatRepo {
	dbName := fmt.Sprint("chat_db_", time.Now().UnixNano())
	r, _ := NewChatRepo(
		getEnv("MONGO-TEST-USER", "developer"),
		getEnv("MONGO-TEST-PASS", "Qwerty123"),
		dbName,
		[]string{getEnv("MONGO-TEST-ADDR", "192.168.1.13:27017")},
	)
	return r
}

func conversation(name string, users ...string) *model.Conversation {
	c := &model.Conversation{
		Name:         name,
		Participants: make([]*model.Participant, len(users)),
	}
	for i, user := range users {
		c.Participants[i] = &model.Participant{Id: user}
	}
	return c
}

func message(user string, conv *model.Conversation, t int64) *model.Message {
	return &model.Message{
		SenderId:       user,
		ConversationId: conv.Id.Hex(),
		Text:           fmt.Sprintf("%s %s %d", user, conv.Name, t),
		Timestamp:      t,
	}
}

func reply(userId, title, text string) *model.Reply {
	return &model.Reply{
		Title:   title,
		Text:    text,
		OwnerId: userId,
	}
}
func prepareDB(r *ChatRepo) {
	conv1, _ = r.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))
	conv2, _ = r.CreateConversation(nil, conversation("c2", "user1", "user2"))

	r.SendMessage(nil, message("user1", conv1, 1))
	r.SendMessage(nil, message("user1", conv1, 3))
	r.SendMessage(nil, message("user2", conv1, 5))
	r.SendMessage(nil, message("user3", conv1, 7))

	r.CreateReply(nil, reply("user1", "first reply", "this is my reply"))
	r.CreateReply(nil, reply("user1", "second reply", "this is my second reply"))
}

func dropDB(r *ChatRepo) {
	r.db.DropDatabase()
	r.Close()
}

func TestChatRepo_FindConversationsByParticipants(t *testing.T) {

	tester := func(users []*model.Participant, expectedLength int) func(t *testing.T) {
		return func(t *testing.T) {
			conversations, err := repo.FindConversationsByParticipants(nil, users)
			if err != nil {
				t.Error("Returned error: ", err.Error())
			}
			if len(conversations) != expectedLength {
				t.Errorf("Should return %d conversations, returned %d", expectedLength, len(conversations))
			}
		}
	}

	t.Run("subset of participants", tester([]*model.Participant{{Id: "user1"}, {Id: "user2"}}, 2))

	t.Run("all participants", tester([]*model.Participant{{Id: "user1"}, {Id: "user2"}, {Id: "user3"}}, 1))

	t.Run("other participant", tester([]*model.Participant{{Id: "user1"}, {Id: "other"}}, 0))
}

func TestChatRepo_FindConversationsWithOnlyParticipants(t *testing.T) {
	t.Run("two participants", func(t *testing.T) {
		conversations, err := repo.FindConversationsWithOnlyParticipants(nil, []*model.Participant{{Id: "user1"}, {Id: "user2"}})
		if err != nil {
			t.Error("Returned error: ", err.Error())
		}
		if len(conversations) != 1 {
			t.Error("Should return 1 conversations, returned ", len(conversations))
		}
	})
}

func TestChatRepo_GetConversationsOf(t *testing.T) {

	tester := func(userId string, expectedLength int) func(t *testing.T) {
		return func(t *testing.T) {
			conversations, err := repo.GetConversationsOf(nil, userId)
			if err != nil {
				t.Error("Returned error: ", err.Error())
			}
			if len(conversations) != expectedLength {
				t.Errorf("Should return %d conversations, returned %d", expectedLength, len(conversations))
			}
		}
	}

	t.Run("user1", tester("user1", 2))

	t.Run("user3", tester("user3", 1))

	t.Run("other user", tester("other", 0))
}

func TestChatRepo_GetMessages(t *testing.T) {
	tester := func(convId string, from, to int64, expectedLength int) func(t *testing.T) {
		return func(t *testing.T) {
			messages, err := repo.GetMessages(nil, convId, from, to)
			if err != nil {
				t.Error("Returned error: ", err.Error())
			}
			if len(messages) != expectedLength {
				t.Errorf("Should return %d messages, returned %d", expectedLength, len(messages))
			}
		}
	}

	t.Run("for normal user", tester(conv1.Id.Hex(), 0, math.MaxInt64, 4))
	t.Run("for deleted user", tester(conv1.Id.Hex(), 4, math.MaxInt64, 2))
	t.Run("for left user", tester(conv1.Id.Hex(), 0, 6, 3))
	t.Run("for left & deleted user", tester(conv1.Id.Hex(), 2, 4, 1))
}

func TestChatRepo_UpdateParticipant(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))
	t.Run("update existing participant", func(t *testing.T) {
		err := repo.UpdateParticipant(nil, conv1.Id.Hex(), &model.Participant{
			Id:    "user1",
			Muted: true,
		})
		if err != nil {
			t.Error("Returned error: ", err.Error())
		}
		c, _ := repo.GetConversation(nil, conv1.Id.Hex())
		for _, p := range c.Participants {
			if p.Id == "user1" && !p.Muted {
				t.Error("Participant was not updated")
			}
		}
	})

	t.Run("update unexisting participant", func(t *testing.T) {
		err := repo.UpdateParticipant(nil, conv1.Id.Hex(), &model.Participant{
			Id:       "other",
			Archived: true,
		})
		if err == nil {
			t.Error("Updated unexisting participant")
		}
	})
	dropDB(repo)
}

func TestChatRepo_SetUnreadFor(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))
	t.Run("existing participant", func(t *testing.T) {
		err := repo.SetUnreadFor(nil, conv1.Id.Hex(), "user1", true)
		if err != nil {
			t.Error("Returned error: ", err.Error())
		}
		c, _ := repo.GetConversation(nil, conv1.Id.Hex())
		for _, p := range c.Participants {
			if p.Id == "user1" && !p.Unread {
				t.Error("Participant was not updated")
			}
		}
	})

	t.Run("unexisting participant", func(t *testing.T) {
		err := repo.SetUnreadFor(nil, conv1.Id.Hex(), "other", true)
		if err == nil {
			t.Error("Updated unexisting participant")
		}
	})
	dropDB(repo)
}

func TestChatRepo_SetUnreadForAll(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))

	err := repo.SetUnreadForAll(nil, conv1.Id.Hex(), true)
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}
	c, _ := repo.GetConversation(nil, conv1.Id.Hex())
	for _, p := range c.Participants {
		if !p.Unread {
			t.Error("Participant was updated incorrectly")
		}
	}
	dropDB(repo)
}

func TestChatRepo_MuteConversation(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))

	err := repo.MuteConversation(nil, "user1", conv1.Id.Hex(), true)
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}
	c, _ := repo.GetConversation(nil, conv1.Id.Hex())
	for _, p := range c.Participants {
		if p.Id == "user1" && !p.Muted {
			t.Error("Participant was updated incorrectly")
		}
	}
	dropDB(repo)
}

func TestChatRepo_ArchiveConversation(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))

	err := repo.ArchiveConversation(nil, "user1", conv1.Id.Hex(), true)
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}
	c, _ := repo.GetConversation(nil, conv1.Id.Hex())
	for _, p := range c.Participants {
		if p.Id == "user1" && !p.Archived {
			t.Error("Participant was updated incorrectly")
		}
	}
	dropDB(repo)
}

func TestChatRepo_BlockConversation(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))

	err := repo.BlockConversation(nil, "user1", conv1.Id.Hex(), true)
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}
	c, _ := repo.GetConversation(nil, conv1.Id.Hex())
	for _, p := range c.Participants {
		if p.Id == "user1" && !p.Blocked {
			t.Error("Participant was updated incorrectly")
		}
	}
	dropDB(repo)
}

func TestChatRepo_RenameConversation(t *testing.T) {
	repo := createDB()
	conv1, _ := repo.CreateConversation(nil, conversation("c1", "user1", "user2", "user3"))

	err := repo.RenameConversation(nil, conv1.Id.Hex(), "changed name")
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}
	c, _ := repo.GetConversation(nil, conv1.Id.Hex())
	if c.Name != "changed name" {
		t.Error("Name was not changed")
	}
	dropDB(repo)
}

func TestChatRepo_GetUserReplies(t *testing.T) {
	t.Run("green path", func(t *testing.T) {
		replies, err := repo.GetUserReplies(nil, "user1")
		if err != nil {
			t.Error("Returned error: ", err.Error())
		}
		if len(replies) != 2 {
			t.Error("Expected 2 replies, returned ", len(replies))
		}
	})

	t.Run("no replies", func(t *testing.T) {
		replies, err := repo.GetUserReplies(nil, "user2")
		if err != nil {
			t.Error("Returned error: ", err.Error())
		}
		if len(replies) > 0 {
			t.Error("This user should not have replies")
		}
	})
}

func TestChatRepo_DeleteReply(t *testing.T) {
	repo := createDB()
	reply, err := repo.CreateReply(nil, reply("user1", "my reply", "text"))
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}
	if reply.Id.Hex() == "" {
		t.Error("Id was not populated")
	}

	err = repo.DeleteReply(nil, reply.Id.Hex())
	if err != nil {
		t.Error("Returned error: ", err.Error())
	}

	r, _ := repo.GetUserReplies(nil, "user1")
	if len(r) > 0 {
		t.Error("Reply was not deleted from db")
	}

	dropDB(repo)
}
