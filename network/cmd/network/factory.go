package main

import (
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	handlers "gitlab.lan/Rightnao-site/microservices/network/pkg/grpc-handlers"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/auth"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/chat"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/company"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/mq"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/network"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/repo/user"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/services"
	"gitlab.lan/Rightnao-site/microservices/network/pkg/validators"
	"gitlab.lan/Rightnao-site/microservices/shared/grpc/server"
	"google.golang.org/grpc"
)

func createGrpcServer(config *Configuration) (*server.Server, error) {
	h, err := createGrpcHandlers(config)
	if err != nil {
		return nil, err
	}

	grpcServer := server.NewServer(func(s *grpc.Server) {
		networkRPC.RegisterNetworkServiceServer(s, h)
	})

	return grpcServer, nil
}

func createGrpcHandlers(config *Configuration) (*handlers.GrpcHandlers, error) {
	service, err := createNetworkService(config)
	if err != nil {
		return nil, err
	}
	return handlers.NewGrpcHandlers(service), nil
}

func createNetworkService(config *Configuration) (*services.NetworkService, error) {
	repo, err := network.NewNetworkRepo(config.ArangoUser, config.ArangoPass, config.ArangoAddresses)
	if err != nil {
		return nil, err
	}
	validator := validators.NewNetworkValidator()

	auth, err := auth.NewAuthClient(config.AuthServiceAddr)
	if err != nil {
		return nil, err
	}

	user, err := user.NewUserClient(config.UserServiceAddr)
	if err != nil {
		return nil, err
	}

	company, err := company.NewCompanyClient(config.CompanyServiceAddr)
	if err != nil {
		return nil, err
	}

	chat, err := chat.NewChatClient(config.ChatServiceAddr)
	if err != nil {
		return nil, err
	}

	mq, err := mq.NewPublisher(mq.Config{
		URL:  getEnv("ADDR_RABBITMQ", "localhost:5672"),
		User: getEnv("USER_RABBITMQ", ""),
		Pass: getEnv("PASS_RABBITMQ", ""),
	},
	)
	if err != nil {
		panic(err)
	}
	// defer mq.Close()

	return services.NewNetworkService(repo, validator, auth, user, company, chat, mq), nil
}
