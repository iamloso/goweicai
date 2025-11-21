package gowencai

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetValue safely gets a nested value from a map
func GetValue(data map[string]interface{}, path string) interface{} {
	keys := strings.Split(path, ".")
	var current interface{} = data

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
			if current == nil {
				return nil
			}
		case []interface{}:
			// Handle array index like "0"
			if idx := parseIndex(key); idx >= 0 && idx < len(v) {
				current = v[idx]
			} else {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}

// parseIndex converts string to int for array indexing
func parseIndex(s string) int {
	var idx int
	if _, err := fmt.Sscanf(s, "%d", &idx); err == nil {
		return idx
	}
	return -1
}

// ParseURLParams parses URL query parameters
func ParseURLParams(urlStr string) map[string]string {
	params := make(map[string]string)
	if urlStr == "" {
		return params
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return params
	}

	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params
}

// xuanguTableV1Handler handles xuangu_tableV1 type components
func xuanguTableV1Handler(comp map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"condition": GetValue(comp, "data.meta.extra.condition"),
		"comp_id":   comp["cid"],
		"uuid":      comp["puuid"],
	}
}

// commonHandler handles common type components
func commonHandler(comp map[string]interface{}) interface{} {
	datas := GetValue(comp, "data.datas")
	if datas != nil {
		if dataList, ok := datas.([]interface{}); ok {
			return dataList
		}
	}
	return GetValue(comp, "data")
}

// showTypeHandler processes different show_type components
func showTypeHandler(comp map[string]interface{}, comps []interface{}) interface{} {
	showType, _ := comp["show_type"].(string)

	switch showType {
	// Add specific handlers here if needed
	default:
		return commonHandler(comp)
	}
}

// getKey extracts the key for a component
func getKey(comp map[string]interface{}) string {
	if h1 := GetValue(comp, "title_config.data.h1"); h1 != nil {
		if str, ok := h1.(string); ok && str != "" {
			return str
		}
	}
	if title := GetValue(comp, "config.title"); title != nil {
		if str, ok := title.(string); ok && str != "" {
			return str
		}
	}
	if showType := comp["show_type"]; showType != nil {
		if str, ok := showType.(string); ok {
			return str
		}
	}
	return ""
}

// multiShowTypeHandler handles multiple show_type components
func multiShowTypeHandler(components []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, comp := range components {
		compMap, ok := comp.(map[string]interface{})
		if !ok {
			continue
		}

		key := getKey(compMap)
		value := showTypeHandler(compMap, components)

		if key != "" && value != nil {
			result[key] = value
		}
	}

	return result
}

// Convert processes the robot data response
func Convert(responseBody []byte) (*RobotDataParams, error) {
	var robotResp RobotResponse
	if err := json.Unmarshal(responseBody, &robotResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract content
	var content map[string]interface{}
	if len(robotResp.Data.Answer) > 0 && len(robotResp.Data.Answer[0].Txt) > 0 {
		contentRaw := robotResp.Data.Answer[0].Txt[0].Content

		// If content is a string, parse it as JSON
		if contentStr, ok := contentRaw.(string); ok {
			if err := json.Unmarshal([]byte(contentStr), &content); err != nil {
				return nil, fmt.Errorf("failed to parse content string: %w", err)
			}
		} else if contentMap, ok := contentRaw.(map[string]interface{}); ok {
			content = contentMap
		} else {
			return nil, fmt.Errorf("unexpected content type")
		}
	} else {
		return nil, fmt.Errorf("no answer data found")
	}

	// Extract components
	componentsRaw, ok := content["components"]
	if !ok {
		return nil, fmt.Errorf("no components found")
	}

	components, ok := componentsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("components is not an array")
	}

	params := &RobotDataParams{
		URLParams: make(map[string]string),
	}

	// Check if single xuangu_tableV1 component
	if len(components) == 1 {
		firstComp, ok := components[0].(map[string]interface{})
		if ok {
			showType, _ := firstComp["show_type"].(string)
			if showType == "xuangu_tableV1" {
				urlStr := ""
				if u := GetValue(firstComp, "config.other_info.footer_info.url"); u != nil {
					urlStr, _ = u.(string)
				}

				rowCount := 0
				if rc := GetValue(firstComp, "data.meta.extra.row_count"); rc != nil {
					if rcFloat, ok := rc.(float64); ok {
						rowCount = int(rcFloat)
					}
				}

				params.Data = xuanguTableV1Handler(firstComp)
				params.RowCount = rowCount
				params.URL = urlStr
				params.URLParams = ParseURLParams(urlStr)
				return params, nil
			}
		}
	}

	// Multiple components or not xuangu_tableV1
	var urlStr string
	if len(components) > 0 {
		if firstComp, ok := components[0].(map[string]interface{}); ok {
			if u := GetValue(firstComp, "config.other_info.footer_info.url"); u != nil {
				urlStr, _ = u.(string)
			}
		}
	}

	params.Data = multiShowTypeHandler(components)
	params.URL = urlStr
	params.URLParams = ParseURLParams(urlStr)

	return params, nil
}
