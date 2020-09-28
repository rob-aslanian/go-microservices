package rpc

import (
	"context"
	"log"
	"os"
	"sync"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/chatRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/groupsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/networkRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/notificationsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/searchRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/shopRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/statisticsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var (
	userAddress          = getEnv("ADDR_GRPC_USER", "127.0.0.1:8822")
	authAddress          = getEnv("ADDR_GRPC_AUTH", "127.0.0.1:50000")
	infoAddress          = getEnv("ADDR_GRPC_INFO", "127.0.0.1:8804")
	companyAddress       = getEnv("ADDR_GRPC_COMPANY", "127.0.0.1:50000")
	networkAddress       = getEnv("ADDR_GRPC_NETWORK", "localhost:8806")
	searchAddress        = getEnv("ADDR_GRPC_SEARCH", "127.0.0.1:8810")
	jobsAddress          = getEnv("ADDR_GRPC_JOBS", "192.168.1.13:8812")
	chatAddress          = getEnv("ADDR_GRPC_CHAT", "localhost:8808")
	advertAddress        = getEnv("ADDR_GRPC_ADVERT", "127.0.0.1:8021")
	stuffAddress         = getEnv("ADDR_GRPC_STUFF", "127.0.0.1:8022")
	statisticsAddress    = getEnv("ADDR_GRPC_STAT", "192.168.1.111:8814")
	notificationsAddress = getEnv("ADDR_GRPC_NOTTIFICATIONS", "192.168.1.111:8814")
	servicesAddress      = getEnv("ADDR_GRPC_SERVICES", "192.168.1.111:6969")
	newsfeedAddress      = getEnv("ADDR_GRPC_NEWSFEED", "192.168.1.111:8814")
	groupsAddress        = getEnv("ADDR_GRPC_GROUPS", "192.168.1.111:8814")
	shopAddress          = getEnv("ADDR_GRPC_SHOP", "192.168.1.111:8814")
	rentalAddress        = getEnv("ADDR_GRPC_RENTAL", "192.168.1.111:8814")
)

var (
	instance *grpcClient
	once     sync.Once
)

type grpcClient struct {
	userConn   *grpc.ClientConn
	UserClient userRPC.UserServiceClient

	authConn   *grpc.ClientConn
	AuthClient authRPC.AuthServiceClient

	infoConn   *grpc.ClientConn
	InfoClient infoRPC.InfoServiceClient

	companyConn   *grpc.ClientConn
	CompanyClient companyRPC.CompanyServiceClient

	networkConn   *grpc.ClientConn
	NetworkClient networkRPC.NetworkServiceClient

	searchConn   *grpc.ClientConn
	SearchClient searchRPC.SearchServiceClient

	jobsConn   *grpc.ClientConn
	JobsClient jobsRPC.JobsServiceClient

	chatConn   *grpc.ClientConn
	ChatClient chatRPC.ChatServiceClient

	advertConn   *grpc.ClientConn
	AdvertClient advertRPC.AdvertServiceClient

	stuffConn   *grpc.ClientConn
	StuffClient stuffRPC.StuffServiceClient

	statisticsConn   *grpc.ClientConn
	StatisticsClient statisticsRPC.StatisticsClient

	notificationsConn   *grpc.ClientConn
	NotificationsClient notificationsRPC.NotificationsServiceClient

	servicesConn   *grpc.ClientConn
	ServicesClient servicesRPC.ServicesServiceClient

	newsfeedConn   *grpc.ClientConn
	NewsfeedClient newsfeedRPC.NewsfeedServiceClient

	groupsConn   *grpc.ClientConn
	GroupsClient groupsRPC.GroupsServiceClient

	shopConn   *grpc.ClientConn
	ShopClient shopRPC.ShopServiceClient

	rentalConn   *grpc.ClientConn
	RentalClient rentalRPC.RentalServiceClient

	// ___Conn * grpc.ClientConn
	// __Client __RPC.__Client
}

func (c *grpcClient) start() (err error) {
	c.userConn, err = grpc.DialContext(
		context.Background(),
		userAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.UserClient = userRPC.NewUserServiceClient(c.userConn)

	c.authConn, err = grpc.DialContext(
		context.Background(),
		authAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.AuthClient = authRPC.NewAuthServiceClient(c.authConn)

	c.infoConn, err = grpc.DialContext(
		context.Background(),
		infoAddress, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		log.Fatal(err)
	}
	c.InfoClient = infoRPC.NewInfoServiceClient(c.infoConn)

	c.companyConn, err = grpc.DialContext(
		context.Background(),
		companyAddress, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		log.Fatal(err)
	}
	c.CompanyClient = companyRPC.NewCompanyServiceClient(c.companyConn)

	c.networkConn, err = grpc.DialContext(
		context.Background(),
		networkAddress, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		log.Fatal(err)
	}
	c.NetworkClient = networkRPC.NewNetworkServiceClient(c.networkConn)

	c.searchConn, err = grpc.DialContext(
		context.Background(),
		searchAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		log.Fatal(err)
	}
	c.SearchClient = searchRPC.NewSearchServiceClient(c.searchConn)

	c.jobsConn, err = grpc.DialContext(
		context.Background(),
		jobsAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.JobsClient = jobsRPC.NewJobsServiceClient(c.jobsConn)

	c.chatConn, err = grpc.DialContext(
		context.Background(),
		chatAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.ChatClient = chatRPC.NewChatServiceClient(c.chatConn)

	c.statisticsConn, err = grpc.DialContext(
		context.Background(),
		statisticsAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.StatisticsClient = statisticsRPC.NewStatisticsClient(c.statisticsConn)

	c.advertConn, err = grpc.DialContext(
		context.Background(),
		advertAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.AdvertClient = advertRPC.NewAdvertServiceClient(c.advertConn)

	c.stuffConn, err = grpc.DialContext(
		context.Background(),
		stuffAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.StuffClient = stuffRPC.NewStuffServiceClient(c.stuffConn)

	c.notificationsConn, err = grpc.DialContext(
		context.Background(),
		notificationsAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.NotificationsClient = notificationsRPC.NewNotificationsServiceClient(c.notificationsConn)

	c.servicesConn, err = grpc.DialContext(
		context.Background(),
		servicesAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.ServicesClient = servicesRPC.NewServicesServiceClient(c.servicesConn)

	c.newsfeedConn, err = grpc.DialContext(
		context.Background(),
		newsfeedAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.NewsfeedClient = newsfeedRPC.NewNewsfeedServiceClient(c.newsfeedConn)

	c.groupsConn, err = grpc.DialContext(
		context.Background(),
		groupsAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.GroupsClient = groupsRPC.NewGroupsServiceClient(c.groupsConn)

	c.shopConn, err = grpc.DialContext(
		context.Background(),
		shopAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.ShopClient = shopRPC.NewShopServiceClient(c.shopConn)

	c.rentalConn, err = grpc.DialContext(
		context.Background(),
		rentalAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.RentalClient = rentalRPC.NewRentalServiceClient(c.rentalConn)

	//...

	return
}

func (c *grpcClient) Close() {
	c.userConn.Close()
	c.authConn.Close()
	c.infoConn.Close()
	c.companyConn.Close()
	c.networkConn.Close()
	c.searchConn.Close()
	c.jobsConn.Close()
	c.chatConn.Close()
	c.advertConn.Close()
	c.stuffConn.Close()
	c.statisticsConn.Close()
	c.notificationsConn.Close()
	c.newsfeedConn.Close()
	c.groupsConn.Close()
	c.shopConn.Close()
	c.rentalConn.Close()
	//...
}

func GetGrpcClient() *grpcClient {
	once.Do(func() {
		instance = &grpcClient{}
		instance.start()
	})

	return instance
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
