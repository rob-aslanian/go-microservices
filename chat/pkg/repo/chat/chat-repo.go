package chat

import "github.com/globalsign/mgo"

type ChatRepo struct {
	session       *mgo.Session
	db            *mgo.Database
	conversations *mgo.Collection
	messages      *mgo.Collection
	labels        *mgo.Collection
	replies       *mgo.Collection
	reports       *mgo.Collection

	fs *mgo.GridFS
}

func NewChatRepo(user, pass, dbName string, addresses []string) (*ChatRepo, error) {
	dialInfo := mgo.DialInfo{
		Addrs:    addresses,
		Username: user,
		Password: pass,
	}
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, err
	}

	db := session.DB(dbName)

	conversations := db.C(CONVERSATIONS_COLLECTION_NAME)
	conversations.EnsureIndex(mgo.Index{Name: CONVERSATIONS_INDEX_NAME, Key: []string{CONVERSATIONS_INDEX_KEY}})

	messages := db.C(MESSAGES_COLLECTION_NAME)
	messages.EnsureIndex(mgo.Index{Name: MESSAGES_INDEX_NAME, Key: []string{MESSAGES_INDEX_KEY}})

	labels := db.C(LABELS_COLLECTION_NAME)
	labels.EnsureIndex(mgo.Index{Name: LABELS_INDEX_NAME, Key: []string{LABELS_INDEX_KEY}})

	replies := db.C(REPLIES_COLLECTIONS_NAME)
	replies.EnsureIndex(mgo.Index{Name: REPLIES_INDEX_NAME, Key: []string{REPLIES_INDEX_KEY}})

	reports := db.C(REPORTS_COLLECTIONS_NAME)

	return &ChatRepo{
		session:       session,
		db:            db,
		conversations: conversations,
		messages:      messages,
		replies:       replies,
		labels:        labels,
		reports:       reports,

		fs: db.GridFS(FS_NAME),
	}, nil
}

func (m *ChatRepo) Close() error {
	m.session.Close()
	return nil
}
