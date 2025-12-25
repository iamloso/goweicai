-- 添加流通市值字符串字段

ALTER TABLE `base_info`
    ADD COLUMN `circulation_market_value_str` varchar(50) DEFAULT NULL COMMENT '流通市值（如：12万、1.3亿）' AFTER `circulation_market_value`;
