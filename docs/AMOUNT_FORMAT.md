# 金额字段格式说明

## 字段列表

以下字段存储格式为带单位的字符串：

- `auction_unmatched_amount_str`: 竞价未匹配金额
- `morning_auction_amount_str`: 早盘竞价金额
- `turnover_str`: 成交额

## 存储格式

金额以**万**或**亿**为单位存储，格式示例：

- `12万` - 金额 < 1亿时使用万为单位
- `284.11万` - 保留2位小数（自动去除尾部0）
- `1.3亿` - 金额 >= 1亿时使用亿为单位
- `45.67亿` - 保留2位小数

## 转换规则

1. **单位选择**：
   - 原始金额 >= 1亿（100,000,000）→ 使用**亿**为单位
   - 原始金额 < 1亿 → 使用**万**为单位

2. **数值格式**：
   - 保留最多2位小数
   - 自动去除尾部的0和小数点
   - 示例：`12.00万` → `12万`，`1.30亿` → `1.3亿`

3. **原始数据处理**：
   - 支持普通数字：`123456.78`
   - 支持科学计数法：`4.598915e+08`
   - 支持逗号分隔符：`1,234,567.89`

## 代码实现

格式化函数位于 `internal/service/utils.go`：

```go
// formatAmountStr 从 map 中获取金额值并格式化为带单位的字符串
// 返回格式如：12万、1.3亿
func formatAmountStr(m map[string]interface{}, key string) string {
    value, unit := getAmountValueWithUnit(m, key)
    if value == 0 {
        return ""
    }
    // 保留2位小数，去除尾部的0
    str := fmt.Sprintf("%.2f", value)
    // 去除尾部的0和小数点
    str = strings.TrimRight(str, "0")
    str = strings.TrimRight(str, ".")
    return str + unit
}
```

## 使用示例

### 存储示例

| 原始金额 | 转换后 | 存储值 |
|---------|--------|--------|
| 123456.78 | 12.3456万 | `12.35万` |
| 1234567.89 | 123.4567万 | `123.46万` |
| 123456789 | 1.23456789亿 | `1.23亿` |
| 459891500 | 4.598915亿 | `4.6亿` |
| 10000 | 1万 | `1万` |
| 100000000 | 1亿 | `1亿` |

### 查询示例

```sql
-- 查看格式化后的金额
SELECT 
    stock_name AS '股票简称',
    auction_unmatched_amount_str AS '竞价未匹配金额',
    morning_auction_amount_str AS '竞价金额',
    turnover_str AS '成交额'
FROM base_info
WHERE trade_date = CURDATE()
LIMIT 10;
```
