#!/bin/bash
# 启动所有服务（不构建）

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "========== 启动所有服务 =========="
cd "$SCRIPT_DIR"
docker compose up -d

echo ""
echo "========== 启动主从复制配置 (如需要) =========="
docker compose --profile init up -d || true

echo ""
echo "========== 完成！=========="
echo "服务状态："
docker ps --format "table {{.Names}}\t{{.Status}}"
