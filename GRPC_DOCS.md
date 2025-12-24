# gRPC 服务文档

## 概述

GoWencai 提供了 gRPC API 接口用于高性能的股票数据查询和管理。

**gRPC 地址**: `localhost:9000`

## Proto 文件

位置: `api/proto/stock.proto`

## 服务定义

```protobuf
service StockService {
  // QueryStocks 查询股票
  rpc QueryStocks(QueryStocksRequest) returns (QueryStocksResponse);
  
  // GetLatestStocks 获取最新股票
  rpc GetLatestStocks(GetLatestStocksRequest) returns (GetLatestStocksResponse);
  
  // TriggerFetch 手动触发数据抓取
  rpc TriggerFetch(TriggerFetchRequest) returns (TriggerFetchResponse);
}
```

## 接口说明

### 1. QueryStocks - 查询股票

**请求**: `QueryStocksRequest`

```protobuf
message QueryStocksRequest {
  string code = 1;         // 股票代码
  string start_date = 2;   // 开始日期 YYYY-MM-DD
  string end_date = 3;     // 结束日期 YYYY-MM-DD
  int32 page = 4;          // 页码
  int32 page_size = 5;     // 每页数量
}
```

**响应**: `QueryStocksResponse`

```protobuf
message QueryStocksResponse {
  int32 total = 1;           // 总数
  int32 page = 2;            // 当前页
  int32 size = 3;            // 每页大小
  repeated Stock stocks = 4;  // 股票列表
}
```

### 2. GetLatestStocks - 获取最新股票

**请求**: `GetLatestStocksRequest`

```protobuf
message GetLatestStocksRequest {
  int32 limit = 1;  // 返回数量限制
}
```

**响应**: `GetLatestStocksResponse`

```protobuf
message GetLatestStocksResponse {
  int32 total = 1;           // 总数
  repeated Stock stocks = 2;  // 股票列表
}
```

### 3. TriggerFetch - 手动触发数据抓取

**请求**: `TriggerFetchRequest`

```protobuf
message TriggerFetchRequest {}
```

**响应**: `TriggerFetchResponse`

```protobuf
message TriggerFetchResponse {
  bool success = 1;        // 是否成功
  string message = 2;      // 消息
  int32 fetched_count = 3; // 抓取数量
}
```

## Stock 消息定义

```protobuf
message Stock {
  int64 id = 1;
  string code = 2;
  string name = 3;
  string date = 4;
  double price = 5;
  double change_percent = 6;
  double turnover = 7;
  double market_value = 8;
  int32 board_days = 9;
  string limit_reason = 10;
  string highlights = 11;
  string industry = 12;
  string concepts = 13;
  double pre_close = 14;
  double open = 15;
  double high = 16;
  double low = 17;
  double close = 18;
  double volume = 19;
  double bid_amount = 20;
  double bid_unmatched = 21;
  double bid_change = 22;
  double seal_amount = 23;
  string created_at = 24;
  string updated_at = 25;
}
```

## 使用示例

### Go 客户端

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    pb "github.com/iamloso/goweicai/api/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // 连接 gRPC 服务器
    conn, err := grpc.NewClient("localhost:9000", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer conn.Close()

    client := pb.NewStockServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 1. 获取最新股票
    latestResp, err := client.GetLatestStocks(ctx, &pb.GetLatestStocksRequest{
        Limit: 10,
    })
    if err != nil {
        log.Fatalf("调用失败: %v", err)
    }
    fmt.Printf("获取到 %d 条股票\n", len(latestResp.Stocks))

    // 2. 查询特定股票
    queryResp, err := client.QueryStocks(ctx, &pb.QueryStocksRequest{
        Code:     "000001",
        Page:     1,
        PageSize: 20,
    })
    if err != nil {
        log.Fatalf("调用失败: %v", err)
    }
    fmt.Printf("查询到 %d 条记录\n", len(queryResp.Stocks))

    // 3. 触发数据抓取
    fetchResp, err := client.TriggerFetch(ctx, &pb.TriggerFetchRequest{})
    if err != nil {
        log.Fatalf("调用失败: %v", err)
    }
    fmt.Printf("抓取结果: %v, %s\n", fetchResp.Success, fetchResp.Message)
}
```

### Python 客户端

首先安装依赖：

```bash
pip install grpcio grpcio-tools
```

生成 Python 代码：

```bash
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. api/proto/stock.proto
```

使用示例：

```python
import grpc
from api.proto import stock_pb2
from api.proto import stock_pb2_grpc

