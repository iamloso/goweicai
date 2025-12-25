# GoWencai API 使用示例

本文档提供了使用 GoWencai HTTP API 的实际示例代码。

## 启动服务

```bash
# 编译项目
make build

# 后台启动
make daemon

# 查看状态
make status

# 查看日志
make logs
```

## cURL 示例

### 1. 健康检查

```bash
curl http://localhost:8000/health
```

**响应**:
```json
{"status":"ok","time":"2024-12-24T10:00:00+08:00"}
```

### 2. 获取最新股票（最多 10 条）

```bash
curl "http://localhost:8000/api/stocks/latest?limit=10"
```

### 3. 查询特定股票（POST）

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

### 4. 查询特定股票（GET）

```bash
curl "http://localhost:8000/api/stocks/query?code=000001&start_date=2024-01-01&end_date=2024-12-31&page=1&page_size=20"
```

## Python 示例

```python
#!/usr/bin/env python3
import requests
import json

BASE_URL = "http://localhost:8000"

# 1. 健康检查
def health_check():
    response = requests.get(f"{BASE_URL}/health")
    print("健康检查:", response.json())

# 2. 获取最新股票
def get_latest_stocks(limit=10):
    response = requests.get(f"{BASE_URL}/api/stocks/latest", params={"limit": limit})
    data = response.json()
    print(f"获取最新 {limit} 条股票:")
    print(json.dumps(data, indent=2, ensure_ascii=False))

# 3. 查询股票（POST）
def query_stocks(code=None, start_date=None, end_date=None, page=1, page_size=20):
    payload = {
        "page": page,
        "page_size": page_size
    }
    if code:
        payload["code"] = code
    if start_date:
        payload["start_date"] = start_date
    if end_date:
        payload["end_date"] = end_date
    
    response = requests.post(f"{BASE_URL}/api/stocks/query", json=payload)
    data = response.json()
    print("查询结果:")
    print(json.dumps(data, indent=2, ensure_ascii=False))

# 4. 查询股票（GET）
def query_stocks_get(code=None, page=1, page_size=20):
    params = {"page": page, "page_size": page_size}
    if code:
        params["code"] = code
    
    response = requests.get(f"{BASE_URL}/api/stocks/query", params=params)
    data = response.json()
    print("查询结果:")
    print(json.dumps(data, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    print("=== GoWencai API 测试 ===\n")
    
    # 测试健康检查
    health_check()
    print()
    
    # 获取最新股票
    get_latest_stocks(5)
    print()
    
    # 查询特定股票
    query_stocks(code="000001", page=1, page_size=10)
    print()
    
    # GET 方式查询
    query_stocks_get(code="000001", page=1, page_size=10)
```

保存为 `test_api.py` 并运行：

```bash
python3 test_api.py
```

## JavaScript/Node.js 示例

```javascript
const axios = require('axios');

const BASE_URL = 'http://localhost:8000';

// 1. 健康检查
async function healthCheck() {
  try {
    const response = await axios.get(`${BASE_URL}/health`);
    console.log('健康检查:', response.data);
  } catch (error) {
    console.error('错误:', error.message);
  }
}

// 2. 获取最新股票
async function getLatestStocks(limit = 10) {
  try {
    const response = await axios.get(`${BASE_URL}/api/stocks/latest`, {
      params: { limit }
    });
    console.log(`获取最新 ${limit} 条股票:`, JSON.stringify(response.data, null, 2));
  } catch (error) {
    console.error('错误:', error.message);
  }
}

// 3. 查询股票（POST）
async function queryStocks(options = {}) {
  const { code, startDate, endDate, page = 1, pageSize = 20 } = options;
  
  const payload = { page, page_size: pageSize };
  if (code) payload.code = code;
  if (startDate) payload.start_date = startDate;
  if (endDate) payload.end_date = endDate;
  
  try {
    const response = await axios.post(`${BASE_URL}/api/stocks/query`, payload);
    console.log('查询结果:', JSON.stringify(response.data, null, 2));
  } catch (error) {
    console.error('错误:', error.message);
  }
}

// 4. 查询股票（GET）
async function queryStocksGet(options = {}) {
  const { code, page = 1, pageSize = 20 } = options;
  
  const params = { page, page_size: pageSize };
  if (code) params.code = code;
  
  try {
    const response = await axios.get(`${BASE_URL}/api/stocks/query`, { params });
    console.log('查询结果:', JSON.stringify(response.data, null, 2));
  } catch (error) {
    console.error('错误:', error.message);
  }
}

// 主函数
async function main() {
  console.log('=== GoWencai API 测试 ===\n');
  
  await healthCheck();
  console.log();
  
  await getLatestStocks(5);
  console.log();
  
  await queryStocks({ code: '000001', page: 1, pageSize: 10 });
  console.log();
  
  await queryStocksGet({ code: '000001', page: 1, pageSize: 10 });
}

main();
```

