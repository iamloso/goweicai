# Base Info 表防重复数据方案

## 问题描述

`base_info` 表可能出现重复数据，即相同的 `(code, trade_date)` 组合存在多条记录。

## 解决方案

### 1. 数据库层面保护

数据库表已设置唯一索引：

```sql
UNIQUE KEY `uk_code_date` (`code`, `trade_date`)
```

此索引确保同一股票代码和交易日期的组合在数据库中是唯一的。

### 2. 应用层面保护

#### 2.1 BatchSave 方法

`BatchSave` 方法使用 GORM 的 `ON CONFLICT` 子句实现 **upsert** 操作：

- 当 `(code, trade_date)` 组合已存在时，更新该记录
- 当 `(code, trade_date)` 组合不存在时，插入新记录

代码示例：
```go
err := r.data.gormDB.WithContext(ctx).
    Clauses(clause.OnConflict{
        Columns: []clause.Column{{Name: "code"}, {Name: "trade_date"}},
        DoUpdates: clause.AssignmentColumns([]string{
            "stock_name", "latest_price", "auction_change_rate", 
            // ... 其他需要更新的字段
        }),
    }).
    CreateInBatches(models, 100).Error
```

#### 2.2 Save 方法

`Save` 方法也已修改为使用 upsert 逻辑，不再使用简单的 `Create` 操作。

#### 2.3 Update 方法

`Update` 方法现在基于 `code` 和 `trade_date` 进行更新：

```go
result := r.data.gormDB.WithContext(ctx).
    Where("code = ? AND trade_date = ?", model.Code, model.TradeDate).
    Updates(model)
```

### 3. 清理已有重复数据

#### 3.1 使用自动化脚本

运行清理脚本：

```bash
./scripts/fix_duplicate_baseinfo.sh
```

此脚本会：
1. 检查重复数据
2. 显示重复记录详情
3. 询问是否清理（保留 ID 最大的记录）
4. 验证唯一索引是否存在

#### 3.2 手动执行 SQL

```bash
mysql -u root -p goweicai < migrations/004_fix_duplicate_base_info.sql
```

#### 3.3 检查重复数据

查询重复记录：

```sql
SELECT code, trade_date, COUNT(*) as count, GROUP_CONCAT(id) as ids
FROM base_info
GROUP BY code, trade_date
HAVING count > 1
ORDER BY count DESC;
```

### 4. 验证修复

#### 4.1 检查唯一索引

```sql
SHOW INDEX FROM base_info WHERE Key_name = 'uk_code_date';
```

#### 4.2 尝试插入重复数据

测试唯一索引是否生效：

```sql
-- 这应该会失败（如果数据已存在）
INSERT INTO base_info (code, trade_date, stock_name, create_time)
VALUES ('000001', '2025-12-26', '测试股票', NOW());

-- 第二次执行相同的 INSERT 应该报错：
-- ERROR 1062 (23000): Duplicate entry '000001-2025-12-26' for key 'uk_code_date'
```

## 最佳实践

### 1. 数据插入

- **推荐**：使用 `BatchSave` 批量插入/更新数据
- **避免**：直接使用 `Create` 方法插入数据（除非确定数据不存在）

### 2. 数据更新

- **推荐**：使用 `Update` 方法或 `BatchSave` 方法
- **注意**：确保提供正确的 `code` 和 `trade_date`

### 3. 定期检查

建议定期执行以下查询检查数据一致性：

```sql
-- 检查是否有重复数据
SELECT COUNT(*) FROM (
    SELECT code, trade_date, COUNT(*) as cnt
    FROM base_info
    GROUP BY code, trade_date
    HAVING cnt > 1
) as duplicates;
```

## 修改记录

| 日期 | 修改内容 | 修改文件 |
|------|---------|---------|
| 2025-12-26 | 修正 BaseInfoModel GORM 标签 | internal/data/baseinfo.go |
| 2025-12-26 | Save 方法改为 upsert 逻辑 | internal/data/baseinfo.go |
| 2025-12-26 | Update 方法优化 | internal/data/baseinfo.go |
| 2025-12-26 | 创建清理脚本 | scripts/fix_duplicate_baseinfo.sh |
| 2025-12-26 | 创建清理 SQL | migrations/004_fix_duplicate_base_info.sql |

## 相关文件

- [internal/data/baseinfo.go](../../internal/data/baseinfo.go) - 数据访问层
- [migrations/001_create_base_info.sql](../../migrations/001_create_base_info.sql) - 表结构定义
- [migrations/004_fix_duplicate_base_info.sql](../../migrations/004_fix_duplicate_base_info.sql) - 清理重复数据
- [scripts/fix_duplicate_baseinfo.sh](../../scripts/fix_duplicate_baseinfo.sh) - 自动化清理脚本
