package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	gowencai "github.com/fenghuang/gowencai"
	_ "github.com/go-sql-driver/mysql"
)

var queryZpjj string = "早盘竞价未匹配资金排序,包含竞价涨幅;不含ST"
var query string = queryZpjj
var cookie string = `other_uid=Ths_iwencai_Xuangu_8lcskmb3gni2zd4etfbg1wkwnup0v634; ta_random_userid=sybmbrbrj7; cid=6c489adc15242a39ac7e243149197e521763607556; u_ukey=A10702B8689642C6BE607730E11E6E4A; u_uver=1.0.0; u_dpass=FfSUBGVVjZnSzvgF2N%2B39xCLBedVnnNrvixvHGqK5KsphhqUy8o1q4QIg2l7Zez0Hi80LrSsTFH9a%2B6rtRvqGg%3D%3D; u_did=A7005F415DF440FE8F8AC699543FC6AA; u_ttype=WEB; ttype=WEB; user=MDppYW1sb3NvOjpOb25lOjUwMDoxNjQ1NzAwMjI6NywxMTExMTExMTExMSw0MDs0NCwxMSw0MDs2LDEsNDA7NSwxLDQwOzEsMTAxLDQwOzIsMSw0MDszLDEsNDA7NSwxLDQwOzgsMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDEsNDA7MTAyLDEsNDA6MjU6OjoxNTQ1NzAwMjI6MTc2MzYwNzU4Njo6OjEzNjEwMjMxNDA6NjA0ODAwOjA6MTgxZWVhODllMTNjNDc4MDhkNTM3NjQ5Y2YzNTBkZmU1OmRlZmF1bHRfNTow; userid=154570022; u_name=iamloso; escapename=iamloso; ticket=a479fb4f0c4ba40c96cd826fc079ca41; user_status=0; utk=d104680c9471ede714de89c677bb5a53; sess_tk=eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiIsImtpZCI6InNlc3NfdGtfMSIsImJ0eSI6InNlc3NfdGsifQ.eyJqdGkiOiJlNWRmNTBmMzljNjQzN2Q1MDg3OGM0MTM5ZWE4ZWU4MTEiLCJpYXQiOjE3NjM2MDc1ODYsImV4cCI6MTc2NDIxMjM4Niwic3ViIjoiMTU0NTcwMDIyIiwiaXNzIjoidXBhc3MuaXdlbmNhaS5jb20iLCJhdWQiOiIyMDIwMTExODUyODg5MDcyIiwiYWN0Ijoib2ZjIiwiY3VocyI6IjY2ZGI2YmQzMTQyNmJjY2ZkMGMxNjBkNTVlNWY2YTQ4NjMwZmUzYTBmNDhjNjQ5MTQ5YjI3ZWMwZTY1YTA5Y2QifQ.Gza8-0nqW_cpgczyXYK4zHc2TH5ZRs7jFnTF5FwLN7GsSsao0cItkIh1OgrF3BP74bcotIEu2EZ70d0WFlVa9A; cuc=gkhkxir81hv3; RT="z=1&dm=iwencai.com&si=5f513b9c-0114-482f-a335-da7bce145f64&ss=mi6ueja9&sl=2&tt=aif&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=q47&ul=12g1&hd=12i4"; v=A6j36q_zVVUf_XnRlv9iHTeRcJ2_0Qzb7jXgX2LZ9CMWvUaDCuHcaz5FsOix` // 请替换为你的实际cookie值

// 数据库配置
const (
	dbUser     = "root"
	dbPassword = "eauDx15FxO83lS"
	dbHost     = "localhost"
	dbPort     = "3306"
	dbName     = "wc"
)

func main() {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库ping失败: %v", err)
	}
	log.Println("数据库连接成功")
	// 示例1: 基本查询
	fmt.Println("=== 示例1: 基本查询 ===")
	result1, err := gowencai.Get(&gowencai.QueryOptions{
		Query:  query,
		Cookie: cookie, // 必须填入你的cookie
		Log:    true,
		Loop:   true,
	})
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		printResult(result1)
		// 写入数据库
		if err := saveToDatabase(db, result1); err != nil {
			log.Printf("写入数据库失败: %v", err)
		} else {
			log.Println("数据写入成功")
		}
	}
}

func printResult(result interface{}) {
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("格式化结果失败: %v", err)
		return
	}
	fmt.Println(string(jsonData))

	// 限制输出长度
	// if len(jsonData) > 1000 {
	// 	fmt.Printf("%s\n... (结果太长，已截断)\n", jsonData[:1000])
	// } else {
	// 	fmt.Println(string(jsonData))
	// }
}

