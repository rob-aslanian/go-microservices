package main

import (
	"log"
	"os"
	"strings"

	repo "gitlab.lan/Rightnao-site/microservices/authenticator/internal/repository/mmdb"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/repository/sessions"
	authtoken "gitlab.lan/Rightnao-site/microservices/authenticator/internal/repository/tokens"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/service"
	"gitlab.lan/Rightnao-site/microservices/authenticator/internal/tracer"
)

func main() {
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Authenticator",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	// session repository
	sessRepo, err := sessions.NewRepository(
		&sessions.Settings{
			DBAddresses:            strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:                   getEnv("USER_MONGO", "developer"),
			Password:               getEnv("PASS_MONGO", ""),
			DBName:                 "users_db",
			SessionsCollectionName: "sessions",
		},
	)
	if err != nil {
		panic(err)
	}

	// ips repository
	ipsRepo, err := repo.NewIPsRepo(`./data/GeoIP2-City.mmdb`)
	if err != nil {
		panic(err)
	}

	// ips repository
	tokensRepo, err := authtoken.NewRepo(
		authtoken.Settings{
			Address:  getEnv("ADDR_REDIS", "192.168.1.13:6379"),
			Password: getEnv("PASS_REDIS", ""),
			DB:       0,
		},
	)
	if err != nil {
		panic(err)
	}

	// -----------------
	//service
	service := service.NewService(
		&service.Settings{
			Tracer: tracer,
			IPDB:   ipsRepo,
			Auth:   sessRepo,
			Tokens: tokensRepo,
		},
	)

	//gRPC server
	servergRPC, err := serverRPC.NewRPCServer(serverRPC.Settings{
		Address: getEnv("ADDR_GRPC_SERVER", "0.0.0.0:8803"),
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
