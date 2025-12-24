package mainpackage grpcclient

package main

import (

	"context"import (

	"flag"	"context"

	"fmt"	"flag"

	"log"	"fmt"

	"time"	"log"

	"time"

	pb "github.com/iamloso/goweicai/api/proto"

	"google.golang.org/grpc"	pb "github.com/iamloso/goweicai/api/proto"

	"google.golang.org/grpc/credentials/insecure"	"google.golang.org/grpc"

)	"google.golang.org/grpc/credentials/insecure"

)

var (

	addr = flag.String("addr", "localhost:9000", "gRPC server address")var (

)	addr = flag.String("addr", "localhost:9000", "gRPC server address")

)

func main() {

	flag.Parse()func main() {

	flag.Parse()

	// 连接 gRPC 服务器

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))	// 连接 gRPC 服务器

	if err != nil {	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

		log.Fatalf("连接失败: %v", err)	if err != nil {

	}		log.Fatalf("连接失败: %v", err)

	defer conn.Close()	}

	defer conn.Close()

	client := pb.NewStockServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)	client := pb.NewStockServiceClient(conn)

	defer cancel()	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	fmt.Println("=== GoWencai gRPC 客户端测试 ===\n")

	fmt.Println("=== GoWencai gRPC 客户端测试 ===\n")

	// 测试 1: 获取最新股票

	fmt.Println("1. 获取最新股票（limit=5）")	// 测试 1: 获取最新股票

	fmt.Println("----------------------------------------")	fmt.Println("1. 获取最新股票（limit=5）")

	latestResp, err := client.GetLatestStocks(ctx, &pb.GetLatestStocksRequest{	fmt.Println("----------------------------------------")

		Limit: 5,	latestResp, err := client.GetLatestStocks(ctx, &pb.GetLatestStocksRequest{

	})		Limit: 5,

	if err != nil {	})

		log.Fatalf("调用失败: %v", err)	if err != nil {

	}		log.Fatalf("调用失败: %v", err)

	fmt.Printf("Total: %d\n", latestResp.Total)	}

	fmt.Printf("Stocks: %d 条\n", len(latestResp.Stocks))	fmt.Printf("Total: %d\n", latestResp.Total)

	fmt.Println()	fmt.Printf("Stocks: %d 条\n", len(latestResp.Stocks))

	fmt.Println()

	// 测试 2: 查询股票

	fmt.Println("2. 查询股票（code=000001）")	// 测试 2: 查询股票

	fmt.Println("----------------------------------------")	fmt.Println("2. 查询股票（code=000001）")

	queryResp, err := client.QueryStocks(ctx, &pb.QueryStocksRequest{	fmt.Println("----------------------------------------")

		Code:     "000001",	queryResp, err := client.QueryStocks(ctx, &pb.QueryStocksRequest{

		Page:     1,		Code:     "000001",

		PageSize: 10,		Page:     1,

	})		PageSize: 10,

	if err != nil {	})

		log.Fatalf("调用失败: %v", err)	if err != nil {

	}		log.Fatalf("调用失败: %v", err)

	fmt.Printf("Total: %d\n", queryResp.Total)	}

	fmt.Printf("Page: %d\n", queryResp.Page)	fmt.Printf("Total: %d\n", queryResp.Total)

	fmt.Printf("Size: %d\n", queryResp.Size)	fmt.Printf("Page: %d\n", queryResp.Page)

	fmt.Printf("Stocks: %d 条\n", len(queryResp.Stocks))	fmt.Printf("Size: %d\n", queryResp.Size)

	fmt.Println()	fmt.Printf("Stocks: %d 条\n", len(queryResp.Stocks))

	fmt.Println()

	// 测试 3: 触发数据抓取

	fmt.Println("3. 触发数据抓取")	// 测试 3: 触发数据抓取

	fmt.Println("----------------------------------------")	fmt.Println("3. 触发数据抓取")

	fetchResp, err := client.TriggerFetch(ctx, &pb.TriggerFetchRequest{})	fmt.Println("----------------------------------------")

	if err != nil {	fetchResp, err := client.TriggerFetch(ctx, &pb.TriggerFetchRequest{})

		log.Fatalf("调用失败: %v", err)	if err != nil {

	}		log.Fatalf("调用失败: %v", err)

	fmt.Printf("Success: %v\n", fetchResp.Success)	}

	fmt.Printf("Message: %s\n", fetchResp.Message)	fmt.Printf("Success: %v\n", fetchResp.Success)

	fmt.Printf("Fetched: %d 条\n", fetchResp.FetchedCount)	fmt.Printf("Message: %s\n", fetchResp.Message)

	fmt.Println()	fmt.Printf("Fetched: %d 条\n", fetchResp.FetchedCount)

	fmt.Println()

	fmt.Println("=== 所有测试完成 ===")

}	fmt.Println("=== 所有测试完成 ===")

}