// saveToDatabase 将查询结果保存到数据库
func saveToDatabase(db *sql.DB, result interface{}) error {
	// 将结果转换为 []map[string]interface{}
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
		return fmt.Errorf("不支持的数据类型")
	}

	if len(items) == 0 {
		return fmt.Errorf("没有数据需要保存")
	}

	// 获取当前日期作为交易日期
	tradeDate := time.Now().Format("2006-01-02")
	tradeDateKey := time.Now().Format("20060102") // 用于字段名的日期格式，例如：20251121
	log.Printf("当前交易日期: %s, 字段日期键: %s", tradeDate, tradeDateKey)

	// 准备SQL语句
	checkSQL := `SELECT id FROM zp_jj WHERE code = ? AND trade_date = ?`

	insertSQL := `INSERT INTO zp_jj (
		code, market_code, latest_price, latest_change_rate,
		auction_unmatched_amount, auction_unmatched_amount_str, auction_unmatched_amount_rank,
		auction_unmatched_amount_rank_num,
		auction_change_rate, morning_auction_amount, morning_auction_amount_str, 
		turnover, turnover_str, circulation_market_value,
		stock_code, stock_name, trade_date
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		stock_name = ?
	WHERE code = ? AND trade_date = ?` // 开启事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback()

	// 批量插入/更新
	insertCount := 0
	updateCount := 0

	for _, item := range items {
		code := getStringValue(item, "code")
		if code == "" {
			continue
		}

		// 检查记录是否存在
		var existingID int64
		err := tx.QueryRow(checkSQL, code, tradeDate).Scan(&existingID)

		// 准备数据（使用动态日期字段）
		marketCode := getStringValue(item, "market_code")
		latestPrice := getFloatValue(item, "最新价")
		latestChangeRate := getFloatValue(item, "最新涨跌幅")

		// 动态构建字段名
		auctionAmountKey := fmt.Sprintf("竞价未匹配金额[%s]", tradeDateKey)
		auctionAmountRankKey := fmt.Sprintf("竞价未匹配金额排名[%s]", tradeDateKey)
		auctionChangeRateKey := fmt.Sprintf("竞价涨幅[%s]", tradeDateKey)

		auctionAmount := getIntValue(item, auctionAmountKey)
		auctionAmountStr := formatAmount(auctionAmount) // 格式化为万或亿
		auctionAmountRank := getStringValue(item, auctionAmountRankKey)
		auctionAmountRankNum := parseRankNumber(auctionAmountRank) // 解析排名数字
		auctionChangeRate := getFloatValue(item, auctionChangeRateKey)
		morningAuctionAmount := int64(0)      // 待添加
		morningAuctionAmountStr := ""         // 格式化的早盘竞价金额
		turnover := 0.0                       // 待添加
		turnoverStr := ""                     // 格式化的成交额
		circulationMarketValue := 0.0         // 待添加
		stockCode := getStringValue(item, "股票代码")
		stockName := getStringValue(item, "股票简称")

		if err == sql.ErrNoRows {
			// 记录不存在，执行插入
			_, err = tx.Exec(insertSQL,
				code, marketCode, latestPrice, latestChangeRate,
				auctionAmount, auctionAmountStr, auctionAmountRank, auctionAmountRankNum,
				auctionChangeRate, morningAuctionAmount, morningAuctionAmountStr, 
				turnover, turnoverStr, circulationMarketValue,
				stockCode, stockName, tradeDate,
			)
			if err != nil {
				log.Printf("插入数据失败 [%s]: %v", code, err)
				continue
			}
			insertCount++
			log.Printf("插入新记录: %s - %s (排名: %d)", code, stockName, auctionAmountRankNum)
		} else if err == nil {
			// 记录存在，执行更新
			_, err = tx.Exec(updateSQL,
				marketCode, latestPrice, latestChangeRate,
				auctionAmount, auctionAmountStr, auctionAmountRank, auctionAmountRankNum,
				auctionChangeRate, morningAuctionAmount, morningAuctionAmountStr, 
				turnover, turnoverStr, circulationMarketValue,
				stockName, code, tradeDate,
			)
			if err != nil {
				log.Printf("更新数据失败 [%s]: %v", code, err)
				continue
			}
			updateCount++
			log.Printf("更新记录: %s - %s (排名: %d)", code, stockName, auctionAmountRankNum)
		} else {
			log.Printf("查询记录失败 [%s]: %v", code, err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	log.Printf("操作完成 - 新插入: %d 条, 更新: %d 条", insertCount, updateCount)
	return nil
}

// 辅助函数：从map中获取字符串值
func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// 辅助函数：从map中获取浮点数值
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

// 辅助函数：从map中获取整数值
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
			var i int64
			fmt.Sscanf(val, "%d", &i)
			return i
		}
	}
	return 0
}

// formatAmount 将金额转换为带单位的字符串（万或亿）
func formatAmount(amount int64) string {
	if amount == 0 {
		return "" // 为0时返回空字符串
	}

	// 转换为浮点数以便计算
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

// parseRankNumber 从排名字符串中解析出排名数字（/前面的部分）
// 例如："64/5458" -> 64
func parseRankNumber(rankStr string) int {
	if rankStr == "" {
		return 0
	}

	// 按 / 分割
	parts := strings.Split(rankStr, "/")
	if len(parts) > 0 {
		var rank int
		fmt.Sscanf(parts[0], "%d", &rank)
		return rank
	}

	return 0
}
