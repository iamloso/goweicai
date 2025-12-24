package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	gowencai "github.com/iamloso/goweicai"
	"github.com/iamloso/goweicai/internal/biz"
	"github.com/iamloso/goweicai/internal/conf"
	"github.com/iamloso/goweicai/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

// WencaiService 问财服务
type WencaiService struct {
	uc     *biz.StockUsecase
	config *conf.Wencai
	log    *log.Helper
}

// NewWencaiService 创建问财服务
func NewWencaiService(uc *biz.StockUsecase, c *conf.Wencai, logger log.Logger) *WencaiService {
	return &WencaiService{
		uc:     uc,
		config: c,
		log:    log.NewHelper(log.With(logger, "module", "service/wencai")),
	}
}

// FetchAndSaveStocks 获取并保存股票数据
func (s *WencaiService) FetchAndSaveStocks(ctx context.Context) error {
	s.log.Info("开始查询股票数据...")

	result, err := gowencai.Get(&gowencai.QueryOptions{
		Query:  s.config.Query,
		Cookie: s.config.Cookie,
		Log:    true,
		Loop:   true,
	})
	if err != nil {
		return fmt.Errorf("查询失败: %w", err)
	}

	stocks, err := s.parseResult(result)
	if err != nil {
		return fmt.Errorf("解析结果失败: %w", err)
	}

	s.log.Infof("查询到 %d 条股票数据", len(stocks))

	if err := s.uc.SaveStocks(ctx, stocks); err != nil {
		return fmt.Errorf("保存数据失败: %w", err)
	}

	s.log.Info("数据保存成功")
	return nil
}

// parseResult 解析查询结果
func (s *WencaiService) parseResult(result interface{}) ([]*biz.Stock, error) {
	var items []map[string]interface{}

	switch v := result.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				items = append(items, m)
			}
		}
	case []map[string]interface{}:
		items = v
	default:
		return nil, fmt.Errorf("不支持的数据类型")
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("没有数据")
	}

	tradeDate := time.Now()
	tradeDateKey := tradeDate.Format("20060102")

	stocks := make([]*biz.Stock, 0, len(items))

	for _, item := range items {
		code := getStringValue(item, "code")
		if code == "" {
			continue
		}

		// 准备数据
		marketCode := getStringValue(item, "market_code")
		latestPrice := getFloatValue(item, "最新价")

		// 涨跌幅字段
		changeRateKey := fmt.Sprintf("涨跌幅:前复权[%s]", tradeDateKey)
		latestChangeRate := getFloatValue(item, changeRateKey)

		// 动态构建字段名
		auctionAmountKey := fmt.Sprintf("竞价未匹配金额[%s]", tradeDateKey)
		auctionAmountRankKey := fmt.Sprintf("竞价未匹配金额排名[%s]", tradeDateKey)
		auctionChangeRateKey := fmt.Sprintf("竞价涨幅[%s]", tradeDateKey)

		auctionAmount := getIntValue(item, auctionAmountKey)
		auctionAmountStr := data.FormatAmount(auctionAmount)
		auctionAmountRank := getStringValue(item, auctionAmountRankKey)
		auctionAmountRankNum := data.ParseRankNumber(auctionAmountRank)
		auctionChangeRate := getFloatValue(item, auctionChangeRateKey)

		// 竞价金额、成交金额、流通市值
		morningAuctionAmountKey := fmt.Sprintf("竞价金额[%s]", tradeDateKey)
		morningAuctionAmount := getIntValue(item, morningAuctionAmountKey)
		morningAuctionAmountStr := data.FormatAmount(morningAuctionAmount)

		turnoverKey := fmt.Sprintf("成交额[%s]", tradeDateKey)
		turnoverInt := getIntValue(item, turnoverKey)
		turnover := float64(turnoverInt)
		turnoverStr := data.FormatAmount(turnoverInt)

		circulationMarketValueKey := fmt.Sprintf("a股市值(不含限售股)[%s]", tradeDateKey)
		circulationMarketValue := getFloatValue(item, circulationMarketValueKey)

		stockCode := getStringValue(item, "股票代码")
		stockName := getStringValue(item, "股票简称")

		// 新增字段
		limitUpReasonKey := fmt.Sprintf("涨停原因类别[%s]", tradeDateKey)
		limitUpReason := getStringValue(item, limitUpReasonKey)

		companyHighlights := getStringValue(item, "公司亮点")
		industryCategory := getStringValue(item, "所属同花顺行业")
		conceptTheme := getStringValue(item, "所属概念")

		limitUpSealAmountKey := fmt.Sprintf("涨停封单额[%s]", tradeDateKey)
		limitUpSealAmount := getIntValue(item, limitUpSealAmountKey)
		limitUpSealAmountStr := data.FormatAmount(limitUpSealAmount)

		consecutiveLimitDaysKey := fmt.Sprintf("连续涨停天数[%s]", tradeDateKey)
		consecutiveLimitDays := int(getIntValue(item, consecutiveLimitDaysKey))

		stock := &biz.Stock{
			Code:                      code,
			MarketCode:                marketCode,
			StockCode:                 stockCode,
			StockName:                 stockName,
			LatestPrice:               latestPrice,
			LatestChangeRate:          latestChangeRate,
			AuctionUnmatchedAmount:    auctionAmount,
			AuctionUnmatchedAmountStr: auctionAmountStr,
			AuctionUnmatchedAmountRank: auctionAmountRank,
			AuctionUnmatchedAmountRankNum: auctionAmountRankNum,
			AuctionChangeRate:         auctionChangeRate,
			MorningAuctionAmount:      morningAuctionAmount,
			MorningAuctionAmountStr:   morningAuctionAmountStr,
			Turnover:                  turnover,
			TurnoverStr:               turnoverStr,
			CirculationMarketValue:    circulationMarketValue,
			LimitUpReason:             limitUpReason,
			CompanyHighlights:         companyHighlights,
			IndustryCategory:          industryCategory,
			ConceptTheme:              conceptTheme,
			LimitUpSealAmount:         limitUpSealAmount,
			LimitUpSealAmountStr:      limitUpSealAmountStr,
			ConsecutiveLimitDays:      consecutiveLimitDays,
			TradeDate:                 tradeDate,
		}

		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// 辅助函数
func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getFloatValue(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok && v != nil {
		switch val := v.(type) {
		case float64:
			return val
		case float32:
			return float64(val)
		case int:
			return float64(val)
		case int64:
			return float64(val)
		case string:
			var f float64
			fmt.Sscanf(val, "%f", &f)
			return f
		}
	}
	return 0
}

func getIntValue(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok && v != nil {
		switch val := v.(type) {
		case int64:
			return val
		case int:
			return int64(val)
		case float64:
			return int64(val)
		case float32:
			return int64(val)
		case string:
			// 处理科学计数法
			var f float64
			if strings.Contains(val, "E") || strings.Contains(val, "e") {
				fmt.Sscanf(val, "%e", &f)
				return int64(f)
			}
			var i int64
			fmt.Sscanf(val, "%d", &i)
			return i
		}
	}
	return 0
}
