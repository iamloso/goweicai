#!/bin/bash

# 应用 zt_day 表创建迁移

DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_NAME="${DB_NAME:-goweicai}"

echo "==================================="
echo "创建 zt_day 表"
echo "==================================="

if ! command -v mysql &> /dev/null; then
    echo "错误: 未找到 mysql 命令，请先安装 MySQL 客户端"
    exit 1
fi

MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER}"
if [ -n "$DB_PASS" ]; then
    MYSQL_CMD="$MYSQL_CMD -p${DB_PASS}"
fi
MYSQL_CMD="$MYSQL_CMD ${DB_NAME}"

echo ""
echo "执行迁移脚本..."

$MYSQL_CMD < migrations/007_create_zt_day.sql

if [ $? -eq 0 ]; then
    echo "✓ zt_day 表创建成功！"
    
    TABLE_EXISTS=$(echo "SHOW TABLES LIKE 'zt_day';" | $MYSQL_CMD -N | wc -l)
    
    if [ "$TABLE_EXISTS" -eq 1 ]; then
        echo "✓ 表已验证存在"
        
        echo ""
        echo "表结构："
        echo "DESC zt_day;" | $MYSQL_CMD
        
        echo ""
        echo "索引信息："
        echo "SHOW INDEX FROM zt_day;" | $MYSQL_CMD
    fi
else
    echo "✗ 表创建失败"
    exit 1
fi

echo ""
echo "==================================="
echo "迁移完成！"
echo "==================================="
