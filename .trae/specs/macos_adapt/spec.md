# 第二版系统完善计划

## 为什么
第一版已完成docker和数据库基础框架，但存在以下问题需要修复：
1. 配置文件使用Windows路径格式，需要适配macOS环境
2. Docker中MySQL密码需要从root改为dtw258989971
3. 本地3309端口MySQL密码也需要改为dtw258989971（与Docker一致）
4. 需要确认minio是否需要配置
5. 没有数据库自动映射（表结构需要手动创建或重建）
6. 需要验证数据库连接和调用是否成功，必要时重建数据库

## 什么变化
- 修改.env配置为macOS路径格式
- 修改docker-compose.yaml中MySQL密码为dtw258989971
- 修改所有config.yaml中的数据库密码为dtw258989971
- 验证docker-compose.yaml的macOS兼容性
- 检查并修复各微服务配置文件
- 确认minio配置状态
- 测试数据库连接和调用
- 必要时重建数据库

## 影响
- 受影响的服务：project-api, project-project, project-user, project-common
- 受影响的配置：.env, docker-compose.yaml, 各微服务config.yaml
- 关键文件：
  - /Users/daitingwei/Desktop/第二版/web/.env
  - /Users/daitingwei/Desktop/第二版/web/docker-compose.yaml
  - /Users/daitingwei/Desktop/第二版/web/project-project/config/config.yaml
  - /Users/daitingwei/Desktop/第二版/web/project-user/config/config.yaml
  - /Users/daitingwei/Desktop/第二版/web/project-api/config/config.yaml

## 新增需求
### 需求：macOS环境适配
系统应能在macOS环境下正常运行，包括：
- docker容器路径正确映射
- 各微服务配置正确连接

#### 场景：Docker服务启动
- **给定**用户运行docker-compose
- **当**所有配置文件正确
- **则**所有服务正常启动

### 需求：MySQL密码统一修改
Docker和本地3309端口的MySQL密码都需要从root改为dtw258989971

#### 场景：修改MySQL密码
- **当**用户修改.env中MySQL密码并重启
- **则**docker-compose.yaml中MYSQL_ROOT_PASSWORD同步更新为dtw258989971
- **则**各微服务config.yaml中数据库密码同步更新为dtw258989971

### 需求：数据库连接验证
修改密码后需要验证数据库连接和调用是否成功

#### 场景：数据库连接测试
- **给定**所有配置文件已修改
- **当**启动后端服务
- **则**服务能成功连接数据库
- **或则**如果连接失败，删除并重建数据库

### 需求：MinIO配置确认
确认MinIO在系统中的用途，并进行相应配置

#### 场景：MinIO使用确认
- **给定**系统需要文件存储功能
- **则**MinIO配置保持不变
- **或则**如果不需要，可禁用MinIO相关功能

## 修改需求
### 需求：环境配置文件
- 将.env中Windows路径改为macOS路径格式
- 将MySQL密码从root改为dtw258989971
- 验证config.yaml中的配置与.env一致

## 移除需求
### 需求：旧的Windows配置
- 移除所有硬编码的Windows路径
