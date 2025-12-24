package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/iamloso/goweicai/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// HTTPService HTTP 服务
type HTTPService struct {
	stockUc *biz.StockUsecase
	log     *log.Helper
}

// NewHTTPService 创建 HTTP 服务
func NewHTTPService(stockUc *biz.StockUsecase, logger log.Logger) *HTTPService {
	return &HTTPService{
		stockUc: stockUc,
		log:     log.NewHelper(log.With(logger, "module", "service/http")),
	}
}

// StockQueryRequest 股票查询请求
type StockQueryRequest struct {
	Code      string `json:"code"`       // 股票代码
	StartDate string `json:"start_date"` // 开始日期 YYYY-MM-DD
	EndDate   string `json:"end_date"`   // 结束日期 YYYY-MM-DD
	Page      int    `json:"page"`       // 页码，默认 1
	PageSize  int    `json:"page_size"`  // 每页数量，默认 20
}

// StockQueryResponse 股票查询响应
type StockQueryResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *StockData    `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// StockData 股票数据
type StockData struct {
	Total  int          `json:"total"`
	Page   int          `json:"page"`
	Size   int          `json:"size"`
	Stocks []*biz.Stock `json:"stocks"`
}

// RegisterRoutes 注册路由
func (s *HTTPService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/stocks/query", s.handleQueryStocks)
	mux.HandleFunc("/api/stocks/latest", s.handleGetLatestStocks)
	mux.HandleFunc("/health", s.handleHealth)
}

// handleQueryStocks 处理股票查询
func (s *HTTPService) handleQueryStocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "只支持 GET 或 POST 请求")
		return
	}

	var req StockQueryRequest

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "请求参数格式错误: "+err.Error())
			return
		}
	} else {
		// GET 请求从 query string 获取参数
		req.Code = r.URL.Query().Get("code")
		req.StartDate = r.URL.Query().Get("start_date")
		req.EndDate = r.URL.Query().Get("end_date")
		if page := r.URL.Query().Get("page"); page != "" {
			req.Page, _ = strconv.Atoi(page)
		}
		if pageSize := r.URL.Query().Get("page_size"); pageSize != "" {
			req.PageSize, _ = strconv.Atoi(pageSize)
		}
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大每页数量
	}
	
	// 这里暂时返回空数据，等待实现具体的查询逻辑
	// TODO: 实现从数据库查询股票数据
	s.log.Infof("查询股票: code=%s, start_date=%s, end_date=%s, page=%d, size=%d",
		req.Code, req.StartDate, req.EndDate, req.Page, req.PageSize)

	stocks := make([]*biz.Stock, 0)
	
	s.respondSuccess(w, &StockData{
		Total:  0,
		Page:   req.Page,
		Size:   req.PageSize,
		Stocks: stocks,
	})
}

// handleGetLatestStocks 获取最新股票数据
func (s *HTTPService) handleGetLatestStocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "只支持 GET 请求")
		return
	}

	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}
	
	s.log.Infof("获取最新股票数据: limit=%d", limit)

	// TODO: 实现从数据库查询最新股票数据
	stocks := make([]*biz.Stock, 0)

	s.respondSuccess(w, &StockData{
		Total:  0,
		Page:   1,
		Size:   limit,
		Stocks: stocks,
	})
}

// handleHealth 健康检查
func (s *HTTPService) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// respondSuccess 返回成功响应
func (s *HTTPService) respondSuccess(w http.ResponseWriter, data interface{}) {
	response := StockQueryResponse{
		Code:    0,
		Message: "success",
		Data:    data.(*StockData),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// respondError 返回错误响应
func (s *HTTPService) respondError(w http.ResponseWriter, statusCode int, message string) {
	response := StockQueryResponse{
		Code:    statusCode,
		Message: "error",
		Error:   message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
