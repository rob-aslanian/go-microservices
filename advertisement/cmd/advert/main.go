package main

import (
	"log"
	"os"
	"strings"

	repository "gitlab.lan/Rightnao-site/microservices/advertisement/pkg/repository/mongo"
	clientRPC "gitlab.lan/Rightnao-site/microservices/advertisement/pkg/rpc/client"
	serverRPC "gitlab.lan/Rightnao-site/microservices/advertisement/pkg/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/advertisement/pkg/service"
	tracing "gitlab.lan/Rightnao-site/microservices/advertisement/pkg/tracer"
)

func main() {
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Advert",
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
	user := clientRPC.NewUserClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_USER", ":8803"),
		},
	)
	// ------------------------
	repo, err := repository.NewRepository(
		repository.Settings{
			Addresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:      getEnv("USER_MONGO", "developer"),
			Password:  getEnv("PASS_MONGO", "Qwerty123"),

			Database:                    "adverts",
			AdvertismentCollection:      "adverts",
			GalleryCollectionCollection: "gallery",
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	serv, err := service.NewService(
		service.Settings{
			Repository: repo,
			AuthRPC:    auth,
			UserRPC:    user,
			Tracer:     tracer,
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	rpcServer, err := serverRPC.NewRPCServer(
		serverRPC.Settings{
			Address: getEnv("ADDR_GRPC_SERVER", ":8822"),
			Service: serv,
		},
	)
	if err != nil {
		panic(err)
	}

	if err := rpcServer.Launch(); err != nil {
		panic(err)
	}
	// ------------------------
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
