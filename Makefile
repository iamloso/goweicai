.PHONY: build run clean test daemon stop status

# 编译
build:
	cd cmd/goweicai && go build -o goweicai

# 前台运行
run:
	cd cmd/goweicai && ./goweicai

# 后台运行（守护进程模式）
daemon:
	cd cmd/goweicai && nohup ./goweicai > goweicai.log 2>&1 &
	@echo "定时任务已在后台启动，日志文件: cmd/goweicai/goweicai.log"
	@echo "PID: $$(pgrep -f './goweicai')"

# 停止后台进程
stop:
	@pkill -f './goweicai' || echo "没有运行中的进程"

# 查看运行状态
status:
	@pgrep -fa './goweicai' || echo "服务未运行"

# 查看日志
logs:
	tail -f cmd/goweicai/goweicai.log

# 编译并运行
all: build run

# 清理
clean:
	rm -f cmd/goweicai/goweicai
	rm -f cmd/goweicai/goweicai.log

# 测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...
