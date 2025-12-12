# 更新说明

## 更新内容

### 1. 查询字段更新
- **query**: `"竞价未匹配金额；竞价金额；竞价涨幅；涨幅；成交金额；流通市值；连板天数；不含ST; 涨停原因；公司亮点；行业分类；概念题材；"`
- **query1**: query + `"涨停封单金额；"` (第二次查询包含涨停封单金额)

### 2. 新增字段（7个）

| 字段名 | 数据类型 | 说明 | 格式化字段 |
|--------|---------|------|-----------|
| limit_up_reason | TEXT | 涨停原因 | - |
| company_highlights | TEXT | 公司亮点 | - |
| industry_category | VARCHAR(100) | 行业分类 | - |
| concept_theme | TEXT | 概念题材 | - |
| limit_up_seal_amount | BIGINT | 涨停封单金额(元) | limit_up_seal_amount_str |
| consecutive_limit_days | INT | 连板天数 | - |

### 3. 更新已有字段映射

| 字段名 | API字段名 | 说明 |
|--------|----------|------|
| morning_auction_amount | 竞价金额 | 早盘竞价金额 |
| turnover | 成交金额 | 成交额 |
| circulation_market_value | 流通市值 | 流通市值 |

### 4. 执行流程

程序会执行两次查询：

1. **第一次查询** (query): 获取基础字段，不包含涨停封单金额
2. **第二次查询** (query1): 获取包含涨停封单金额的完整数据

两次查询的结果都会写入数据库，第二次查询会更新第一次查询的数据。

## 数据库迁移

执行以下命令添加新字段：

```bash
mysql -uroot -peauDx15FxO83lS wc < /home/administrator/workplace/gowencai/cmd/add_new_fields.sql
```

## 运行程序

```bash
cd /home/administrator/workplace/gowencai/cmd
go run main.go
```

## 字段格式化

以下字段会自动格式化为"万/亿"单位：

- `auction_unmatched_amount_str` - 竞价未匹配金额
- `morning_auction_amount_str` - 竞价金额
- `turnover_str` - 成交金额
- `limit_up_seal_amount_str` - 涨停封单金额

格式化规则：
- >= 1亿: `X.XX亿`
- >= 1万: `X.XX万`
- < 1万: 原始数值

## 日志输出

插入/更新记录时会输出连板天数和排名信息：
```
插入新记录: 000001 - 平安银行 (连板: 3天, 排名: 64)
更新记录: 000002 - 万科A (连板: 1天, 排名: 128)
```
