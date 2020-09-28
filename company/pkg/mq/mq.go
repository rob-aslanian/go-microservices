package mq

import (
	"encoding/json"
	"log"
	"time"

	notmes "gitlab.lan/Rightnao-site/microservices/company/pkg/notification_messages"

	"github.com/streadway/amqp"
)

// RabbitmqConnection ...
type RabbitmqConnection struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string

	notificationQueue       amqp.Queue // queue for sending notification to client
	notificationRecordQueue amqp.Queue // queue for saving in db

	emailQueue amqp.Queue // queue for sending emails
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

	err = mq.createEmailSendingQueue()
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

func (r *RabbitmqConnection) createEmailSendingQueue() error {
	q, err := r.ch.QueueDeclare(
		"sending_emails",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{ // arguments
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return err
	}

	err = r.ch.QueueBind(
		"sending_emails",
		"sending_emails",
		r.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.emailQueue = q

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

// SendNewCompanyReview ...
func (r *RabbitmqConnection) SendNewCompanyReview(companyID string, message *notmes.NewCompanyReview) error {
	message.Type = notmes.TypeNewReview
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

// SendNewFounderRequest ...
func (r *RabbitmqConnection) SendNewFounderRequest(companyID string, message *notmes.NewFounderRequest) error {
	message.Type = notmes.TypeNewFounderRequest
	message.GenerateID()
	message.ReceiverID = message.Founder
	message.CompanyID = companyID
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

// SendEmail ...
func (r *RabbitmqConnection) SendEmail(emailAddress string, bodyMessage string) error {
	log.Println("sending email")
	message := make(map[string]string, 2)
	message["address"] = emailAddress
	message["message"] = bodyMessage

	j, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.emit(r.emailQueue.Name, j)
	if err != nil {
		return err
	}

	return nil
}