def main():
    # 连接 gRPC 服务器
    channel = grpc.insecure_channel('localhost:9000')
    client = stock_pb2_grpc.StockServiceStub(channel)

    # 1. 获取最新股票
    latest_resp = client.GetLatestStocks(
        stock_pb2.GetLatestStocksRequest(limit=10)
    )
    print(f"获取到 {len(latest_resp.stocks)} 条股票")

    # 2. 查询特定股票
    query_resp = client.QueryStocks(
        stock_pb2.QueryStocksRequest(
            code="000001",
            page=1,
            page_size=20
        )
    )
    print(f"查询到 {len(query_resp.stocks)} 条记录")

    # 3. 触发数据抓取
    fetch_resp = client.TriggerFetch(
        stock_pb2.TriggerFetchRequest()
    )
    print(f"抓取结果: {fetch_resp.success}, {fetch_resp.message}")

if __name__ == '__main__':
    main()
```

### grpcurl 测试

安装 grpcurl:

```bash
# macOS
brew install grpcurl

# Linux
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

使用示例：

```bash
# 1. 列出所有服务
grpcurl -plaintext localhost:9000 list

# 2. 列出服务的所有方法
grpcurl -plaintext localhost:9000 list stock.v1.StockService

# 3. 获取最新股票
grpcurl -plaintext -d '{"limit": 10}' \
  localhost:9000 stock.v1.StockService/GetLatestStocks

# 4. 查询股票
grpcurl -plaintext -d '{"code": "000001", "page": 1, "page_size": 20}' \
  localhost:9000 stock.v1.StockService/QueryStocks

# 5. 触发数据抓取
grpcurl -plaintext -d '{}' \
  localhost:9000 stock.v1.StockService/TriggerFetch
```

## 配置说明

在 `configs/config.yaml` 中配置 gRPC 服务器：

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 60s
  grpc:
    addr: 0.0.0.0:9000  # gRPC 监听地址
    timeout: 60s
```

## 代码生成

修改 proto 文件后，重新生成代码：

```bash
./generate.sh
```

或手动生成：

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/stock.proto
```

## 性能优势

相比 HTTP REST API，gRPC 具有以下优势：

1. **更高效的序列化**: 使用 Protocol Buffers，比 JSON 更紧凑
2. **HTTP/2 支持**: 多路复用、头部压缩、流控制
3. **强类型**: 编译时类型检查，减少运行时错误
4. **双向流**: 支持客户端流、服务端流、双向流
5. **更好的性能**: 更低的延迟和更高的吞吐量

## 错误处理

gRPC 使用标准的状态码：

| 状态码 | 说明 |
|--------|------|
| OK | 成功 |
| CANCELLED | 操作被取消 |
| UNKNOWN | 未知错误 |
| INVALID_ARGUMENT | 无效参数 |
| DEADLINE_EXCEEDED | 超时 |
| NOT_FOUND | 未找到 |
| ALREADY_EXISTS | 已存在 |
| PERMISSION_DENIED | 权限拒绝 |
| RESOURCE_EXHAUSTED | 资源耗尽 |
| FAILED_PRECONDITION | 前置条件失败 |
| ABORTED | 中止 |
| OUT_OF_RANGE | 超出范围 |
| UNIMPLEMENTED | 未实现 |
| INTERNAL | 内部错误 |
| UNAVAILABLE | 不可用 |
| DATA_LOSS | 数据丢失 |
| UNAUTHENTICATED | 未认证 |

## 注意事项

1. **安全性**: 当前使用不安全的连接（insecure），生产环境建议启用 TLS
2. **认证**: 建议添加认证机制（如 Token、mTLS）
3. **限流**: 建议添加请求限流
4. **日志**: 所有请求都会记录在应用日志中
5. **监控**: 建议添加 gRPC 监控指标

## 相关文档

- [Protocol Buffers](https://protobuf.dev/)
- [gRPC Documentation](https://grpc.io/docs/)
- [gRPC-Go](https://github.com/grpc/grpc-go)

## 文件结构

```
api/
└── proto/
    ├── stock.proto          # Proto 定义文件
    ├── stock.pb.go          # 生成的消息代码
    └── stock_grpc.pb.go     # 生成的服务代码

internal/
└── service/
    └── grpc.go              # gRPC 服务实现

cmd/
└── goweicai/
    └── main.go              # gRPC 服务器启动代码
```

## 后续优化

- [ ] 添加 TLS 支持
- [ ] 添加认证和授权
- [ ] 添加请求拦截器
- [ ] 添加监控指标
- [ ] 添加流式接口
- [ ] 添加反射服务
- [ ] 性能调优
