#!/bin/bash

# 启动所有后端服务
echo "Starting project-user..."
cd /Users/daitingwei/Desktop/第二版/后端/project-user
nohup ./project-user > /tmp/project-user.log 2>&1 &
echo "project-user PID: $!"

sleep 2

echo "Starting project-project..."
cd /Users/daitingwei/Desktop/第二版/后端/project-project
nohup ./project-project > /tmp/project-project.log 2>&1 &
echo "project-project PID: $!"

sleep 2

echo "Starting project-api..."
cd /Users/daitingwei/Desktop/第二版/后端/project-api
nohup ./project-api > /tmp/project-api.log 2>&1 &
echo "project-api PID: $!"

echo ""
echo "All services started!"
echo "Logs:"
echo "  project-user:  /tmp/project-user.log"
echo "  project-project: /tmp/project-project.log"
echo "  project-api: /tmp/project-api.log"
