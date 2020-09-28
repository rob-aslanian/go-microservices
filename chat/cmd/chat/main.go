package main

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.lan/Rightnao-site/microservices/chat/pkg/tracer"
)

func main() {
	conf := GetConfig()

	// tracer
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("TRACING-ADDR", "192.168.1.13:5775"),
			ServiceName: "Chat",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()

	grpcServer, httpServer, err := createServers(conf, tracer)
	if err != nil {
		log.Fatal(err)
	}
	// grpcServer.AddClosable(closer)

	go func() {
		err = grpcServer.Start(fmt.Sprint(conf.LocalGrpcHost, ":", conf.LocalGrpcPort))
		if err != nil {
			log.Fatal(err)
		}
	}()

	// c := cors.AllowAll() // TODO setup better cors
	server := http.Server{
		Addr:                              fmt.Sprint(conf.LocalWsHost, ":", conf.LocalWsPort),
		Handler:/*c.Handler(*/ httpServer, /*)*/
	}
	log.Fatal(server.ListenAndServe())

	defer grpcServer.Close()
}
