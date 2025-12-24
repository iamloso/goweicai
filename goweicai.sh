#!/bin/bash

# GoWenCai 定时任务服务管理脚本

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
BINARY="$PROJECT_ROOT/cmd/goweicai/goweicai"
PID_FILE="$PROJECT_ROOT/goweicai.pid"
LOG_FILE="$PROJECT_ROOT/cmd/goweicai/goweicai.log"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否正在运行
is_running() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$PID_FILE"
            return 1
        fi
    fi
    return 1
}

# 编译
build() {
    print_info "开始编译..."
    cd "$PROJECT_ROOT" || exit 1
    if make build; then
        print_info "编译成功"
        return 0
    else
        print_error "编译失败"
        return 1
    fi
}

# 启动服务
start() {
    if is_running; then
        print_warn "服务已在运行中 (PID: $(cat "$PID_FILE"))"
        return 1
    fi

    print_info "启动服务..."
    cd "$PROJECT_ROOT/cmd/goweicai" || exit 1
    
    nohup ./goweicai > "$LOG_FILE" 2>&1 &
    PID=$!
    echo "$PID" > "$PID_FILE"
    
    sleep 2
    
    if is_running; then
        print_info "服务启动成功 (PID: $PID)"
        print_info "日志文件: $LOG_FILE"
        return 0
    else
        print_error "服务启动失败，请查看日志: $LOG_FILE"
        rm -f "$PID_FILE"
        return 1
    fi
}

# 停止服务
stop() {
    if ! is_running; then
        print_warn "服务未运行"
        return 1
    fi

    PID=$(cat "$PID_FILE")
    print_info "正在停止服务 (PID: $PID)..."
    
    kill "$PID"
    
    # 等待最多 10 秒
    for i in {1..10}; do
        if ! ps -p "$PID" > /dev/null 2>&1; then
            rm -f "$PID_FILE"
            print_info "服务已停止"
            return 0
        fi
        sleep 1
    done
    
    # 强制停止
    print_warn "尝试强制停止..."
    kill -9 "$PID" 2>/dev/null
    rm -f "$PID_FILE"
    print_info "服务已强制停止"
}

# 重启服务
restart() {
    print_info "重启服务..."
    stop
    sleep 2
    start
}

# 查看状态
status() {
    if is_running; then
        PID=$(cat "$PID_FILE")
        print_info "服务运行中"
        echo "  PID: $PID"
        echo "  运行时间: $(ps -p "$PID" -o etime= | tr -d ' ')"
        echo "  内存使用: $(ps -p "$PID" -o rss= | awk '{printf "%.2f MB\n", $1/1024}')"
        echo "  CPU 使用: $(ps -p "$PID" -o %cpu= | tr -d ' ')%"
    else
        print_warn "服务未运行"
    fi
}

# 查看日志
logs() {
    if [ ! -f "$LOG_FILE" ]; then
        print_warn "日志文件不存在"
        return 1
    fi
    
    if [ "$1" == "-f" ]; then
        print_info "实时查看日志 (Ctrl+C 退出)..."
        tail -f "$LOG_FILE"
    else
        tail -n 50 "$LOG_FILE"
    fi
}

# 部署（编译 + 重启）
deploy() {
    print_info "开始部署..."
    
    if ! build; then
        return 1
    fi
    
    if is_running; then
        restart
    else
        start
    fi
}

# 使用说明
usage() {
    cat << EOF
GoWenCai 定时任务服务管理脚本

用法: $0 {build|start|stop|restart|status|logs|deploy}

命令:
  build       编译项目
  start       启动服务
  stop        停止服务
  restart     重启服务
  status      查看服务状态
  logs        查看日志 (最近 50 行)
  logs -f     实时查看日志
  deploy      部署 (编译 + 重启)

示例:
  $0 start          # 启动服务
  $0 status         # 查看状态
  $0 logs -f        # 实时查看日志
  $0 deploy         # 编译并部署

EOF
}

# 主函数
main() {
    case "$1" in
        build)
            build
            ;;
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        status)
            status
            ;;
        logs)
            logs "$2"
            ;;
        deploy)
            deploy
            ;;
        *)
            usage
            exit 1
            ;;
    esac
}

main "$@"
