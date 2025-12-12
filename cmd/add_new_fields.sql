-- 添加新字段：涨停原因、公司亮点、行业分类、概念题材、涨停封单金额、连板天数

-- 检查并添加字段
ALTER TABLE zp_jj 
ADD COLUMN limit_up_reason TEXT COMMENT '涨停原因',
ADD COLUMN company_highlights TEXT COMMENT '公司亮点',
ADD COLUMN industry_category VARCHAR(100) DEFAULT '' COMMENT '行业分类',
ADD COLUMN concept_theme TEXT COMMENT '概念题材',
ADD COLUMN limit_up_seal_amount BIGINT DEFAULT 0 COMMENT '涨停封单金额(元)',
ADD COLUMN limit_up_seal_amount_str VARCHAR(50) DEFAULT '' COMMENT '涨停封单金额(万/亿)',
ADD COLUMN consecutive_limit_days INT DEFAULT 0 COMMENT '连板天数';

-- 添加索引以便查询
CREATE INDEX idx_consecutive_days ON zp_jj(consecutive_limit_days);
CREATE INDEX idx_industry ON zp_jj(industry_category);
