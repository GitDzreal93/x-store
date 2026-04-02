#!/bin/bash

# X-Store 开发环境启动脚本（不依赖 Docker）
# 使用方式: ./dev-start.sh [backend|frontend|admin|docs|all|stop|restart|check]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  X-Store 开发环境启动脚本${NC}"
echo -e "${BLUE}  (不依赖 Docker)${NC}"
echo -e "${BLUE}========================================${NC}"

# 项目根目录（无论在什么路径执行脚本，都以脚本所在目录为根）
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 日志目录
LOG_DIR="${ROOT_DIR}/logs"
mkdir -p "${LOG_DIR}"

# 简单日志轮转：保留最近 5 份
rotate_log() {
    local file="$1"
    if [ -f "$file" ]; then
        for i in 5 4 3 2 1; do
            if [ -f "${file}.${i}" ]; then
                mv "${file}.${i}" "${file}.$((i + 1))" 2>/dev/null || true
            fi
        done
        mv "$file" "${file}.1"
    fi
}

# 获取本机 IP
get_local_ip() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1
    else
        # Linux
        hostname -I | awk '{print $1}'
    fi
}

# 启动后端
start_backend() {
    echo -e "${GREEN}启动后端服务...${NC}"
    cd "${ROOT_DIR}/backend"
    
    # 检查是否需要编译
    if [ ! -f "x-store-backend" ] || [ "cmd/main.go" -nt "x-store-backend" ]; then
        echo -e "${YELLOW}编译后端...${NC}"
        go build -o x-store-backend ./cmd/main.go
    fi
    
    LOCAL_IP=$(get_local_ip)
    
    echo -e "${GREEN}✓ 后端服务启动成功${NC}"
    echo -e "${YELLOW}📍 本地访问: ${CYAN}http://localhost:8082${NC}"
    echo -e "${YELLOW}🌐 外网访问: ${CYAN}http://${LOCAL_IP}:8082${NC}"
    echo -e "${YELLOW}📖 API 测试: ${CYAN}http://localhost:8082/api/categories${NC}"
    echo -e "${YELLOW}📝 日志文件: ${CYAN}${LOG_DIR}/backend.log${NC}"
    echo ""
    
    rotate_log "${LOG_DIR}/backend.log"
    # 前台运行（由调用方决定是否放到后台）
    ./x-store-backend >> "${LOG_DIR}/backend.log" 2>&1
}

# 启动 C 端前台
start_frontend() {
    echo -e "${GREEN}启动 C 端前台...${NC}"
    cd "${ROOT_DIR}/frontend-store"
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装依赖...${NC}"
        npm install
    fi
    
    LOCAL_IP=$(get_local_ip)
    
    echo -e "${GREEN}✓ C 端前台启动成功${NC}"
    echo -e "${YELLOW}📍 本地访问: ${CYAN}http://localhost:3000${NC}"
    echo -e "${YELLOW}🌐 外网访问: ${CYAN}http://${LOCAL_IP}:3000${NC}"
    echo -e "${YELLOW}🛒 商城首页: ${CYAN}http://localhost:3000${NC}"
    echo -e "${YELLOW}📝 日志文件: ${CYAN}${LOG_DIR}/frontend.log${NC}"
    echo ""
    
    rotate_log "${LOG_DIR}/frontend.log"
    npm run dev >> "${LOG_DIR}/frontend.log" 2>&1
}

# 启动管理后台
start_admin() {
    echo -e "${GREEN}启动管理后台...${NC}"
    cd "${ROOT_DIR}/admin-panel"
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装依赖...${NC}"
        npm install
    fi
    
    LOCAL_IP=$(get_local_ip)
    
    echo -e "${GREEN}✓ 管理后台启动成功${NC}"
    echo -e "${YELLOW}📍 本地访问: ${CYAN}http://localhost:5174${NC}"
    echo -e "${YELLOW}🌐 外网访问: ${CYAN}http://${LOCAL_IP}:5174${NC}"
    echo -e "${YELLOW}⚙️  管理面板: ${CYAN}http://localhost:5174${NC}"
    echo -e "${YELLOW}📝 日志文件: ${CYAN}${LOG_DIR}/admin.log${NC}"
    echo ""
    
    rotate_log "${LOG_DIR}/admin.log"
    npm run dev >> "${LOG_DIR}/admin.log" 2>&1
}

# 启动文档站
start_docs() {
    echo -e "${GREEN}启动文档站...${NC}"
    cd "${ROOT_DIR}/docs"
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装依赖...${NC}"
        npm install
    fi
    
    LOCAL_IP=$(get_local_ip)
    
    echo -e "${GREEN}✓ 文档站启动成功${NC}"
    echo -e "${YELLOW}📍 本地访问: ${CYAN}http://localhost:3001${NC}"
    echo -e "${YELLOW}🌐 外网访问: ${CYAN}http://${LOCAL_IP}:3001${NC}"
    echo -e "${YELLOW}📚 项目文档: ${CYAN}http://localhost:3001${NC}"
    echo -e "${YELLOW}📝 日志文件: ${CYAN}${LOG_DIR}/docs.log${NC}"
    echo ""
    
    rotate_log "${LOG_DIR}/docs.log"
    npm start >> "${LOG_DIR}/docs.log" 2>&1
}

