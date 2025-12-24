package main

import (
main.go "context"
main.go "flag"
main.go "fmt"
main.go "log"
main.go "time"

main.go pb "github.com/iamloso/goweicai/api/proto"
main.go "google.golang.org/grpc"
main.go "google.golang.org/grpc/credentials/insecure"
)

var (
main.go addr = flag.String("addr", "localhost:9000", "gRPC server address")
)

func main() {
main.go flag.Parse()

main.go conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
main.go if err != nil {
main.go main.go log.Fatalf("连接失败: %v", err)
main.go }
main.go defer conn.Close()

main.go client := pb.NewStockServiceClient(conn)
main.go ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
main.go defer cancel()

main.go fmt.Println("=== GoWencai gRPC 客户端测试 ===\n")

main.go fmt.Println("1. 获取最新股票（limit=5）")
main.go fmt.Println("----------------------------------------")
main.go latestResp, err := client.GetLatestStocks(ctx, &pb.GetLatestStocksRequest{
main.go main.go Limit: 5,
main.go })
main.go if err != nil {
main.go main.go log.Fatalf("调用失败: %v", err)
main.go }
main.go fmt.Printf("Total: %d\n", latestResp.Total)
main.go fmt.Printf("Stocks: %d 条\n", len(latestResp.Stocks))
main.go fmt.Println()

main.go fmt.Println("2. 查询股票（code=000001）")
main.go fmt.Println("----------------------------------------")
main.go queryResp, err := client.QueryStocks(ctx, &pb.QueryStocksRequest{
main.go main.go Code:     "000001",
main.go main.go Page:     1,
main.go main.go PageSize: 10,
main.go })
main.go if err != nil {
main.go main.go log.Fatalf("调用失败: %v", err)
main.go }
main.go fmt.Printf("Total: %d\n", queryResp.Total)
main.go fmt.Printf("Page: %d\n", queryResp.Page)
main.go fmt.Printf("Size: %d\n", queryResp.Size)
main.go fmt.Printf("Stocks: %d 条\n", len(queryResp.Stocks))
main.go fmt.Println()

main.go fmt.Println("3. 触发数据抓取")
main.go fmt.Println("----------------------------------------")
main.go fetchResp, err := client.TriggerFetch(ctx, &pb.TriggerFetchRequest{})
main.go if err != nil {
main.go main.go log.Fatalf("调用失败: %v", err)
main.go }
main.go fmt.Printf("Success: %v\n", fetchResp.Success)
main.go fmt.Printf("Message: %s\n", fetchResp.Message)
main.go fmt.Printf("Fetched: %d 条\n", fetchResp.FetchedCount)
main.go fmt.Println()

main.go fmt.Println("=== 所有测试完成 ===")
}
