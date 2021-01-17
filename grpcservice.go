package walmart

import (
	"fmt"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// GrpcService is a wrapper for grpc service
type GrpcService interface {
	RegisterHandlers()
	Serve()
}

type grpcChatService struct {
	lis    net.Listener
	server *grpc.Server
}

func (gs *grpcChatService) RegisterHandlers() {
	glog.V(9).Infof("Registerig gRPC Handlers")
	registerServer(gs.server)
}

func (gs *grpcChatService) Serve() {
	if err := gs.server.Serve(gs.lis); err != nil {
		glog.Fatal("grpc server failed to start")
	}
}

// NewGrpcService creates an instance of GrpcService
func NewGrpcChatService(portNum int) GrpcService {
	glog.V(9).Infof("Begin NewGrpcService")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	return &grpcChatService{lis: lis, server: grpcServer}
}
