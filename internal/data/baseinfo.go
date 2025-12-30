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
	StockCode                 string    `gorm:"column:stock_code;type:varchar(20);index:idx_stock_code"`
	TradeDate                 time.Time `gorm:"column:trade_date;type:date;index:idx_trade_date"`
	MarketCode                string    `gorm:"column:market_code;type:varchar(10)"`
	Code                      string    `gorm:"column:code;type:varchar(20);uniqueIndex:uk_code;index:idx_code"`
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

// Save 保存单条数据（使用 upsert 避免重复）
func (r *baseInfoRepo) Save(ctx context.Context, info *biz.BaseInfo) error {
	model := r.toModel(info)
	// 使用 GORM 的 Clauses 进行 upsert 操作，只基于 code 更新
	err := r.data.gormDB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}}, // 唯一键字段（只基于 code）
			DoUpdates: clause.AssignmentColumns([]string{
				"stock_name", "latest_price", "auction_change_rate", "latest_change_rate",
				"auction_unmatched_amount_str", "morning_auction_amount_str", "turnover_str",
				"circulation_market_value", "circulation_market_value_str", "stock_code", "market_code",
				"trade_date", "turnover", "morning_auction_amount", "auction_unmatched_amount",
				"company_highlights", "industry_category", "concept_theme", "consecutive_limit_days", "update_time",
			}), // 冲突时更新的字段（包括 trade_date）
		}).
		Create(model).Error

	if err != nil {
		r.log.Errorf("Save base info error: %v", err)
		return err
	}
	info.ID = model.ID
	return nil
}

// Update 更新数据（使用 upsert，只基于 code 更新）
func (r *baseInfoRepo) Update(ctx context.Context, info *biz.BaseInfo) error {
	model := r.toModel(info)
	// 使用 upsert 逻辑，只基于 code 进行更新
	err := r.data.gormDB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}}, // 唯一键字段（只基于 code）
			DoUpdates: clause.AssignmentColumns([]string{
				"stock_name", "latest_price", "auction_change_rate", "latest_change_rate",
				"auction_unmatched_amount_str", "morning_auction_amount_str", "turnover_str",
				"circulation_market_value", "circulation_market_value_str", "stock_code", "market_code",
				"trade_date", "turnover", "morning_auction_amount", "auction_unmatched_amount",
				"company_highlights", "industry_category", "concept_theme", "consecutive_limit_days", "update_time",
			}), // 冲突时更新的字段（包括 trade_date）
		}).
		Create(model).Error

	if err != nil {
		r.log.Errorf("Update base info error: %v", err)
		return err
	}
	info.ID = model.ID
	return nil
}

