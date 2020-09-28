package main

import (
	"log"
	"os"
	"strings"

	"gitlab.lan/Rightnao-site/microservices/jobs/internal/mq"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/repository/jobs"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/rpc/client"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/service"
	"gitlab.lan/Rightnao-site/microservices/jobs/internal/tracer"
)

func main() {
	// tracer
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
			ServiceName: "Jobs",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	// repository
	repo, err := jobsrepo.NewRepository(
		&jobsrepo.Settings{
			DBAddresses: strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:        getEnv("USER_MONGO", "developer"),
			Password:    getEnv("PASS_MONGO", "Qwerty123"),
			DBName:      "jobs_db",

			JobsCollectionName:             "jobs",
			ProfileCollectionName:          "profiles",
			CompaniesCollectionName:        "companies",
			JobReportsCollectionName:       "job_reports",
			CandidateReportsCollectionName: "candidate_reports",
			JobFiltersCollectionName:       "job_saved_filters",
			CandidateFiltersCollectionName: "candidate_saved_filters",
			PricesCollectionCollectionName: "job_posting_pricings",
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
	},
	)
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

	// service
	service := service.NewService(
		&service.Settings{
			AuthRPC:        authService,
			NetworkRPC:     networkService,
			InfoRPC:        infoService,
			JobsRepository: repo,
			Tracer:         tracer,
			MQ:             mq,
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

	// tasks.ScheduleDailyTask(0, 0, 0, tasks.NewJobExpirationTask(repo))
	// tasks.ScheduleDailyTask(12, 0, 0, tasks.NewJobAlertTask(repo))
	// tasks.ScheduleDailyTask(12, 0, 0, tasks.NewCandidateAlertTask(repo))
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