# 停止服务
stop_service() {
    SERVICE_NAME=$1
    PORT=$2
    PID=$(lsof -ti:$PORT 2>/dev/null || true)
    
    if [ -n "$PID" ]; then
        echo -e "${YELLOW}停止 ${SERVICE_NAME} (PID: $PID)...${NC}"
        kill -TERM $PID 2>/dev/null || true
        
        # 等待进程结束
        for i in {1..10}; do
            if ! kill -0 $PID 2>/dev/null; then
                echo -e "${GREEN}✓ ${SERVICE_NAME} 已停止${NC}"
                return 0
            fi
            sleep 1
        done
        
        # 强制杀死
        echo -e "${YELLOW}强制停止 ${SERVICE_NAME}...${NC}"
        kill -KILL $PID 2>/dev/null || true
        echo -e "${GREEN}✓ ${SERVICE_NAME} 已强制停止${NC}"
    else
        echo -e "${CYAN}ℹ️  ${SERVICE_NAME} 未运行${NC}"
    fi
}

# 停止所有服务
stop_all() {
    echo -e "${YELLOW}停止所有 X-Store 服务...${NC}"
    echo ""
    
    stop_service "后端服务" 8082
    stop_service "C 端前台" 3000
    stop_service "管理后台" 5174
    stop_service "文档站点" 3001
    
    echo ""
    echo -e "${GREEN}所有服务已停止${NC}"
}

# 检查开发环境
check_dev_env() {
    echo -e "${BLUE}检查开发环境...${NC}"
    echo ""
    
    # 检查 Go
    if command -v go >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Go: $(go version)${NC}"
    else
        echo -e "${RED}✗ Go 未安装${NC}"
        echo -e "${CYAN}安装: https://go.dev/dl/${NC}"
    fi
    
    # 检查 Node.js
    if command -v node >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Node.js: $(node --version)${NC}"
    else
        echo -e "${RED}✗ Node.js 未安装${NC}"
        echo -e "${CYAN}安装: https://nodejs.org/${NC}"
    fi
    
    echo ""
    echo -e "${CYAN}💡 提示: PostgreSQL 和 Redis 需要独立安装和启动${NC}"
    echo -e "${CYAN}   数据库配置在 backend/config.yaml 中设置${NC}"
    echo ""
}

# 主逻辑
main() {
    MODE=${1:-all}
    
    case $MODE in
        backend)
            start_backend
            ;;
        frontend)
            start_frontend
            ;;
        admin)
            start_admin
            ;;
        docs)
            start_docs
            ;;
        all)
            echo -e "${YELLOW}一键启动所有服务（后端 + 前台 + 管理后台 + 文档站）...${NC}"
            echo ""

            # 在后台分别启动四个服务，每个服务在自己的子 shell 中运行，不互相干扰
            (
                cd "$ROOT_DIR"
                start_backend
            ) &
            BACKEND_PID=$!

            (
                cd "$ROOT_DIR"
                start_frontend
            ) &
            FRONTEND_PID=$!

            (
                cd "$ROOT_DIR"
                start_admin
            ) &
            ADMIN_PID=$!

            (
                cd "$ROOT_DIR"
                start_docs
            ) &
            DOCS_PID=$!

            echo ""
            echo -e "${GREEN}已在后台启动所有服务：${NC}"
            echo -e "  🖥️  后端服务 (Go + Gin)        PID: ${CYAN}${BACKEND_PID}${NC}  端口: ${CYAN}8082${NC}"
            echo -e "  🛍️  C 端商城 (Next.js)         PID: ${CYAN}${FRONTEND_PID}${NC}  端口: ${CYAN}3000${NC}"
            echo -e "  ⚙️  管理后台 (React + Antd)    PID: ${CYAN}${ADMIN_PID}${NC}  端口: ${CYAN}5174${NC}"
            echo -e "  📚 文档站 (Docusaurus)        PID: ${CYAN}${DOCS_PID}${NC}  端口: ${CYAN}3001${NC}"
            echo ""
            echo -e "${YELLOW}提示：如需停止所有服务，可运行: ${GREEN}./dev-start.sh stop${NC}"
            ;;
        stop)
            stop_all
            ;;
        restart)
            echo -e "${YELLOW}重启所有服务...${NC}"
            echo ""
            stop_all
            echo ""
            echo -e "${YELLOW}等待 3 秒后重新启动...${NC}"
            sleep 3
            echo ""
            echo -e "${YELLOW}现在请在新终端运行以下命令启动服务:${NC}"
            echo ""
            echo -e "${CYAN}🖥️  后端服务:${NC} ${GREEN}./dev-start.sh backend${NC}"
            echo -e "${CYAN}🛍️  C 端商城:${NC} ${GREEN}./dev-start.sh frontend${NC}"
            echo -e "${CYAN}⚙️  管理后台:${NC} ${GREEN}./dev-start.sh admin${NC}"
            echo -e "${CYAN}📚 项目文档:${NC} ${GREEN}./dev-start.sh docs${NC}"
            ;;
        check)
            check_dev_env
            ;;
        *)
            echo -e "${RED}未知参数: $MODE${NC}"
            echo -e "${YELLOW}使用方式: ./dev-start.sh [backend|frontend|admin|docs|all|stop|restart|check]${NC}"
            echo -e "${YELLOW}  backend/frontend/admin/docs - 启动单个服务${NC}"
            echo -e "${YELLOW}  all - 显示所有服务启动命令${NC}"
            echo -e "${YELLOW}  stop - 停止所有服务${NC}"
            echo -e "${YELLOW}  restart - 重启所有服务${NC}"
            echo -e "${YELLOW}  check - 检查开发环境${NC}"
            exit 1
            ;;
    esac
}

main "$@"
