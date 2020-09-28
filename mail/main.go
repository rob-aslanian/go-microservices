package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/smtp"
	"os"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/mailRPC"
	"gitlab.lan/Rightnao-site/microservices/mail/opentracing"
	"google.golang.org/grpc"
)

var (
	grcpPort string

	smtpMailServerAddress string
	loginMailServer       string
	passwordMailServer    string
)

func parseArgs() {
	// smtpAddressPtr := flag.String("smtpAddr", "192.168.1.66:25", "smtp address")
	// mailUserPtr := flag.String("mailUser", "vova@test.com", "mail user")
	// mailPassPtr := flag.String("mailpass", "12345678", "mail password")
	//
	// serverAddrPtr := flag.String("rpcAddr", ":8805", "local http address")
	//
	// flag.Parse()
	//
	// grcpPort = *serverAddrPtr
	// smtpMailServerAddress = *smtpAddressPtr
	// loginMailServer = *mailUserPtr
	// passwordMailServer = *mailPassPtr

	grcpPort = getEnv("ADDR_GRPC_SERVER", ":8805")
	smtpMailServerAddress = getEnv("ADDR_SMTP_SERVER", "192.168.1.66:25")
	loginMailServer = getEnv("USER_SMTP_SERVER", "vova@test.com")
	passwordMailServer = getEnv("PASS_SMTP_SERVER", "12345678")
	// hostTarget = getEnv("HOST", "localhost:8123")
}

type EmailServer struct{}

func sendEmail(sender, receiver, data string) error {
	host, _, _ := net.SplitHostPort(smtpMailServerAddress)

	// m := email.NewMessage("title", data)
	// m.From = mail.Address{
	// 	Name:    "From",
	// 	Address: sender,
	// }
	//
	// m.To = []string{receiver}
	//
	// auth := smtp.PlainAuth("", loginMailServer, passwordMailServer, host)
	// if err := email.Send(smtpMailServerAddress, auth, m); err != nil {
	// 	return err
	// }

	// // c, err := smtp.Dial(smtpMailServerAddress)
	// conn, err := tls.Dial(
	// 	"tcp",
	// 	smtpMailServerAddress,
	// 	&tls.Config{
	// 		InsecureSkipVerify: true,
	// 		ServerName:         host,
	// 	},
	// )
	// if err != nil {
	// 	return err
	// }

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	c, err := smtp.Dial(smtpMailServerAddress)
	if err != nil {
		log.Println("error: dial:", err)
		return err
	}

	// conn, err := tls.Dial("tcp", smtpMailServerAddress, tlsConfig)
	// if err != nil {
	// 	log.Println("error: dial:", err)
	// 	return err
	// }

	err = c.StartTLS(tlsConfig)
	if err != nil {
		log.Println("error: start tls:", err)
		return err
	}

	// c, err = smtp.NewClient(conn, host)
	// if err != nil {
	// 	log.Println("error: new connection", err)
	// }

	if err = c.Auth(
		smtp.PlainAuth("", loginMailServer, passwordMailServer, host),
	); err != nil {
		log.Println("error: auth", err)
		return err
	}

	// err = c.Hello(host)
	// if err != nil {
	// 	log.Println("error: hello", err)
	// 	return err
	// }

	err = c.Mail(sender)
	if err != nil {
		log.Println("error: mail", err)
		return err
	}

	err = c.Rcpt(receiver)
	if err != nil {
		log.Println("error: rcpt", err)
		return err
	}

	wc, err := c.Data()
	if err != nil {
		log.Println("error: data", err)
		return err
	}
	defer wc.Close()

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := bytes.NewBufferString(mime + data)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Println("error: writing", err)
		return err
	}
	return nil
}

func (es *EmailServer) SendMail(ctx context.Context, data *mailRPC.SendMailRequest) (*mailRPC.Empty, error) {
	err := sendEmail(loginMailServer, data.Receiver, data.Data)
	if err != nil {
		return nil, err
	}
	return &mailRPC.Empty{}, nil
}

func main() {
	parseArgs()

	closer, err := tracer.Create()
	if err != nil {
		// TODO: handle error
		log.Println(err)
	}
	defer closer.Close()

	ES := EmailServer{}

	lis, err := net.Listen("tcp", grcpPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
	)

	mailRPC.RegisterMailServiceServer(grpcServer, &ES)

	// rabitMQ
	ch, err := connectMQ(
		getEnv("USER_RABBITMQ", ""),
		getEnv("PASS_RABBITMQ", ""),
		getEnv("ADDR_RABBITMQ", "localhost:5672"),
	)
	if err != nil {
		panic(err)
	}
	log.Println("connected to RabbitMQ")

	go func() {
		forever := make(<-chan struct{})

		msgs, err := getMessages(ch)
		if err != nil {
			log.Println("error: get message:", err)
		}

		for m := range msgs {
			log.Println("got message")

			j := make(map[string]string, 2)

			err := json.Unmarshal(m.Body, &j)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("email sent:", j["address"], j["message"])

			if j["address"] == "" {
				log.Println("email address is empty")
				continue
			}

			err = sendEmail(loginMailServer, j["address"], j["message"])
			if err != nil {
				log.Println("error: send email:", err)
			}
		}
		log.Println("email has been send")

		<-forever
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}

func connectMQ(user string, pass string, url string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(`amqp://` + user + `:` + pass + `@` + url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	exchangeName := "notifications"

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	emailQueue := "sending_emails"

	_, err = ch.QueueDeclare(
		emailQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		amqp.Table{ // arguments
			"x-message-ttl": int32(300000), // 5 minutes?
		},
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		emailQueue,   // queue name
		emailQueue,   // routing key
		exchangeName, // exhange
		false,        // no wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func getMessages(ch *amqp.Channel) (<-chan amqp.Delivery, error) {
	emailQueue := "sending_emails"

	msgs, err := ch.Consume(
		emailQueue,
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
