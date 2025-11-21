# 数据库配置说明

## 使用步骤

### 1. 创建数据库表

```sql
CREATE TABLE zp_jj (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '自增主键',
    code VARCHAR(20) NOT NULL COMMENT '股票代码（不含后缀）',
    market_code VARCHAR(10) COMMENT '市场代码',
    latest_price DECIMAL(10, 2) COMMENT '最新价',
    latest_change_rate DECIMAL(10, 2) COMMENT '最新涨跌幅(%)',
    auction_unmatched_amount BIGINT COMMENT '竞价未匹配金额（原始值）',
    auction_unmatched_amount_str VARCHAR(50) COMMENT '竞价未匹配金额（带单位，如：284.11万、2.84亿）',
    auction_unmatched_amount_rank VARCHAR(50) COMMENT '竞价未匹配金额排名',
    auction_change_rate DECIMAL(10, 3) COMMENT '竞价涨幅(%)',
    morning_auction_amount BIGINT COMMENT '早盘竞价金额',
    turnover DECIMAL(20, 2) COMMENT '成交额',
    circulation_market_value DECIMAL(20, 2) COMMENT '流通市值',
    stock_code VARCHAR(20) NOT NULL COMMENT '股票代码（含后缀）',
    stock_name VARCHAR(100) NOT NULL COMMENT '股票简称',
    trade_date DATE NOT NULL COMMENT '交易日期',
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_code (code),
    INDEX idx_stock_code (stock_code),
    INDEX idx_trade_date (trade_date),
    UNIQUE KEY uk_code_date (code, trade_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='竞价数据表';
```

### 2. 修改数据库连接配置

在 `main.go` 文件中修改以下常量：

```go
const (
    dbUser     = "root"              // 数据库用户名
    dbPassword = "your_password"     // 数据库密码
    dbHost     = "localhost"         // 数据库主机
    dbPort     = "3306"              // 数据库端口
    dbName     = "your_database"     // 数据库名称
)
```

### 3. 运行程序

```bash
cd /home/administrator/workplace/gowencai/cmd
go run main.go
```

## 功能说明

程序会：
1. 从问财API查询数据
2. 打印查询结果到控制台
3. 自动将结果保存到数据库
4. 如果数据已存在（相同股票代码和日期），则更新数据

## 数据字段映射

| API字段 | 数据库字段 |
|--------|-----------|
| code | code |
| market_code | market_code |
| 最新价 | latest_price |
| 最新涨跌幅 | latest_change_rate |
| 竞价未匹配金额[20251120] | auction_unmatched_amount |
| 竞价未匹配金额排名[20251120] | auction_unmatched_amount_rank |
| 竞价涨幅[20251120] | auction_change_rate |
| 股票代码 | stock_code |
| 股票简称 | stock_name |

## 注意事项

1. 确保MySQL服务正在运行
2. 数据库用户需要有相应的读写权限
3. 日期字段默认使用当前日期
4. 使用 `ON DUPLICATE KEY UPDATE` 避免重复插入
