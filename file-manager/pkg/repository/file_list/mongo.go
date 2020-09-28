package fileList

import (
	"github.com/globalsign/mgo"
)

type Mongo struct {
	db         *mgo.Database
	session    *mgo.Session
	collection map[string]*mgo.Collection
}

type Index mgo.Index

type Collection struct {
	*mgo.Collection
}

func Connect(conf Configuration) (*Mongo, error) {
	session, err := mgo.DialWithInfo(getDialInfo(conf))
	if err != nil {
		return nil, err
	}

	db := session.DB(conf.GetDatabase())

	return &Mongo{
		session: session,
		db:      db,
	}, nil
}

func (m *Mongo) Close() error {
	m.session.Close()
	return nil
}

func (m *Mongo) SetDatabase(dbName string) {
	m.db = m.session.DB(dbName)
}

func (m *Mongo) AddCollection(collection string) Collection {
	return Collection{m.db.C(collection)}
}

func (c *Collection) EnsureIndexes(index Index) {
	c.EnsureIndex(mgo.Index(index))
}

func getDialInfo(conf Configuration) *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    conf.GetAddress(),
		Username: conf.GetUser(),
		Password: conf.GetPassword(),
	}
}
