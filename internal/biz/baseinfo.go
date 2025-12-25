package biz

import (
	"context"
	"time"
)

// BaseInfo 基础数据模型
type BaseInfo struct {
	ID                        int64     `json:"id"`
	StockName                 string    `json:"stock_name"`
	LatestPrice               float64   `json:"latest_price"`
	AuctionChangeRate         float64   `json:"auction_change_rate"`
	LatestChangeRate          float64   `json:"latest_change_rate"`
	AuctionUnmatchedAmountStr string    `json:"auction_unmatched_amount_str"` // 带单位字符串，如：12万、1.3亿
	MorningAuctionAmountStr   string    `json:"morning_auction_amount_str"`   // 带单位字符串，如：12万、1.3亿
	TurnoverStr               string    `json:"turnover_str"`                 // 带单位字符串，如：12万、1.3亿
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
	CompanyHighlights         string    `json:"company_highlights"`
	IndustryCategory          string    `json:"industry_category"`
	ConceptTheme              string    `json:"concept_theme"`
	ConsecutiveLimitDays      int       `json:"consecutive_limit_days"`
}

// BaseInfoRepo 基础数据仓库接口
type BaseInfoRepo interface {
	// Save 保存基础数据
	Save(ctx context.Context, info *BaseInfo) error
	// Update 更新基础数据
	Update(ctx context.Context, info *BaseInfo) error
	// FindByCode 根据代码查询
	FindByCode(ctx context.Context, code string) (*BaseInfo, error)
	// BatchSave 批量保存
	BatchSave(ctx context.Context, infos []*BaseInfo) error
}

// BaseInfoUsecase 基础数据用例
type BaseInfoUsecase struct {
	repo BaseInfoRepo
}

// NewBaseInfoUsecase 创建基础数据用例
func NewBaseInfoUsecase(repo BaseInfoRepo) *BaseInfoUsecase {
	return &BaseInfoUsecase{
		repo: repo,
	}
}

// SaveBaseInfos 保存基础数据
func (uc *BaseInfoUsecase) SaveBaseInfos(ctx context.Context, infos []*BaseInfo) error {
	return uc.repo.BatchSave(ctx, infos)
}
