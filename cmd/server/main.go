package main

import (
	"flag"
	"walmart"

	"github.com/golang/glog"
)

func main() {
	port := flag.Int("port", 13000, "The grpc server port")
	flag.Parse()

	grpcService := walmart.NewGrpcChatService(*port)
	grpcService.RegisterHandlers()
	glog.Info("Starting gRPC server")
	grpcService.Serve()
}
