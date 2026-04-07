package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/biz"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MarketStatisticsModel 市场统计数据模型
type MarketStatisticsModel struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement"`
	TradeDate           time.Time `gorm:"type:date;uniqueIndex;not null"`
	LimitUpCount        int       `gorm:"not null;default:0"`
	LimitDownCount      int       `gorm:"not null;default:0"`
	BrokenCount         int       `gorm:"not null;default:0"`
	MaxConsecutiveDays  int       `gorm:"not null;default:0"`
	TwoConsecutiveCount int       `gorm:"not null;default:0"`
	CreateTime          time.Time `gorm:"autoCreateTime"`
	UpdateTime          time.Time `gorm:"autoUpdateTime"`
}

func (MarketStatisticsModel) TableName() string {
	return "market_statistics"
}

// MarketStatisticsRepo 市场统计数据仓库实现
type MarketStatisticsRepo struct {
	data *Data
	log  *log.Helper
}

// NewMarketStatisticsRepo 创建市场统计数据仓库
func NewMarketStatisticsRepo(data *Data, logger log.Logger) biz.MarketStatisticsRepo {
	return &MarketStatisticsRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/market_statistics")),
	}
}

// Save 保存单条数据
func (r *MarketStatisticsRepo) Save(ctx context.Context, stat *biz.MarketStatistics) error {
	model := &MarketStatisticsModel{
		TradeDate:           stat.TradeDate,
		LimitUpCount:        stat.LimitUpCount,
		LimitDownCount:      stat.LimitDownCount,
		BrokenCount:         stat.BrokenCount,
		MaxConsecutiveDays:  stat.MaxConsecutiveDays,
		TwoConsecutiveCount: stat.TwoConsecutiveCount,
	}

	result := r.data.gormDB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "trade_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"limit_up_count", "limit_down_count", "broken_count", "max_consecutive_days", "two_consecutive_count", "update_time"}),
	}).Create(model)

	if result.Error != nil {
		r.log.Errorf("保存市场统计数据失败: %v", result.Error)
		return result.Error
	}

	r.log.Infof("保存市场统计数据成功，日期: %s", stat.TradeDate.Format("2006-01-02"))
	return nil
}

// BatchSave 批量保存数据
func (r *MarketStatisticsRepo) BatchSave(ctx context.Context, stats []*biz.MarketStatistics) error {
	if len(stats) == 0 {
		return nil
	}

	models := make([]*MarketStatisticsModel, 0, len(stats))
	for _, stat := range stats {
		models = append(models, &MarketStatisticsModel{
			TradeDate:           stat.TradeDate,
			LimitUpCount:        stat.LimitUpCount,
			LimitDownCount:      stat.LimitDownCount,
			BrokenCount:         stat.BrokenCount,
			MaxConsecutiveDays:  stat.MaxConsecutiveDays,
			TwoConsecutiveCount: stat.TwoConsecutiveCount,
		})
	}

	result := r.data.gormDB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "trade_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"limit_up_count", "limit_down_count", "broken_count", "max_consecutive_days", "two_consecutive_count", "update_time"}),
	}).Create(&models)

	if result.Error != nil {
		r.log.Errorf("批量保存市场统计数据失败: %v", result.Error)
		return result.Error
	}

	r.log.Infof("批量保存市场统计数据成功，共 %d 条", len(stats))
	return nil
}

// GetByDate 根据日期获取数据
func (r *MarketStatisticsRepo) GetByDate(ctx context.Context, date time.Time) (*biz.MarketStatistics, error) {
	var model MarketStatisticsModel
	result := r.data.gormDB.WithContext(ctx).Where("trade_date = ?", date.Format("2006-01-02")).First(&model)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("查询市场统计数据失败: %v", result.Error)
		return nil, result.Error
	}

	return &biz.MarketStatistics{
		ID:                  model.ID,
		TradeDate:           model.TradeDate,
		LimitUpCount:        model.LimitUpCount,
		LimitDownCount:      model.LimitDownCount,
		BrokenCount:         model.BrokenCount,
		MaxConsecutiveDays:  model.MaxConsecutiveDays,
		TwoConsecutiveCount: model.TwoConsecutiveCount,
		CreateTime:          model.CreateTime,
		UpdateTime:          model.UpdateTime,
	}, nil
}

// List 获取数据列表
func (r *MarketStatisticsRepo) List(ctx context.Context, startDate, endDate time.Time) ([]*biz.MarketStatistics, error) {
	var models []MarketStatisticsModel
	query := r.data.gormDB.WithContext(ctx)

	if !startDate.IsZero() {
		query = query.Where("trade_date >= ?", startDate.Format("2006-01-02"))
	}

	if !endDate.IsZero() {
		query = query.Where("trade_date <= ?", endDate.Format("2006-01-02"))
	}

	result := query.Order("trade_date DESC").Find(&models)

	if result.Error != nil {
		r.log.Errorf("查询市场统计数据列表失败: %v", result.Error)
		return nil, result.Error
	}

	stats := make([]*biz.MarketStatistics, 0, len(models))
	for _, model := range models {
		stats = append(stats, &biz.MarketStatistics{
			ID:                  model.ID,
			TradeDate:           model.TradeDate,
			LimitUpCount:        model.LimitUpCount,
			LimitDownCount:      model.LimitDownCount,
			BrokenCount:         model.BrokenCount,
			MaxConsecutiveDays:  model.MaxConsecutiveDays,
			TwoConsecutiveCount: model.TwoConsecutiveCount,
			CreateTime:          model.CreateTime,
			UpdateTime:          model.UpdateTime,
		})
	}

	return stats, nil
}
