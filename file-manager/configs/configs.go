package configs

import (
	"log"
	"os"
	"strings"
)

// HTTP server

type HttpServerConfig struct {
	address string
}

func NewHttpServerConfig() *HttpServerConfig {
	return &HttpServerConfig{
		address: getEnv("HTTP_SERVER_ADDRESS", ":2090"), //os.Getenv("HTTP_SERVER_ADDRESS"),
	}
}

func (c *HttpServerConfig) GetAddress() string {
	return c.address
}

// Opentracing

type OpentracingConfig struct {
	address     string
	serviceName string
}

func NewOpentracingConfig() *OpentracingConfig {
	return &OpentracingConfig{
		address:     getEnv("TRACING_SERVER_ADDRESS", "192.168.1.13:5775"), //os.Getenv("TRACING_SERVER_ADDRESS"),
		serviceName: "File_manager",
	}
}

func (c *OpentracingConfig) GetAddress() string {
	return c.address
}

func (c *OpentracingConfig) GetServiceName() string {
	return c.serviceName
}

// RPC client

type RPCClientConfig struct {
	authAddress     string
	userAddress     string
	companyAddress  string
	jobsAddress     string
	advertAddress   string
	servicesAddress string
	newsfeedAddress string
	stuffAddress    string
}

func NewRPCClientConfig() *RPCClientConfig {
	return &RPCClientConfig{
		authAddress:     getEnv("RPC_CLIENT_AUTH_ADDRESS", "127.0.0.1:8803"), //os.Getenv("RPC_CLIENT_AUTH_ADDRESS"),
		userAddress:     getEnv("RPC_CLIENT_USER_ADDRESS", "127.0.0.1:8824"), //os.Getenv("RPC_CLIENT_USER_ADDRESS"),
		companyAddress:  getEnv("RPC_CLIENT_COMPANY_ADDRESS", "127.0.0.1:8824"),
		jobsAddress:     getEnv("RPC_CLIENT_JOBS_ADDRESS", "127.0.0.1:8824"),
		advertAddress:   getEnv("RPC_CLIENT_ADVERT_ADDRESS", "127.0.0.1:8824"),
		servicesAddress: getEnv("RPC_CLIENT_SERVICES_ADDRESS", "127.0.0.1:8824"),
		newsfeedAddress: getEnv("RPC_CLIENT_NEWSFEED_ADDRESS", "127.0.0.1:8824"),
		stuffAddress:    getEnv("RPC_CLIENT_STUFF_ADDRESS", "127.0.0.1:8824"),
	}
}

func (c *RPCClientConfig) GetAuthAddress() string {
	return c.authAddress
}

func (c *RPCClientConfig) GetUserAddress() string {
	return c.userAddress
}

func (c *RPCClientConfig) GetCompanyAddress() string {
	return c.companyAddress
}

func (c *RPCClientConfig) GetJobAddress() string {
	return c.jobsAddress
}

func (c *RPCClientConfig) GetAdvertAddress() string {
	return c.advertAddress
}

func (c *RPCClientConfig) GetServicesAddress() string {
	return c.servicesAddress
}

func (c *RPCClientConfig) GetNewsfeedAddress() string {
	return c.newsfeedAddress
}

func (c *RPCClientConfig) GetStuffAddress() string {
	return c.stuffAddress
}

// MongoDB

type FileListConfig struct {
	address         []string
	user            string
	password        string
	database        string
	collectionFiles string
}

func NewFileListRepositoryConfig() *FileListConfig {
	// var addresses []string
	// for count := 0; ; count++ {
	// 	if address := os.Getenv("DB_MONGO_ADDRESS_USER_" + strconv.Itoa(count)); address != "" {
	// 		addresses = append(addresses, address)
	// 	} else {
	// 		break
	// 	}
	// }

	return &FileListConfig{
		address:         strings.Split(getEnv("DB_MONGO_ADDRESS", "192.168.1.13:27017"), ","), //addresses,
		user:            getEnv("DB_MONGO_USER_USER", "developer"),                            //os.Getenv("DB_MONGO_USER_USER"),
		password:        getEnv("DB_MONGO_PASSWORD_USER", "Qwerty123"),                        //os.Getenv("DB_MONGO_PASSWORD_USER"),
		database:        "file",                                                               //os.Getenv("DB_MONGO_DATABASE_FILE"),
		collectionFiles: "files",                                                              //os.Getenv("DB_MONGO_COLLECTION_FILES"),
	}
}

func (c *FileListConfig) GetAddress() []string {
	return c.address
}

func (c *FileListConfig) GetUser() string {
	return c.user
}

func (c *FileListConfig) GetPassword() string {
	return c.password
}

func (c *FileListConfig) GetDatabase() string {
	return c.database
}

func (c *FileListConfig) GetCollectionFiles() string {
	return c.collectionFiles
}

// Ceph

type FileStorageConfig struct {
	path string
}

func (c *FileStorageConfig) GetPath() string {
	return c.path
}

func NewFileStorageRepositoryConfig() *FileStorageConfig {
	return &FileStorageConfig{
		path: "data/",
	}
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
