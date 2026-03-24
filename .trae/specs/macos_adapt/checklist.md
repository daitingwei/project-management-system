# Checklist - macOS环境适配与配置修复验证

## 配置验证

- [ ] .env文件使用macOS路径格式
- [ ] .env文件包含MySQL密码配置 MYSQL_PASSWORD=dtw258989971
- [ ] docker-compose.yaml中MySQL密码为dtw258989971
- [ ] docker-compose.yaml volumes配置正确
- [ ] project-project/config/config.yaml 数据库密码为dtw258989971
- [ ] project-project/config/config.yaml Redis配置正确
- [ ] project-project/config/config.yaml gRPC配置正确
- [ ] project-user/config/config.yaml 数据库密码为dtw258989971
- [ ] project-user/config/config.yaml Redis配置正确
- [ ] project-user/config/config.yaml gRPC配置正确
- [ ] project-api/config/config.yaml 配置正确
- [ ] 前端.env文件配置正确
- [ ] 前端vue.config.js配置正确

## MinIO验证

- [ ] 确认代码中是否使用MinIO
- [ ] 确认是否需要启用MinIO

## Docker服务验证

- [ ] MySQL容器正常启动
- [ ] Redis容器正常启动
- [ ] 其他容器正常启动

## 数据库验证

- [ ] 数据库连接成功
- [ ] 数据库调用成功
- [ ] 必要时重建数据库成功

## 目录验证

- [ ] MySQL数据目录存在
- [ ] Redis配置目录存在
- [ ] 其他数据目录存在
