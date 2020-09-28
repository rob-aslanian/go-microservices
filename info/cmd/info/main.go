package main

import (
	"log"
	"os"

	"gitlab.lan/Rightnao-site/microservices/info/pkg/db/maxmind"
	"gitlab.lan/Rightnao-site/microservices/info/pkg/db/postgres"
	"gitlab.lan/Rightnao-site/microservices/info/pkg/rpc/server"
	"gitlab.lan/Rightnao-site/microservices/info/pkg/tracing"
)

func main() {
	tracer, closer, err := tracing.NewTracer(
		tracing.Settings{
			Address:     getEnv("ADDR_TRACE_SERVER", "127.0.0.1:5775"),
			ServiceName: "Information",
		},
	)
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()
	// ------------------------
	db, err := postgres.Connect(&postgres.ConnectionInfo{
		User:     getEnv("USER_POSTGRE", "developer"),
		Password: getEnv("PASS_POSTGRE", "Qwerty123"),
		Name:     getEnv("DB_POSTGRE", "rightnao"),
		Host:     getEnv("ADDR_POSTGRE", "localhost"),
		Port:     getEnv("PORT_POSTGRE", "5432"),
		SSLMode:  false,
		Tracer:   tracer,
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// test
	// cityByID, err := db.GetCityInfoByIDNew(context.Background(), "ru", "611717")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("cityByID(ru, 611717):", cityByID)
	//
	// countryCodeByID, err := db.GetCountryCodeByIDNew(context.Background(), 88)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("countryCodeByID(88):", countryCodeByID)
	//
	// countryCodes, err := db.GetListOfCountryCodesNew(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("countryCodes:", countryCodes)
	//
	// cityByLetters, err := db.GetListOfAllCitiesNew(context.Background(), "ru", "kir", 5, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("cityByLetters(ru, kir, 5, 0):", cityByLetters)
	//
	// cityByLettersInCountry, err := db.GetListOfCitiesNew(context.Background(), "ru", "RU", "kir", 5, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("cityByLettersInCountry(ru, RU, kir, 5, 0):", cityByLettersInCountry)
	// test

	// ------------------------
	mDB, err := maxmind.Open("./data/GeoIP2-City.mmdb")
	if err != nil {
		panic(err)
	}
	defer mDB.Close()
	// ------------------------
	rpcServer, err := server.NewRPCServer(
		server.Settings{
			Address: getEnv("ADDR_GRPC_SERVER", ":8804"),
			Tracer:  tracer,
			DB:      db,
			Mdb:     mDB,
		},
	)
	if err != nil {
		panic(err)
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
