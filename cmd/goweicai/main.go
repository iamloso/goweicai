package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/iamloso/goweicai/api/proto"
	"github.com/iamloso/goweicai/internal/biz"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/data"
	"github.com/iamloso/goweicai/internal/service"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs/config.yaml", "config path")
}

func main() {
	flag.Parse()

	// 初始化日志
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
	helper := log.NewHelper(logger)

	// 加载配置
	cfg := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	defer cfg.Close()

	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		panic(err)
	}

	helper.Info("配置加载成功")

	// 初始化数据层
	d, cleanup, err := data.NewData(bc.Data, logger)
	if err != nil {
		panic(fmt.Sprintf("failed to create data: %v", err))
	}
	defer cleanup()

	// 初始化仓库
	stockRepo := data.NewStockRepo(d, logger)

	// 初始化业务层
	stockUc := biz.NewStockUsecase(stockRepo)

	// 初始化服务层
	wencaiSvc := service.NewWencaiService(stockUc, bc.Wencai, logger)
	httpSvc := service.NewHTTPService(stockUc, logger)
	grpcSvc := service.NewGRPCService(stockUc, wencaiSvc, logger)

	// 启动 gRPC 服务器
	grpcAddr := ":9000"
	if bc.Server != nil && bc.Server.Grpc != nil && bc.Server.Grpc.Addr != "" {
		grpcAddr = bc.Server.Grpc.Addr
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		helper.Fatalf("gRPC 监听失败: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStockServiceServer(grpcServer, grpcSvc)

	go func() {
		helper.Infof("gRPC 服务器启动在 %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			helper.Errorf("gRPC 服务器启动失败: %v", err)
		}
	}()

	// 启动 HTTP 服务器
	mux := http.NewServeMux()
	httpSvc.RegisterRoutes(mux)
	
	httpAddr := ":8000"
	if bc.Server != nil && bc.Server.Http != nil && bc.Server.Http.Addr != "" {
		httpAddr = bc.Server.Http.Addr
	}
	
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	
	go func() {
		helper.Infof("HTTP 服务器启动在 %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			helper.Errorf("HTTP 服务器启动失败: %v", err)
		}
	}()

	// 创建定时任务
	c := cron.New(cron.WithSeconds()) // 支持秒级别的 cron 表达式
	ctx := context.Background()

	// 定义任务函数
	job := func() {
		helper.Info("开始执行定时任务...")
		if err := wencaiSvc.FetchAndSaveStocks(ctx); err != nil {
			helper.Errorf("定时任务执行失败: %v", err)
		} else {
			helper.Info("定时任务执行成功")
		}
	}

	// 添加定时任务
	cronExpr := bc.Scheduler.Cron
	if cronExpr == "" {
		cronExpr = "0 0 9 * * *" // 默认每天 9:00
	}
	
	_, err = c.AddFunc(cronExpr, job)
	if err != nil {
		helper.Fatalf("添加定时任务失败: %v", err)
	}

	helper.Infof("定时任务已配置，Cron 表达式: %s", cronExpr)

	// 如果配置了启动时立即执行
	if bc.Scheduler.RunOnStart {
		helper.Info("启动时立即执行一次任务...")
		job()
	}

	// 启动定时任务
	c.Start()
	helper.Info("定时任务调度器已启动")

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	helper.Info("收到退出信号，正在关闭...")
	
	// 停止定时任务
	c.Stop()
	helper.Info("定时任务调度器已停止")
	
	// 优雅关闭 gRPC 服务器
	grpcServer.GracefulStop()
	helper.Info("gRPC 服务器已停止")
	
	// 优雅关闭 HTTP 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		helper.Errorf("HTTP 服务器关闭失败: %v", err)
	} else {
		helper.Info("HTTP 服务器已停止")
	}
}
