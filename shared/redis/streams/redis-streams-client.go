package streams

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisStreamsClient redis.Client

func NewRedisStreamsClient(options *redis.Options) *RedisStreamsClient {
	options.ReadTimeout = -1
	c := RedisStreamsClient(*redis.NewClient(options))
	return &c
}

func (c *RedisStreamsClient) ListenStream(streamName, group, consumer string, listener func(message RedisStreamsMessage)) {
	go func() {
		for {
			cmd := redis.NewXStreamSliceCmd("XREADGROUP", "GROUP", group, consumer, "BLOCK", 0, "COUNT", 1, "STREAMS", streamName, ">")
			c.Process(cmd)
			stream, err := cmd.Result()
			if err != nil {
				fmt.Println(err)
			} else {
				listener(RedisStreamsMessage{
					Stream:   streamName,
					Group:    group,
					client:   c,
					XMessage: stream[0].Messages[0],
				})
			}
		}
	}()
}

func (c *RedisStreamsClient) Publish(stream string, data map[string]interface{}) (string, error) {
	return c.XAdd(&redis.XAddArgs{
		Stream:       stream,
		Values:       data,
		MaxLenApprox: 1000,
	}).Result()
}
