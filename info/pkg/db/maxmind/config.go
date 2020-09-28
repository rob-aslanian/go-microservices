package maxmind

import (
	maxminddb "github.com/oschwald/maxminddb-golang"
)

type DB struct {
	*maxminddb.Reader
}

func Open(path string) (*DB, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return &DB{}, err
	}

	return &DB{db}, nil
}
