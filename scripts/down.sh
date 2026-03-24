#!/bin/bash
# 停止并清理项目服务

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "========== 停止服务 =========="
cd "$SCRIPT_DIR"
docker compose down

echo ""
echo "========== 清理已完成 =========="
