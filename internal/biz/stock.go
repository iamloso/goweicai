package biz

import (
	"context"
	"time"
)

// Stock 股票数据模型
type Stock struct {
	ID                        int64     `json:"id"`
	Code                      string    `json:"code"`
	MarketCode                string    `json:"market_code"`
	StockCode                 string    `json:"stock_code"`
	StockName                 string    `json:"stock_name"`
	LatestPrice               float64   `json:"latest_price"`
	LatestChangeRate          float64   `json:"latest_change_rate"`
	AuctionUnmatchedAmount    int64     `json:"auction_unmatched_amount"`
	AuctionUnmatchedAmountStr string    `json:"auction_unmatched_amount_str"`
	AuctionUnmatchedAmountRank string   `json:"auction_unmatched_amount_rank"`
	AuctionUnmatchedAmountRankNum int   `json:"auction_unmatched_amount_rank_num"`
	AuctionChangeRate         float64   `json:"auction_change_rate"`
	MorningAuctionAmount      int64     `json:"morning_auction_amount"`
	MorningAuctionAmountStr   string    `json:"morning_auction_amount_str"`
	Turnover                  float64   `json:"turnover"`
	TurnoverStr               string    `json:"turnover_str"`
	CirculationMarketValue    float64   `json:"circulation_market_value"`
	LimitUpReason             string    `json:"limit_up_reason"`
	CompanyHighlights         string    `json:"company_highlights"`
	IndustryCategory          string    `json:"industry_category"`
	ConceptTheme              string    `json:"concept_theme"`
	LimitUpSealAmount         int64     `json:"limit_up_seal_amount"`
	LimitUpSealAmountStr      string    `json:"limit_up_seal_amount_str"`
	ConsecutiveLimitDays      int       `json:"consecutive_limit_days"`
	TradeDate                 time.Time `json:"trade_date"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// StockRepo 股票数据仓库接口
type StockRepo interface {
	Save(context.Context, *Stock) error
	BatchSave(context.Context, []*Stock) error
	FindByCodeAndDate(context.Context, string, time.Time) (*Stock, error)
	Update(context.Context, *Stock) error
}

// StockUsecase 股票业务用例
type StockUsecase struct {
	repo StockRepo
}

// NewStockUsecase 创建股票业务用例
func NewStockUsecase(repo StockRepo) *StockUsecase {
	return &StockUsecase{repo: repo}
}

// SaveStocks 批量保存股票数据
func (uc *StockUsecase) SaveStocks(ctx context.Context, stocks []*Stock) error {
	return uc.repo.BatchSave(ctx, stocks)
}
