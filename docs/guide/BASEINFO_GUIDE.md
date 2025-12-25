# BaseInfo 基础数据功能说明

## 概述

新增了基础数据采集和存储功能，用于保存股票的基础数据信息。该功能与原有的涨停股票数据采集功能并行运行。

## 数据库表结构

### 表名：`base_info`

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | BIGINT | 主键，自增 |
| stock_name | VARCHAR(100) | 股票简称 |
| latest_price | DECIMAL(10,2) | 最新价 |
| auction_change_rate | DECIMAL(10,2) | 集合竞价涨跌幅 |
| latest_change_rate | DECIMAL(10,2) | 最新涨跌幅 |
| auction_unmatched_amount_str | VARCHAR(50) | 竞价未匹配量(字符串) |
| morning_auction_amount_str | VARCHAR(50) | 早盘竞价成交额(字符串) |
| turnover_str | VARCHAR(50) | 成交额(字符串) |
| circulation_market_value | DECIMAL(20,2) | 流通市值 |
| stock_code | VARCHAR(20) | 股票代码 |
| trade_date | DATE | 交易日期 |
| market_code | VARCHAR(10) | 市场代码 |
| code | VARCHAR(50) | 代码 |
| turnover | DECIMAL(20,2) | 成交额 |
| morning_auction_amount | BIGINT | 早盘竞价成交额 |
| auction_unmatched_amount | BIGINT | 竞价未匹配量 |
| create_time | DATETIME | 创建时间 |
| update_time | DATETIME | 更新时间 |
| company_highlights | TEXT | 公司亮点 |
| industry_category | VARCHAR(100) | 所属行业 |
| concept_theme | VARCHAR(500) | 概念题材 |
| consecutive_limit_days | INT | 连续涨停天数 |

### 索引

- **主键索引**: `id`
- **复合索引**: `idx_stock_date` (stock_code, trade_date) - 用于按股票代码和日期查询
- **日期索引**: `idx_trade_date` (trade_date) - 用于按日期查询
- **唯一索引**: `uk_code_date` (code, trade_date) - 确保同一股票同一天只有一条记录

## 配置说明

在 `configs/config.yaml` 中添加了 `baseinfo` 配置段：

```yaml
baseinfo:
  query: "竞价未匹配金额；竞价金额；竞价涨幅；涨幅；成交金额；流通市值；连板天数；公司亮点；行业分类；概念题材；"
  cookie: "your_cookie_here"
```

- `query`: 问财查询条件，定义需要获取的数据字段
- `cookie`: 同花顺问财的 Cookie，用于 API 认证

## 定时任务配置

系统现在有两个独立的定时任务：

### 1. 股票涨停数据任务
- **执行时间**: 每天 9:00 (由 `scheduler.cron` 配置)
- **功能**: 查询并保存涨停股票数据到 `zp_jj` 表

### 2. 基础数据任务
- **执行时间**: 每天 9:15
- **功能**: 查询并保存股票基础数据到 `base_info` 表

> **注意**: 两个任务间隔 15 分钟，避免请求过于密集。

## 代码架构

项目遵循 DDD (领域驱动设计) 分层架构：

### 1. Biz 层 (业务层)
**文件**: `internal/biz/baseinfo.go`
- `BaseInfo`: 基础数据业务模型
- `BaseInfoRepo`: 仓库接口定义
- `BaseInfoUsecase`: 业务用例，包含保存逻辑

### 2. Data 层 (数据层)
**文件**: `internal/data/baseinfo.go`
- `baseInfoRepo`: 仓库接口实现
- `Save()`: 保存单条记录
- `Update()`: 更新记录
- `FindByCodeAndDate()`: 按代码和日期查询
- `BatchSave()`: 批量保存，使用 MySQL 的 `ON DUPLICATE KEY UPDATE` 实现 upsert

### 3. Service 层 (服务层)
**文件**: `internal/service/baseinfo.go`
- `BaseInfoService`: 基础数据服务
- `FetchAndSaveBaseInfo()`: 从问财 API 获取数据并保存
- `parseResult()`: 解析问财 API 返回的数据

