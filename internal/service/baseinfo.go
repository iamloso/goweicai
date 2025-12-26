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
	if err != nil {
		return fmt.Errorf("查询失败: %w", err)
	}

	// 根据配置决定是否打印调试信息（控制输出长度）
	if s.bootstrap.Debug {
		resultStr := fmt.Sprintf("%v", result)
		maxLen := 500 // 最多显示 500 字符
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

	s.log.Infof("查询到 %d 条基础数据", len(infos))

	// 保存到 base_info 表（原逻辑不变，只根据 code 更新）
	if err := s.uc.SaveBaseInfos(ctx, infos); err != nil {
		return fmt.Errorf("保存数据失败: %w", err)
	}
	s.log.Info("base_info 表数据保存成功")

	// 保存到 base_info_day 表（根据 code 和 trade_date 更新）
	if err := s.uc.SaveBaseInfosDay(ctx, infos); err != nil {
		return fmt.Errorf("保存每日数据失败: %w", err)
	}
	s.log.Info("base_info_day 表数据保存成功")

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
