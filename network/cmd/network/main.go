package main

import (
	"fmt"
	"log"

	"gitlab.lan/Rightnao-site/microservices/shared/opentracing"
)

func main() {

	conf := GetConfig()

	closer, err := opentracing.Create(conf.OpentracingTitle, conf.OpentracingAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer, err := createGrpcServer(conf)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.AddClosable(closer)
	err = grpcServer.Start(fmt.Sprint(conf.LocalGrpcHost /*, ":", conf.LocalGrpcPort*/))
	if err != nil {
		log.Fatal(err)
	}
	defer grpcServer.Close()
}
