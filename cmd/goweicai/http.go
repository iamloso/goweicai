package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/service"
)

// HTTPServer HTTP 服务器
type HTTPServer struct {
	server *http.Server
	log    *log.Helper
}

// NewHTTPServer 创建 HTTP 服务器
func NewHTTPServer(httpSvc *service.HTTPService, c *conf.Server, logger log.Logger) *HTTPServer {
	helper := log.NewHelper(log.With(logger, "module", "http"))

	mux := http.NewServeMux()
	httpSvc.RegisterRoutes(mux)

	httpAddr := ":8000"
	if c != nil && c.Http != nil && c.Http.Addr != "" {
		httpAddr = c.Http.Addr
	}

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	return &HTTPServer{
		server: server,
		log:    helper,
	}
}

// Start 启动 HTTP 服务器
func (s *HTTPServer) Start() error {
	s.log.Infof("HTTP 服务器启动在 %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Errorf("HTTP 服务器启动失败: %v", err)
		return err
	}
	return nil
}

// Stop 停止 HTTP 服务器
func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Errorf("HTTP 服务器关闭失败: %v", err)
		return err
	}
	
	s.log.Info("HTTP 服务器已停止")
	return nil
}
