package main

import (
	"context"
	"log"
	"os"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/user/pkg/http-server"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/mq"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/repository/arango"
	cacheRepository "gitlab.lan/Rightnao-site/microservices/user/pkg/repository/cache"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/repository/geoip"
	usersRepository "gitlab.lan/Rightnao-site/microservices/user/pkg/repository/users"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/rpc/client"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/service"
	"gitlab.lan/Rightnao-site/microservices/user/pkg/tracer"
)

func main() {
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Users",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()
	// ------------------------
	authService := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":8803"),
		},
	)
	// ------------------------
	// mailService := clientRPC.NewMailClient(
	// 	clientRPC.Settings{
	// 		Address: getEnv("ADDR_GRPC_MAIL", ":8805"),
	// 	},
	// )
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
	companyService := clientRPC.NewCompanyClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_COMPANY", ":50000"),
		},
	)
	// ------------------------
	stuffService := clientRPC.NewStuffClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_STUFF", ":50000"),
		},
	)
	// ------------------------
	chatService := clientRPC.NewChatClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_CHAT", ":50000"),
		},
	)
	// ------------------------
	geoipRepo, err := geoip.NewRepository("./data/GeoIP2-City.mmdb")
	if err != nil {
		panic(err)
	}
	defer geoipRepo.Close()
	// ------------------------
	usersRepo, err := usersRepository.NewRepository(
		usersRepository.Settings{
			Addresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:      getEnv("USER_MONGO", "developer"),
			Password:  getEnv("PASS_MONGO", "Qwerty123"),
			Database:  "users_db",
		},
	)
	if err != nil {
		panic(err)
	}
	// ------------------------
	cacheRepo, err := cacheRepository.NewRepository(
		cacheRepository.Settings{
			Address:  getEnv("ADDR_REDIS", "192.168.1.13:6379"),
			Database: 0, // os.Getenv("DB_REDIS_DATABASE_CACHE"),
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
			Host:            getEnv("HOST", ""),
			UsersRepository: usersRepo,
			CacheRepository: cacheRepo,
			GeoIPRepository: geoipRepo,
			AuthRPC:         authService,
			ArrangoRepo:     arrRepo,
			// MailRPC:         mailService,
			InfoRPC:    infoService,
			NetworkRPC: networkService,
			CompanyRPC: companyService,
			ChatRPC:    chatService,
			StuffRPC:   stuffService,
			Tracer:     tracer,
			MQ:         mq,
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
