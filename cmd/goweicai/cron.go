package main

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/service"
	"github.com/robfig/cron/v3"
)

// CronScheduler 定时任务调度器
type CronScheduler struct {
	cron        *cron.Cron
	wencaiSvc   *service.WencaiService
	baseInfoSvc *service.BaseInfoService
	ztInfoSvc   *service.ZtInfoService
	config      *conf.Scheduler
	log         *log.Helper
}

// NewCronScheduler 创建定时任务调度器
func NewCronScheduler(
	wencaiSvc *service.WencaiService,
	baseInfoSvc *service.BaseInfoService,
	ztInfoSvc *service.ZtInfoService,
	c *conf.Scheduler,
	logger log.Logger,
) *CronScheduler {
	helper := log.NewHelper(log.With(logger, "module", "cron"))

	return &CronScheduler{
		cron:        cron.New(cron.WithSeconds()), // 支持秒级别的 cron 表达式
		wencaiSvc:   wencaiSvc,
		baseInfoSvc: baseInfoSvc,
		ztInfoSvc:   ztInfoSvc,
		config:      c,
		log:         helper,
	}
}

// Start 启动定时任务调度器
func (s *CronScheduler) Start() error {
	ctx := context.Background()

	// 定义股票涨停数据任务函数
	stockJob := func() {
		s.log.Info("开始执行股票涨停数据定时任务...")
		if err := s.wencaiSvc.FetchAndSaveStocks(ctx); err != nil {
			s.log.Errorf("股票涨停数据定时任务执行失败: %v", err)
		} else {
			s.log.Info("股票涨停数据定时任务执行成功")
		}
	}

	// 定义基础数据任务函数
	baseInfoJob := func() {
		s.log.Info("开始执行基础数据定时任务...")
		if err := s.baseInfoSvc.FetchAndSaveBaseInfo(ctx); err != nil {
			s.log.Errorf("基础数据定时任务执行失败: %v", err)
		} else {
			s.log.Info("基础数据定时任务执行成功")
		}
	}

	// 定义涨停数据任务函数
	ztInfoJob := func() {
		s.log.Info("开始执行涨停数据定时任务...")
		if err := s.ztInfoSvc.FetchAndSaveZtInfo(ctx); err != nil {
			s.log.Errorf("涨停数据定时任务执行失败: %v", err)
		} else {
			s.log.Info("涨停数据定时任务执行成功")
		}
	}

	// 添加股票涨停数据定时任务
	if s.config.Stock != nil && s.config.Stock.Enabled {
		cronExpr := s.config.Stock.Cron
		if cronExpr == "" {
			cronExpr = "0 0 9 * * *" // 默认每天 9:00
		}

		if _, err := s.cron.AddFunc(cronExpr, stockJob); err != nil {
			s.log.Errorf("添加股票涨停数据定时任务失败: %v", err)
			return err
		}
		s.log.Infof("股票涨停数据定时任务已配置，Cron 表达式: %s", cronExpr)
	} else {
		s.log.Info("股票涨停数据定时任务已禁用")
	}

	// 添加基础数据定时任务
	if s.config.BaseInfo != nil && s.config.BaseInfo.Enabled {
		cronExpr := s.config.BaseInfo.Cron
		if cronExpr == "" {
			cronExpr = "0 15 9 * * *" // 默认每天 9:15
		}

		if _, err := s.cron.AddFunc(cronExpr, baseInfoJob); err != nil {
			s.log.Errorf("添加基础数据定时任务失败: %v", err)
			return err
		}
		s.log.Infof("基础数据定时任务已配置，Cron 表达式: %s", cronExpr)
	} else {
		s.log.Info("基础数据定时任务已禁用")
	}

	// 添加涨停数据定时任务
	if s.config.ZtInfo != nil && s.config.ZtInfo.Enabled {
		cronExpr := s.config.ZtInfo.Cron
		if cronExpr == "" {
			cronExpr = "0 */1 * * * *" // 默认每分钟执行一次（秒 分 时 日 月 周）
		}

		if _, err := s.cron.AddFunc(cronExpr, ztInfoJob); err != nil {
			s.log.Errorf("添加涨停数据定时任务失败: %v", err)
			return err
		}
		s.log.Infof("涨停数据定时任务已配置，Cron 表达式: %s", cronExpr)
	} else {
		s.log.Info("涨停数据定时任务已禁用")
	}

	// 如果配置了启动时立即执行
	if s.config.RunOnStart {
		s.log.Info("启动时立即执行一次任务...")

		// 只执行已启用的任务
		if s.config.Stock != nil && s.config.Stock.Enabled {
			go stockJob()
			time.Sleep(2 * time.Second) // 等待2秒，避免请求过于密集
		}

		if s.config.BaseInfo != nil && s.config.BaseInfo.Enabled {
			go baseInfoJob()
			time.Sleep(2 * time.Second) // 等待2秒
		}

		if s.config.ZtInfo != nil && s.config.ZtInfo.Enabled {
			go ztInfoJob()
		}
	}

	// 启动定时任务
	s.cron.Start()
	s.log.Info("定时任务调度器已启动")

	return nil
}

// Stop 停止定时任务调度器
func (s *CronScheduler) Stop() {
	s.cron.Stop()
	s.log.Info("定时任务调度器已停止")
}
