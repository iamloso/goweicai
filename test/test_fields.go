package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	gowencai "github.com/fenghuang/gowencai"
)

var query string = "竞价未匹配金额；竞价金额；竞价涨幅；涨幅；成交金额；流通市值；连板天数；不含ST; 涨停原因；公司亮点；行业分类；概念题材；"
var cookie string = `other_uid=Ths_iwencai_Xuangu_8lcskmb3gni2zd4etfbg1wkwnup0v634; ta_random_userid=sybmbrbrj7; cid=6c489adc15242a39ac7e243149197e521763607556; u_ukey=A10702B8689642C6BE607730E11E6E4A; u_uver=1.0.0; u_dpass=FfSUBGVVjZnSzvgF2N%2B39xCLBedVnnNrvixvHGqK5KsphhqUy8o1q4QIg2l7Zez0Hi80LrSsTFH9a%2B6rtRvqGg%3D%3D; u_did=A7005F415DF440FE8F8AC699543FC6AA; u_ttype=WEB; ttype=WEB; user=MDppYW1sb3NvOjpOb25lOjUwMDoxNjQ1NzAwMjI6NywxMTExMTExMTExMSw0MDs0NCwxMSw0MDs2LDEsNDA7NSwxLDQwOzEsMTAxLDQwOzIsMSw0MDszLDEsNDA7NSwxLDQwOzgsMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDEsNDA7MTAyLDEsNDA6MjU6OjoxNTQ1NzAwMjI6MTc2MzYwNzU4Njo6OjEzNjEwMjMxNDA6NjA0ODAwOjA6MTgxZWVhODllMTNjNDc4MDhkNTM3NjQ5Y2YzNTBkZmU1OmRlZmF1bHRfNTow; userid=154570022; u_name=iamloso; escapename=iamloso; ticket=a479fb4f0c4ba40c96cd826fc079ca41; user_status=0; utk=d104680c9471ede714de89c677bb5a53; sess_tk=eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiIsImtpZCI6InNlc3NfdGtfMSIsImJ0eSI6InNlc3NfdGsifQ.eyJqdGkiOiJlNWRmNTBmMzljNjQzN2Q1MDg3OGM0MTM5ZWE4ZWU4MTEiLCJpYXQiOjE3NjM2MDc1ODYsImV4cCI6MTc2NDIxMjM4Niwic3ViIjoiMTU0NTcwMDIyIiwiaXNzIjoidXBhc3MuaXdlbmNhaS5jb20iLCJhdWQiOiIyMDIwMTExODUyODg5MDcyIiwiYWN0Ijoib2ZjIiwiY3VocyI6IjY2ZGI2YmQzMTQyNmJjY2ZkMGMxNjBkNTVlNWY2YTQ4NjMwZmUzYTBmNDhjNjQ5MTQ5YjI3ZWMwZTY1YTA5Y2QifQ.Gza8-0nqW_cpgczyXYK4zHc2TH5ZRs7jFnTF5FwLN7GsSsao0cItkIh1OgrF3BP74bcotIEu2EZ70d0WFlVa9A; cuc=gkhkxir81hv3; RT="z=1&dm=iwencai.com&si=5f513b9c-0114-482f-a335-da7bce145f64&ss=mi6ueja9&sl=2&tt=aif&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=q47&ul=12g1&hd=12i4"; v=A6j36q_zVVUf_XnRlv9iHTeRcJ2_0Qzb7jXgX2LZ9CMWvUaDCuHcaz5FsOix`

func main() {
	tradeDateKey := time.Now().Format("20060102")
	log.Printf("当前日期键: %s", tradeDateKey)

	// 查询第一页数据
	result, err := gowencai.Get(&gowencai.QueryOptions{
		Query:  query,
		Cookie: cookie,
		Log:    false,
		Loop:   false, // 只查询第一页
	})
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

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
	}

	if len(items) == 0 {
		log.Fatal("没有数据")
	}

	// 打印第一条记录的所有字段
	log.Println("\n=== 第一条记录的所有字段 ===")
	firstItem := items[0]
	for key, value := range firstItem {
		log.Printf("%s = %v", key, value)
	}

	// 测试字段提取
	log.Println("\n=== 测试字段提取 ===")
	
	code := getStringValue(firstItem, "code")
	stockName := getStringValue(firstItem, "股票简称")
	
	changeRateKey := fmt.Sprintf("涨跌幅:前复权[%s]", tradeDateKey)
	latestChangeRate := getFloatValue(firstItem, changeRateKey)
	
	turnoverKey := fmt.Sprintf("成交额[%s]", tradeDateKey)
	turnover := getIntValue(firstItem, turnoverKey)
	
	circulationMarketValueKey := fmt.Sprintf("a股市值(不含限售股)[%s]", tradeDateKey)
	circulationMarketValue := getFloatValue(firstItem, circulationMarketValueKey)
	
	industryCategory := getStringValue(firstItem, "所属同花顺行业")
	conceptTheme := getStringValue(firstItem, "所属概念")
	
	limitUpReasonKey := fmt.Sprintf("涨停原因类别[%s]", tradeDateKey)
	limitUpReason := getStringValue(firstItem, limitUpReasonKey)
	
	limitUpSealAmountKey := fmt.Sprintf("涨停封单额[%s]", tradeDateKey)
	limitUpSealAmount := getIntValue(firstItem, limitUpSealAmountKey)
	
	consecutiveLimitDaysKey := fmt.Sprintf("连续涨停天数[%s]", tradeDateKey)
	consecutiveLimitDays := getIntValue(firstItem, consecutiveLimitDaysKey)
	
	log.Printf("\n股票代码: %s", code)
	log.Printf("股票简称: %s", stockName)
	log.Printf("涨跌幅: %.2f%%", latestChangeRate)
	log.Printf("成交额: %d", turnover)
	log.Printf("流通市值: %.2f", circulationMarketValue)
	log.Printf("行业分类: %s", industryCategory)
	log.Printf("概念题材: %s", conceptTheme)
	log.Printf("涨停原因: %s", limitUpReason)
	log.Printf("涨停封单额: %d", limitUpSealAmount)
	log.Printf("连板天数: %d", consecutiveLimitDays)

	// 打印前3条记录的JSON
	log.Println("\n=== 前3条记录 ===")
	for i := 0; i < 3 && i < len(items); i++ {
		jsonData, _ := json.MarshalIndent(items[i], "", "  ")
		fmt.Printf("\n记录 %d:\n%s\n", i+1, string(jsonData))
	}
}

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
			var i int64
			fmt.Sscanf(val, "%d", &i)
			return i
		}
	}
	return 0
}
