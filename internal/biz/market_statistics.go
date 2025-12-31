package biz

import (
	"context"
	"time"
)

// MarketStatistics 市场统计数据模型
type MarketStatistics struct {
	ID                 int64     `json:"id"`
	TradeDate          time.Time `json:"trade_date"`          // 交易日期
	LimitUpCount       int       `json:"limit_up_count"`      // 涨停数
	LimitDownCount     int       `json:"limit_down_count"`    // 跌停数
	BrokenCount        int       `json:"broken_count"`        // 开板数
	MaxConsecutiveDays int       `json:"max_consecutive_days"` // 连板高度
	TwoConsecutiveCount int      `json:"two_consecutive_count"` // 二连板数量
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time"`
}

// MarketStatisticsRepo 市场统计数据仓库接口
type MarketStatisticsRepo interface {
	// Save 保存单条数据
	Save(ctx context.Context, stat *MarketStatistics) error
	// BatchSave 批量保存数据
	BatchSave(ctx context.Context, stats []*MarketStatistics) error
	// GetByDate 根据日期获取数据
	GetByDate(ctx context.Context, date time.Time) (*MarketStatistics, error)
	// List 获取数据列表
	List(ctx context.Context, startDate, endDate time.Time) ([]*MarketStatistics, error)
}

// MarketStatisticsUsecase 市场统计用例
type MarketStatisticsUsecase struct {
	repo MarketStatisticsRepo
}

// NewMarketStatisticsUsecase 创建市场统计用例
func NewMarketStatisticsUsecase(repo MarketStatisticsRepo) *MarketStatisticsUsecase {
	return &MarketStatisticsUsecase{
		repo: repo,
	}
}

// SaveMarketStatistics 保存市场统计数据
func (uc *MarketStatisticsUsecase) SaveMarketStatistics(ctx context.Context, stat *MarketStatistics) error {
	return uc.repo.Save(ctx, stat)
}

// BatchSaveMarketStatistics 批量保存市场统计数据
func (uc *MarketStatisticsUsecase) BatchSaveMarketStatistics(ctx context.Context, stats []*MarketStatistics) error {
	return uc.repo.BatchSave(ctx, stats)
}

// GetByDate 根据日期获取市场统计数据
func (uc *MarketStatisticsUsecase) GetByDate(ctx context.Context, date time.Time) (*MarketStatistics, error) {
	return uc.repo.GetByDate(ctx, date)
}

// List 获取市场统计数据列表
func (uc *MarketStatisticsUsecase) List(ctx context.Context, startDate, endDate time.Time) ([]*MarketStatistics, error) {
	return uc.repo.List(ctx, startDate, endDate)
}
