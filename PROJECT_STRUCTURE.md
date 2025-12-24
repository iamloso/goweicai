# 项目结构说明

```
gowencai/
├── cmd/                          # 应用程序入口
│   └── goweicai/
│       ├── main.go              # 主程序（定时任务调度）
│       ├── goweicai             # 编译后的可执行文件
│       └── goweicai.log         # 运行日志（后台模式）
│
├── internal/                     # 内部代码（DDD 分层架构）
│   ├── biz/                     # 业务逻辑层
│   │   └── stock.go             # 股票业务模型和用例
│   ├── data/                    # 数据访问层
│   │   └── stock.go             # 股票数据仓库实现
│   ├── service/                 # 服务层
│   │   └── wencai.go            # 问财服务（调用 API）
│   └── conf/                    # 配置结构
│       └── conf.go              # 配置模型定义
│
├── configs/                      # 配置文件
│   └── config.yaml              # 主配置文件
│       ├── scheduler            # 定时任务配置
│       ├── data                 # 数据库配置
│       └── wencai               # 问财 API 配置
│
├── pywencai/                     # Python 原始实现（参考）
│   ├── hexin-v.js               # 同花顺 token 生成脚本
│   └── ...
│
├── example/                      # 示例代码
│   └── main.go                  # 库模式使用示例
│
├── goweicai.sh                  # 服务管理脚本
├── Makefile                     # 构建工具
│
├── go.mod                       # Go 模块定义
├── go.sum                       # 依赖版本锁定
│
├── wencai.go                    # 问财 API 封装（库模式）
├── types.go                     # 公共类型定义
├── headers.go                   # HTTP 请求头
├── convert.go                   # 数据转换工具
│
└── 文档/
    ├── README.md                # 项目主文档
    ├── SCHEDULER_GUIDE.md       # 定时任务使用指南
    ├── SCHEDULER_REFACTOR.md    # 定时任务改造总结
    ├── README_KRATOS.md         # Kratos 框架说明
    ├── REFACTOR_SUMMARY.md      # 重构总结
    ├── QUICKSTART.md            # 快速开始
    ├── MIGRATION_GUIDE.md       # 迁移指南
    └── LICENSE                  # 开源协议
```

## 核心组件说明

### 1. 主程序 (cmd/goweicai/main.go)

**职责**：
- 加载配置文件
- 初始化依赖（数据库、服务等）
- 创建定时任务调度器
- 监听系统退出信号
- 优雅关闭

**关键代码**：
```go
// 创建 Cron 调度器
c := cron.New(cron.WithSeconds())

// 添加定时任务
c.AddFunc(cronExpr, job)

// 启动调度器
c.Start()

// 等待退出信号
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
```

### 2. 业务逻辑层 (internal/biz/)

**职责**：
- 定义业务模型（Stock）
- 定义业务接口（StockRepo）
- 实现业务用例（StockUsecase）

**特点**：
- 不依赖具体实现（依赖倒置）
- 纯业务逻辑，无技术细节

### 3. 数据访问层 (internal/data/)

**职责**：
- 实现 StockRepo 接口
- 数据库 CRUD 操作
- 数据格式化（金额、排名等）

**关键功能**：
- `Save()`: 保存单条记录
- `Update()`: 更新记录
- `FindByCodeAndDate()`: 查询记录
- `BatchSave()`: 批量保存（去重）
- `FormatAmount()`: 金额格式化
- `ParseRankNumber()`: 排名解析

### 4. 服务层 (internal/service/)

**职责**：
- 调用问财 API 获取数据
- 解析 API 返回结果
- 调用业务层保存数据

**流程**：
```
API 调用 → 数据解析 → 字段映射 → 业务层保存
```

### 5. 配置管理 (configs/config.yaml)

**结构**：
```yaml
scheduler:      # 定时任务配置
  cron: "..."   # Cron 表达式
  run_on_start: true

data:           # 数据库配置
  database:
    driver: "mysql"
    source: "..."

wencai:         # 问财 API 配置
  query: "..."  # 查询语句
  cookie: "..." # 认证 Cookie
```

## 依赖关系

```
main.go
  ↓
service (WencaiService)
  ↓
biz (StockUsecase) ← 依赖接口
  ↓
data (stockRepo)   ← 实现接口
  ↓
MySQL 数据库
```

## 数据流向

```
1. 定时触发
   ↓
2. WencaiService.FetchAndSaveStocks()
   ↓
3. 调用问财 API (wencai.Get)
   ↓
4. 解析返回数据
   ↓
5. StockUsecase.SaveStocks()
   ↓
6. stockRepo.BatchSave()
   ↓
7. 写入 MySQL
```

## 技术栈

| 层级 | 技术选型 |
|------|---------|
| 框架 | Kratos v2.9.2 |
| 调度 | robfig/cron/v3 |
| 数据库 | MySQL |
| 驱动 | go-sql-driver/mysql |
| 配置 | YAML (gopkg.in/yaml.v3) |
| 日志 | Kratos log |

## 设计模式

1. **依赖注入 (DI)**
   - 通过构造函数注入依赖
   - 降低耦合度

2. **仓库模式 (Repository)**
   - 抽象数据访问层
   - 易于测试和替换

3. **分层架构 (Layered Architecture)**
   - 职责分离
   - 易于维护

4. **策略模式 (Strategy)**
   - 数据格式化可扩展

## 配置说明

### Scheduler 配置

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| cron | string | Cron 表达式 | `"0 0 9 * * *"` |
| run_on_start | bool | 启动时执行 | `true` |

### Data 配置

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| driver | string | 数据库驱动 | `"mysql"` |
| source | string | 连接字符串 | `"user:pass@tcp(...)"` |

### Wencai 配置

| 字段 | 类型 | 说明 | 必填 |
|------|------|------|------|
| query | string | 查询语句 | ✅ |
| cookie | string | 认证 Cookie | ✅ |

## 运行模式

### 开发模式
```bash
make run
# 前台运行，Ctrl+C 退出
```

### 生产模式
```bash
make daemon
# 后台运行，日志输出到文件
```

### 脚本管理
```bash
./goweicai.sh start    # 启动
./goweicai.sh stop     # 停止
./goweicai.sh restart  # 重启
./goweicai.sh status   # 状态
./goweicai.sh logs -f  # 查看日志
```

## 扩展指南

### 添加新的数据字段

1. 修改 `internal/biz/stock.go` - 添加字段到 Stock 结构体
2. 修改 `internal/service/wencai.go` - 添加字段映射
3. 修改数据库表结构
4. 更新 `internal/data/stock.go` - 调整 SQL 语句

### 添加新的数据源

1. 在 `internal/service/` 创建新服务
2. 在 `internal/biz/` 定义业务接口
3. 在 `internal/data/` 实现数据仓库
4. 在 `main.go` 中注入依赖

### 自定义调度策略

1. 修改 `configs/config.yaml` 中的 cron 表达式
2. 或在代码中动态添加多个任务：
   ```go
   c.AddFunc("0 0 9 * * *", job1)
   c.AddFunc("0 0 15 * * *", job2)
   ```

## 最佳实践

1. **配置管理**
   - 敏感信息使用环境变量
   - 不同环境使用不同配置文件

2. **日志管理**
   - 使用日志轮转（logrotate）
   - 定期清理旧日志

3. **监控告警**
   - 添加 Prometheus metrics
   - 配置任务失败告警

4. **高可用**
   - 使用 systemd 管理进程
   - 配置自动重启
   - 多实例部署需加分布式锁

5. **性能优化**
   - 批量操作使用事务
   - 合理设置连接池
   - 定期分析慢查询
