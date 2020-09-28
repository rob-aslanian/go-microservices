package serverRPC

import (
	"log"
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/stuffRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents gRPC server
type Server struct {
	service Service
	address string
}

// Settings for RPC server
type Settings struct {
	Address string
	Service Service
}

// NewRPCServer creates new gRPC server instance
func NewRPCServer(settings Settings) (*Server, error) {
	return &Server{
		service: settings.Service,
		address: settings.Address,
	}, nil
}

// Launch starts gRPC server
func (s *Server) Launch() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Panicln("error while luanching gRPC server:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
	)

	stuffRPC.RegisterStuffServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
