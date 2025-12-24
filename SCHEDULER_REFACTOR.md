# 定时任务改造总结

## 改造内容

### 1. 依赖引入
- 引入 `github.com/robfig/cron/v3` 定时任务调度库
- 版本：v3.0.1

### 2. 配置文件更新

**文件**: `configs/config.yaml`

新增 `scheduler` 配置段：
```yaml
scheduler:
  cron: "0 0 9 * * *"  # Cron 表达式
  run_on_start: true    # 启动时立即执行
```

### 3. 配置结构更新

**文件**: `internal/conf/conf.go`

新增 `Scheduler` 结构体：
```go
type Bootstrap struct {
    Server    *Server    `json:"server"`
    Scheduler *Scheduler `json:"scheduler"`  // 新增
    Data      *Data      `json:"data"`
    Wencai    *Wencai    `json:"wencai"`
}

type Scheduler struct {
    Cron       string `json:"cron"`
    RunOnStart bool   `json:"run_on_start"`
}
```

### 4. 主程序改造

**文件**: `cmd/goweicai/main.go`

主要变更：
- 从一次性执行改为持续运行的定时任务
- 支持优雅退出（SIGINT/SIGTERM 信号）
- 支持启动时立即执行一次

核心逻辑：
```go
// 创建定时任务调度器
c := cron.New(cron.WithSeconds())

// 添加定时任务
c.AddFunc(cronExpr, job)

// 可选：启动时立即执行
if bc.Scheduler.RunOnStart {
    job()
}

// 启动调度器
c.Start()

// 等待退出信号
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

// 优雅停止
c.Stop()
```

### 5. Makefile 增强

**文件**: `Makefile`

新增命令：
- `make daemon`: 后台运行
- `make stop`: 停止后台进程
- `make status`: 查看运行状态
- `make logs`: 查看实时日志

## 功能特性

### ✅ 灵活的调度配置
- 支持标准 Cron 表达式（秒级别精度）
- 通过配置文件轻松调整执行时间
- 默认每天 9:00 执行

### ✅ 启动时立即执行
- 可配置 `run_on_start: true`
- 不必等待第一次定时触发
- 适合需要立即同步数据的场景

### ✅ 优雅退出
- 监听系统信号（SIGINT/SIGTERM）
- 正确停止调度器
- 清理资源

### ✅ 后台运行支持
- 支持守护进程模式
- 日志输出到文件
- 进程管理命令

### ✅ 完善的日志
- 记录任务启动、执行、完成状态
- 错误信息完整输出
- 支持实时查看

## 使用示例

### 开发环境
```bash
# 测试配置：每 10 秒执行一次
# 修改 configs/config.yaml:
#   cron: "*/10 * * * * *"

make build
make run
```

### 生产环境
```bash
# 配置：每天 9:00 执行
# configs/config.yaml:
#   cron: "0 0 9 * * *"
#   run_on_start: true

make build
make daemon      # 后台启动
make status      # 检查状态
make logs        # 查看日志
make stop        # 停止服务
```

## Cron 表达式示例

| 需求 | Cron 表达式 |
|------|-------------|
| 每天 9:00 | `0 0 9 * * *` |
| 每 30 分钟 | `0 */30 * * * *` |
| 每天 9:00 和 15:00 | `0 0 9,15 * * *` |
| 工作日 9:00 | `0 0 9 * * 1-5` |
| 每小时整点 | `0 0 * * * *` |
| 每 5 秒（测试） | `*/5 * * * * *` |

## 文件清单

### 新增文件
- `SCHEDULER_GUIDE.md`: 定时任务使用指南
- `SCHEDULER_REFACTOR.md`: 本文件，改造总结

### 修改文件
- `configs/config.yaml`: 新增 scheduler 配置
- `internal/conf/conf.go`: 新增 Scheduler 结构
- `cmd/goweicai/main.go`: 改造为定时任务模式
- `Makefile`: 新增后台运行相关命令
- `go.mod`: 新增 cron/v3 依赖

## 技术栈

- **框架**: Kratos v2.9.2
- **调度器**: robfig/cron v3.0.1
- **数据库**: MySQL
- **配置**: YAML

## 后续优化建议

1. **监控告警**
   - 集成 Prometheus metrics
   - 添加任务执行时长监控
   - 失败告警通知

2. **日志管理**
   - 集成日志轮转（logrotate）
   - 结构化日志（JSON 格式）
   - 日志级别动态调整

3. **高可用**
   - 使用 systemd 管理进程
   - 配置自动重启
   - 健康检查接口

4. **分布式锁**
   - 如果部署多实例，需要添加分布式锁
   - 防止任务重复执行
   - 可使用 Redis 或 etcd

5. **任务管理界面**
   - 提供 HTTP API 查看任务状态
   - 手动触发任务
   - 暂停/恢复调度

## 测试检查清单

- [x] 编译通过
- [x] 配置文件正确加载
- [x] 定时任务可以正常调度
- [x] 启动时立即执行功能正常
- [x] 优雅退出功能正常
- [x] 后台运行正常
- [x] 日志输出正常
- [ ] 长时间运行稳定性测试（建议运行 24 小时）
- [ ] 异常情况处理测试（数据库断开等）

## 版本信息

- 改造日期：2024-12-24
- Go 版本：1.23.1
- Kratos 版本：v2.9.2
- Cron 版本：v3.0.1
