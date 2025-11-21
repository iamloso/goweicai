package gowencai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents the Wencai API client
type Client struct {
	httpClient *http.Client
	logger     *log.Logger
}

// NewClient creates a new Wencai client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: log.New(io.Discard, "[gowencai] ", log.LstdFlags),
	}
}

// SetLogger enables logging for the client
func (c *Client) SetLogger(logger *log.Logger) {
	c.logger = logger
}

// whileDo retries an operation with exponential backoff
func (c *Client) whileDo(do func() (interface{}, error), retry int, sleepSec int, logEnabled bool) (interface{}, error) {
	var lastErr error

	for count := 0; count < retry; count++ {
		if count > 0 {
			time.Sleep(time.Duration(sleepSec) * time.Second)
		}

		result, err := do()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if logEnabled {
			c.logger.Printf("第%d次尝试失败: %v", count+1, err)
		}
	}

	return nil, lastErr
}

// getRobotData fetches the condition data
func (c *Client) getRobotData(opts *QueryOptions) (*RobotDataParams, error) {
	if opts.Log {
		c.logger.Println("获取condition开始")
	}

	queryType := opts.QueryType
	if queryType == "" {
		queryType = "stock"
	}

	data := map[string]interface{}{
		"add_info":         `{"urp":{"scene":1,"company":1,"business":1},"contentType":"json","searchInfo":true}`,
		"perpage":          "10",
		"page":             1,
		"source":           "Ths_iwencai_Xuangu",
		"log_info":         `{"input_type":"click"}`,
		"version":          "2.0",
		"secondary_intent": queryType,
		"question":         opts.Query,
	}

	if opts.Pro {
		data["iwcpro"] = 1
	}

	do := func() (interface{}, error) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}

		req, err := http.NewRequest("POST", "http://www.iwencai.com/customized/chart/get-robot-data", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		headers, err := Headers(opts.Cookie, opts.UserAgent)
		if err != nil {
			return nil, fmt.Errorf("failed to generate headers: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		params, err := Convert(body)
		if err != nil {
			return nil, fmt.Errorf("failed to convert response: %w", err)
		}

		if opts.Log {
			c.logger.Println("获取get_robot_data成功")
		}

		return params, nil
	}

	retry := opts.Retry
	if retry == 0 {
		retry = 10
	}

	result, err := c.whileDo(do, retry, opts.Sleep, opts.Log)
	if err != nil {
		if opts.Log {
			c.logger.Printf("获取get_robot_data失败: %v", err)
		}
		return nil, err
	}

	return result.(*RobotDataParams), nil
}

// getPage fetches data for a single page
func (c *Client) getPage(urlParams map[string]string, opts *QueryOptions) ([]map[string]interface{}, error) {
	if opts.Log {
		c.logger.Printf("第%d页开始", opts.Page)
	}

	var targetURL string
	var path string
	formData := url.Values{}

	// Add URL params
	for key, value := range urlParams {
		formData.Set(key, value)
	}

	perpage := opts.PerPage
	if perpage == 0 {
		perpage = 100
	}
	formData.Set("perpage", strconv.Itoa(perpage))
	formData.Set("page", strconv.Itoa(opts.Page))

	if len(opts.Find) == 0 {
		targetURL = "http://www.iwencai.com/gateway/urp/v7/landing/getDataList"
		if opts.Pro {
			targetURL += "?iwcpro=1"
		}
		path = "answer.components.0.data.datas"
	} else {
		find := strings.Join(opts.Find, ",")
		formData.Set("query_type", opts.QueryType)
		formData.Set("question", find)
		targetURL = "http://www.iwencai.com/unifiedwap/unified-wap/v2/stock-pick/find"
		path = "data.data.datas"
	}

	do := func() (interface{}, error) {
		req, err := http.NewRequest("POST", targetURL, strings.NewReader(formData.Encode()))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		headers, err := Headers(opts.Cookie, opts.UserAgent)
		if err != nil {
			return nil, fmt.Errorf("failed to generate headers: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		dataList := GetValue(result, path)
		if dataList == nil {
			return nil, fmt.Errorf("第%d页返回空", opts.Page)
		}

		dataArray, ok := dataList.([]interface{})
		if !ok || len(dataArray) == 0 {
			return nil, fmt.Errorf("第%d页返回空", opts.Page)
		}

		// Convert to []map[string]interface{}
		var dataListMap []map[string]interface{}
		for _, item := range dataArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				dataListMap = append(dataListMap, itemMap)
			}
		}

		if opts.Log {
			c.logger.Printf("第%d页成功", opts.Page)
		}

		return dataListMap, nil
	}

	retry := opts.Retry
	if retry == 0 {
		retry = 10
	}

	result, err := c.whileDo(do, retry, opts.Sleep, opts.Log)
	if err != nil {
		if opts.Log {
			c.logger.Printf("第%d页失败: %v", opts.Page, err)
		}
		return nil, err
	}

	return result.([]map[string]interface{}), nil
}

// loopPage fetches multiple pages
func (c *Client) loopPage(loop interface{}, rowCount int, urlParams map[string]string, opts *QueryOptions) ([]map[string]interface{}, error) {
	perpage := opts.PerPage
	if perpage == 0 {
		perpage = 100
	}

	maxPage := int(math.Ceil(float64(rowCount) / float64(perpage)))
	
	var loopCount int
	switch v := loop.(type) {
	case bool:
		if v {
			loopCount = maxPage
		} else {
			loopCount = 1
		}
	case int:
		loopCount = v
	case float64:
		loopCount = int(v)
	default:
		loopCount = 1
	}

	if opts.Page == 0 {
		opts.Page = 1
	}
	initPage := opts.Page

	var result []map[string]interface{}

	for count := 0; count < loopCount; count++ {
		opts.Page = initPage + count
		pageResult, err := c.getPage(urlParams, opts)
		if err != nil {
			return result, err
		}

		result = append(result, pageResult...)
	}

	return result, nil
}

// Get performs the main query
func (c *Client) Get(opts *QueryOptions) (interface{}, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	if opts.Query == "" {
		return nil, fmt.Errorf("query is required")
	}

	if opts.Cookie == "" {
		return nil, fmt.Errorf("cookie is required")
	}

	// Set defaults
	if opts.QueryType == "" {
		opts.QueryType = "stock"
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if opts.PerPage == 0 {
		opts.PerPage = 100
	}
	if opts.Retry == 0 {
		opts.Retry = 10
	}

	params, err := c.getRobotData(opts)
	if err != nil {
		return nil, err
	}

	// Check if data contains condition
	if dataMap, ok := params.Data.(map[string]interface{}); ok {
		if condition := dataMap["condition"]; condition != nil {
			// Has condition, fetch page data
			if opts.Loop != nil && opts.Loop != false && len(opts.Find) == 0 {
				return c.loopPage(opts.Loop, params.RowCount, params.URLParams, opts)
			}
			return c.getPage(params.URLParams, opts)
		}
	}

	// No condition, return data directly
	if opts.NoDetail {
		return nil, nil
	}

	return params.Data, nil
}

// Get is a convenience function that creates a client and performs a query
func Get(opts *QueryOptions) (interface{}, error) {
	client := NewClient()
	if opts.Log {
		client.SetLogger(log.New(log.Writer(), "[gowencai] ", log.LstdFlags))
	}
	return client.Get(opts)
}
