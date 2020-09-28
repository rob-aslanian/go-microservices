package serverRPC

import (
	"log"
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"gitlab.lan/Rightnao-site/microservices/grpc-proto/rentalRPC"
	tracing "gitlab.lan/Rightnao-site/microservices/rental/internal/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents gRPC server
type Server struct {
	house   House
	address string
	tracer  *tracing.Tracer
}

// Settings for RPC server
type Settings struct {
	Address string
	House   House
	Tracer  *tracing.Tracer
}

// NewRPCServer creates new gRPC server instance
func NewRPCServer(settings Settings) (*Server, error) {
	return &Server{
		house:   settings.House,
		address: settings.Address,
		tracer:  settings.Tracer,
	}, nil
}

// Launch starts gRPC server
func (s *Server) Launch() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Panicln("error while luanching gRPC server:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.tracer.GetTracer()),
		),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(s.tracer.GetTracer()),
		),
	)

	rentalRPC.RegisterRentalServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
