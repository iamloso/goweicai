# GoWencai 快速参考

## 服务端口

| 服务 | 端口 | 协议 |
|------|------|------|
| HTTP API | 8000 | HTTP/1.1 |
| gRPC API | 9000 | HTTP/2 |

## 快速命令

```bash
# 编译
make build

# 运行
make run              # 前台运行
make daemon           # 后台运行
make stop             # 停止服务
make status           # 查看状态
make logs             # 查看日志

# 代码生成
make proto            # 生成 gRPC 代码

# 测试
make grpc-test        # 测试 gRPC 服务
./test_api.sh         # 测试 HTTP API
```

## HTTP API 示例

```bash
# 健康检查
curl http://localhost:8000/health

# 获取最新股票
curl "http://localhost:8000/api/stocks/latest?limit=10"

# 查询股票
curl -X POST http://localhost:8000/api/stocks/query \
  -H "Content-Type: application/json" \
  -d '{"code":"000001","page":1,"page_size":20}'
```

## gRPC API 示例

```bash
# 安装 grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 列出服务
grpcurl -plaintext localhost:9000 list

# 获取最新股票
grpcurl -plaintext -d '{"limit": 10}' \
  localhost:9000 stock.v1.StockService/GetLatestStocks

# 查询股票
grpcurl -plaintext -d '{"code": "000001", "page": 1, "page_size": 20}' \
  localhost:9000 stock.v1.StockService/QueryStocks

# 触发数据抓取
grpcurl -plaintext -d '{}' \
  localhost:9000 stock.v1.StockService/TriggerFetch
```

## Cron 表达式示例

```yaml
# 每天 9:00
cron: "0 0 9 * * *"

# 每 30 分钟
cron: "0 */30 * * * *"

# 每天 9:00 和 15:00
cron: "0 0 9,15 * * *"

# 工作日 9:00
cron: "0 0 9 * * 1-5"

# 每 5 秒（测试用）
cron: "*/5 * * * * *"
```

## 配置文件位置

- 主配置: `configs/config.yaml`
- Cookie 获取: 浏览器登录问财后获取

## 文档索引

- [README.md](./README.md) - 项目总览
- [API_DOCS.md](./API_DOCS.md) - HTTP API 文档
- [API_EXAMPLES.md](./API_EXAMPLES.md) - API 示例代码
- [GRPC_DOCS.md](./GRPC_DOCS.md) - gRPC 文档
- [SCHEDULER_GUIDE.md](./SCHEDULER_GUIDE.md) - 定时任务指南
- [README_KRATOS.md](./README_KRATOS.md) - Kratos 框架说明

## 故障排查

### 服务无法启动

```bash
# 检查端口占用
lsof -i :8000
lsof -i :9000

# 查看日志
tail -f cmd/goweicai/goweicai.log

# 检查配置
cat configs/config.yaml
```

### 数据库连接失败

```bash
# 检查 MySQL 服务
systemctl status mysql

# 测试连接
mysql -u root -p

# 检查配置
grep "source:" configs/config.yaml
```

### Proto 代码生成失败

```bash
# 检查 protoc
protoc --version

# 重新安装工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 重新生成
./generate.sh
```

## 项目结构

```
gowencai/
├── api/proto/           # Proto 定义和生成代码
├── cmd/
│   └── goweicai/       # 主程序
├── configs/            # 配置文件
├── internal/
│   ├── biz/           # 业务逻辑层
│   ├── data/          # 数据访问层
│   ├── service/       # 服务层
│   └── conf/          # 配置结构
├── pywencai/          # Python 相关文件
├── Makefile           # 构建脚本
└── *.md               # 文档文件
```

## 开发工作流

```bash
# 1. 克隆项目
git clone https://github.com/iamloso/goweicai.git
cd goweicai

# 2. 安装依赖
go mod download

# 3. 配置
cp configs/config.yaml.example configs/config.yaml
vim configs/config.yaml  # 修改数据库和 cookie

# 4. 编译
make build

# 5. 运行
make run

# 6. 测试
./test_api.sh
make grpc-test
```

## 生产部署

```bash
# 1. 编译
make build

# 2. 配置 systemd
sudo cp goweicai.service /etc/systemd/system/
sudo systemctl enable goweicai
sudo systemctl start goweicai

# 3. 查看状态
sudo systemctl status goweicai
sudo journalctl -u goweicai -f
```

## 性能调优

### HTTP

- 调整连接超时: `server.http.timeout`
- 启用 GZIP 压缩
- 添加缓存层

### gRPC

- 调整连接池大小
- 启用连接复用
- 使用流式 RPC

### 数据库

- 添加索引
- 使用连接池
- 开启查询缓存

## 监控指标

建议监控：
- HTTP 请求量和延迟
- gRPC 调用量和延迟
- 数据库连接数
- 定时任务执行状态
- 内存和 CPU 使用率

## 安全建议

1. 使用 HTTPS/TLS
2. 添加 API 认证
3. 限制请求频率
4. 定期更新 cookie
5. 数据库访问控制
6. 日志脱敏处理

## 许可证

MIT License - 详见 [LICENSE](./LICENSE)

## 贡献

欢迎提交 Issue 和 Pull Request！
