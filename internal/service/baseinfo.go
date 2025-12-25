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

// BaseInfoService 基础数据服务
type BaseInfoService struct {
	uc        *biz.BaseInfoUsecase
	config    *conf.BaseInfo
	bootstrap *conf.Bootstrap
	log       *log.Helper
}

// NewBaseInfoService 创建基础数据服务
func NewBaseInfoService(uc *biz.BaseInfoUsecase, c *conf.BaseInfo, bc *conf.Bootstrap, logger log.Logger) *BaseInfoService {
	return &BaseInfoService{
		uc:        uc,
		config:    c,
		bootstrap: bc,
		log:       log.NewHelper(log.With(logger, "module", "service/baseinfo")),
	}
}

// FetchAndSaveBaseInfo 获取并保存基础数据
func (s *BaseInfoService) FetchAndSaveBaseInfo(ctx context.Context) error {
	s.log.Info("开始查询基础数据...")

	result, err := gowencai.Get(&gowencai.QueryOptions{
		Query:  s.config.Query,
		Cookie: s.config.Cookie,
		Log:    true,
		Loop:   true,
	})
	// 根据配置决定是否打印调试信息
	if s.bootstrap.Debug {
		s.log.Infof("查询结果: %v", result)
	}
	if err != nil {
		return fmt.Errorf("查询失败: %w", err)
	}

	infos, err := s.parseResult(result)
	if err != nil {
		return fmt.Errorf("解析结果失败: %w", err)
	}

	s.log.Infof("查询到 %d 条基础数据", len(infos))

	if err := s.uc.SaveBaseInfos(ctx, infos); err != nil {
		return fmt.Errorf("保存数据失败: %w", err)
	}

	s.log.Info("基础数据保存成功")
	return nil
}

// parseResult 解析查询结果
func (s *BaseInfoService) parseResult(result interface{}) ([]*biz.BaseInfo, error) {
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

	infos := make([]*biz.BaseInfo, 0, len(items))

	for _, item := range items {
		code := getStringValue(item, "code")
		if code == "" {
			continue
		}

		// 解析数据 - 根据实际返回字段映射
		info := &biz.BaseInfo{
			StockName:                 getStringValue(item, "股票简称"),
			LatestPrice:               getFloatValue(item, "最新价"),
			AuctionChangeRate:         getFloatValue(item, fmt.Sprintf("竞价涨幅[%s]", tradeDateKey)),
			LatestChangeRate:          getFloatValue(item, fmt.Sprintf("涨跌幅:前复权[%s]", tradeDateKey)),
			AuctionUnmatchedAmountStr: formatAmountStr(item, fmt.Sprintf("竞价未匹配金额[%s]", tradeDateKey)),
			MorningAuctionAmountStr:   formatAmountStr(item, fmt.Sprintf("竞价金额[%s]", tradeDateKey)),
			TurnoverStr:               formatAmountStr(item, fmt.Sprintf("成交额[%s]", tradeDateKey)),
			CirculationMarketValue:    getFloatValue(item, fmt.Sprintf("a股市值(不含限售股)[%s]", tradeDateKey)),
			CirculationMarketValueStr: formatAmountStr(item, fmt.Sprintf("a股市值(不含限售股)[%s]", tradeDateKey)),
			StockCode:                 getStringValue(item, "股票代码"),
			TradeDate:                 tradeDate,
			MarketCode:                getStringValue(item, "market_code"),
			Code:                      code,
			Turnover:                  getFloatValue(item, fmt.Sprintf("成交额[%s]", tradeDateKey)),
			MorningAuctionAmount:      getIntValue(item, fmt.Sprintf("竞价量[%s]", tradeDateKey)),
			AuctionUnmatchedAmount:    getIntValue(item, fmt.Sprintf("竞价未匹配量[%s]", tradeDateKey)),
			CompanyHighlights:         getStringValue(item, "公司亮点"),
			IndustryCategory:          getStringValue(item, "所属同花顺行业"),
			ConceptTheme:              getStringValue(item, "所属概念"),
			ConsecutiveLimitDays:      int(getIntValue(item, fmt.Sprintf("连续涨停天数[%s]", tradeDateKey))),
		}

		infos = append(infos, info)
	}

	return infos, nil
}
