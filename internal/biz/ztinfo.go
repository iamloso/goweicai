package biz

import (
	"context"
	"time"
)

// ZtInfo 涨停数据模型
type ZtInfo struct {
	ID                        int64     `json:"id"`
	StockName                 string    `json:"stock_name"`
	LatestPrice               float64   `json:"latest_price"`
	AuctionChangeRate         float64   `json:"auction_change_rate"`
	LatestChangeRate          float64   `json:"latest_change_rate"`
	AuctionUnmatchedAmountStr string    `json:"auction_unmatched_amount_str"` // 带单位字符串
	MorningAuctionAmountStr   string    `json:"morning_auction_amount_str"`   // 带单位字符串
	TurnoverStr               string    `json:"turnover_str"`                 // 带单位字符串
	CirculationMarketValue    float64   `json:"circulation_market_value"`
	StockCode                 string    `json:"stock_code"`
	TradeDate                 time.Time `json:"trade_date"`
	MarketCode                string    `json:"market_code"`
	Code                      string    `json:"code"`
	Turnover                  float64   `json:"turnover"`
	MorningAuctionAmount      int64     `json:"morning_auction_amount"`
	AuctionUnmatchedAmount    int64     `json:"auction_unmatched_amount"`
	CreateTime                time.Time `json:"create_time"`
	UpdateTime                time.Time `json:"update_time"`
	LimitUpReason             string    `json:"limit_up_reason"` // 涨停原因
	CompanyHighlights         string    `json:"company_highlights"`
	IndustryCategory          string    `json:"industry_category"`
	ConceptTheme              string    `json:"concept_theme"`
	LimitUpSealAmount         int64     `json:"limit_up_seal_amount"`     // 涨停封单金额(元)
	LimitUpSealAmountStr      string    `json:"limit_up_seal_amount_str"` // 涨停封单金额(万/亿)
	ConsecutiveLimitDays      int       `json:"consecutive_limit_days"`
}

// ZtInfoRepo 涨停数据仓库接口
type ZtInfoRepo interface {
	// BatchSaveDay 批量保存每日数据
	BatchSaveDay(ctx context.Context, infos []*ZtInfo) error
}

// ZtInfoUsecase 涨停数据用例
type ZtInfoUsecase struct {
	repo ZtInfoRepo
}

// NewZtInfoUsecase 创建涨停数据用例
func NewZtInfoUsecase(repo ZtInfoRepo) *ZtInfoUsecase {
	return &ZtInfoUsecase{
		repo: repo,
	}
}

// SaveZtInfosDay 保存每日涨停数据
func (uc *ZtInfoUsecase) SaveZtInfosDay(ctx context.Context, infos []*ZtInfo) error {
	return uc.repo.BatchSaveDay(ctx, infos)
}
