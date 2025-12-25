# GoWeicai - Kratos 版本

基于 Kratos 框架重构的问财股票数据采集工具。

## 项目结构

```
.
├── cmd/
│   └── goweicai/          # 主程序入口
│       └── main.go
├── configs/               # 配置文件
│   └── config.yaml
├── internal/              # 内部代码
│   ├── biz/              # 业务逻辑层
│   │   └── stock.go
│   ├── conf/             # 配置定义
│   │   └── conf.go
│   ├── data/             # 数据访问层
│   │   └── stock.go
│   └── service/          # 服务层
│       └── wencai.go
├── convert.go            # 类型转换工具
├── headers.go            # 请求头生成
├── types.go              # 类型定义
└── wencai.go             # 问财API客户端

```

## 架构说明

项目采用标准的 Kratos DDD（领域驱动设计）分层架构：

### 1. Service 层（服务层）
- 位置：`internal/service/`
- 职责：对外提供服务接口，组织业务逻辑调用
- 文件：`wencai.go` - 问财数据获取服务

### 2. Biz 层（业务逻辑层）
- 位置：`internal/biz/`
- 职责：核心业务逻辑，定义业务模型和用例
- 文件：`stock.go` - 股票业务模型和用例

### 3. Data 层（数据访问层）
- 位置：`internal/data/`
- 职责：数据持久化，实现仓库接口
- 文件：`stock.go` - 股票数据仓库实现

### 4. Conf 层（配置层）
- 位置：`internal/conf/`
- 职责：配置结构定义
- 文件：`conf.go` - 配置模型

## 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 60s

data:
  database:
    driver: mysql
    source: root:password@tcp(localhost:3306)/wc?charset=utf8mb4&parseTime=True&loc=Local

wencai:
  query: "竞价未匹配金额；竞价金额；竞价涨幅；涨幅；成交金额；流通市值；连板天数；不含ST; 涨停原因；公司亮点；行业分类；概念题材； 封单金额；"
  cookie: "your_cookie_here"
```

## 编译运行

### 编译

```bash
cd cmd/goweicai
go build -o goweicai
```

### 运行

```bash
# 使用默认配置文件
./goweicai

# 指定配置文件
./goweicai -conf /path/to/config.yaml
```

## 数据库表结构

```sql
CREATE TABLE `zp_jj` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL COMMENT '股票代码',
  `market_code` varchar(10) DEFAULT NULL COMMENT '市场代码',
  `stock_code` varchar(20) DEFAULT NULL COMMENT '完整股票代码',
  `stock_name` varchar(50) DEFAULT NULL COMMENT '股票名称',
  `latest_price` decimal(10,2) DEFAULT NULL COMMENT '最新价',
  `latest_change_rate` decimal(10,2) DEFAULT NULL COMMENT '涨跌幅',
  `auction_unmatched_amount` bigint DEFAULT '0' COMMENT '竞价未匹配金额(元)',
  `auction_unmatched_amount_str` varchar(50) DEFAULT '' COMMENT '竞价未匹配金额(万/亿)',
  `auction_unmatched_amount_rank` varchar(50) DEFAULT '' COMMENT '竞价未匹配金额排名',
  `auction_unmatched_amount_rank_num` int DEFAULT '0' COMMENT '竞价未匹配金额排名数字',
  `auction_change_rate` decimal(10,2) DEFAULT NULL COMMENT '竞价涨幅',
  `morning_auction_amount` bigint DEFAULT '0' COMMENT '早盘竞价金额(元)',
  `morning_auction_amount_str` varchar(50) DEFAULT '' COMMENT '早盘竞价金额(万/亿)',
  `turnover` decimal(20,2) DEFAULT NULL COMMENT '成交额',
  `turnover_str` varchar(50) DEFAULT '' COMMENT '成交额(万/亿)',
  `circulation_market_value` decimal(20,2) DEFAULT NULL COMMENT '流通市值',
  `limit_up_reason` text COMMENT '涨停原因',
  `company_highlights` text COMMENT '公司亮点',
  `industry_category` varchar(100) DEFAULT '' COMMENT '行业分类',
  `concept_theme` text COMMENT '概念题材',
  `limit_up_seal_amount` bigint DEFAULT '0' COMMENT '涨停封单金额(元)',
  `limit_up_seal_amount_str` varchar(50) DEFAULT '' COMMENT '涨停封单金额(万/亿)',
  `consecutive_limit_days` int DEFAULT '0' COMMENT '连板天数',
  `trade_date` date NOT NULL COMMENT '交易日期',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code_date` (`code`,`trade_date`),
  KEY `idx_consecutive_days` (`consecutive_limit_days`),
  KEY `idx_industry` (`industry_category`),
  KEY `idx_trade_date` (`trade_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='早盘竞价数据';
```

## 功能特性

1. **数据采集**
   - 从问财网获取股票数据
   - 支持动态日期字段解析
   - 自动处理科学计数法格式

2. **数据处理**
   - 自动格式化金额为万/亿单位
   - 解析排名数字（如 "64/5458" → 64）
   - 支持去重和增量更新

3. **数据存储**
   - MySQL 数据库持久化
   - 事务支持
   - 自动判断插入或更新

4. **日志记录**
   - 结构化日志输出
   - 操作详情跟踪
   - 错误信息记录

## 依赖管理

主要依赖：

- `github.com/go-kratos/kratos/v2` - Kratos 微服务框架
- `github.com/go-sql-driver/mysql` - MySQL 驱动
- `gopkg.in/yaml.v3` - YAML 配置解析

## 开发说明

### 添加新功能

1. 在 `biz` 层定义业务模型和接口
2. 在 `data` 层实现数据访问
3. 在 `service` 层实现业务逻辑
4. 在 `main.go` 中组装依赖

### 修改配置

编辑 `configs/config.yaml`，重启程序即可生效。

## License

MIT License
