package offlinerepo

import (
	"log"

	"github.com/go-redis/redis"
)

// SetOffline set user online status permanently offline
func (r Repository) SetOffline(id string, isOffline bool) error {
	err := r.client.Set(id+"_online", isOffline, 0).Err()
	if err != nil {
		log.Println("saving status error:", err)
		return err
	}

	return nil
}

// IsOffline ...
func (r Repository) IsOffline(id string) (bool, error) {
	cmd := r.client.Get(id + "_online")
	if cmd.Err() != nil {
		// if not found
		if cmd.Err() == redis.Nil {
			return false, nil
		}
		return false, cmd.Err()
	}

	b := false

	if cmd.Val() == "1" {
		b = true
	}

	return b, nil
}

func (r Repository) DeleteOffline(id string) error {
	cmd := r.client.Del(id + "_online")
	if cmd.Err() != nil {
		// if not found
		if cmd.Err() == redis.Nil {
			return nil
		}
		return cmd.Err()
	}

	return nil
}
