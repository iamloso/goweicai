#!/bin/bash

# 检查并清理 base_info 表的重复数据

# 从配置文件读取数据库连接信息
CONFIG_FILE="configs/config.yaml"

# 提取数据库连接信息（需要安装 yq 或使用其他方法）
# 这里假设使用默认的本地 MySQL 配置
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASS="${DB_PASS:-}"
DB_NAME="${DB_NAME:-goweicai}"

echo "==================================="
echo "Base Info 表重复数据检查和清理工具"
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
echo "步骤 1: 检查重复数据..."

# 检查是否存在重复数据
DUPLICATE_COUNT=$(echo "
SELECT COUNT(*) as dup_count
FROM (
    SELECT code, trade_date, COUNT(*) as cnt
    FROM base_info
    GROUP BY code, trade_date
    HAVING cnt > 1
) as duplicates;
" | $MYSQL_CMD -N)

echo "发现 ${DUPLICATE_COUNT} 组重复的 (code, trade_date) 记录"

if [ "$DUPLICATE_COUNT" -eq 0 ]; then
    echo "✓ 没有重复数据，数据库状态良好！"
    exit 0
fi

# 显示重复数据详情
echo ""
echo "重复数据详情："
echo "
SELECT code, trade_date, COUNT(*) as count, GROUP_CONCAT(id) as ids
FROM base_info
GROUP BY code, trade_date
HAVING count > 1
ORDER BY count DESC
LIMIT 10;
" | $MYSQL_CMD

echo ""
read -p "是否要清理重复数据？这将保留每组中 ID 最大的记录 (y/n): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "操作已取消"
    exit 0
fi

echo ""
echo "步骤 2: 开始清理重复数据..."

# 执行清理脚本
$MYSQL_CMD < migrations/004_fix_duplicate_base_info.sql

if [ $? -eq 0 ]; then
    echo "✓ 重复数据清理完成！"
    
    # 再次检查
    NEW_DUPLICATE_COUNT=$(echo "
    SELECT COUNT(*) as dup_count
    FROM (
        SELECT code, trade_date, COUNT(*) as cnt
        FROM base_info
        GROUP BY code, trade_date
        HAVING cnt > 1
    ) as duplicates;
    " | $MYSQL_CMD -N)
    
    echo "清理后剩余重复记录组数: ${NEW_DUPLICATE_COUNT}"
    
    if [ "$NEW_DUPLICATE_COUNT" -eq 0 ]; then
        echo "✓ 所有重复数据已成功清理！"
    else
        echo "⚠ 警告: 仍有重复数据，请手动检查"
    fi
else
    echo "✗ 清理失败，请检查错误信息"
    exit 1
fi

echo ""
echo "步骤 3: 验证唯一索引..."

# 检查唯一索引是否存在
INDEX_EXISTS=$(echo "
SHOW INDEX FROM base_info WHERE Key_name = 'uk_code_date';
" | $MYSQL_CMD -N | wc -l)

if [ "$INDEX_EXISTS" -gt 0 ]; then
    echo "✓ 唯一索引 uk_code_date 存在"
else
    echo "⚠ 警告: 唯一索引 uk_code_date 不存在，需要创建"
    read -p "是否创建唯一索引？(y/n): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "ALTER TABLE base_info ADD UNIQUE KEY uk_code_date (code, trade_date);" | $MYSQL_CMD
        if [ $? -eq 0 ]; then
            echo "✓ 唯一索引创建成功"
        else
            echo "✗ 唯一索引创建失败"
        fi
    fi
fi

echo ""
echo "==================================="
echo "处理完成！"
echo "==================================="
