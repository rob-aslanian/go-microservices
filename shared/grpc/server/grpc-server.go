package server

import (
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type Server struct {
	registrator Registrator
	closables   []Closable
}

type Closable interface {
	Close() error
}

type Registrator func(*grpc.Server)

func NewServer(registrator Registrator, resources ...Closable) *Server {
	return &Server{
		registrator: registrator,
		closables:   resources,
	}
}

func (s *Server) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
	)

	s.registrator(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) AddClosable(resource Closable) {
	s.closables = append(s.closables, resource)
}

func (s *Server) Close() error {
	for _, c := range s.closables {
		c.Close()
	}
	return nil // maybe better to return error which wrappes multiple errors
}
