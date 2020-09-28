package main

import (
	"os"
	"strings"
	"sync"
)

type Configuration struct {
	LocalGrpcHost string
	LocalGrpcPort string

	LocalWsHost string
	LocalWsPort string

	MongoUser     string
	MongoPassword string
	MongoDb       string
	MongoAddr     []string

	AuthServiceAddr    string
	NetworkServiceAddr string

	OpentracingTitle   string
	OpentracingAddress string
}

var config *Configuration
var once sync.Once

func GetConfig() *Configuration {
	once.Do(func() {
		config = loadConfig()
	})
	return config
}

func loadConfig() *Configuration {
	return &Configuration{
		LocalGrpcHost: getEnv("GRPC-HOST", ""),
		LocalGrpcPort: getEnv("GRPC-PORT", "8808"),

		LocalWsHost: getEnv("WS-HOST", ""),
		LocalWsPort: getEnv("WS-PORT", "8809"),

		MongoAddr:     strings.Split(getEnv("MONGO-ADDR", "192.168.1.13:27017"), ","),
		MongoUser:     getEnv("MONGO-USER", "developer"),
		MongoPassword: getEnv("MONGO-PASS", "Qwerty123"),
		MongoDb:       getEnv("MONGO-DB", "chat_db"),

		AuthServiceAddr:    getEnv("AUTH-ADDR", "127.0.0.1:8803"),
		NetworkServiceAddr: getEnv("NET-ADDR", "127.0.0.1:8806"),

		OpentracingTitle:   getEnv("TRACING-TITLE", "Chat"),
		OpentracingAddress: getEnv("TRACING-ADDR", "192.168.1.13:5775"),
	}
}

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return def
}
