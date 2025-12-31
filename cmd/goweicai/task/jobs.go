package task

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/service"
)

// JobHandler 任务处理器接口
type JobHandler interface {
	Execute(ctx context.Context) error
	Name() string
}

// StockJobHandler 股票任务处理器
type StockJobHandler struct {
	svc *service.WencaiService
	log *log.Helper
}

func (h *StockJobHandler) Execute(ctx context.Context) error {
	return h.svc.FetchAndSaveStocks(ctx)
}

func (h *StockJobHandler) Name() string {
	return "stock"
}

// BaseInfoJobHandler 基础数据任务处理器
type BaseInfoJobHandler struct {
	svc *service.BaseInfoService
	log *log.Helper
}

func (h *BaseInfoJobHandler) Execute(ctx context.Context) error {
	return h.svc.FetchAndSaveBaseInfo(ctx)
}

func (h *BaseInfoJobHandler) Name() string {
	return "base_info"
}

// ZtInfoJobHandler 涨停数据任务处理器
type ZtInfoJobHandler struct {
	svc *service.ZtInfoService
	log *log.Helper
}

func (h *ZtInfoJobHandler) Execute(ctx context.Context) error {
	return h.svc.FetchAndSaveZtInfo(ctx)
}

func (h *ZtInfoJobHandler) Name() string {
	return "zt_info"
}
