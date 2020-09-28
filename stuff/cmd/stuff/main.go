package main

import (
	"log"
	"os"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	repository "gitlab.lan/Rightnao-site/microservices/stuff/pkg/repository/mongo"
	clientRPC "gitlab.lan/Rightnao-site/microservices/stuff/pkg/rpc/client"
	serverRPC "gitlab.lan/Rightnao-site/microservices/stuff/pkg/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/stuff/pkg/service"
)

func main() {
	tracer, closer, err := config.Configuration{
		ServiceName: "Stuff",
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  getEnv("ADDR_TRACE_SERVER", "192.168.1.13:5775"),
		},
	}.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		// ...
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)
	// ------------------------
	auth := clientRPC.NewAuthClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_AUTH", ":8803"),
		},
	)
	// ------------------------
	user := clientRPC.NewUserClient(
		clientRPC.Settings{
			Address: getEnv("ADDR_GRPC_USER", ":8803"),
		},
	)
	// ------------------------
	repo, err := repository.NewRepository(
		repository.Settings{
			Addresses:   strings.Split(getEnv("ADDR_MONGO", "192.168.1.13:27017"), ","),
			User:        getEnv("USER_MONGO", "developer"),
			Password:    getEnv("PASS_MONGO", "Qwerty123"),
			Collections: []string{"feedback", "additional_feedback" , "wallet" , "coming_soon"}, //
			Database:    "stuff",                                     // os.Getenv("DB_MONGO_DATABASE_USER"),
		},
	)
	if err != nil {
		// ...
	}

	// ------------------------
	serv, err := service.NewService(
		service.Settings{
			Repository: repo,
			AuthRPC:    auth,
			UserRPC:    user,
		},
	)
	if err != nil {
		// ...
	}

	// ------------------------
	rpcServer, err := serverRPC.NewRPCServer(
		serverRPC.Settings{
			Address: getEnv("ADDR_GRPC_SERVER", ":8822"),
			Service: serv,
		},
	)


	if err != nil {
		// ...
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
