-- 修改唯一索引：从 (code, trade_date) 改为只用 code
-- 这样每个股票代码只保留一条最新记录

-- 1. 删除旧的唯一索引
ALTER TABLE base_info DROP INDEX uk_code_date;

-- 2. 创建新的唯一索引（只基于 code）
ALTER TABLE base_info ADD UNIQUE KEY uk_code (code);

-- 3. 保留 trade_date 的普通索引用于查询
-- idx_trade_date 应该已存在，如果不存在则创建
-- ALTER TABLE base_info ADD INDEX idx_trade_date (trade_date);
