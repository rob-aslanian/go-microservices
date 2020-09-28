package mq

import (
	"encoding/json"
	"time"

	notmes "gitlab.lan/Rightnao-site/microservices/jobs/internal/notification_messages"

	"github.com/streadway/amqp"
)

// RabbitmqConnection ...
type RabbitmqConnection struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string

	notificationQueue       amqp.Queue // queue for sending notification to client
	notificationRecordQueue amqp.Queue // queue for saving in db
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
		"notifications",
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

	r.exchangeName = "notifications"
	return nil
}

func (r *RabbitmqConnection) createNotificationQueue() error {
	q, err := r.ch.QueueDeclare(
		"notifications",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"notifications",
		"notifications",
		r.exchangeName,
		false,
		nil,
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

// SendNewInvitation ...
func (r *RabbitmqConnection) SendNewInvitation(userID string, message *notmes.NewInvitation) error {
	message.Type = notmes.TypeNewInvitation
	message.ReceiverID = userID
	message.CreatedAt = time.Now()

	j, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.emit(r.notificationQueue.Name, j)
	if err != nil {
		return err
	}

	err = r.emit(r.notificationRecordQueue.Name, j)
	if err != nil {
		return err
	}

	return nil
}

// SendNewJobApplicant ...
func (r *RabbitmqConnection) SendNewJobApplicant(companyID string, message *notmes.NewJobApplicant) error {
	message.Type = notmes.TypeNewApplicant
	message.GenerateID()
	message.ReceiverID = companyID
	message.CreatedAt = time.Now()

	j, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.emit(r.notificationQueue.Name, j)
	if err != nil {
		return err
	}

	err = r.emit(r.notificationRecordQueue.Name, j)
	if err != nil {
		return err
	}

	return nil
}
