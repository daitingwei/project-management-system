#!/bin/bash
# 启动所有服务的脚本

echo "正在启动所有服务..."

# 创建必要的数据目录
mkdir -p data/mysql/data data/mysql/data-slave1 data/mysql/data-slave2 data/mysql/conf data/mysql/conf-slave1 data/mysql/conf-slave2 data/mysql/logs data/mysql/logs-slave1 data/mysql/logs-slave2
mkdir -p data/redis data/etcd/data data/nacos/data data/nacos/logs data/nacos/conf data/kafka/data data/jaeger data/minio/data

# 启动所有服务
docker-compose -f all-services.yml up -d

echo "所有服务启动完成！"
echo "检查服务状态："
docker-compose -f all-services.yml ps

echo ""
echo "服务端口映射："
echo "MySQL 主库: 3309 -> 3306"
echo "MySQL 从库1: 3310 -> 3306"
echo "MySQL 从库2: 3311 -> 3306"
echo "Redis: 6379 -> 6379"
echo "Etcd: 2379 -> 2379"
echo "Nacos: 8848 -> 8848"
echo "Kafka: 9092 -> 9092"
echo "Jaeger: 16686 -> 16686"
echo "MinIO: 9009 -> 9000"