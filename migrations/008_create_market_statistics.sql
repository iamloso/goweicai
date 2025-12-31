-- 创建市场统计表
CREATE TABLE IF NOT EXISTS `market_statistics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `trade_date` date NOT NULL COMMENT '交易日期',
  `limit_up_count` int NOT NULL DEFAULT '0' COMMENT '涨停数',
  `limit_down_count` int NOT NULL DEFAULT '0' COMMENT '跌停数',
  `broken_count` int NOT NULL DEFAULT '0' COMMENT '开板数',
  `max_consecutive_days` int NOT NULL DEFAULT '0' COMMENT '连板高度',
  `two_consecutive_count` int NOT NULL DEFAULT '0' COMMENT '二连板数量',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_trade_date` (`trade_date`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='市场统计数据表';
