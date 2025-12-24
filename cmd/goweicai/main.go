package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
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

	// 执行任务
	helper.Info("开始获取股票数据...")
	ctx := context.Background()
	if err := wencaiSvc.FetchAndSaveStocks(ctx); err != nil {
		helper.Errorf("failed to fetch and save stocks: %v", err)
		os.Exit(1)
	}

	helper.Info("任务完成")
}
