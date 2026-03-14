#!/bin/bash

# X-Store 项目启动脚本
# 使用方式: ./start.sh [backend|frontend|admin|all]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  X-Store 项目启动脚本${NC}"
echo -e "${GREEN}========================================${NC}"

# 检查 PostgreSQL 连接
check_postgres() {
    echo -e "${YELLOW}检查 PostgreSQL 连接...${NC}"
    if docker run --rm -e PGPASSWORD='Postgres@2026' postgres:15 psql -h host.docker.internal -p 5432 -U admin -d x_store -c "SELECT 1" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PostgreSQL 连接成功${NC}"
        return 0
    else
        echo -e "${RED}✗ PostgreSQL 连接失败${NC}"
        echo -e "${YELLOW}请确保 PostgreSQL 已启动并配置正确${NC}"
        return 1
    fi
}

# 检查 Redis 连接
check_redis() {
    echo -e "${YELLOW}检查 Redis 连接...${NC}"
    if redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Redis 连接成功${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠ Redis 未启动，部分功能可能受限${NC}"
        return 1
    fi
}

# 启动后端
start_backend() {
    echo -e "${GREEN}启动后端服务...${NC}"
    cd backend
    
    # 检查是否需要编译
    if [ ! -f "x-store-backend" ] || [ "cmd/main.go" -nt "x-store-backend" ]; then
        echo -e "${YELLOW}编译后端...${NC}"
        go build -o x-store-backend ./cmd/main.go
    fi
    
    # 启动后端
    echo -e "${GREEN}后端服务启动在 http://localhost:8082${NC}"
    ./x-store-backend
}

# 启动 C 端前台
start_frontend() {
    echo -e "${GREEN}启动 C 端前台...${NC}"
    cd frontend-store
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装依赖...${NC}"
        npm install
    fi
    
    echo -e "${GREEN}C 端前台启动在 http://localhost:3000${NC}"
    npm run dev
}

# 启动管理后台
start_admin() {
    echo -e "${GREEN}启动管理后台...${NC}"
    cd admin-panel
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}安装依赖...${NC}"
        npm install
    fi
    
    echo -e "${GREEN}管理后台启动在 http://localhost:5174${NC}"
    npm run dev
}

# 主逻辑
main() {
    MODE=${1:-all}
    
    # 检查依赖
    check_postgres || exit 1
    check_redis
    
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
        all)
            echo -e "${YELLOW}启动所有服务（需要多个终端）${NC}"
            echo -e "${YELLOW}请在不同终端分别运行:${NC}"
            echo -e "  ${GREEN}./start.sh backend${NC}  - 启动后端"
            echo -e "  ${GREEN}./start.sh frontend${NC} - 启动 C 端"
            echo -e "  ${GREEN}./start.sh admin${NC}    - 启动管理后台"
            ;;
        *)
            echo -e "${RED}未知参数: $MODE${NC}"
            echo -e "使用方式: ./start.sh [backend|frontend|admin|all]"
            exit 1
            ;;
    esac
}

main "$@"
