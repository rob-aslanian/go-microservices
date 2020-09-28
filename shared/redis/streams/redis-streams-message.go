package streams

import "github.com/go-redis/redis"

type RedisStreamsMessage struct {
	redis.XMessage
	client *RedisStreamsClient

	Group, Stream string
}

func (this *RedisStreamsMessage) Ack() error {
	return this.client.XAck(this.Stream, this.Group, this.ID).Err()
}
