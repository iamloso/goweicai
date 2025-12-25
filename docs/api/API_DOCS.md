# HTTP API 文档

## 概述

GoWencai 提供了 RESTful HTTP API 接口用于查询股票数据。

**Base URL**: `http://localhost:8000`

## 接口列表

### 1. 健康检查

检查服务是否正常运行。

**接口**: `GET /health`

**响应示例**:
```json
{
  "status": "ok",
  "time": "2024-12-24T10:30:00+08:00"
}
```

---

### 2. 查询股票数据

根据条件查询股票数据。

**接口**: `POST /api/stocks/query` 或 `GET /api/stocks/query`

#### 请求参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| code | string | 否 | 股票代码 | "000001" |
| start_date | string | 否 | 开始日期 | "2024-01-01" |
| end_date | string | 否 | 结束日期 | "2024-12-31" |
| page | int | 否 | 页码，默认 1 | 1 |
| page_size | int | 否 | 每页数量，默认 20，最大 100 | 20 |

#### POST 请求示例

```bash
curl -X POST http://localhost:8000/api/stocks/query \
  -H "Content-Type: application/json" \
  -d '{
    "code": "000001",
    "start_date": "2024-01-01",
    "end_date": "2024-12-31",
    "page": 1,
    "page_size": 20
  }'
```

#### GET 请求示例

```bash
curl "http://localhost:8000/api/stocks/query?code=000001&start_date=2024-01-01&end_date=2024-12-31&page=1&page_size=20"
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "size": 20,
    "stocks": [
      {
        "id": 1,
        "code": "000001",
        "name": "平安银行",
        "date": "2024-12-24",
        "price": 12.34,
        "change_percent": 1.23,
        "turnover": 1234567890.00,
        "market_value": 9876543210.00,
        "created_at": "2024-12-24T10:00:00+08:00",
        "updated_at": "2024-12-24T10:00:00+08:00"
      }
    ]
  }
}
```

---

### 3. 获取最新股票数据

获取最新的股票数据列表。

**接口**: `GET /api/stocks/latest`

#### 请求参数

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| limit | int | 否 | 返回数量，默认 50，最大 200 | 50 |

#### 请求示例

```bash
curl "http://localhost:8000/api/stocks/latest?limit=50"
```

#### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "page": 1,
    "size": 50,
    "stocks": [
      {
        "id": 1,
        "code": "000001",
        "name": "平安银行",
        "date": "2024-12-24",
        "price": 12.34,
        "change_percent": 1.23,
        "turnover": 1234567890.00,
        "market_value": 9876543210.00,
        "created_at": "2024-12-24T10:00:00+08:00",
        "updated_at": "2024-12-24T10:00:00+08:00"
      }
    ]
  }
}
```

---

## 响应格式

所有接口统一使用以下响应格式：

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 具体数据
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "error",
  "error": "错误详细信息"
}
```

## 状态码说明

| HTTP 状态码 | 说明 |
|------------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 405 | 请求方法不允许 |
| 500 | 服务器内部错误 |

## 股票数据字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 记录 ID |
| code | string | 股票代码 |
| name | string | 股票名称 |
| date | string | 交易日期 |
| price | float | 股价 |
| change_percent | float | 涨跌幅 (%) |
| turnover | float | 成交额 |
| market_value | float | 流通市值 |
| board_days | int | 连板天数 |
| limit_reason | string | 涨停原因 |
| highlights | string | 公司亮点 |
| industry | string | 行业分类 |
| concepts | string | 概念题材 |
| pre_close | float | 昨收价 |
| open | float | 开盘价 |
| high | float | 最高价 |
| low | float | 最低价 |
| close | float | 收盘价 |
| volume | float | 成交量 |
| bid_amount | float | 竞价金额 |
| bid_unmatched | float | 竞价未匹配金额 |
| bid_change | float | 竞价涨幅 |
| seal_amount | float | 封单金额 |
| created_at | string | 创建时间 |
| updated_at | string | 更新时间 |

## 使用示例

### Python 示例

```python
import requests

# 健康检查
response = requests.get('http://localhost:8000/health')
print(response.json())

# 查询股票
data = {
    'code': '000001',
    'start_date': '2024-01-01',
    'end_date': '2024-12-31',
    'page': 1,
    'page_size': 20
}
response = requests.post('http://localhost:8000/api/stocks/query', json=data)
print(response.json())

# 获取最新数据
response = requests.get('http://localhost:8000/api/stocks/latest?limit=50')
print(response.json())
```

### JavaScript 示例

```javascript
// 健康检查
fetch('http://localhost:8000/health')
  .then(response => response.json())
  .then(data => console.log(data));

// 查询股票
fetch('http://localhost:8000/api/stocks/query', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    code: '000001',
    start_date: '2024-01-01',
    end_date: '2024-12-31',
    page: 1,
    page_size: 20
  })
})
  .then(response => response.json())
  .then(data => console.log(data));

// 获取最新数据
fetch('http://localhost:8000/api/stocks/latest?limit=50')
  .then(response => response.json())
  .then(data => console.log(data));
```

### Go 示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type QueryRequest struct {
    Code      string `json:"code"`
    StartDate string `json:"start_date"`
    EndDate   string `json:"end_date"`
    Page      int    `json:"page"`
    PageSize  int    `json:"page_size"`
}

func main() {
    // 健康检查
    resp, _ := http.Get("http://localhost:8000/health")
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
    resp.Body.Close()

    // 查询股票
    req := QueryRequest{
        Code:      "000001",
        StartDate: "2024-01-01",
        EndDate:   "2024-12-31",
        Page:      1,
        PageSize:  20,
    }
    data, _ := json.Marshal(req)
    resp, _ = http.Post("http://localhost:8000/api/stocks/query", 
        "application/json", bytes.NewReader(data))
    body, _ = io.ReadAll(resp.Body)
    fmt.Println(string(body))
    resp.Body.Close()

    // 获取最新数据
    resp, _ = http.Get("http://localhost:8000/api/stocks/latest?limit=50")
    body, _ = io.ReadAll(resp.Body)
    fmt.Println(string(body))
    resp.Body.Close()
}
```

## 配置说明

HTTP 服务器监听地址在 `configs/config.yaml` 中配置：

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 60s
```

## 注意事项

1. **CORS**: 当前版本未启用 CORS，如需跨域访问请自行添加中间件
2. **认证**: 当前版本无认证机制，生产环境建议添加 API Token 认证
3. **限流**: 建议在生产环境添加限流中间件
4. **日志**: 所有请求都会记录在应用日志中
5. **性能**: 建议使用分页查询，避免一次性查询大量数据

## 后续优化

- [ ] 添加数据库查询实现
- [ ] 添加更多筛选条件
- [ ] 添加数据导出功能（CSV/Excel）
- [ ] 添加 WebSocket 实时推送
- [ ] 添加 API 认证和授权
- [ ] 添加请求限流
- [ ] 添加 CORS 支持
- [ ] 添加 API 文档页面（Swagger）
