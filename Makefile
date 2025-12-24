.PHONY: build run clean test

# 编译
build:
	cd cmd/goweicai && go build -o goweicai

# 运行
run:
	cd cmd/goweicai && ./goweicai

# 编译并运行
all: build run

# 清理
clean:
	rm -f cmd/goweicai/goweicai

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
