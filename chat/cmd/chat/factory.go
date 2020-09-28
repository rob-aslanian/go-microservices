package main

import (
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/tracer"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/server"
	"google.golang.org/grpc"

	handlers "gitlab.lan/Rightnao-site/microservices/chat/pkg/grpc-handlers"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/http-server"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/live-connections"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/repo/auth"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/repo/chat"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/repo/network"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/repo/offline-users"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/services"
	"gitlab.lan/Rightnao-site/microservices/chat/pkg/validation"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
)

func createServers(config *Configuration, tracer *tracing.Tracer) (*server.Server, *http_server.HttpServer, error) {
	repo, err := chat.NewChatRepo(config.MongoUser, config.MongoPassword, config.MongoDb, config.MongoAddr)
	if err != nil {
		return nil, nil, err
	}
	validator := validation.NewChatValidator()

	auth, err := auth.NewAuthClient(config.AuthServiceAddr)
	if err != nil {
		return nil, nil, err
	}

	offRepo, err := offlinerepo.NewRepository(
		offlinerepo.Settings{
			Address: getEnv("ADDR_REDIS", "localhost:6379"),
			// Database: "",
			Password: getEnv("PASS_REDIS", ""),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	liveConnections := live_connections.NewLiveConnections(tracer, offRepo)

	adminProvider, err := network.NewAdminLevelProvider(config.NetworkServiceAddr)
	if err != nil {
		return nil, nil, err
	}

	service := services.NewChatService(repo, validator, auth, liveConnections, adminProvider, tracer)

	h := handlers.NewGrpcHandlers(service)

	liveConnections.SetChatService(service)

	httpServer := http_server.NewHttpServer(auth, service, liveConnections, tracer)

	grpcServer := server.NewServer(func(s *grpc.Server) {
		chatRPC.RegisterChatServiceServer(s, h)
	})

	return grpcServer, httpServer, nil
}
