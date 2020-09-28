package cacheRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// CreateTemporaryCodeForEmailActivation generates temporary code and saves data (email and user id) during 24 hours.
func (r Repository) CreateTemporaryCodeForEmailActivation(ctx context.Context, userID string, email string) (string, error) {
	tmpCode := generateString(36)

	data := emailActivation{
		Email:  email,
		UserID: userID,
	}
	js, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	cmd := r.client.Set(tmpCode, js, 24*time.Hour)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return tmpCode, nil
}

func generateString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

// CheckTemporaryCodeForEmailActivation ...
func (r Repository) CheckTemporaryCodeForEmailActivation(ctx context.Context, companyID string, code string) (matched bool, email string, err error) {
	fmt.Println("start checking")

	cmd := r.client.Get(code)
	// if cmd.Err() != nil && cmd.Err() != redis.Nil {
	// 	fmt.Println("error happend", cmd.Err())
	// 	return false, cmd.Err()
	// }

	if cmd.Err() != nil {
		log.Printf("type of err in CheckTemporaryCodeForEmailActivation is %T\n", cmd.Err())
		return false, "", cmd.Err()
	}

	act := emailActivation{}
	err = json.Unmarshal([]byte(cmd.Val()), &act)
	if err != nil {
		return false, "", err
	}

	log.Printf("Inputs:\n\tUserID: %s\n\tcode: %s\nCache:\n\temail: %s\n\tuserID: %s\n", companyID, code, act.Email, act.UserID)

	if companyID != act.UserID {
		fmt.Println("value", cmd.Val())
		return false, "", nil
	}

	fmt.Println("good end")
	return true, act.Email, nil
}

// Remove any cached data by key
func (r Repository) Remove(ctx context.Context, key string) error {
	return r.client.Del(key).Err()
}