保存为 `test_api.js`，安装依赖并运行：

```bash
npm install axios
node test_api.js
```

## Go 示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const baseURL = "http://localhost:8000"

type QueryRequest struct {
    Code      string `json:"code,omitempty"`
    StartDate string `json:"start_date,omitempty"`
    EndDate   string `json:"end_date,omitempty"`
    Page      int    `json:"page"`
    PageSize  int    `json:"page_size"`
}

type Response struct {
    Code    int             `json:"code"`
    Message string          `json:"message"`
    Data    json.RawMessage `json:"data,omitempty"`
    Error   string          `json:"error,omitempty"`
}

// 1. 健康检查
func healthCheck() error {
    resp, err := http.Get(baseURL + "/health")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println("健康检查:", string(body))
    return nil
}

// 2. 获取最新股票
func getLatestStocks(limit int) error {
    url := fmt.Sprintf("%s/api/stocks/latest?limit=%d", baseURL, limit)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var result Response
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }

    fmt.Printf("获取最新 %d 条股票:\n", limit)
    fmt.Printf("%+v\n", result)
    return nil
}

// 3. 查询股票（POST）
func queryStocks(req QueryRequest) error {
    data, _ := json.Marshal(req)
    resp, err := http.Post(
        baseURL+"/api/stocks/query",
        "application/json",
        bytes.NewReader(data),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var result Response
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }

    fmt.Println("查询结果:")
    fmt.Printf("%+v\n", result)
    return nil
}

// 4. 查询股票（GET）
func queryStocksGet(code string, page, pageSize int) error {
    url := fmt.Sprintf("%s/api/stocks/query?code=%s&page=%d&page_size=%d",
        baseURL, code, page, pageSize)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var result Response
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return err
    }

    fmt.Println("查询结果:")
    fmt.Printf("%+v\n", result)
    return nil
}

func main() {
    fmt.Println("=== GoWencai API 测试 ===\n")

    // 健康检查
    if err := healthCheck(); err != nil {
        fmt.Println("错误:", err)
        return
    }
    fmt.Println()

    // 获取最新股票
    if err := getLatestStocks(5); err != nil {
        fmt.Println("错误:", err)
        return
    }
    fmt.Println()

    // 查询特定股票
    req := QueryRequest{
        Code:     "000001",
        Page:     1,
        PageSize: 10,
    }
    if err := queryStocks(req); err != nil {
        fmt.Println("错误:", err)
        return
    }
    fmt.Println()

    // GET 方式查询
    if err := queryStocksGet("000001", 1, 10); err != nil {
        fmt.Println("错误:", err)
        return
    }
}
```

保存为 `test_api_example.go` 并运行：

```bash
go run test_api_example.go
```

## 测试脚本

项目提供了自动化测试脚本：

```bash
./test_api.sh
```

该脚本会自动测试所有 API 接口并显示结果。

## 注意事项

1. **空数据**: 当前 API 返回空数据是正常的，因为查询逻辑尚未实现（标记为 TODO）
2. **端口配置**: 默认端口为 8000，可在 `configs/config.yaml` 中修改
3. **错误处理**: 所有接口都有完善的错误处理，返回标准格式的错误信息
4. **分页**: 查询接口支持分页，默认每页 20 条，最大 100 条

## 更多信息

- 完整 API 文档: [API_DOCS.md](./API_DOCS.md)
- 项目 README: [README.md](./README.md)
- 定时任务指南: [SCHEDULER_GUIDE.md](./SCHEDULER_GUIDE.md)
