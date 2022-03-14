package main

import (
	"context"
	"net"

	h "user-unary/handler"
	svc "user-unary/service"

	mongo "grpc-applications/crud-mongodb"
	upb "grpc-applications/protoc/protobuf-user"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	unaryServiceHost = "127.0.0.1:8080"
)

func main() {
	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(0))

	logger.Debug("Internal dependencies creation...")
	mongoCl, err := mongo.NewMongoConnection(context.Background(), logger, "127.0.0.1", "27017")
	if err != nil {
		logger.With(zap.Error(err)).Fatal("failed to create mongoDB client")
	}
	service := svc.NewUnaryService(logger, mongoCl)
	handler := h.NewUnaryHandler(logger, service)

	logger.Debug("gRPC server creation...")
	listener, err := net.Listen("tcp", unaryServiceHost)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("failed to allocate TCP port")
	}
	server := grpc.NewServer()
	upb.RegisterUserServiceServer(server, handler)
	if server.Serve(listener) != nil {
		logger.With(zap.Error(err)).Fatal("failed to serve gRPC server")
	}
}
