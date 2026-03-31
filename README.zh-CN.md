# RESTful-to-MCP

[English](./README.md) | [中文](./README.zh-CN.md)

将现有 RESTful API 转换为 [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) 服务，让 AI 客户端可以通过统一的 MCP 网关调用你的接口。

## 功能特性

- 支持导入 Swagger 2.0 / OpenAPI 3.x 文档，自动生成 MCP Server 和 Tools。
- 支持通过前端管理界面手动创建 MCP Server 和 Tool。
- 提供统一的 MCP 网关 `/gateway/:serverToken/mcp`。
- 支持工具测试、鉴权配置、Server 详情查看与连接信息展示。
- 支持 PostgreSQL / MySQL 和 Redis。
- 可选接入 Nacos 服务注册。

## 界面截图

### Server 列表

![MCP Server 列表](./docs/images/mcp_server.png)

### Server 详情

![MCP Server 详情](./docs/images/mcp_server_detail.png)

### Tool 详情

![MCP Server Tool 详情](./docs/images/mcp_server_tools_detail.png)

### OpenAPI 导入

![OpenAPI 导入](./docs/images/openapi_import.png)

### Postman 测试示例

![Postman Test](./docs/images/postman_test.png)

## 架构

```text
RESTful API / OpenAPI
        |
        v
RESTful-to-MCP 管理 API (/v1/...)
        |
        +--> OpenAPI 解析 / 表单配置
        |
        +--> MCP Server 与 Tool 注册
        |
        v
MCP Gateway (/gateway/:serverToken/mcp)
        |
        v
MCP 客户端 (Claude Desktop、Cursor、Cherry Studio 等)
```

## 快速开始

### 环境要求

- Go 1.24+
- PostgreSQL 14+ 或 MySQL 8+
- Redis 7+

### 本地运行

```bash
git clone https://github.com/lengyue1084/restful-to-mcp.git
cd restful-to-mcp

cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml

go install github.com/google/wire/cmd/wire@latest
make wire
make run
```

默认后端地址：

```text
http://localhost:9002
```

### 前端运行

仓库中包含 React + Vite 管理前端，位于 [`web/`](./web)。

```bash
cd web
npm install
npm run dev
```

默认前端地址：

```text
http://localhost:3000
```

## MCP 调用方式

网关入口为：

```text
POST /gateway/:serverToken/mcp
```

其中 `serverToken` 目前同时支持：

- server UUID
- 页面生成的 `connectToken`

注意：

- 不要直接在浏览器里打开 MCP Endpoint。
- MCP 网关要求使用 `POST` JSON-RPC。
- 浏览器直接 `GET` 访问通常会返回 `404`。

`tools/list` 示例：

```bash
curl -X POST "http://localhost:9002/gateway/{serverToken}/mcp" \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'
```

## 主要接口

| 接口 | 说明 |
| --- | --- |
| `POST /v1/openapi/upload` | 导入 OpenAPI 文档 |
| `POST /v1/openapi/updateForAuth` | 更新鉴权配置 |
| `POST /v1/mcpServer/list` | 获取 MCP Server 列表 |
| `POST /v1/mcpServer/createByForm` | 手动创建 Server |
| `POST /v1/mcpServer/updateMcpServerByForm` | 更新手动创建的 Server |
| `POST /v1/mcpServer/getMcpServerInfoByUUID` | 获取 Server 详情 |
| `POST /v1/mcpServer/getMcpConnectTokenByUUID` | 生成 connectToken |
| `POST /v1/mcpServer/getMcpServerToolsByUUID` | 获取某个 Server 下的工具 |
| `POST /v1/mcpServer/createMcpServerTool` | 创建工具 |
| `POST /v1/mcpServer/updateMcpServerTool` | 更新工具 |
| `POST /v1/mcpServer/testMcpServerTool` | 测试工具请求 |
| `POST /gateway/:serverToken/mcp` | MCP 网关入口 |

## 项目结构

```text
api/                  HTTP API 层
cmd/                  程序入口与 wire
configs/              配置文件
docs/                 文档和截图
internal/             业务逻辑、数据层、MCP 运行时
middleware/           中间件
pkg/                  公共工具
router/               路由注册
web/                  前端管理界面
```

## License

[MIT](./LICENSE)
