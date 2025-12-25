package service

import (
	"fmt"
	"strings"
)

// getStringValue 从 map 中获取字符串值
func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// getFloatValue 从 map 中获取浮点值
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

// getIntValue 从 map 中获取整数值
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

// getAmountValue 从 map 中获取金额值并转换为万或亿为单位
// 规则：金额 >= 1亿 时用亿为单位，< 1亿时用万为单位
// 返回：转换后的金额值和单位（"万" 或 "亿"）
func getAmountValueWithUnit(m map[string]interface{}, key string) (float64, string) {
	if v, ok := m[key]; ok && v != nil {
		var amount float64

		switch val := v.(type) {
		case float64:
			amount = val
		case float32:
			amount = float64(val)
		case int:
			amount = float64(val)
		case int64:
			amount = float64(val)
		case string:
			// 处理科学计数法和普通数字字符串
			if strings.Contains(val, "E") || strings.Contains(val, "e") {
				fmt.Sscanf(val, "%e", &amount)
			} else {
				// 移除可能的逗号分隔符
				cleaned := strings.ReplaceAll(val, ",", "")
				fmt.Sscanf(cleaned, "%f", &amount)
			}
		}

		// 转换为万或亿为单位
		if amount >= 100000000 { // >= 1亿
			return amount / 100000000, "亿" // 转换为亿
		} else {
			return amount / 10000, "万" // 转换为万
		}
	}
	return 0, "万"
}

// getAmountValue 从 map 中获取金额值并转换为万或亿为单位（仅返回数值）
// 规则：金额 >= 1亿 时用亿为单位，< 1亿时用万为单位
func getAmountValue(m map[string]interface{}, key string) float64 {
	value, _ := getAmountValueWithUnit(m, key)
	return value
}

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
