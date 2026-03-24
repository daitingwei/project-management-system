#!/bin/bash

# 停止所有后端服务
pkill -f project-user
pkill -f project-project
pkill -f project-api

echo "All services stopped!"
