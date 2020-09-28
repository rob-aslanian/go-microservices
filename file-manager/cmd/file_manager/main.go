package main

import (
	"context"
	"log"

	"gitlab.lan/Rightnao-site/microservices/file-manager/configs"
	"gitlab.lan/Rightnao-site/microservices/file-manager/pkg/http_server"
	tracer "gitlab.lan/Rightnao-site/microservices/file-manager/pkg/opentracing"
	fileList "gitlab.lan/Rightnao-site/microservices/file-manager/pkg/repository/file_list"
	files "gitlab.lan/Rightnao-site/microservices/file-manager/pkg/repository/file_storage"
	clientRPC "gitlab.lan/Rightnao-site/microservices/file-manager/pkg/rpc/client"
	"gitlab.lan/Rightnao-site/microservices/file-manager/pkg/service"
)

func main() {
	// read configs
	httpServerConfig := configs.NewHttpServerConfig()
	rpcClientConfig := configs.NewRPCClientConfig()
	opentracingConfig := configs.NewOpentracingConfig()
	fileListConfig := configs.NewFileListRepositoryConfig()
	fileStorageConfig := configs.NewFileStorageRepositoryConfig()

	// initiate Opentracing
	_, closer, err := tracer.Create(opentracingConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer closer.Close()

	// initiate repositories
	fileStorageRepo, err := files.NewFileStorageRepository(fileStorageConfig)
	if err != nil {
		log.Fatalln(err)
	}

	fileListRepo, err := fileList.NewFileListRepository(fileListConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// prepare rpc client
	rpcClient := clientRPC.NewRPCClient(rpcClientConfig)
	defer rpcClient.Close()

	// create service
	fileService, err := service.NewFileService(rpcClient, fileListRepo, fileStorageRepo)
	if err != nil {
		log.Fatalln(err)
	}

	// launch HTTP server
	httpServer := http_server.NewHttpServer(context.Background(), httpServerConfig, fileService)
	log.Fatalln(httpServer.Start())
}
