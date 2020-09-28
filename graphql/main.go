package main

import (
	"log"
	"os"

	"gitlab.lan/Rightnao-site/microservices/graphql/mq"
	tracer "gitlab.lan/Rightnao-site/microservices/graphql/opentracing"
	"gitlab.lan/Rightnao-site/microservices/graphql/resolver"
	"gitlab.lan/Rightnao-site/microservices/graphql/server"
)

func main() {
	closer, err := tracer.Create()
	if err != nil {
		// TODO: handle error
		log.Println(err)
	}
	defer closer.Close()

	// ----------
	mqSub, err := mq.NewSubscriber(mq.Config{
		URL:  getEnv("ADDR_RABBITMQ", "localhost:5672"),
		User: getEnv("USER_RABBITMQ", ""),
		Pass: getEnv("PASS_RABBITMQ", ""),
	},
	)
	if err != nil {
		panic(err)
	}

	defer mqSub.Close()
	// ----------

	res := resolver.Resolver{}
	res.Init()

	// optionsCors := cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowCredentials: true,
	// 	// AllowedMethods:   []string{"*"},
	// 	// AllowedHeaders:   []string{"*"},
	// 	Debug: true,
	// }

	// -------
	go mq.ListenNewsPostEvents(mqSub, res.AddedPostEvents)
	go mq.ListenNewsCommentsEvents(mqSub, res.AddedCommentPostEvents)
	go mq.ListenNewsLikesEvents(mqSub, res.AddedLikePostEvents)
	// -------

	s := server.NewGqlServer(&res, getEnv("HTTP_SERVER_ADDR", ":8000"), nil /*&optionsCors*/)

	log.Fatal(s.Serve()) // TODO: run as https
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
