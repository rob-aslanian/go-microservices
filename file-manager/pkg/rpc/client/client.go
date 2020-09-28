package clientRPC

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/advertRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/authRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/companyRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/jobsRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/newsfeedRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/servicesRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"google.golang.org/grpc"
)

var (
	instance *grpcClient
)

type grpcClient struct {
	authConn   *grpc.ClientConn
	AuthClient authRPC.AuthServiceClient

	userConn   *grpc.ClientConn
	UserClient userRPC.UserServiceClient

	companyConn   *grpc.ClientConn
	CompanyClient companyRPC.CompanyServiceClient

	jobsConn   *grpc.ClientConn
	JobsClient jobsRPC.JobsServiceClient

	advertConn   *grpc.ClientConn
	AdvertClient advertRPC.AdvertServiceClient

	servicesConn   *grpc.ClientConn
	ServicesClient servicesRPC.ServicesServiceClient

	newsfeedConn   *grpc.ClientConn
	NewsfeedClient newsfeedRPC.NewsfeedServiceClient

	stuffConn   *grpc.ClientConn
	StuffClient stuffRPC.StuffServiceClient

	// ___Conn * grpc.ClientConn
	// __Client __RPC.__Client
}

func (c *grpcClient) start(config Configuration) (err error) {
	c.authConn, err = grpc.DialContext(
		context.Background(),
		config.GetAuthAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.AuthClient = authRPC.NewAuthServiceClient(c.authConn)

	c.userConn, err = grpc.DialContext(
		context.Background(),
		config.GetUserAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.UserClient = userRPC.NewUserServiceClient(c.userConn)

	c.companyConn, err = grpc.DialContext(
		context.Background(),
		config.GetCompanyAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.CompanyClient = companyRPC.NewCompanyServiceClient(c.companyConn)

	c.jobsConn, err = grpc.DialContext(
		context.Background(),
		config.GetJobAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.JobsClient = jobsRPC.NewJobsServiceClient(c.jobsConn)

	c.advertConn, err = grpc.DialContext(
		context.Background(),
		config.GetAdvertAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.AdvertClient = advertRPC.NewAdvertServiceClient(c.advertConn)

	c.servicesConn, err = grpc.DialContext(
		context.Background(),
		config.GetServicesAddress(),
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
		config.GetNewsfeedAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.NewsfeedClient = newsfeedRPC.NewNewsfeedServiceClient(c.newsfeedConn)

	c.stuffConn, err = grpc.DialContext(
		context.Background(),
		config.GetStuffAddress(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.StuffClient = stuffRPC.NewStuffServiceClient(c.stuffConn)

	// c.___Conn, err = grpc.DialContext(context.Background(), config.Get___Address(), grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.___Client = ___RPC.New___ServiceClient(c.___Conn)

	return
}

func (c *grpcClient) Close() {
	c.authConn.Close()
	c.userConn.Close()
	c.companyConn.Close()
	c.jobsConn.Close()
	c.advertConn.Close()
	c.servicesConn.Close()
	c.newsfeedConn.Close()
	c.stuffConn.Close()
	//...
}

func NewRPCClient(config Configuration) *grpcClient {
	instance = &grpcClient{}
	instance.start(config)
	return instance
}

func GetRPCClient() *grpcClient {
	return instance
}

func (c *grpcClient) Auth() authRPC.AuthServiceClient {
	return c.AuthClient
}

func (c *grpcClient) User() userRPC.UserServiceClient {
	return c.UserClient
}

func (c *grpcClient) Company() companyRPC.CompanyServiceClient {
	return c.CompanyClient
}

func (c *grpcClient) Jobs() jobsRPC.JobsServiceClient {
	return c.JobsClient
}

func (c *grpcClient) Advert() advertRPC.AdvertServiceClient {
	return c.AdvertClient
}

func (c *grpcClient) Services() servicesRPC.ServicesServiceClient {
	return c.ServicesClient
}

func (c *grpcClient) Newsfeed() newsfeedRPC.NewsfeedServiceClient {
	return c.NewsfeedClient
}

func (c *grpcClient) Stuff() stuffRPC.StuffServiceClient {
	return c.StuffClient
}
