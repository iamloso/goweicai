# GoWencai

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Kratos](https://img.shields.io/badge/Framework-Kratos%20v2.9-brightgreen)](https://go-kratos.dev/)
[![Cron](https://img.shields.io/badge/Scheduler-Cron%20v3-orange)](https://github.com/robfig/cron)

Go语言实现的同花顺问财数据获取工具，从 [pywencai](https://github.com/zsrl/pywencai) 移植而来。

⚠️ **注意**：由于问财登录策略调整，目前**必须提供cookie参数**才能使用。

## ✨ 特性

- 🚀 基于 Kratos 微服务框架的 DDD 分层架构
- ⏰ 支持 Cron 表达式的定时任务调度
- 💾 MySQL 数据持久化存储
- 🔧 YAML 配置文件，灵活易用
- 📊 完整的股票数据字段支持
- 🛡️ 优雅退出和错误处理

## 声明

1. gowencai为开源社区开发，并非同花顺官方提供的工具。
2. 该工具只是效率工具，为了便于通过Go获取问财数据，用于量化研究和学习，其原理与登录网页获取数据方式一致。
3. 建议低频使用，反对高频调用，高频调用会被问财屏蔽，请自行评估技术和法律风险。
4. 项目代码遵循MIT开源协议，但不赞成商用，商用请自行评估法律风险。
5. 感谢问财提供免费接口和数据分享。

## 环境依赖

**需要 Node.js v12+**

虽然原计划使用纯 Go 实现，但由于 JavaScript 代码的复杂性（使用了异步生成器等特性），当前版本需要 Node.js 来生成 token。

### 安装 Node.js

#### Ubuntu/Debian
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### macOS
```bash
brew install node
```

#### Windows
下载安装：https://nodejs.org/

详细说明请参考 [NODEJS_REQUIRED.md](./NODEJS_REQUIRED.md)

## 🚀 快速开始（定时任务模式）

### 1. 配置文件

编辑 `configs/config.yaml`：

```yaml
scheduler:
  # Cron 表达式：每天 9:00 执行
  cron: "0 0 9 * * *"
  # 启动时立即执行一次
  run_on_start: true

data:
  database:
    driver: mysql
    source: root:password@tcp(localhost:3306)/wc?charset=utf8mb4&parseTime=True&loc=Local

wencai:
  query: "竞价未匹配金额；竞价金额；竞价涨幅；涨幅；成交金额；流通市值；连板天数；不含ST"
  cookie: "your_cookie_here"  # 必填
```

### 2. 编译并运行

```bash
# 使用 Makefile
make build
make run          # 前台运行（开发测试）
make daemon       # 后台运行（生产环境）

# 或使用管理脚本
./goweicai.sh build
./goweicai.sh start    # 启动服务
./goweicai.sh status   # 查看状态
./goweicai.sh logs -f  # 实时查看日志
./goweicai.sh stop     # 停止服务
```

### 3. 常用 Cron 表达式

| 需求 | Cron 表达式 |
|------|-------------|
| 每天 9:00 | `0 0 9 * * *` |
| 每 30 分钟 | `0 */30 * * * *` |
| 每天 9:00 和 15:00 | `0 0 9,15 * * *` |
| 工作日 9:00 | `0 0 9 * * 1-5` |
| 每 5 秒（测试） | `*/5 * * * * *` |

📚 **详细文档**：
- [定时任务使用指南](./SCHEDULER_GUIDE.md)
- [定时任务改造总结](./SCHEDULER_REFACTOR.md)
- [Kratos 框架说明](./README_KRATOS.md)

---

## 💡 库模式使用（编程调用）

### 安装

```bash
go get github.com/iamloso/goweicai
```

### 示例代码

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/fenghuang/gowencai/gowencai"
)

func main() {
    result, err := gowencai.Get(&gowencai.QueryOptions{
        Query:     "退市股票",
        SortKey:   "退市@退市日期",
        SortOrder: "asc",
        Cookie:    "your_cookie_here", // 必填
        Log:       true,
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("%+v\n", result)
}
```

## API 文档

### Get(opts *QueryOptions) (interface{}, error)

根据问财语句查询结果。

#### QueryOptions 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|-----|------|------|--------|------|
| Query | string | ✅ | - | 查询问句 |
| Cookie | string | ✅ | - | Cookie值（获取方法见下文） |
| SortKey | string | ❌ | - | 排序字段，值为返回结果的列名 |
| SortOrder | string | ❌ | - | 排序规则：`asc`(升序) 或 `desc`(降序) |
| Page | int | ❌ | 1 | 查询的页号 |
| PerPage | int | ❌ | 100 | 每页数据条数，最大100 |
| Loop | interface{} | ❌ | false | 是否循环分页：`false`/`true`/数字 |
| QueryType | string | ❌ | stock | 查询类型，见下表 |
| Retry | int | ❌ | 10 | 请求失败重试次数 |
| Sleep | int | ❌ | 0 | 循环请求间隔秒数 |
| Log | bool | ❌ | false | 是否打印日志 |
| Pro | bool | ❌ | false | 是否使用付费版 |
| NoDetail | bool | ❌ | false | 详情类问题返回nil而非字典 |
| Find | []string | ❌ | nil | 指定股票代码列表 |
| UserAgent | string | ❌ | 随机 | 自定义User-Agent |

#### QueryType 取值

| 值 | 含义 |
|----|------|
| stock | 股票 |
| zhishu | 指数 |
| fund | 基金 |
| hkstock | 港股 |
| usstock | 美股 |
| threeboard | 新三板 |
| conbond | 可转债 |
| insurance | 保险 |
| futures | 期货 |
| lccp | 理财 |
| foreign_exchange | 外汇 |

### Cookie 获取方法

1. 打开浏览器访问 [问财网站](http://www.iwencai.com)
2. 登录你的账号
3. 打开浏览器开发者工具（F12）
4. 切换到 Network 标签
5. 在问财进行任意查询
6. 找到请求，复制请求头中的 `Cookie` 字段值

![cookie获取示例](./cookie.png)

## 使用示例

### 示例1：基本查询

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:  "退市股票",
    Cookie: "your_cookie_here",
})
```

### 示例2：带排序的查询

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:     "退市股票",
    SortKey:   "退市@退市日期",
    SortOrder: "asc",
    Cookie:    "your_cookie_here",
})
```

### 示例3：分页查询

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:   "昨日涨幅",
    Page:    2,
    PerPage: 50,
    Cookie:  "your_cookie_here",
})
```

### 示例4：循环分页获取所有数据

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:  "昨日涨幅",
    Loop:   true, // 获取所有页
    Cookie: "your_cookie_here",
})
```

### 示例5：限制循环页数

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:  "昨日涨幅",
    Loop:   3, // 只获取3页
    Cookie: "your_cookie_here",
})
```

### 示例6：查询不同类型数据

```go
// 查询基金
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:     "基金规模大于10亿",
    QueryType: "fund",
    Cookie:    "your_cookie_here",
})

// 查询港股
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:     "恒生指数成分股",
    QueryType: "hkstock",
    Cookie:    "your_cookie_here",
})
```

### 示例7：指定股票代码查询

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:  "最新价",
    Find:   []string{"600519", "000001"},
    Cookie: "your_cookie_here",
})
```

### 示例8：使用Client进行多次查询

```go
client := gowencai.NewClient()
client.SetLogger(log.New(os.Stdout, "[gowencai] ", log.LstdFlags))

