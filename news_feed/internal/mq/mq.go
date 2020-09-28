package mq

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"gitlab.lan/Rightnao-site/microservices/news_feed/internal/post"
)

// RabbitmqConnection ...
type RabbitmqConnection struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string

	notificationQueue    amqp.Queue // queue for sending notification to client
	newsfeedPostQueue    amqp.Queue
	newsfeedCommentQueue amqp.Queue
	newsfeedLikeQueue    amqp.Queue
}

// Config ...
type Config struct {
	URL  string
	User string
	Pass string
}

// NewPublisher ...
func NewPublisher(conf Config) (*RabbitmqConnection, error) {
	conn, ch, err := connect(conf)
	if err != nil {
		return nil, err
	}
	mq := &RabbitmqConnection{
		conn: conn,
		ch:   ch,
	}

	err = mq.createExchange()
	if err != nil {
		return nil, err
	}

	err = mq.createPostQueue()
	if err != nil {
		return nil, err
	}

	err = mq.createCommentQueue()
	if err != nil {
		return nil, err
	}

	err = mq.createLikeQueue()
	if err != nil {
		return nil, err
	}

	return mq, nil
}

func connect(conf Config) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(`amqp://` + conf.User + `:` + conf.Pass + `@` + conf.URL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func (r *RabbitmqConnection) createExchange() error {
	err := r.ch.ExchangeDeclare(
		"newsfeed",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.exchangeName = "newsfeed"
	return nil
}

func (r *RabbitmqConnection) createPostQueue() error {
	q, err := r.ch.QueueDeclare(
		"post", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		amqp.Table{
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"post",         // queue name
		"post",         // routing key
		r.exchangeName, // exhange
		false,          // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	r.newsfeedPostQueue = q

	return nil
}

func (r *RabbitmqConnection) createCommentQueue() error {
	q, err := r.ch.QueueDeclare(
		"comment", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"comment",      // queue name
		"comment",      // routing key
		r.exchangeName, // exhange
		false,          // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	r.newsfeedCommentQueue = q

	return nil
}

func (r *RabbitmqConnection) createLikeQueue() error {
	q, err := r.ch.QueueDeclare(
		"like", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		amqp.Table{
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"like",         // queue name
		"like",         // routing key
		r.exchangeName, // exhange
		false,          // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	r.newsfeedLikeQueue = q

	return nil
}

// Close ...
func (r *RabbitmqConnection) Close() {
	r.ch.Close()
	r.conn.Close()
}

func (r *RabbitmqConnection) emit(queueName string, message []byte) error {
	err := r.ch.Publish(
		r.exchangeName,
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// -------

// SendNewPostEvent ...
func (r *RabbitmqConnection) SendNewPostEvent(p *post.Post) error {
	j, err := json.Marshal(p)
	if err != nil {
		return err
	}

	err = r.emit(r.newsfeedPostQueue.Name, j)
	if err != nil {
		return err
	}

	return nil
}

// SendNewCommentEvent ...
func (r *RabbitmqConnection) SendNewCommentEvent(c *post.Comment) error {
	j, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = r.emit(r.newsfeedCommentQueue.Name, j)
	if err != nil {
		return err
	}

	return nil
}

// SendNewLikeEvent ...
func (r *RabbitmqConnection) SendNewLikeEvent(c *post.Like) error {
	j, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = r.emit(r.newsfeedLikeQueue.Name, j)
	if err != nil {
		return err
	}

	return nil
}
