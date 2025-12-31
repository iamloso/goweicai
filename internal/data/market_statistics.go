package data
package data

import (
	"context"
	"time"

	"github.com/go-kraken/kratos/v2/log"
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





































































































































}	return stats, nil	}		})			UpdateTime:          model.UpdateTime,			CreateTime:          model.CreateTime,			TwoConsecutiveCount: model.TwoConsecutiveCount,			MaxConsecutiveDays:  model.MaxConsecutiveDays,			BrokenCount:         model.BrokenCount,			LimitDownCount:      model.LimitDownCount,			LimitUpCount:        model.LimitUpCount,			TradeDate:           model.TradeDate,			ID:                  model.ID,		stats = append(stats, &biz.MarketStatistics{	for _, model := range models {	stats := make([]*biz.MarketStatistics, 0, len(models))	}		return nil, result.Error		r.log.Errorf("查询市场统计数据列表失败: %v", result.Error)	if result.Error != nil {	result := query.Order("trade_date DESC").Find(&models)	}		query = query.Where("trade_date <= ?", endDate.Format("2006-01-02"))	if !endDate.IsZero() {	}		query = query.Where("trade_date >= ?", startDate.Format("2006-01-02"))	if !startDate.IsZero() {	query := r.data.gormDB.WithContext(ctx)	var models []MarketStatisticsModelfunc (r *MarketStatisticsRepo) List(ctx context.Context, startDate, endDate time.Time) ([]*biz.MarketStatistics, error) {// List 获取数据列表}	}, nil		UpdateTime:          model.UpdateTime,		CreateTime:          model.CreateTime,		TwoConsecutiveCount: model.TwoConsecutiveCount,		MaxConsecutiveDays:  model.MaxConsecutiveDays,		BrokenCount:         model.BrokenCount,		LimitDownCount:      model.LimitDownCount,		LimitUpCount:        model.LimitUpCount,		TradeDate:           model.TradeDate,		ID:                  model.ID,	return &biz.MarketStatistics{	}		return nil, result.Error		r.log.Errorf("查询市场统计数据失败: %v", result.Error)		}			return nil, nil		if result.Error == gorm.ErrRecordNotFound {	if result.Error != nil {	result := r.data.gormDB.WithContext(ctx).Where("trade_date = ?", date.Format("2006-01-02")).First(&model)	var model MarketStatisticsModelfunc (r *MarketStatisticsRepo) GetByDate(ctx context.Context, date time.Time) (*biz.MarketStatistics, error) {// GetByDate 根据日期获取数据}	return nil	r.log.Infof("批量保存市场统计数据成功，共 %d 条", len(stats))	}		return result.Error		r.log.Errorf("批量保存市场统计数据失败: %v", result.Error)	if result.Error != nil {	}).Create(&models)		DoUpdates: clause.AssignmentColumns([]string{"limit_up_count", "limit_down_count", "broken_count", "max_consecutive_days", "two_consecutive_count", "update_time"}),		Columns:   []clause.Column{{Name: "trade_date"}},	result := r.data.gormDB.WithContext(ctx).Clauses(clause.OnConflict{	}		})			TwoConsecutiveCount: stat.TwoConsecutiveCount,			MaxConsecutiveDays:  stat.MaxConsecutiveDays,			BrokenCount:         stat.BrokenCount,			LimitDownCount:      stat.LimitDownCount,			LimitUpCount:        stat.LimitUpCount,			TradeDate:           stat.TradeDate,		models = append(models, &MarketStatisticsModel{	for _, stat := range stats {	models := make([]*MarketStatisticsModel, 0, len(stats))	}		return nil	if len(stats) == 0 {func (r *MarketStatisticsRepo) BatchSave(ctx context.Context, stats []*biz.MarketStatistics) error {// BatchSave 批量保存数据}	return nil	r.log.Infof("保存市场统计数据成功，日期: %s", stat.TradeDate.Format("2006-01-02"))	}		return result.Error		r.log.Errorf("保存市场统计数据失败: %v", result.Error)	if result.Error != nil {	}).Create(model)		DoUpdates: clause.AssignmentColumns([]string{"limit_up_count", "limit_down_count", "broken_count", "max_consecutive_days", "two_consecutive_count", "update_time"}),		Columns:   []clause.Column{{Name: "trade_date"}},	result := r.data.gormDB.WithContext(ctx).Clauses(clause.OnConflict{	}		TwoConsecutiveCount: stat.TwoConsecutiveCount,		MaxConsecutiveDays:  stat.MaxConsecutiveDays,		BrokenCount:         stat.BrokenCount,		LimitDownCount:      stat.LimitDownCount,		LimitUpCount:        stat.LimitUpCount,		TradeDate:           stat.TradeDate,	model := &MarketStatisticsModel{func (r *MarketStatisticsRepo) Save(ctx context.Context, stat *biz.MarketStatistics) error {// Save 保存单条数据}	}		log:  log.NewHelper(log.With(logger, "module", "data/market_statistics")),		data: data,	return &MarketStatisticsRepo{func NewMarketStatisticsRepo(data *Data, logger log.Logger) biz.MarketStatisticsRepo {// NewMarketStatisticsRepo 创建市场统计数据仓库}	log  *log.Helper	data *Datatype MarketStatisticsRepo struct {// MarketStatisticsRepo 市场统计数据仓库实现