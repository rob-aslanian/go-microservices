package main

import (
	"log"
	"os"
	"strings"

	repository "gitlab.lan/Rightnao-site/microservices/shop/internal/repository/mongo"
	clientRPC "gitlab.lan/Rightnao-site/microservices/shop/internal/rpc/client"
	serverRPC "gitlab.lan/Rightnao-site/microservices/shop/internal/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/shop/internal/service"
	tracing "gitlab.lan/Rightnao-site/microservices/shop/internal/tracer"
)

func main() {
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Shop",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()
	// ------------------------
	auth := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":8803"),
		},
	)
	// ------------------------
	network := clientRPC.NewNetworkClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_NETWORK", ":8803"),
		},
	)
	// ------------------------
	// info RPC
	infoService := clientRPC.NewInfoClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_INFO", ":8804"),
		},
	)
	// ------------------------
	repo, err := repository.NewRepository(
		repository.Settings{
			Addresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:      getEnv("USER_MONGO", "developer"),
			Password:  getEnv("PASS_MONGO", "Qwerty123"),

			Database:           "shop",
			ShopsCollection:    "shops",
			ProductsCollection: "products",
			OrdersCollection:   "orders",
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	// mq, err := mq.NewPublisher(mq.Config{
	// 	URL:  getEnv("ADDR_RABBITMQ", "localhost:5672"),
	// 	User: getEnv("USER_RABBITMQ", ""),
	// 	Pass: getEnv("PASS_RABBITMQ", ""),
	// },
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// defer mq.Close()
	// ------------------------
	serv, err := service.NewService(
		service.Settings{
			Repository: repo,
			AuthRPC:    auth,
			NetworkRPC: network,
			InfoRPC:    infoService,
			Tracer:     tracer,
			// MQ:         mq,
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	servergRPC, err := serverRPC.NewRPCServer(serverRPC.Settings{
		Address: getEnv("ADDR_GRPC_SERVER", ":8822"),
		Service: serv,
		Tracer:  tracer,
	})
	if err != nil {
		panic(err)
	}
	panic(servergRPC.Launch())
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}

	return value
}
