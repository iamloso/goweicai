package main

import (
	"net"

	pb "github.com/iamloso/goweicai/api/proto"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

// GRPCServer gRPC 服务器
type GRPCServer struct {
	server *grpc.Server
	lis    net.Listener
	log    *log.Helper
}

// NewGRPCServer 创建 gRPC 服务器
func NewGRPCServer(grpcSvc *service.GRPCService, c *conf.Server, logger log.Logger) (*GRPCServer, error) {
	helper := log.NewHelper(log.With(logger, "module", "grpc"))

	grpcAddr := ":9000"
	if c != nil && c.Grpc != nil && c.Grpc.Addr != "" {
		grpcAddr = c.Grpc.Addr
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		helper.Errorf("gRPC 监听失败: %v", err)
		return nil, err
	}

	server := grpc.NewServer()
	pb.RegisterStockServiceServer(server, grpcSvc)

	return &GRPCServer{
		server: server,
		lis:    lis,
		log:    helper,
	}, nil
}

// Start 启动 gRPC 服务器
func (s *GRPCServer) Start() error {
	s.log.Infof("gRPC 服务器启动在 %s", s.lis.Addr().String())
	if err := s.server.Serve(s.lis); err != nil {
		s.log.Errorf("gRPC 服务器启动失败: %v", err)
		return err
	}
	return nil
}

// Stop 停止 gRPC 服务器
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
	s.log.Info("gRPC 服务器已停止")
}
