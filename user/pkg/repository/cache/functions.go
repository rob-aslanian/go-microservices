package cacheRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// TestFunctionInCache ...
// func (r Repository) TestFunctionInCache() {
//
// sub := r.client.Subscribe("channel_test")
// s, err := sub.Receive()
// if err != nil {
// 	fmt.Println(err)
// } else {
// 	cmd := r.client.Publish("channel_test", "TEST")
// 	if cmd.Err() != redis.Nil {
// 		fmt.Println(cmd.Err())
// 	}
//
// 	fmt.Printf("Recivied: %#v\n", s)
//
// }
// }

// CreateTemporaryCodeForEmailActivation generates temporary code and saves data (email and user id) during 24 hours.
func (r Repository) CreateTemporaryCodeForEmailActivation(ctx context.Context, userID string, email string) (string, error) {
	tmpCode := generateString(6)

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

// CreateTemporaryCodeForNotActivatedUser generates temporary token and saves it during 2 hours.
func (r Repository) CreateTemporaryCodeForNotActivatedUser(id string) (string, error) {
	tmpToken := generateString(56)
	cmd := r.client.Set(id, tmpToken, 24*time.Hour)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return tmpToken, nil
}

// CreateTemporaryCodeForRecoveryByEmail generates temporary code and saves it during 24 hours.
func (r Repository) CreateTemporaryCodeForRecoveryByEmail(id string) (string, error) {
	tmpCode := generateString(36)
	cmd := r.client.Set(id, tmpCode, 24*time.Hour)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return tmpCode, nil
}

// CheckTemporaryCodeForNotActivatedUser returns true if code matched
func (r Repository) CheckTemporaryCodeForNotActivatedUser(userID string, code string) (bool, error) {
	cmd := r.client.Get(userID)
	// if cmd.Err() != nil && cmd.Err() != redis.Nil {
	// 	fmt.Println("error happend", cmd.Err())
	// 	return false, cmd.Err()
	// }

	if code != cmd.Val() {
		return false, nil
	}

	return true, nil
}

// CheckTemporaryCodeForEmailActivation returns true if code matched
func (r Repository) CheckTemporaryCodeForEmailActivation(ctx context.Context, userID string, code string) (bool, string, error) {
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
	err := json.Unmarshal([]byte(cmd.Val()), &act)
	if err != nil {
		return false, "", err
	}

	log.Printf("Inputs:\n\tUserID: %s\n\tcode: %s\nCache:\n\temail: %s\n\tuserID: %s\n", userID, code, act.Email, act.UserID)

	if userID != act.UserID {
		fmt.Println("value", cmd.Val())
		return false, "", nil
	}

	fmt.Println("good end")
	return true, act.Email, nil
}

// CheckTemporaryCodeForRecoveryByEmail ...
func (r Repository) CheckTemporaryCodeForRecoveryByEmail(userID string, code string) (bool, error) {
	fmt.Println("start checking")
	cmd := r.client.Get(userID)
	// if cmd.Err() != nil && cmd.Err() != redis.Nil {
	// 	fmt.Println("error happend", cmd.Err())
	// 	return false, cmd.Err()
	// }

	if code != cmd.Val() {
		fmt.Println("value", cmd.Val())
		return false, nil
	}

	fmt.Println("good end")
	return true, nil
}

// Remove any cached data by key
func (r Repository) Remove(key string) error {
	return r.client.Del(key).Err()
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