### 4. Conf 层 (配置层)
**文件**: `internal/conf/conf.go`
- `BaseInfo`: 基础数据配置结构

### 5. 主程序
**文件**: `cmd/goweicai/main.go`
- 初始化所有组件
- 配置两个定时任务
- 服务器启动和优雅关闭

## 数据流程

```
1. 定时任务触发 (每天 9:15)
   ↓
2. BaseInfoService.FetchAndSaveBaseInfo()
   ↓
3. 调用 gowencai.Get() 查询问财 API
   ↓
4. parseResult() 解析返回数据
   ↓
5. BaseInfoUsecase.SaveBaseInfos()
   ↓
6. BaseInfoRepo.BatchSave()
   ↓
7. 使用 ON DUPLICATE KEY UPDATE 批量插入/更新
   ↓
8. 数据保存到 base_info 表
```

## Upsert 实现

使用 MySQL 的 `ON DUPLICATE KEY UPDATE` 语法实现 upsert（插入或更新）：

```sql
INSERT INTO base_info (...) VALUES (...)
ON DUPLICATE KEY UPDATE
    stock_name = VALUES(stock_name),
    latest_price = VALUES(latest_price),
    ...
```

- 基于唯一键 `uk_code_date` (code, trade_date)
- 如果记录不存在，则插入
- 如果记录已存在（相同的 code 和 trade_date），则更新

## 使用说明

### 1. 创建数据库表

```bash
mysql -u root -p wc < migrations/001_create_base_info.sql
```

### 2. 编译程序

```bash
cd /home/administrator/workplace/gowencai
go build -o cmd/goweicai/goweicai cmd/goweicai/main.go
```

### 3. 运行程序

```bash
cd cmd/goweicai
./goweicai
```

或使用 Makefile：

```bash
make run
```

### 4. 查看数据

```sql
-- 查看今天的基础数据
SELECT * FROM base_info 
WHERE trade_date = CURDATE() 
ORDER BY create_time DESC 
LIMIT 10;

-- 查看某只股票的历史数据
SELECT * FROM base_info 
WHERE stock_code = '000001' 
ORDER BY trade_date DESC;

-- 统计每天的数据量
SELECT trade_date, COUNT(*) as count 
FROM base_info 
GROUP BY trade_date 
ORDER BY trade_date DESC;
```

## 日志示例

程序运行时会输出详细日志：

```
INFO ts=2025-01-19T09:15:00+08:00 msg=开始执行基础数据定时任务...
INFO ts=2025-01-19T09:15:00+08:00 module=service/baseinfo msg=开始查询基础数据...
INFO ts=2025-01-19T09:15:02+08:00 module=service/baseinfo msg=查询到 4500 条基础数据
INFO ts=2025-01-19T09:15:03+08:00 module=data/baseinfo msg=BatchSave success, count: 4500
INFO ts=2025-01-19T09:15:03+08:00 module=service/baseinfo msg=基础数据保存成功
INFO ts=2025-01-19T09:15:03+08:00 msg=基础数据定时任务执行成功
```

## 配置修改

如需修改查询条件或执行时间，请编辑 `configs/config.yaml`:

```yaml
baseinfo:
  # 修改查询条件以获取不同的数据字段
  query: "你的查询条件；"
  cookie: "your_cookie"
```

如需修改基础数据任务的执行时间，请修改 `cmd/goweicai/main.go` 中的：

```go
baseInfoCron := "0 15 9 * * *"  // 修改为你想要的时间
```

## 注意事项

1. **Cookie 有效期**: 同花顺问财的 Cookie 可能会过期，需要定期更新
2. **请求频率**: 两个任务间隔执行，避免被服务器限流
3. **数据去重**: 通过唯一索引确保同一天同一股票只有一条记录
4. **启动执行**: `scheduler.run_on_start: true` 时，程序启动会立即执行两个任务
5. **数据库连接**: 确保数据库连接配置正确，程序会自动重连

## 未来扩展

可以考虑的功能扩展：

1. 添加 HTTP API 查询基础数据
2. 添加 gRPC 接口支持
3. 数据分析和统计功能
4. 数据导出功能（CSV, Excel）
5. 实时数据更新（Websocket）
6. 数据可视化（图表展示）
