package mq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"gitlab.lan/Rightnao-site/microservices/graphql/resolver"
)

// RabbitmqConnection ...
type RabbitmqConnection struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	exchangeNewsfeedName string
	PostNewsfeedQueue    amqp.Queue
	CommentNewsfeedQueue amqp.Queue
	LikeNewsfeedQueue    amqp.Queue
}

// Config ...
type Config struct {
	URL  string
	User string
	Pass string
}

// NewSubscriber ...
func NewSubscriber(conf Config) (*RabbitmqConnection, error) {
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

// ListenNewsPostEvents ...
func ListenNewsPostEvents(mq *RabbitmqConnection, addedPostCh chan *resolver.NewsfeedPostResolverCustom) (err error) {
	forever := make(<-chan struct{})

	msgs, qmErr := mq.GetNewsfeedPostsEvents()
	if qmErr != nil {
		return qmErr
	}

	go func() {
		for m := range msgs {
			var j resolver.NewsfeedPostCustom

			err := json.Unmarshal(m.Body, &j)
			if err != nil {
				log.Println(err)
				continue
			}

			addedPostCh <- &resolver.NewsfeedPostResolverCustom{
				R: &j,
			}
		}
	}()

	<-forever
	return err
}

// ListenNewsCommentsEvents ...
func ListenNewsCommentsEvents(mq *RabbitmqConnection, addedPostCh chan *resolver.NewsfeedPostCommentResolverCustom) (err error) {
	forever := make(<-chan struct{})

	msgs, qmErr := mq.GetNewsfeedCommentsEvents()
	if qmErr != nil {
		return qmErr
	}

	go func() {
		for m := range msgs {
			var j resolver.NewsfeedPostCommentCustom

			err := json.Unmarshal(m.Body, &j)
			if err != nil {
				log.Println(err)
				continue
			}

			addedPostCh <- &resolver.NewsfeedPostCommentResolverCustom{
				R: &j,
			}
		}
	}()

	<-forever
	return err
}

// ListenNewsLikesEvents ...
func ListenNewsLikesEvents(mq *RabbitmqConnection, addedPostCh chan *resolver.LikeResolverCustom) (err error) {
	forever := make(<-chan struct{})

	msgs, qmErr := mq.GetNewsfeedLikeEvents()
	if qmErr != nil {
		return qmErr
	}

	go func() {
		for m := range msgs {
			log.Println(string(m.Body))
			var j resolver.LikeCustom

			err := json.Unmarshal(m.Body, &j)
			if err != nil {
				log.Println(err)
				continue
			}

			addedPostCh <- &resolver.LikeResolverCustom{
				R: &j,
			}
		}
	}()

	<-forever
	return err
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

	r.exchangeNewsfeedName = "newsfeed"
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
		"post",                 // queue name
		"post",                 // routing key
		r.exchangeNewsfeedName, // exhange
		false,                  // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	r.PostNewsfeedQueue = q

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
		"comment",              // queue name
		"comment",              // routing key
		r.exchangeNewsfeedName, // exhange
		false,                  // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	r.CommentNewsfeedQueue = q

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
		"like",                 // queue name
		"like",                 // routing key
		r.exchangeNewsfeedName, // exhange
		false,                  // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	r.LikeNewsfeedQueue = q

	return nil
}

// Close ...
func (r *RabbitmqConnection) Close() {
	r.ch.Close()
	r.conn.Close()
}

/// -------

// GetNewsfeedPostsEvents ...
func (r *RabbitmqConnection) GetNewsfeedPostsEvents() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.PostNewsfeedQueue.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// GetNewsfeedCommentsEvents ...
func (r *RabbitmqConnection) GetNewsfeedCommentsEvents() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.CommentNewsfeedQueue.Name, // queue
		"",                          // consumer
		true,                        // auto-ack
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// GetNewsfeedLikeEvents ...
func (r *RabbitmqConnection) GetNewsfeedLikeEvents() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.LikeNewsfeedQueue.Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
