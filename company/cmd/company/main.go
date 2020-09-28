package main

import (
	"context"
	"log"
	"os"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/company/pkg/http-server"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/mq"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/repository/arango"
	cacheRepository "gitlab.lan/Rightnao-site/microservices/company/pkg/repository/cache"
	companiesRepository "gitlab.lan/Rightnao-site/microservices/company/pkg/repository/companies"
	reviewsRepository "gitlab.lan/Rightnao-site/microservices/company/pkg/repository/reviews"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/rpc/client"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/service"
	"gitlab.lan/Rightnao-site/microservices/company/pkg/tracer"
)

func main() {
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Company",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()
	// ------------------------
	authService := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":50000"),
		},
	)
	// ------------------------
	mailService := clientRPC.NewMailClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_MAIL", ":8805"),
		},
	)
	// ------------------------
	stuffService := clientRPC.NewStuffClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_STUFF", ":8805"),
		},
	)
	// ------------------------
	infoService := clientRPC.NewInfoClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_INFO", ":8804"),
		},
	)
	// ------------------------
	networkService := clientRPC.NewNetworkClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_NETWORK", ":8806"),
		},
	)
	// ------------------------
	chatService := clientRPC.NewChatClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_CHAT", ":8806"),
		},
	)
	// ------------------------
	userService := clientRPC.NewUserClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_USER", ":50000"),
		},
	)
	// ------------------------
	jobsService := clientRPC.NewJobsClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_JOBS", ":50000"),
		},
	)
	// ------------------------
	compRepo, err := companiesRepository.NewRepository(
		companiesRepository.Settings{
			Addresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:      getEnv("USER_MONGO", "developer"),
			Password:  getEnv("PASS_MONGO", "Qwerty123"),
			Database:  "companies-db",
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	reviewsRepo, err := reviewsRepository.NewRepository(
		reviewsRepository.Settings{
			Addresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:      getEnv("USER_MONGO", "developer"),
			Password:  getEnv("PASS_MONGO", "Qwerty123"),
			Database:  "reviews",
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	cacheRepo, err := cacheRepository.NewRepository(
		cacheRepository.Settings{
			Address:  getEnv("ADDR_REDIS", "192.168.1.13:6379"),
			Database: 1, // os.Getenv("DB_REDIS_DATABASE_CACHE"),
			Password: getEnv("PASS_REDIS", ""),
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	mq, err := mq.NewPublisher(mq.Config{
		URL:  getEnv("ADDR_RABBITMQ", "localhost:5672"),
		User: getEnv("USER_RABBITMQ", ""),
		Pass: getEnv("PASS_RABBITMQ", ""),
	},
	)
	if err != nil {
		panic(err)
	}
	defer mq.Close()
	// ------------------------
	arrRepo, err := arangorepo.NewNetworkRepo(
		getEnv("ARANGO_USER", "root"),
		getEnv("ARANGO_PASS", "pass"),
		strings.Split(getEnv("ARANGO_ADDR", "http://192.168.1.13:8529"), ","),
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	serv, err := service.NewService(
		service.Settings{
			Host:              getEnv("HOST", ""),
			CompanyRepository: compRepo,
			CacheRepository:   cacheRepo,
			ReviewsRepository: reviewsRepo,
			AuthRPC:           authService,
			MailRPC:           mailService,
			InfoRPC:           infoService,
			NetworkRPC:        networkService,
			ChatRPC:           chatService,
			UserRPC:           userService,
			JobsRPC:           jobsService,
			StuffRPC:		   stuffService,
			Tracer:            tracer,
			MQ:                mq,
			ArangoRepo:        arrRepo,
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	httpServer, _ := serverHTTP.NewHTTPServer(context.Background(), serverHTTP.Settings{
		Address: getEnv("HTTP_SERVER_ADDR", ":8123"),
		Service: serv,
		Tracer:  tracer,
	})
	go func() {
		if errHTTP := httpServer.Launch(); errHTTP != nil {
			panic(errHTTP)
		}
	}()
	// ------------------------
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

	if err := rpcServer.Launch(); err != nil {
		log.Println(err)
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
