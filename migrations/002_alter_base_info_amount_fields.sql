-- 修改金额字段为字符串类型，存储带单位的格式（如：12万、1.3亿）

ALTER TABLE `base_info`
    MODIFY COLUMN `auction_unmatched_amount_str` varchar(50) DEFAULT NULL COMMENT '竞价未匹配金额（如：12万、1.3亿）',
    MODIFY COLUMN `morning_auction_amount_str` varchar(50) DEFAULT NULL COMMENT '早盘竞价金额（如：12万、1.3亿）',
    MODIFY COLUMN `turnover_str` varchar(50) DEFAULT NULL COMMENT '成交额（如：12万、1.3亿）';
