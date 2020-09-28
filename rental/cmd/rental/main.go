package main

import (
	"log"
	"os"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/rental/internal/houses"
	clientRPC "gitlab.lan/Rightnao-site/microservices/rental/internal/rpc/client"
	serverRPC "gitlab.lan/Rightnao-site/microservices/rental/internal/rpc/server"

	mq "gitlab.lan/Rightnao-site/microservices/rental/internal/mq"

	housesrepo "gitlab.lan/Rightnao-site/microservices/rental/internal/repository"
	tracing "gitlab.lan/Rightnao-site/microservices/rental/internal/tracer"
)

func main() {
	// tracer
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Rental",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	// ------------------
	repo, err := housesrepo.NewRepository(
		&housesrepo.Settings{
			DBAddresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:        getEnv("USER_MONGO", "developer"),
			Password:    getEnv("PASS_MONGO", "Qwerty123"),
			DBName:      "houses_db",

			HousesCollectionName: "houses",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	// MQ
	mq, err := mq.NewPublisher(mq.Config{
		URL:  getEnv("ADDR_RABBITMQ", "localhost:5672"),
		User: getEnv("USER_RABBITMQ", ""),
		Pass: getEnv("PASS_RABBITMQ", ""),
	})
	if err != nil {
		panic(err)
	}
	defer mq.Close()

	// auth RPC

	authService := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":8803"),
		},
	)

	// network RPC
	networkService := clientRPC.NewNetworkClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_NETWORK", ":8806"),
		},
	)

	// info RPC
	infoService := clientRPC.NewInfoClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_INFO", ":8804"),
		},
	)

	// houses
	houses := houses.NewHouseRental(
		houses.Settings{
			AuthRPC:          authService,
			NetworkRPC:       networkService,
			InfoRPC:          infoService,
			HousesRepository: repo,
			Tracer:           tracer,
			MQ:               mq,
		},
	)

	// gRPC server
	servergRPC, err := serverRPC.NewRPCServer(serverRPC.Settings{
		Address: getEnv("ADDR_GRPC_SERVER", ":8822"),
		House:   houses,
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
