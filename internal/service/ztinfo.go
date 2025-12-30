package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	gowencai "github.com/iamloso/goweicai"
	"github.com/iamloso/goweicai/internal/biz"
	"github.com/iamloso/goweicai/internal/conf"
)

// ZtInfoService 涨停数据服务
type ZtInfoService struct {
	uc        *biz.ZtInfoUsecase
	config    *conf.ZtInfo
	bootstrap *conf.Bootstrap
	log       *log.Helper
}

// NewZtInfoService 创建涨停数据服务
func NewZtInfoService(uc *biz.ZtInfoUsecase, c *conf.ZtInfo, bc *conf.Bootstrap, logger log.Logger) *ZtInfoService {
	return &ZtInfoService{
		uc:        uc,
		config:    c,
		bootstrap: bc,
		log:       log.NewHelper(log.With(logger, "module", "service/ztinfo")),
	}
}

// FetchAndSaveZtInfo 获取并保存涨停数据
func (s *ZtInfoService) FetchAndSaveZtInfo(ctx context.Context) error {
	s.log.Info("开始查询涨停数据...")

	result, err := gowencai.Get(&gowencai.QueryOptions{
		Query:  s.config.Query,
		Cookie: s.config.Cookie,
		Log:    true,
		Loop:   true,
	})
	if err != nil {
		return fmt.Errorf("查询失败: %w", err)
	}

	// 根据配置决定是否打印调试信息
	if s.bootstrap.Debug {
		resultStr := fmt.Sprintf("%v", result)
		maxLen := 5000
		if len(resultStr) > maxLen {
			s.log.Infof("查询结果（截断）: %s... [总长度: %d]", resultStr[:maxLen], len(resultStr))
		} else {
			s.log.Infof("查询结果: %s", resultStr)
		}
	}

	infos, err := s.parseResult(result)
	if err != nil {
		return fmt.Errorf("解析结果失败: %w", err)
	}

	s.log.Infof("查询到 %d 条涨停数据", len(infos))

	// 保存到 zt_day 表
	if err := s.uc.SaveZtInfosDay(ctx, infos); err != nil {
		return fmt.Errorf("保存涨停数据失败: %w", err)
	}
	s.log.Info("zt_day 表数据保存成功")

	return nil
}

// parseResult 解析查询结果
func (s *ZtInfoService) parseResult(result interface{}) ([]*biz.ZtInfo, error) {
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

	infos := make([]*biz.ZtInfo, 0, len(items))

	for i, item := range items {
		code := getStringValue(item, "code")
		if code == "" {
			continue
		}

		// 打印第一条记录的所有字段名，用于调试
		if i == 0 {
			s.log.Info("第一条记录的字段名:")
			for key := range item {
				s.log.Infof("  字段: %s", key)
			}
		}

		// 解析数据
		info := &biz.ZtInfo{
			StockName:                 getStringValue(item, "股票简称"),
			LatestPrice:               getFloatValue(item, "最新价"),
			AuctionChangeRate:         getFloatValue(item, fmt.Sprintf("竞价涨幅[%s]", tradeDateKey)),
			LatestChangeRate:          getFloatValue(item, fmt.Sprintf("涨跌幅:前复权[%s]", tradeDateKey)),
			AuctionUnmatchedAmountStr: formatAmountStr(item, fmt.Sprintf("竞价未匹配金额[%s]", tradeDateKey)),
			MorningAuctionAmountStr:   formatAmountStr(item, fmt.Sprintf("竞价金额[%s]", tradeDateKey)),
			TurnoverStr:               formatAmountStr(item, fmt.Sprintf("成交额[%s]", tradeDateKey)),
			CirculationMarketValue:    getFloatValue(item, fmt.Sprintf("a股市值(不含限售股)[%s]", tradeDateKey)),
			StockCode:                 getStringValue(item, "股票代码"),
			TradeDate:                 tradeDate,
			MarketCode:                getStringValue(item, "market_code"),
			Code:                      code,
			Turnover:                  getFloatValue(item, fmt.Sprintf("成交额[%s]", tradeDateKey)),
			MorningAuctionAmount:      int64(getFloatValue(item, fmt.Sprintf("竞价金额[%s]", tradeDateKey))),
			AuctionUnmatchedAmount:    int64(getFloatValue(item, fmt.Sprintf("竞价未匹配金额[%s]", tradeDateKey))),
			LimitUpReason:             getStringValue(item, fmt.Sprintf("涨停原因类别[%s]", tradeDateKey)),
			CompanyHighlights:         getStringValue(item, "公司亮点"),
			IndustryCategory:          getStringValue(item, "所属同花顺行业"),
			ConceptTheme:              getStringValue(item, "所属概念"),
			ConsecutiveLimitDays:      int(getIntValue(item, fmt.Sprintf("连续涨停天数[%s]", tradeDateKey))),
		}

		// 优先取涨停封单额，如果没有则取跌停封单额
		limitUpSealKey := fmt.Sprintf("涨停封单额[%s]", tradeDateKey)
		limitDownSealKey := fmt.Sprintf("跌停封单额[%s]", tradeDateKey)

		sealAmount := getFloatValue(item, limitUpSealKey)
		sealAmountStr := formatAmountStr(item, limitUpSealKey)
		if sealAmount == 0 {
			sealAmount = getFloatValue(item, limitDownSealKey)
			sealAmountStr = formatAmountStr(item, limitDownSealKey)
		}

		info.LimitUpSealAmount = int64(sealAmount)
		info.LimitUpSealAmountStr = sealAmountStr

		infos = append(infos, info)
	}

	return infos, nil
}
