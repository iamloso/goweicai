package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/iamloso/goweicai/internal/biz"
)

type ztInfoRepo struct {
	data *Data
	log  *log.Helper
}

// NewZtInfoRepo 创建涨停数据仓库
func NewZtInfoRepo(data *Data, logger log.Logger) biz.ZtInfoRepo {
	return &ztInfoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/ztinfo")),
	}
}

// ZtDayModel 涨停每日数据库模型
type ZtDayModel struct {
	ID                        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	StockName                 string    `gorm:"column:stock_name;type:varchar(100)"`
	LatestPrice               float64   `gorm:"column:latest_price;type:decimal(10,2)"`
	AuctionChangeRate         float64   `gorm:"column:auction_change_rate;type:decimal(10,3)"`
	LatestChangeRate          float64   `gorm:"column:latest_change_rate;type:decimal(10,2)"`
	AuctionUnmatchedAmountStr string    `gorm:"column:auction_unmatched_amount_str;type:varchar(50)"`
	MorningAuctionAmountStr   string    `gorm:"column:morning_auction_amount_str;type:varchar(50)"`
	TurnoverStr               string    `gorm:"column:turnover_str;type:varchar(50)"`
	CirculationMarketValue    float64   `gorm:"column:circulation_market_value;type:decimal(20,2)"`
	StockCode                 string    `gorm:"column:stock_code;type:varchar(20);index:idx_stock_code"`
	TradeDate                 time.Time `gorm:"column:trade_date;type:date;uniqueIndex:uk_code_date;index:idx_trade_date"`
	MarketCode                string    `gorm:"column:market_code;type:varchar(10)"`
	Code                      string    `gorm:"column:code;type:varchar(20);uniqueIndex:uk_code_date;index:idx_code"`
	Turnover                  float64   `gorm:"column:turnover;type:decimal(20,2)"`
	MorningAuctionAmount      int64     `gorm:"column:morning_auction_amount;type:bigint"`
	AuctionUnmatchedAmount    int64     `gorm:"column:auction_unmatched_amount;type:bigint"`
	CreateTime                time.Time `gorm:"column:create_time;type:datetime;autoCreateTime"`
	UpdateTime                time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime"`
	LimitUpReason             string    `gorm:"column:limit_up_reason;type:text"`
	CompanyHighlights         string    `gorm:"column:company_highlights;type:text"`
	IndustryCategory          string    `gorm:"column:industry_category;type:varchar(100)"`
	ConceptTheme              string    `gorm:"column:concept_theme;type:text"`
	LimitUpSealAmount         int64     `gorm:"column:limit_up_seal_amount;type:bigint"`
	LimitUpSealAmountStr      string    `gorm:"column:limit_up_seal_amount_str;type:varchar(50)"`
	ConsecutiveLimitDays      int       `gorm:"column:consecutive_limit_days;type:int"`
}

func (ZtDayModel) TableName() string {
	return "zt_day"
}

// BatchSaveDay 批量保存每日数据（根据 code 和 trade_date 查询，存在则更新，不存在则插入）
func (r *ztInfoRepo) BatchSaveDay(ctx context.Context, infos []*biz.ZtInfo) error {
	if len(infos) == 0 {
		return nil
	}

	var toCreate []*ZtDayModel
	var toUpdate []*ZtDayModel

	// 收集所有的 code 和 trade_date 组合
	type CodeDateKey struct {
		Code      string
		TradeDate string // 使用字符串格式 YYYY-MM-DD
	}
	codeToInfo := make(map[CodeDateKey]*biz.ZtInfo)
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
	var existingModels []*ZtDayModel
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
	existingMap := make(map[CodeDateKey]*ZtDayModel)
	for _, model := range existingModels {
		// 只保留日期部分
		dateOnly := model.TradeDate.Format("2006-01-02")
		key := CodeDateKey{Code: model.Code, TradeDate: dateOnly}
		existingMap[key] = model
	}

	// 分类处理：更新 vs 插入
	for _, info := range infos {
		model := r.toModel(info)
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
				Model(&ZtDayModel{}).
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

// toModel 转换为数据库模型
func (r *ztInfoRepo) toModel(info *biz.ZtInfo) *ZtDayModel {
	return &ZtDayModel{
		ID:                        info.ID,
		StockName:                 info.StockName,
		LatestPrice:               info.LatestPrice,
		AuctionChangeRate:         info.AuctionChangeRate,
		LatestChangeRate:          info.LatestChangeRate,
		AuctionUnmatchedAmountStr: info.AuctionUnmatchedAmountStr,
		MorningAuctionAmountStr:   info.MorningAuctionAmountStr,
		TurnoverStr:               info.TurnoverStr,
		CirculationMarketValue:    info.CirculationMarketValue,
		StockCode:                 info.StockCode,
		TradeDate:                 info.TradeDate,
		MarketCode:                info.MarketCode,
		Code:                      info.Code,
		Turnover:                  info.Turnover,
		MorningAuctionAmount:      info.MorningAuctionAmount,
		AuctionUnmatchedAmount:    info.AuctionUnmatchedAmount,
		CreateTime:                info.CreateTime,
		UpdateTime:                info.UpdateTime,
		LimitUpReason:             info.LimitUpReason,
		CompanyHighlights:         info.CompanyHighlights,
		IndustryCategory:          info.IndustryCategory,
		ConceptTheme:              info.ConceptTheme,
		LimitUpSealAmount:         info.LimitUpSealAmount,
		LimitUpSealAmountStr:      info.LimitUpSealAmountStr,
		ConsecutiveLimitDays:      info.ConsecutiveLimitDays,
	}
}
