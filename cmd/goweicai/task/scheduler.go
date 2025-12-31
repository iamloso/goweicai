package task

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/service"
	"github.com/robfig/cron/v3"
)

// JobConfig 任务配置
type JobConfig struct {
	Name    string
	Cron    string
	Enabled bool
	Handler JobHandler
}

// Scheduler 动态任务调度器
type Scheduler struct {
	cron        *cron.Cron
	jobs        []*JobConfig
	config      *conf.Scheduler
	log         *log.Helper
	ctx         context.Context
	wencaiSvc   *service.WencaiService
	baseInfoSvc *service.BaseInfoService
	ztInfoSvc   *service.ZtInfoService
}

// NewScheduler 创建动态任务调度器
func NewScheduler(
	wencaiSvc *service.WencaiService,
	baseInfoSvc *service.BaseInfoService,
	ztInfoSvc *service.ZtInfoService,
	c *conf.Scheduler,
	logger log.Logger,
) *Scheduler {
	helper := log.NewHelper(log.With(logger, "module", "scheduler"))
	ctx := context.Background()

	scheduler := &Scheduler{
		cron:        cron.New(cron.WithSeconds()),
		jobs:        make([]*JobConfig, 0),
		config:      c,
		log:         helper,
		ctx:         ctx,
		wencaiSvc:   wencaiSvc,
		baseInfoSvc: baseInfoSvc,
		ztInfoSvc:   ztInfoSvc,
	}

	// 注册所有任务
	scheduler.registerJobs()

	return scheduler
}

// registerJobs 注册所有任务
func (s *Scheduler) registerJobs() {
	// 注册股票涨停数据任务
	if s.config.Stock != nil {
		s.registerJob(&JobConfig{
			Name:    "stock",
			Cron:    s.getCronExpr(s.config.Stock.Cron, "0 */5 9-16 * * 1-5"),
			Enabled: s.config.Stock.Enabled,
			Handler: &StockJobHandler{svc: s.wencaiSvc, log: s.log},
		})
	}

	// 注册基础数据任务
	if s.config.BaseInfo != nil {
		s.registerJob(&JobConfig{
			Name:    "base_info",
			Cron:    s.getCronExpr(s.config.BaseInfo.Cron, "0 */5 9-16 * * 1-5"),
			Enabled: s.config.BaseInfo.Enabled,
			Handler: &BaseInfoJobHandler{svc: s.baseInfoSvc, log: s.log},
		})
	}

	// 注册涨停数据任务
	if s.config.ZtInfo != nil {
		s.registerJob(&JobConfig{
			Name:    "zt_info",
			Cron:    s.getCronExpr(s.config.ZtInfo.Cron, "0 */1 9-16 * * 1-5"),
			Enabled: s.config.ZtInfo.Enabled,
			Handler: &ZtInfoJobHandler{svc: s.ztInfoSvc, log: s.log},
		})
	}
}

// registerJob 注册单个任务
func (s *Scheduler) registerJob(job *JobConfig) {
	s.jobs = append(s.jobs, job)
}

// getCronExpr 获取 cron 表达式，如果为空则使用默认值
func (s *Scheduler) getCronExpr(cronExpr, defaultExpr string) string {
	if cronExpr == "" {
		return defaultExpr
	}
	return cronExpr
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	// 动态添加所有已启用的任务
	for _, job := range s.jobs {
		if !job.Enabled {
			s.log.Infof("[%s] 任务已禁用", job.Name)
			continue
		}

		// 动态生成任务函数
		jobFunc := s.createJobFunc(job)

		// 动态添加到 cron
		if _, err := s.cron.AddFunc(job.Cron, jobFunc); err != nil {
			s.log.Errorf("[%s] 添加任务失败: %v", job.Name, err)
			return fmt.Errorf("添加任务 %s 失败: %w", job.Name, err)
		}

		s.log.Infof("[%s] 任务已配置，Cron: %s", job.Name, job.Cron)
	}

	// 如果配置了启动时立即执行
	if s.config.RunOnStart {
		s.log.Info("启动时立即执行一次所有已启用任务...")
		s.runOnStartJobs()
	}

	// 启动定时任务
	s.cron.Start()
	s.log.Info("动态任务调度器已启动")

	return nil
}

// createJobFunc 动态创建任务执行函数
func (s *Scheduler) createJobFunc(job *JobConfig) func() {
	return func() {
		s.log.Infof("[%s] 开始执行任务...", job.Name)
		if err := job.Handler.Execute(s.ctx); err != nil {
			s.log.Errorf("[%s] 任务执行失败: %v", job.Name, err)
		} else {
			s.log.Infof("[%s] 任务执行成功", job.Name)
		}
	}
}

// runOnStartJobs 启动时立即执行所有已启用任务
func (s *Scheduler) runOnStartJobs() {
	for _, job := range s.jobs {
		if job.Enabled {
			jobFunc := s.createJobFunc(job)
			go jobFunc()
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.log.Info("动态任务调度器已停止")
}
