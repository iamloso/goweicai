#!/bin/bash

# GoWencai HTTP API 测试脚本

set -e

API_BASE="http://localhost:8000"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================================"
echo "  GoWencai HTTP API 测试"
echo "================================================"
echo ""

# 检查服务是否运行
echo -n "检查服务状态... "
if curl -s -f "${API_BASE}/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 服务运行中${NC}"
else
    echo -e "${RED}✗ 服务未运行，请先启动服务${NC}"
    echo ""
    echo "启动命令: make daemon"
    exit 1
fi
echo ""

# 测试 1: 健康检查
echo "测试 1: 健康检查"
echo "----------------------------------------"
echo "请求: GET ${API_BASE}/health"
response=$(curl -s "${API_BASE}/health")
echo "响应: ${response}"
if echo "$response" | grep -q "ok"; then
    echo -e "${GREEN}✓ 测试通过${NC}"
else
    echo -e "${RED}✗ 测试失败${NC}"
fi
echo ""

# 测试 2: 获取最新股票数据
echo "测试 2: 获取最新股票数据"
echo "----------------------------------------"
echo "请求: GET ${API_BASE}/api/stocks/latest?limit=5"
response=$(curl -s "${API_BASE}/api/stocks/latest?limit=5")
echo "响应: ${response}"
if echo "$response" | grep -q "success"; then
    echo -e "${GREEN}✓ 测试通过${NC}"
else
    echo -e "${RED}✗ 测试失败${NC}"
fi
echo ""

# 测试 3: POST 查询股票
echo "测试 3: POST 查询股票"
echo "----------------------------------------"
echo "请求: POST ${API_BASE}/api/stocks/query"
echo '数据: {"code":"000001","page":1,"page_size":10}'
response=$(curl -s -X POST "${API_BASE}/api/stocks/query" \
    -H "Content-Type: application/json" \
    -d '{"code":"000001","page":1,"page_size":10}')
echo "响应: ${response}"
if echo "$response" | grep -q "success"; then
    echo -e "${GREEN}✓ 测试通过${NC}"
else
    echo -e "${RED}✗ 测试失败${NC}"
fi
echo ""

# 测试 4: GET 查询股票
echo "测试 4: GET 查询股票"
echo "----------------------------------------"
echo "请求: GET ${API_BASE}/api/stocks/query?code=000001&page=1&page_size=10"
response=$(curl -s "${API_BASE}/api/stocks/query?code=000001&page=1&page_size=10")
echo "响应: ${response}"
if echo "$response" | grep -q "success"; then
    echo -e "${GREEN}✓ 测试通过${NC}"
else
    echo -e "${RED}✗ 测试失败${NC}"
fi
echo ""

# 测试 5: 错误处理 - 不支持的方法
echo "测试 5: 错误处理（不支持的方法）"
echo "----------------------------------------"
echo "请求: DELETE ${API_BASE}/api/stocks/query"
response=$(curl -s -X DELETE "${API_BASE}/api/stocks/query")
echo "响应: ${response}"
if echo "$response" | grep -q "error"; then
    echo -e "${GREEN}✓ 测试通过（正确返回错误）${NC}"
else
    echo -e "${RED}✗ 测试失败${NC}"
fi
echo ""

echo "================================================"
echo "  所有测试完成"
echo "================================================"
echo ""
echo "提示："
echo "  - 当前 API 返回空数据是正常的（TODO: 实现数据库查询）"
echo "  - 完整 API 文档: API_DOCS.md"
echo ""
