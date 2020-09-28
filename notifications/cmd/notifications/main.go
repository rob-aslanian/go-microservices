package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/jcuga/golongpoll"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/http-server"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/mq"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/repository/notifications"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/rpc/client"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/service"
	"gitlab.lan/Rightnao-site/microservices/notifications/internal/tracer"
)

func main() {
	// tracer
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Notifications",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	// Message Queue
	mqSub, err := mq.NewSubscriber(mq.Config{
		URL:  getEnv("ADDR_RABBITMQ", "localhost:5672"),
		User: getEnv("USER_RABBITMQ", ""),
		Pass: getEnv("PASS_RABBITMQ", ""),
	},
	)
	if err != nil {
		panic(err)
	}
	defer mqSub.Close()

	// auth gRPC client
	authService := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":8803"),
		},
	)

	// network gRPC client
	networkService := clientRPC.NewNetworkClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_NETWORK", ":8803"),
		},
	)

	// Long Poll manager
	lpManager, err := golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled:                 false,
		MaxLongpollTimeoutSeconds:      120,
		MaxEventBufferSize:             100,
		EventTimeToLiveSeconds:         60 * 10, // 10 minutes
		DeleteEventAfterFirstRetrieval: false,
	})
	if err != nil {
		panic(err)
	}

	// repository
	notificationRepo, err := notificationsRepository.NewRepository(
		notificationsRepository.Settings{
			Addresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:      getEnv("USER_MONGO", "developer"),
			Password:  getEnv("PASS_MONGO", "Qwerty123"),
			Database:  "notifications_db",
		},
	)
	if err != nil {
		panic(err)
	}

	// service
	serv, err := service.NewService(
		service.Settings{
			Notifications: notificationRepo,
			AuthRPC:       authService,
			NetworkRPC:    networkService,
			Tracer:        tracer,
		},
	)
	if err != nil {
		panic(err)
	}

	// http server
	httpServer, _ := serverHTTP.NewHTTPServer(context.Background(), serverHTTP.Settings{
		Address:         getEnv("HTTP_SERVER_ADDR", ":8123"),
		Tracer:          tracer,
		Service:         serv,
		AuthRPC:         authService,
		NetworkRPC:      networkService,
		LongPollManager: lpManager,
	})
	go func() {
		if errHTTP := httpServer.Launch(); errHTTP != nil {
			panic(errHTTP)
		}
	}()

	// gRPC server
	rpcServer, err := serverRPC.NewRPCServer(
		serverRPC.Settings{
			Address: getEnv("ADDR_GRPC_SERVER", ":8822"),
			Service: serv,
			Tracer:  tracer,
		},
	)
	if err != nil {
		panic(err)
	}
	go func() {
		err := rpcServer.Launch()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		// listen queue messages for realtime notifications
		err = mq.ListenNotifications(mqSub, lpManager, notificationRepo)
		if err != nil {
			panic(err)
		}
	}()

	// listen queue messages for saving in db
	err = mq.ListenNotificationsRecord(mqSub, notificationRepo)
	if err != nil {
		panic(err)
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
