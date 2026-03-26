# RESTful-to-MCP

将现有的 RESTful API 转换为 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) 服务，让 AI 大模型能够直接调用你的 API。

## 功能特性

- **OpenAPI 文档转换** — 上传 Swagger 2.0 / OpenAPI 3.0 文档，自动解析并生成 MCP Server 和 Tools
- **表单创建** — 通过 API 表单手动创建 MCP Server 和 Tools，无需 OpenAPI 文档
- **MCP 网关** — 统一的 MCP 协议网关，支持 Streamable HTTP 传输方式
- **多数据库支持** — 支持 PostgreSQL 和 MySQL
- **服务注册** — 可选的 Nacos 服务注册与发现
- **Docker 部署** — 提供 Dockerfile 和 docker-compose 一键部署

## 架构

```
┌──────────────────────────────────────────────┐
│                 RESTful-to-MCP               │
├──────────────┬───────────────────────────────┤
│  API Layer   │  MCP Gateway (/gateway/mcp)   │
│  (Gin)       │  Management API (/v1/...)     │
├──────────────┼───────────────────────────────┤
│  Service     │  OpenAPI Parser & Converter   │
│  Layer       │  MCP Server Manager           │
│              │  MCP Tools Manager            │
├──────────────┼───────────────────────────────┤
│  Data Layer  │  PostgreSQL / MySQL + Redis   │
├──────────────┼───────────────────────────────┤
│  MCP Core    │  Streamable HTTP Transport    │
│              │  HTTP Proxy (Tool Execution)  │
└──────────────┴───────────────────────────────┘
```

## 快速开始

### 环境要求

- Go 1.24+
- PostgreSQL 14+ 或 MySQL 8+
- Redis 7+

### 本地运行

```bash
# 克隆项目
git clone https://github.com/lengyue1084/restful-to-mcp.git
cd restful-to-mcp

# 复制配置文件并修改
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml，填入你的数据库和 Redis 连接信息

# 安装 wire（依赖注入代码生成）
go install github.com/google/wire/cmd/wire@latest

# 生成依赖注入代码
make wire

# 运行
make run
```

服务默认启动在 `http://localhost:9002`。

### Docker 部署

```bash
# 一键启动（包含 PostgreSQL + Redis）
docker-compose up -d
```

如需自定义配置，编辑 `docker-compose.yml` 中的环境变量即可。

### 构建

```bash
# Windows
make build

# Linux
make build-linux

# Docker 镜像
docker build -t restful-to-mcp:latest .
```

## 使用方式

### 1. 通过 OpenAPI 文档创建 MCP Server

```bash
# 上传 OpenAPI 文档（支持 JSON/YAML，Swagger 2.0 和 OpenAPI 3.0）
curl -X POST http://localhost:9002/v1/openapi/upload \
  -F "file=@your-api-spec.yaml"
```

上传成功后会返回 MCP Server 的 UUID 和连接信息。

### 2. 通过表单创建 MCP Server

```bash
# 创建 MCP Server
curl -X POST http://localhost:9002/v1/mcpServer/createByForm \
  -H "Content-Type: application/json" \
  -d '{"name": "my-api-server", "description": "My API Server"}'

# 为 Server 添加 Tool
curl -X POST http://localhost:9002/v1/mcpServer/createMcpServerTool \
  -H "Content-Type: application/json" \
  -d '{"server_uuid": "xxx", "name": "get_user", "method": "GET", "url": "https://api.example.com/users/{id}"}'
```

### 3. 连接 MCP Server

使用任意 MCP 客户端（如 Claude Desktop、Cursor 等）连接：

```
MCP Endpoint: http://localhost:9002/gateway/{serverToken}/mcp
```

## API 接口

| 接口 | 说明 |
|------|------|
| `POST /v1/openapi/upload` | 上传 OpenAPI 文档创建 MCP Server |
| `POST /v1/openapi/updateForAuth` | 更新 API 认证信息 |
| `POST /v1/mcpServer/createByForm` | 表单创建 MCP Server |
| `POST /v1/mcpServer/updateByUUID` | 更新 MCP Server |
| `POST /v1/mcpServer/deleteMcpServerByUUID` | 删除 MCP Server |
| `POST /v1/mcpServer/getMcpServerInfoByUUID` | 获取 MCP Server 信息 |
| `POST /v1/mcpServer/getMcpConnectTokenByUUID` | 获取连接 Token |
| `POST /v1/mcpServer/getMcpServerTools` | 获取 Server 的工具列表 |
| `POST /v1/mcpServer/createMcpServerTool` | 创建工具 |
| `POST /v1/mcpServer/updateMcpServerTool` | 更新工具 |
| `POST /v1/mcpServer/testMcpServerTool` | 测试工具连通性 |
| `POST /gateway/mcp` | MCP 网关入口 |
| `POST /gateway/{serverToken}/mcp` | 指定 Server 的 MCP 网关 |

## 配置说明

配置文件位于 `configs/config.yaml`，支持通过环境变量覆盖。

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `MCP_SERVER_HTTP_PORT` | HTTP 服务端口 | 9002 |
| `MCP_DATA_DATABASE_PG_SOURCE` | PostgreSQL 连接串 | - |
| `MCP_DATA_REDIS_ADDR` | Redis 地址 | 127.0.0.1 |
| `MCP_DATA_REDIS_PORT` | Redis 端口 | 6379 |
| `MCP_DATA_REDIS_PASSWORD` | Redis 密码 | - |
| `MCP_LOG_LEVEL` | 日志级别 | info |

完整配置项请参考 `configs/config.yaml.example`。

## 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)
- **依赖注入**: [Wire](https://github.com/google/wire)
- **MCP SDK**: [go-mcp](https://github.com/ThinkInAIXYZ/go-mcp)
- **OpenAPI 解析**: [kin-openapi](https://github.com/getkin/kin-openapi)
- **数据库 ORM**: [GORM](https://gorm.io/)
- **缓存**: Redis

## 项目结构

```
restful-to-mcp/
├── api/                    # API handler 层
├── cmd/                    # 程序入口 & Wire 注入
├── configs/                # 配置文件
├── internal/
│   ├── biz/                # 业务逻辑层
│   ├── conf/               # 配置解析
│   ├── data/               # 数据访问层
│   │   ├── database/       # 数据库连接
│   │   └── model/          # 数据模型
│   ├── mcp/
│   │   ├── proxy/          # HTTP 代理（执行 API 调用）
│   │   ├── server/         # MCP Server 管理
│   │   └── transformer/    # OpenAPI → MCP 转换器
│   ├── pkg/                # 内部工具包
│   └── service/            # 服务层
├── middleware/             # 中间件
├── pkg/                    # 公共工具包
├── router/                 # 路由定义
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

## License

[MIT](LICENSE)
