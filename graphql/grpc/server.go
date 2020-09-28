package rpc

import (
	"log"
	"net"

	// "gitlab.lan/Rightnao-site/microservices/grpc-proto/userRPC"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct{}

func (g *GrpcServer) Launch(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		// TODO: log error
		log.Fatal("Listening gRPC port failed", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
		grpc.StreamInterceptor(otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())),
	)

	// userRPC.RegisterUserServiceServer(grpcServer, g)

	// __RPC.Register__Server(grpcServer, g) // TODO: how to make it in other way?

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve gRPC:", err)
	}
}
