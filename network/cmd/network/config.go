package main

import (
	"os"
	"strings"
	"sync"
)

type Configuration struct {
	ArangoUser      string
	ArangoPass      string
	ArangoAddresses []string

	AuthServiceAddr    string
	UserServiceAddr    string
	CompanyServiceAddr string
	ChatServiceAddr    string

	LocalGrpcHost string
	LocalGrpcPort string

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
		ArangoAddresses: strings.Split(getEnv("ARANGO_ADDR", "http://192.168.1.13:8529"), ","),
		ArangoUser:      getEnv("ARANGO_USER", "root"),
		ArangoPass:      getEnv("ARANGO_PASS", "pass"),

		AuthServiceAddr:    getEnv("AUTH_ADDR", "127.0.0.1:8803"),
		UserServiceAddr:    getEnv("USER_ADDR", "127.0.0.1:8803"),
		CompanyServiceAddr: getEnv("COMPANY_ADDR", "127.0.0.1:8803"),
		ChatServiceAddr:    getEnv("CHAT_ADDR", "127.0.0.1:8803"),

		LocalGrpcHost: getEnv("ADDR_GRPC_SERVER", ":8806"),
		// LocalGrpcHost: getEnv("GRPC-HOST", ""),
		// LocalGrpcPort: getEnv("GRPC-PORT", "8806"),

		// OpentracingTitle:   getEnv("TRACING-TITLE", "Network"),
		OpentracingTitle:   "Network",
		OpentracingAddress: getEnv("TRACING_ADDR", "192.168.1.13:5775"),
	}
}

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return def
}
