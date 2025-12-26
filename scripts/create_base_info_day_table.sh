#!/bin/bash

# 应用 base_info_day 表创建迁移

# 从配置文件读取数据库连接信息
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_NAME="${DB_NAME:-goweicai}"

echo "==================================="
echo "创建 base_info_day 表"
echo "==================================="

# 检查 mysql 命令是否可用
if ! command -v mysql &> /dev/null; then
    echo "错误: 未找到 mysql 命令，请先安装 MySQL 客户端"
    exit 1
fi

# 构建 MySQL 连接参数
MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER}"
if [ -n "$DB_PASS" ]; then
    MYSQL_CMD="$MYSQL_CMD -p${DB_PASS}"
fi
MYSQL_CMD="$MYSQL_CMD ${DB_NAME}"

echo ""
echo "执行迁移脚本..."

# 执行迁移
$MYSQL_CMD < migrations/006_create_base_info_day.sql

if [ $? -eq 0 ]; then
    echo "✓ base_info_day 表创建成功！"
    
    # 验证表是否存在
    TABLE_EXISTS=$(echo "SHOW TABLES LIKE 'base_info_day';" | $MYSQL_CMD -N | wc -l)
    
    if [ "$TABLE_EXISTS" -eq 1 ]; then
        echo "✓ 表已验证存在"
        
        # 显示表结构
        echo ""
        echo "表结构："
        echo "DESC base_info_day;" | $MYSQL_CMD
        
        echo ""
        echo "索引信息："
        echo "SHOW INDEX FROM base_info_day;" | $MYSQL_CMD
    else
        echo "⚠ 警告: 表验证失败"
    fi
else
    echo "✗ 表创建失败，请检查错误信息"
    exit 1
fi

echo ""
echo "==================================="
echo "迁移完成！"
echo "==================================="
