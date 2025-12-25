package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/iamloso/goweicai/internal/biz"
	"github.com/iamloso/goweicai/internal/conf"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data 数据访问层
type Data struct {
	db     *sql.DB
	gormDB *gorm.DB
	log    *log.Helper
}

// NewData 创建数据访问层
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(log.With(logger, "module", "data"))
	
	// 原生 SQL DB（用于 stock）
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}
	
	// GORM DB（用于 baseinfo）
	gormDB, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	
	helper.Info("database connected successfully")
	
	cleanup := func() {
		helper.Info("closing database connection")
		db.Close()
		sqlDB, _ := gormDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
	
	return &Data{db: db, gormDB: gormDB, log: helper}, cleanup, nil
}

// stockRepo 股票数据仓库实现
type stockRepo struct {
	data *Data
	log  *log.Helper
}

// NewStockRepo 创建股票数据仓库
func NewStockRepo(data *Data, logger log.Logger) biz.StockRepo {
	return &stockRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/stock")),
	}
}

// Save 保存单条股票数据
func (r *stockRepo) Save(ctx context.Context, stock *biz.Stock) error {
	insertSQL := `INSERT INTO zp_jj (
		code, market_code, latest_price, latest_change_rate,
		auction_unmatched_amount, auction_unmatched_amount_str, auction_unmatched_amount_rank,
		auction_unmatched_amount_rank_num,
		auction_change_rate, morning_auction_amount, morning_auction_amount_str, 
		turnover, turnover_str, circulation_market_value,
		stock_code, stock_name, 
		limit_up_reason, company_highlights, industry_category, concept_theme,
		limit_up_seal_amount, limit_up_seal_amount_str, consecutive_limit_days,
		trade_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.data.db.ExecContext(ctx, insertSQL,
		stock.Code, stock.MarketCode, stock.LatestPrice, stock.LatestChangeRate,
		stock.AuctionUnmatchedAmount, stock.AuctionUnmatchedAmountStr, stock.AuctionUnmatchedAmountRank,
		stock.AuctionUnmatchedAmountRankNum,
		stock.AuctionChangeRate, stock.MorningAuctionAmount, stock.MorningAuctionAmountStr,
		stock.Turnover, stock.TurnoverStr, stock.CirculationMarketValue,
		stock.StockCode, stock.StockName,
		stock.LimitUpReason, stock.CompanyHighlights, stock.IndustryCategory, stock.ConceptTheme,
		stock.LimitUpSealAmount, stock.LimitUpSealAmountStr, stock.ConsecutiveLimitDays,
		stock.TradeDate,
	)
	
	return err
}

// Update 更新股票数据
func (r *stockRepo) Update(ctx context.Context, stock *biz.Stock) error {
	updateSQL := `UPDATE zp_jj SET
		market_code = ?,
		latest_price = ?,
		latest_change_rate = ?,
		auction_unmatched_amount = ?,
		auction_unmatched_amount_str = ?,
		auction_unmatched_amount_rank = ?,
		auction_unmatched_amount_rank_num = ?,
		auction_change_rate = ?,
		morning_auction_amount = ?,
		morning_auction_amount_str = ?,
		turnover = ?,
		turnover_str = ?,
		circulation_market_value = ?,
		stock_name = ?,
		limit_up_reason = ?,
		company_highlights = ?,
		industry_category = ?,
		concept_theme = ?,
		limit_up_seal_amount = ?,
		limit_up_seal_amount_str = ?,
		consecutive_limit_days = ?
	WHERE code = ? AND trade_date = ?`

	_, err := r.data.db.ExecContext(ctx, updateSQL,
		stock.MarketCode, stock.LatestPrice, stock.LatestChangeRate,
		stock.AuctionUnmatchedAmount, stock.AuctionUnmatchedAmountStr, stock.AuctionUnmatchedAmountRank,
		stock.AuctionUnmatchedAmountRankNum,
		stock.AuctionChangeRate, stock.MorningAuctionAmount, stock.MorningAuctionAmountStr,
		stock.Turnover, stock.TurnoverStr, stock.CirculationMarketValue,
		stock.StockName,
		stock.LimitUpReason, stock.CompanyHighlights, stock.IndustryCategory, stock.ConceptTheme,
		stock.LimitUpSealAmount, stock.LimitUpSealAmountStr, stock.ConsecutiveLimitDays,
		stock.Code, stock.TradeDate,
	)
	
	return err
}

// FindByCodeAndDate 根据代码和日期查找股票
func (r *stockRepo) FindByCodeAndDate(ctx context.Context, code string, date time.Time) (*biz.Stock, error) {
	query := `SELECT id FROM zp_jj WHERE code = ? AND trade_date = ?`
	
	var stock biz.Stock
	err := r.data.db.QueryRowContext(ctx, query, code, date).Scan(&stock.ID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &stock, nil
}

// BatchSave 批量保存股票数据
func (r *stockRepo) BatchSave(ctx context.Context, stocks []*biz.Stock) error {
	tx, err := r.data.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	insertCount := 0
	updateCount := 0

	for _, stock := range stocks {
		existing, err := r.FindByCodeAndDate(ctx, stock.Code, stock.TradeDate)
		if err != nil {
			r.log.Errorf("查询记录失败 [%s]: %v", stock.Code, err)
			continue
		}

		if existing == nil {
			// 插入新记录
			if err := r.Save(ctx, stock); err != nil {
				r.log.Errorf("插入数据失败 [%s]: %v", stock.Code, err)
				continue
			}
			insertCount++
			r.log.Infof("插入新记录: %s - %s (连板: %d天)", stock.Code, stock.StockName, stock.ConsecutiveLimitDays)
		} else {
			// 更新现有记录
			if err := r.Update(ctx, stock); err != nil {
				r.log.Errorf("更新数据失败 [%s]: %v", stock.Code, err)
				continue
			}
			updateCount++
			r.log.Infof("更新记录: %s - %s (连板: %d天)", stock.Code, stock.StockName, stock.ConsecutiveLimitDays)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	r.log.Infof("操作完成 - 新插入: %d 条, 更新: %d 条", insertCount, updateCount)
	return nil
}

// FormatAmount 将金额转换为带单位的字符串（万或亿）
func FormatAmount(amount int64) string {
	if amount == 0 {
		return ""
	}

	fAmount := float64(amount)

	// 1亿 = 100000000
	if amount >= 100000000 {
		yi := fAmount / 100000000
		return fmt.Sprintf("%.2f亿", yi)
	}

	// 1万 = 10000
	if amount >= 10000 {
		wan := fAmount / 10000
		return fmt.Sprintf("%.2f万", wan)
	}

	return fmt.Sprintf("%d", amount)
}

// ParseRankNumber 从排名字符串中解析出排名数字
func ParseRankNumber(rankStr string) int {
	if rankStr == "" {
		return 0
	}

	parts := strings.Split(rankStr, "/")
	if len(parts) > 0 {
		var rank int
		fmt.Sscanf(parts[0], "%d", &rank)
		return rank
	}

	return 0
}
