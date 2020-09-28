package mq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jcuga/golongpoll"
	"github.com/streadway/amqp"
)

// RabbitmqConnection ...
type RabbitmqConnection struct {
	conn                    *amqp.Connection
	ch                      *amqp.Channel
	exchangeName            string
	notificationQueue       amqp.Queue // queue for sending notification to client
	notificationRecordQueue amqp.Queue // queue for saving in db
}

// Config ...
type Config struct {
	URL  string
	User string
	Pass string
}

// Repository ...
type Repository interface {
	SaveNotification(ctx context.Context, not map[string]interface{}) error
	GetMapNotificationsSettings(ctx context.Context, userID string) (map[string]*bool, error)
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

	err = mq.createNotificationQueue()
	if err != nil {
		return nil, err
	}

	err = mq.createNotificationRecordQueue()
	if err != nil {
		return nil, err
	}

	return mq, nil
}

// ListenNotifications ...
func ListenNotifications(mq *RabbitmqConnection, lp *golongpoll.LongpollManager, repo Repository) (err error) {
	forever := make(<-chan struct{})

	msgs, qmErr := mq.GetNotifications()
	if qmErr != nil {
		return qmErr
	}

	go func() {
		for m := range msgs {

			var j map[string]interface{}
			var receiverID string

			err := json.Unmarshal(m.Body, &j)
			if err != nil {
				log.Println(err)
				continue
			}

			if s, isExists := j["receiver_id"]; isExists {
				if v, ok := s.(string); ok {
					receiverID = v
				} else {
					log.Println("not string")
					continue
				}
			} else {
				log.Println("recivier not found")
				continue
			}

			// check notification settings
			set, err := repo.GetMapNotificationsSettings(context.TODO(), receiverID)
			if err != nil {
				log.Println("getting settings error:", err)
			} else {
				// log.Println("settings:", set)
				t, isExits := j["type"].(string)
				if !isExits {
					log.Println("error: recieved notification without type field")
				}
				if set[t] != nil && *set[t] == false {
					continue
				}
			}

			lp.Publish(receiverID, j)
		}
	}()

	<-forever
	return err
}

// ListenNotificationsRecord ...
func ListenNotificationsRecord(mq *RabbitmqConnection, repo Repository) (err error) {
	forever := make(<-chan struct{})

	msgs, qmErr := mq.GetNotificationRecords()
	if qmErr != nil {
		return qmErr
	}

	go func() {
		for m := range msgs {

			var j map[string]interface{}

			err := json.Unmarshal(m.Body, &j)
			if err != nil {
				log.Println(err)
				continue
			}
			j["seen"] = false

			err = repo.SaveNotification(context.TODO(), j)
			if err != nil {
				log.Println("Error saving notification in DB:", err)
				continue
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
		"notifications", // name
		"direct",        // type
		true,            // durable
		false,           // auto-delete
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	r.exchangeName = "notifications"
	return nil
}

func (r *RabbitmqConnection) createNotificationQueue() error {
	q, err := r.ch.QueueDeclare(
		"notifications", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		amqp.Table{ // arguments
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"notifications", // queue name
		"notifications", // routing key
		r.exchangeName,  // exhange
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	r.notificationQueue = q

	return nil
}

func (r *RabbitmqConnection) createNotificationRecordQueue() error {
	q, err := r.ch.QueueDeclare(
		"notifications_record",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"notifications_record",
		"notifications_record",
		r.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.notificationRecordQueue = q

	return nil
}

// Close ...
func (r *RabbitmqConnection) Close() {
	r.ch.Close()
	r.conn.Close()
}

/// -------

// GetNotifications ...
func (r *RabbitmqConnection) GetNotifications() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.notificationQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// GetNotificationRecords ...
func (r *RabbitmqConnection) GetNotificationRecords() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.notificationRecordQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