result1, err := client.Get(&gowencai.QueryOptions{
    Query:  "昨日涨停",
    Cookie: "your_cookie_here",
})

result2, err := client.Get(&gowencai.QueryOptions{
    Query:  "昨日跌停",
    Cookie: "your_cookie_here",
})
```

### 示例9：启用日志

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:  "退市股票",
    Cookie: "your_cookie_here",
    Log:    true, // 打印详细日志
})
```

### 示例10：设置请求间隔

```go
result, err := gowencai.Get(&gowencai.QueryOptions{
    Query:  "昨日涨幅",
    Loop:   true,
    Sleep:  1, // 每次请求间隔1秒
    Cookie: "your_cookie_here",
})
```

## 项目结构

```
gowencai/
├── gowencai/           # 主包目录
│   ├── types.go       # 类型定义
│   ├── headers.go     # HTTP头部和token生成
│   ├── convert.go     # 数据转换和解析
│   └── wencai.go      # 核心查询功能
├── example/           # 示例代码
│   └── main.go
├── pywencai/          # 原Python包（保留用于token生成）
│   └── hexin-v.bundle.js
├── go.mod
└── README.md
```

## 注意事项

1. **Cookie必填**：目前版本必须提供有效的Cookie才能使用
2. **低频使用**：建议合理控制请求频率，避免被封禁
3. **无需Node.js**：使用内置JS引擎，无需安装Node.js
4. **数据格式**：返回的数据格式可能是`[]map[string]interface{}`或`map[string]interface{}`，取决于查询类型
5. **PerPage限制**：每页最多返回100条数据，这是问财接口的限制

## 与Python版本的差异

1. **类型安全**：Go版本提供了更好的类型检查
2. **并发支持**：可以使用goroutine并发查询
3. **性能**：Go版本通常有更好的性能表现
4. **无需Node.js**：Python版本需要Node.js，Go版本内置JS引擎
5. **无依赖pandas**：返回原始的map/slice数据结构，用户可根据需要自行处理

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License

## 致谢

- 原项目：[pywencai](https://github.com/zsrl/pywencai)
- 感谢同花顺问财提供的数据接口
