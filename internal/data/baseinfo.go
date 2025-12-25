package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/biz"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type baseInfoRepo struct {
	data *Data
	log  *log.Helper
}

// NewBaseInfoRepo 创建基础数据仓库
func NewBaseInfoRepo(data *Data, logger log.Logger) biz.BaseInfoRepo {
	return &baseInfoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/baseinfo")),
	}
}

// BaseInfoModel 数据库模型
type BaseInfoModel struct {
	ID                        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	StockName                 string    `gorm:"column:stock_name;type:varchar(100)"`
	LatestPrice               float64   `gorm:"column:latest_price;type:decimal(10,2)"`
	AuctionChangeRate         float64   `gorm:"column:auction_change_rate;type:decimal(10,2)"`
	LatestChangeRate          float64   `gorm:"column:latest_change_rate;type:decimal(10,2)"`
	AuctionUnmatchedAmountStr string    `gorm:"column:auction_unmatched_amount_str;type:varchar(50);comment:竞价未匹配金额（如：12万、1.3亿）"`
	MorningAuctionAmountStr   string    `gorm:"column:morning_auction_amount_str;type:varchar(50);comment:竞价金额（如：12万、1.3亿）"`
	TurnoverStr               string    `gorm:"column:turnover_str;type:varchar(50);comment:成交额（如：12万、1.3亿）"`
	CirculationMarketValue    float64   `gorm:"column:circulation_market_value;type:decimal(20,2)"`
	CirculationMarketValueStr string    `gorm:"column:circulation_market_value_str;type:varchar(50);comment:流通市值（如：12万、1.3亿）"`
	StockCode                 string    `gorm:"column:stock_code;type:varchar(20);index:idx_stock_date"`
	TradeDate                 time.Time `gorm:"column:trade_date;type:date;index:idx_stock_date"`
	MarketCode                string    `gorm:"column:market_code;type:varchar(10)"`
	Code                      string    `gorm:"column:code;type:varchar(50);uniqueIndex:uk_code_date"`
	Turnover                  float64   `gorm:"column:turnover;type:decimal(20,2)"`
	MorningAuctionAmount      int64     `gorm:"column:morning_auction_amount;type:bigint"`
	AuctionUnmatchedAmount    int64     `gorm:"column:auction_unmatched_amount;type:bigint"`
	CreateTime                time.Time `gorm:"column:create_time;type:datetime;autoCreateTime"`
	UpdateTime                time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime"`
	CompanyHighlights         string    `gorm:"column:company_highlights;type:text"`
	IndustryCategory          string    `gorm:"column:industry_category;type:varchar(100)"`
	ConceptTheme              string    `gorm:"column:concept_theme;type:varchar(500)"`
	ConsecutiveLimitDays      int       `gorm:"column:consecutive_limit_days;type:int"`
}

func (BaseInfoModel) TableName() string {
	return "base_info"
}

// Save 保存单条数据
func (r *baseInfoRepo) Save(ctx context.Context, info *biz.BaseInfo) error {
	model := r.toModel(info)
	if err := r.data.gormDB.WithContext(ctx).Create(model).Error; err != nil {
		r.log.Errorf("Save base info error: %v", err)
		return err
	}
	info.ID = model.ID
	return nil
}

// Update 更新数据
func (r *baseInfoRepo) Update(ctx context.Context, info *biz.BaseInfo) error {
	model := r.toModel(info)
	if err := r.data.gormDB.WithContext(ctx).Save(model).Error; err != nil {
		r.log.Errorf("Update base info error: %v", err)
		return err
	}
	return nil
}

// FindByCode 根据代码查询（返回最新的一条）
func (r *baseInfoRepo) FindByCode(ctx context.Context, code string) (*biz.BaseInfo, error) {
	var model BaseInfoModel
	err := r.data.gormDB.WithContext(ctx).
		Where("code = ?", code).
		Order("trade_date DESC").
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.Errorf("FindByCode error: %v", err)
		return nil, err
	}

	return r.toBiz(&model), nil
}

