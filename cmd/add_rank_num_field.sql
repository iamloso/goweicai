-- 添加排名数字字段
ALTER TABLE zp_jj 
ADD COLUMN auction_unmatched_amount_rank_num INT DEFAULT 0 COMMENT '竞价未匹配金额排名数字' 
AFTER auction_unmatched_amount_rank;

-- 添加索引以便按排名查询
CREATE INDEX idx_rank_num ON zp_jj(auction_unmatched_amount_rank_num);
