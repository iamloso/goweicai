-- 添加格式化字段
ALTER TABLE zp_jj 
ADD COLUMN morning_auction_amount_str VARCHAR(50) DEFAULT '' COMMENT '早盘竞价金额(万/亿)' AFTER morning_auction_amount,
ADD COLUMN turnover_str VARCHAR(50) DEFAULT '' COMMENT '成交额(万/亿)' AFTER turnover;