// BatchSave 批量保存（使用 GORM upsert）
func (r *baseInfoRepo) BatchSave(ctx context.Context, infos []*biz.BaseInfo) error {
	if len(infos) == 0 {
		return nil
	}

	models := make([]*BaseInfoModel, 0, len(infos))
	for _, info := range infos {
		models = append(models, r.toModel(info))
	}

	// 使用 GORM 的 Clauses 进行 upsert 操作
	err := r.data.gormDB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}, {Name: "trade_date"}}, // 唯一键字段
			DoUpdates: clause.AssignmentColumns([]string{
				"stock_name", "latest_price", "auction_change_rate", "latest_change_rate",
				"auction_unmatched_amount_str", "morning_auction_amount_str", "turnover_str",
				"circulation_market_value", "circulation_market_value_str", "stock_code", "market_code", "turnover",
				"morning_auction_amount", "auction_unmatched_amount", "company_highlights",
				"industry_category", "concept_theme", "consecutive_limit_days", "update_time",
			}), // 冲突时更新的字段
		}).
		CreateInBatches(models, 100).Error

	if err != nil {
		r.log.Errorf("BatchSave error: %v", err)
		return err
	}

	r.log.Infof("BatchSave success, count: %d", len(infos))
	return nil
}

// toModel 转换为数据库模型
func (r *baseInfoRepo) toModel(info *biz.BaseInfo) *BaseInfoModel {
	return &BaseInfoModel{
		ID:                        info.ID,
		StockName:                 info.StockName,
		LatestPrice:               info.LatestPrice,
		AuctionChangeRate:         info.AuctionChangeRate,
		LatestChangeRate:          info.LatestChangeRate,
		AuctionUnmatchedAmountStr: info.AuctionUnmatchedAmountStr,
		MorningAuctionAmountStr:   info.MorningAuctionAmountStr,
		TurnoverStr:               info.TurnoverStr,
		CirculationMarketValue:    info.CirculationMarketValue,
		CirculationMarketValueStr: info.CirculationMarketValueStr,
		StockCode:                 info.StockCode,
		TradeDate:                 info.TradeDate,
		MarketCode:                info.MarketCode,
		Code:                      info.Code,
		Turnover:                  info.Turnover,
		MorningAuctionAmount:      info.MorningAuctionAmount,
		AuctionUnmatchedAmount:    info.AuctionUnmatchedAmount,
		CreateTime:                info.CreateTime,
		UpdateTime:                info.UpdateTime,
		CompanyHighlights:         info.CompanyHighlights,
		IndustryCategory:          info.IndustryCategory,
		ConceptTheme:              info.ConceptTheme,
		ConsecutiveLimitDays:      info.ConsecutiveLimitDays,
	}
}

// toBiz 转换为业务模型
func (r *baseInfoRepo) toBiz(model *BaseInfoModel) *biz.BaseInfo {
	return &biz.BaseInfo{
		ID:                        model.ID,
		StockName:                 model.StockName,
		LatestPrice:               model.LatestPrice,
		AuctionChangeRate:         model.AuctionChangeRate,
		LatestChangeRate:          model.LatestChangeRate,
		AuctionUnmatchedAmountStr: model.AuctionUnmatchedAmountStr,
		MorningAuctionAmountStr:   model.MorningAuctionAmountStr,
		TurnoverStr:               model.TurnoverStr,
		CirculationMarketValue:    model.CirculationMarketValue,
		CirculationMarketValueStr: model.CirculationMarketValueStr,
		StockCode:                 model.StockCode,
		TradeDate:                 model.TradeDate,
		MarketCode:                model.MarketCode,
		Code:                      model.Code,
		Turnover:                  model.Turnover,
		MorningAuctionAmount:      model.MorningAuctionAmount,
		AuctionUnmatchedAmount:    model.AuctionUnmatchedAmount,
		CreateTime:                model.CreateTime,
		UpdateTime:                model.UpdateTime,
		CompanyHighlights:         model.CompanyHighlights,
		IndustryCategory:          model.IndustryCategory,
		ConceptTheme:              model.ConceptTheme,
		ConsecutiveLimitDays:      model.ConsecutiveLimitDays,
	}
}