// FindByCode 根据代码查询
func (r *baseInfoRepo) FindByCode(ctx context.Context, code string) (*biz.BaseInfo, error) {
	var model BaseInfoModel
	err := r.data.gormDB.WithContext(ctx).
		Where("code = ?", code).
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

// BatchSave 批量保存（先删除所有数据，再批量插入新数据）
func (r *baseInfoRepo) BatchSave(ctx context.Context, infos []*biz.BaseInfo) error {
	if len(infos) == 0 {
		return nil
	}

	// 开启事务
	tx := r.data.gormDB.WithContext(ctx).Begin()
	if tx.Error != nil {
		r.log.Errorf("BatchSave begin transaction error: %v", tx.Error)
		return tx.Error
	}

	// 先删除所有数据
	if err := tx.Exec("DELETE FROM base_info").Error; err != nil {
		tx.Rollback()
		r.log.Errorf("BatchSave delete all records error: %v", err)
		return err
	}
	r.log.Info("BatchSave deleted all existing records")

	// 批量插入新数据
	models := make([]*BaseInfoModel, 0, len(infos))
	for _, info := range infos {
		models = append(models, r.toModel(info))
	}

	if err := tx.CreateInBatches(models, 100).Error; err != nil {
		tx.Rollback()
		r.log.Errorf("BatchSave create error: %v", err)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		r.log.Errorf("BatchSave commit transaction error: %v", err)
		return err
	}

	r.log.Infof("BatchSave success, total: %d records inserted", len(infos))
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

// BaseInfoDayModel 每日基础数据库模型
type BaseInfoDayModel struct {
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
	StockCode                 string    `gorm:"column:stock_code;type:varchar(20);index:idx_stock_code"`
	TradeDate                 time.Time `gorm:"column:trade_date;type:date;uniqueIndex:uk_code_date;index:idx_trade_date"`
	MarketCode                string    `gorm:"column:market_code;type:varchar(10)"`
	Code                      string    `gorm:"column:code;type:varchar(20);uniqueIndex:uk_code_date;index:idx_code"`
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

func (BaseInfoDayModel) TableName() string {
	return "base_info_day"
}

// BatchSaveDay 批量保存每日数据（根据 code 和 trade_date 查询，存在则更新，不存在则插入）
func (r *baseInfoRepo) BatchSaveDay(ctx context.Context, infos []*biz.BaseInfo) error {
	if len(infos) == 0 {
		return nil
	}

	var toCreate []*BaseInfoDayModel
	var toUpdate []*BaseInfoDayModel

	// 收集所有的 code 和 trade_date 组合
	type CodeDateKey struct {
		Code      string
		TradeDate string // 使用字符串格式 YYYY-MM-DD
	}
	codeToInfo := make(map[CodeDateKey]*biz.BaseInfo)
	var conditions []map[string]interface{}

	for _, info := range infos {
		// 只保留日期部分，去掉时分秒
		dateOnly := info.TradeDate.Format("2006-01-02")
		key := CodeDateKey{Code: info.Code, TradeDate: dateOnly}
		codeToInfo[key] = info
		conditions = append(conditions, map[string]interface{}{
			"code":       info.Code,
			"trade_date": dateOnly,
		})
	}

	// 批量查询已存在的记录
	var existingModels []*BaseInfoDayModel
	if len(conditions) > 0 {
		query := r.data.gormDB.WithContext(ctx)
		for i, cond := range conditions {
			if i == 0 {
				query = query.Where("(code = ? AND trade_date = ?)", cond["code"], cond["trade_date"])
			} else {
				query = query.Or("(code = ? AND trade_date = ?)", cond["code"], cond["trade_date"])
			}
		}
		err := query.Find(&existingModels).Error
		if err != nil {
			r.log.Errorf("BatchSaveDay query error: %v", err)
			return err
		}
	}

	// 建立已存在记录的映射
	existingMap := make(map[CodeDateKey]*BaseInfoDayModel)
	for _, model := range existingModels {
		// 只保留日期部分
		dateOnly := model.TradeDate.Format("2006-01-02")
		key := CodeDateKey{Code: model.Code, TradeDate: dateOnly}
		existingMap[key] = model
	}

	// 分类处理：更新 vs 插入
	for _, info := range infos {
		model := r.toDayModel(info)
		// 只保留日期部分
		dateOnly := info.TradeDate.Format("2006-01-02")
		key := CodeDateKey{Code: info.Code, TradeDate: dateOnly}

		if existing, ok := existingMap[key]; ok {
			// 已存在，准备更新（保留原ID）
			model.ID = existing.ID
			toUpdate = append(toUpdate, model)
		} else {
			// 不存在，准备插入
			toCreate = append(toCreate, model)
		}
	}

	// 批量插入新记录
	if len(toCreate) > 0 {
		if err := r.data.gormDB.WithContext(ctx).CreateInBatches(toCreate, 100).Error; err != nil {
			r.log.Errorf("BatchSaveDay create error: %v", err)
			return err
		}
		r.log.Infof("BatchSaveDay created %d new records", len(toCreate))
	}

	// 批量更新已存在的记录
	if len(toUpdate) > 0 {
		for _, model := range toUpdate {
			if err := r.data.gormDB.WithContext(ctx).
				Model(&BaseInfoDayModel{}).
				Where("id = ?", model.ID).
				Updates(model).Error; err != nil {
				r.log.Errorf("BatchSaveDay update error for code %s date %s: %v",
					model.Code, model.TradeDate.Format("2006-01-02"), err)
				return err
			}
		}
		r.log.Infof("BatchSaveDay updated %d existing records", len(toUpdate))
	}

	r.log.Infof("BatchSaveDay success, total: %d (created: %d, updated: %d)", len(infos), len(toCreate), len(toUpdate))
	return nil
}

// toDayModel 转换为每日数据库模型
func (r *baseInfoRepo) toDayModel(info *biz.BaseInfo) *BaseInfoDayModel {
	return &BaseInfoDayModel{
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
