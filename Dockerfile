# 构建阶段
FROM golang:alpine AS builder

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 替换 Alpine 镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 构建阶段只安装必要依赖
RUN apk add --no-cache gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 安装 wire 工具
RUN go install github.com/google/wire/cmd/wire@latest

# 生成 wire 依赖注入代码
RUN wire ./cmd/wire.go

# 构建应用
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o flow-bridge-mcp ./cmd

# 清理构建缓存
RUN go clean -cache -modcache

# 运行阶段
FROM alpine:3.19

# 替换 Alpine 镜像源为阿里云
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装 ca-certificates、时区数据、SSL 支持和网络工具
RUN apk --no-cache add ca-certificates tzdata curl wget net-tools iputils

# 设置时区环境变量
ENV TZ=Asia/Shanghai

# 创建非 root 用户
#RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /app

# 创建日志、临时和配置目录并设置权限
RUN mkdir -p /app/logs /app/tmp /app/configs

# 从构建阶段复制二进制文件
COPY --from=builder /app/flow-bridge-mcp .

# 复制本地 configs 目录到容器
COPY configs /app/configs
RUN ls -la /app/configs/

## 切换到非 root 用户
#USER appuser

# 暴露端口（根据实际应用需求调整）
EXPOSE 9002

# 运行应用，使用 configs 目录下的默认配置
ENTRYPOINT ["./flow-bridge-mcp"]
CMD ["-conf", "./configs"]
