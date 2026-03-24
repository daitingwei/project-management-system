#!/bin/bash
# 构建并启动所有服务

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WEB_DIR="$SCRIPT_DIR/web"

echo "========== 1. 编译 Go 服务 =========="

cd "$WEB_DIR/project-api"
echo "编译 project-api..."
go build -o target/project-api .

cd "$WEB_DIR/project-project"
echo "编译 project-project..."
go build -o target/project-project .

cd "$WEB_DIR/project-user"
echo "编译 project-user..."
go build -o target/project-user .

echo ""
echo "========== 2. 构建 Docker 镜像 =========="

cd "$WEB_DIR/project-api"
echo "构建 project-api 镜像..."
docker build -t project-api:latest .

cd "$WEB_DIR/project-project"
echo "构建 project-project 镜像..."
docker build -t project-project:latest .

cd "$WEB_DIR/project-user"
echo "构建 project-user 镜像..."
docker build -t project-user:latest .

echo ""
echo "========== 3. 启动所有服务 =========="
cd "$SCRIPT_DIR"
docker compose up -d

echo ""
echo "========== 4. 启动主从复制配置 (如需要) =========="
docker compose --profile init up -d || true

echo ""
echo "========== 完成！=========="
echo "服务状态："
docker ps --format "table {{.Names}}\t{{.Status}}"
