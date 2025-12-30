package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamloso/goweicai/internal/biz"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/data"
	"github.com/iamloso/goweicai/internal/service"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/yaml.v3"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "configs/config.yaml", "config path")
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
	d, cleanup, err := data.NewData(bc.Data, &bc, logger)
	if err != nil {
		panic(fmt.Sprintf("failed to create data: %v", err))
	}
	defer cleanup()

	// 初始化仓库
	stockRepo := data.NewStockRepo(d, logger)
	baseInfoRepo := data.NewBaseInfoRepo(d, logger)
	ztInfoRepo := data.NewZtInfoRepo(d, logger)

	// 初始化业务层
	stockUc := biz.NewStockUsecase(stockRepo)
	baseInfoUc := biz.NewBaseInfoUsecase(baseInfoRepo)
	ztInfoUc := biz.NewZtInfoUsecase(ztInfoRepo)

	// 初始化服务层
	wencaiSvc := service.NewWencaiService(stockUc, bc.Wencai, logger)
	baseInfoSvc := service.NewBaseInfoService(baseInfoUc, bc.BaseInfo, &bc, logger)
	ztInfoSvc := service.NewZtInfoService(ztInfoUc, bc.ZtInfo, &bc, logger)
	httpSvc := service.NewHTTPService(stockUc, logger)
	grpcSvc := service.NewGRPCService(stockUc, wencaiSvc, logger)

	// 创建 HTTP 服务器
	httpServer := NewHTTPServer(httpSvc, bc.Server, logger)
	go func() {
		if err := httpServer.Start(); err != nil {
			helper.Errorf("HTTP 服务器错误: %v", err)
		}
	}()

	// 创建 gRPC 服务器
	grpcServer, err := NewGRPCServer(grpcSvc, bc.Server, logger)
	if err != nil {
		helper.Fatalf("创建 gRPC 服务器失败: %v", err)
	}
	go func() {
		if err := grpcServer.Start(); err != nil {
			helper.Errorf("gRPC 服务器错误: %v", err)
		}
	}()

	// 创建定时任务调度器
	cronScheduler := NewCronScheduler(wencaiSvc, baseInfoSvc, ztInfoSvc, bc.Scheduler, logger)
	if err := cronScheduler.Start(); err != nil {
		helper.Fatalf("启动定时任务调度器失败: %v", err)
	}

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	helper.Info("收到退出信号，正在关闭...")

	// 优雅关闭所有服务
	cronScheduler.Stop()
	grpcServer.Stop()
	if err := httpServer.Stop(); err != nil {
		helper.Errorf("HTTP 服务器关闭错误: %v", err)
	}

	helper.Info("所有服务已停止")
}
