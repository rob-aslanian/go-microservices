package model

import (
	"time"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
)

type User struct {
	Id        string `json:"_key"`
	Status    string `json:"status"`
	Url       string `json:"url"`
	Avatar    string `json:"avatar"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Gender    string `json:"gender"`

	PrimaryEmail string `json:"primary_email"`
	PrimaryPhone string `json:"primary_phone"`

	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToRPC() *networkRPC.User {
	return &networkRPC.User{
		Id:        u.Id,
		Status:    u.Status,
		Url:       u.Url,
		Avatar:    u.Avatar,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		// Gender:       u.Gender,
		PrimaryEmail: u.PrimaryEmail,
		PrimaryPhone: u.PrimaryPhone,
	}
}

type UserArr []*User

func (users UserArr) ToRPC() *networkRPC.UserArr {
	networkUsers := make([]*networkRPC.User, len(users))
	for i, user := range users {
		networkUsers[i] = user.ToRPC()
	}
	return &networkRPC.UserArr{Users: networkUsers}
}
