package service

import (
	"context"

	pb "github.com/iamloso/goweicai/api/proto"
	"github.com/iamloso/goweicai/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// GRPCService gRPC 服务
type GRPCService struct {
	pb.UnimplementedStockServiceServer
	stockUc   *biz.StockUsecase
	wencaiSvc *WencaiService
	log       *log.Helper
}

// NewGRPCService 创建 gRPC 服务
func NewGRPCService(stockUc *biz.StockUsecase, wencaiSvc *WencaiService, logger log.Logger) *GRPCService {
	return &GRPCService{
		stockUc:   stockUc,
		wencaiSvc: wencaiSvc,
		log:       log.NewHelper(log.With(logger, "module", "service/grpc")),
	}
}

// QueryStocks 查询股票
func (s *GRPCService) QueryStocks(ctx context.Context, req *pb.QueryStocksRequest) (*pb.QueryStocksResponse, error) {
	s.log.Infof("gRPC 查询股票: code=%s, start_date=%s, end_date=%s, page=%d, size=%d",
		req.Code, req.StartDate, req.EndDate, req.Page, req.PageSize)

	// TODO: 实现从数据库查询
	stocks := make([]*pb.Stock, 0)

	return &pb.QueryStocksResponse{
		Total:  0,
		Page:   req.Page,
		Size:   req.PageSize,
		Stocks: stocks,
	}, nil
}

// GetLatestStocks 获取最新股票
func (s *GRPCService) GetLatestStocks(ctx context.Context, req *pb.GetLatestStocksRequest) (*pb.GetLatestStocksResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	s.log.Infof("gRPC 获取最新股票: limit=%d", limit)

	// TODO: 实现从数据库查询
	stocks := make([]*pb.Stock, 0)

	return &pb.GetLatestStocksResponse{
		Total:  0,
		Stocks: stocks,
	}, nil
}

// TriggerFetch 手动触发数据抓取
func (s *GRPCService) TriggerFetch(ctx context.Context, req *pb.TriggerFetchRequest) (*pb.TriggerFetchResponse, error) {
	s.log.Info("gRPC 触发数据抓取")

	// 调用 WencaiService 抓取数据
	err := s.wencaiSvc.FetchAndSaveStocks(ctx)
	if err != nil {
		s.log.Errorf("数据抓取失败: %v", err)
		return &pb.TriggerFetchResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	s.log.Info("数据抓取成功")
	return &pb.TriggerFetchResponse{
		Success:      true,
		Message:      "数据抓取成功",
		FetchedCount: 0, // TODO: 返回实际抓取数量
	}, nil
}
