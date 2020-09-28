package main

import (
	"log"
	"os"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/search/internal/repository/filters"
	"gitlab.lan/Rightnao-site/microservices/search/internal/repository/search"
	"gitlab.lan/Rightnao-site/microservices/search/internal/rpc/client"
	"gitlab.lan/Rightnao-site/microservices/search/internal/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/search/internal/service"
	"gitlab.lan/Rightnao-site/microservices/search/internal/tracer"
)

func main() {
	// tracer
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Search",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	// repository
	repo := searchrepo.Connect(&searchrepo.Config{
		Addresses: strings.Split(getEnv("ADDR_ES", "127.0.0.1:9200"), ","),
	})

	// auth RPC
	authService := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":8803"),
		},
	)

	// jobs RPC
	jobsService := clientRPC.NewJobsClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_JOBS", ":8803"),
		},
	)

	// network RPC
	networkService := clientRPC.NewNetworkClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_NETWORK", ":8806"),
		},
	)

	// user RPC
	userService := clientRPC.NewUserClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_USER", ":8806"),
		},
	)

	// company RPC
	companyService := clientRPC.NewCompanyClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_COMPANY", ":8806"),
		},
	)

	// company RPC
	infoService := clientRPC.NewInfoClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_INFO", ":8806"),
		},
	)

	// filter RPC
	filterService, err := filters.NewRepository(
		&filters.Settings{
			DBAddresses:           strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:                  getEnv("USER_MONGO", "developer"),
			Password:              getEnv("PASS_MONGO", "Qwerty123"),
			DBName:                "filters_db",
			FiltersCollectionName: "filters",
		},
	)
	if err != nil {
		log.Println(err)
	}

	// service
	service := service.Create(
		&service.Config{
			AuthRPC:          authService,
			NetworkRPC:       networkService,
			UserRPC:          userService,
			CompanyRPC:       &companyService,
			JobsRPC:          &jobsService,
			InfoRPC:          &infoService,
			FilterRepository: filterService,
			Repository:       repo,
			Tracer:           tracer,
		},
	)

	// gRPC server
	servergRPC, err := serverRPC.NewRPCServer(serverRPC.Settings{
		Address: getEnv("ADDR_GRPC_SERVER", ":8822"),
		Service: service,
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
