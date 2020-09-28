package authtoken

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// TokenTimeLimit is a time limit of how long the token should last
const TokenTimeLimit = time.Hour * 24

// SaveToken saves the token in redis through the UserSessionsRepo struct.
func (r UserSessionsRepo) SaveToken(token string, id string) error {
	err := r.Client.Set(token, id, TokenTimeLimit).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserByToken recieves token and gives us back the user ID through the UserSessionsRepo struct.
func (r UserSessionsRepo) GetUserByToken(token string) (string, error) {
	cmd := r.Client.Get(token)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return "", errors.New("Not Authenticated")
		}
		return "", cmd.Err()
	}
	return cmd.Val(), nil
}

// DeleteToken deletes the given token from redis through the UserSessionsRepo struct.
func (r UserSessionsRepo) DeleteToken(token string) error {
	cmd := r.Client.Del(token)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

// DeleteTokens ...
func (r UserSessionsRepo) DeleteTokens(tokens []string) error {
	cmd := r.Client.Del(tokens...)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

// ExpireToken gives the passed token 24 hour time limit, in which case, if the user isn't active it will expire
func (r UserSessionsRepo) ExpireToken(token string) error {
	r.Client.Expire(token, TokenTimeLimit)

	return nil
}
