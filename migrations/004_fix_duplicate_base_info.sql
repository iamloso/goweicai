-- 清理 base_info 表中的重复数据
-- 此脚本保留每个 (code, trade_date) 组合中 id 最大的记录（即最新的记录）

-- 1. 创建临时表存储要保留的记录ID
CREATE TEMPORARY TABLE IF NOT EXISTS temp_keep_ids AS
SELECT MAX(id) as keep_id
FROM base_info
GROUP BY code, trade_date;

-- 2. 删除重复的记录（保留 id 最大的）
DELETE FROM base_info
WHERE id NOT IN (SELECT keep_id FROM temp_keep_ids);

-- 3. 删除临时表
DROP TEMPORARY TABLE IF EXISTS temp_keep_ids;

-- 4. 验证唯一索引存在
-- 如果索引不存在会报错，如果已存在则忽略
-- MySQL 8.0+ 不支持 IF NOT EXISTS for unique index，所以需要手动检查
-- 可以通过查询确认：
-- SHOW INDEX FROM base_info WHERE Key_name = 'uk_code_date';

-- 5. 如果需要重建索引（在索引损坏或不一致的情况下）
-- ALTER TABLE base_info DROP INDEX IF EXISTS uk_code_date;
-- ALTER TABLE base_info ADD UNIQUE KEY uk_code_date (code, trade_date);
