package server

import (
	"log"
	"net"

	"gitlab.lan/Rightnao-site/microservices/grpc-proto/infoRPC"
	"gitlab.lan/Rightnao-site/microservices/info/pkg/db/maxmind"
	"gitlab.lan/Rightnao-site/microservices/info/pkg/db/postgres"
	"gitlab.lan/Rightnao-site/microservices/info/pkg/tracing"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents gRPC server
type Server struct {
	address string
	tracer  *tracing.Tracer
	db      *postgres.PostgresDB
	mDB     *maxmind.DB
}

// Settings for RPC server
type Settings struct {
	Address string
	Tracer  *tracing.Tracer
	DB      *postgres.PostgresDB
	Mdb     *maxmind.DB
}

// NewRPCServer creates new gRPC server instance
func NewRPCServer(settings Settings) (*Server, error) {
	return &Server{
		address: settings.Address,
		tracer:  settings.Tracer,
		db:      settings.DB,
		mDB:     settings.Mdb,
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

	infoRPC.RegisterInfoServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
