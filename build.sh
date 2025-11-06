#!/bin/bash
# 构建镜像
#docker build -t flow-bridge-mcp:latest .
docker build -t your-registry/restful-to-mcp:latest .


#docker login your-registry

 # 推送镜像到阿里云容器镜像仓库
docker push your-registry/restful-to-mcp:latest