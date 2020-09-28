package arangorepo

import (
	"context"
	"time"
)

type Gender struct {
	Gender string  `json:"gender"`
	Type   *string `json:"type"`
}

type User struct {
	ID           string    `json:"_key"`
	Status       string    `json:"status"`
	URL          string    `json:"url"`
	Firstname    string    `json:"first_name"`
	Lastname     string    `json:"last_name"`
	Gender       Gender    `json:"gender"`
	PrimaryEmail string    `json:"primary_email"`
	CreatedAt    time.Time `json:"created_at"`
}

func (n NetworkRepo) SaveUser(ctx context.Context, u *User) error {
	_, err := n.users.CreateDocument(ctx, u)
	if err != nil {
		return err
	}

	return nil
}
