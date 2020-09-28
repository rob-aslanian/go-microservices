package main

import (
	"context"
	"log"
	"os"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/statistics/internal/http_server"
	"gitlab.lan/Rightnao-site/microservices/statistics/internal/repository/statistics"
	"gitlab.lan/Rightnao-site/microservices/statistics/internal/service"
	"gitlab.lan/Rightnao-site/microservices/statistics/internal/tracer"
)

func main() {
	// tracer
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Statistics",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	// repository
	repo, err := statisticsrepo.NewRepository(
		&statisticsrepo.Settings{
			DBAddresses:                 strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:                        getEnv("USER_MONGO", "developer"),
			Password:                    getEnv("PASS_MONGO", "Qwerty123"),
			DBName:                      "statistics",
			CompanyStatisticsCollection: "company",
			UsersStatisticsCollection:   "user",
		},
	)
	if err != nil {
		panic(err)
	}

	// service
	serv := service.NewService(repo)

	// http server
	httpServer, err := httpserver.NewHTTPServer(context.Background(), httpserver.Settings{
		Address: getEnv("HTTP_SERVER_ADDR", ":8123"),
		Service: serv,
		Tracer:  tracer,
	})
	if err != nil {
		panic(err)
	}
	func() {
		if errHTTP := httpServer.Launch(); errHTTP != nil {
			panic(errHTTP)
		}
	}()
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
