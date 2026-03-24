#!/bin/bash
# 停止所有服务的脚本

echo "正在停止所有服务..."
docker-compose -f all-services.yml down

echo "所有服务已停止！"