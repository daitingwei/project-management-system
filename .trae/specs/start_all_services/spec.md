# 启动所有服务组件 Spec

## Why
当前系统只有部分服务运行，需要启动docker-compose.yaml中配置的所有组件，使系统完整运行。

## What Changes
根据docker-compose.yaml，需要启动以下未运行的服务：
- Redis - 缓存服务（端口6379）
- Etcd - 服务注册发现（端口2379）
- Kafka - 消息队列（端口9092）
- Elasticsearch - 搜索引擎（端口9200）
- Kibana - ES可视化（端口5601）
- Logstash - 日志收集（端口9600）
- MinIO - 对象存储（端口9009）

## Impact
- 所有微服务依赖的中间件将完整运行
- 系统功能更完整

## 当前运行状态
- ✅ MySQL主库 (3309)
- ✅ MySQL从库1 (3310)
- ✅ MySQL从库2 (3311)
- ✅ Nacos (8848)
- ✅ Jaeger (16686)
- ❌ Redis (6379) - 未运行
- ❌ Etcd (2379) - 未运行
- ❌ Kafka (9092) - 未运行
- ❌ Elasticsearch (9200) - 未运行
- ❌ Kibana (5601) - 未运行
- ❌ Logstash (9600) - 未运行
- ❌ MinIO (9009) - 未运行

## 需要启动的服务列表
1. Redis - 配置在project-user和project-project的config.yaml中
2. Etcd - 用于gRPC服务注册发现
3. Kafka - 消息队列（可选，取决于是否使用）
4. Elasticsearch - 搜索引擎（可选）
5. Kibana - ES可视化（可选）
6. Logstash - 日志收集（可选）
7. MinIO - 对象存储（代码中有引用）

## 优先级
- 高优先级：Redis, Etcd（代码中直接使用）
- 中优先级：MinIO（代码中有引用）
- 低优先级：Kafka, ES, Kibana, Logstash（根据需要）
