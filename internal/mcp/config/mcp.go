package config

import (
	"flow-bridge-mcp/pkg/tool"
	"time"
)

// MCPStartupPolicy represents the startup policy for MCP servers
type MCPStartupPolicy string

const (
	// PolicyOnStart represents the policy to connect on server start
	PolicyOnStart MCPStartupPolicy = "onStart"
	// PolicyOnDemand represents the policy to connect when needed
	PolicyOnDemand MCPStartupPolicy = "onDemand"
)

// ActionType represents the type of action performed on a configuration
type ActionType string

const (
	// ActionCreate represents a create action
	ActionCreate ActionType = "Create"
	// ActionUpdate represents an update action
	ActionUpdate ActionType = "Update"
	// ActionDelete represents a delete action
	ActionDelete ActionType = "Delete"
	// ActionRevert represents a revert action
	ActionRevert ActionType = "Revert"
)

type AuthMode string

const (
	AuthModeOAuth2 AuthMode = "oauth2"
)

// MCPServer 表示 MCP 服务器的数据结构
type MCPServer struct {
	// 服务器名称
	Name string `json:"name" yaml:"name"`
	// MCP 配置内容
	Content MCPConfig `json:"content" yaml:"content"`
	// 服务器创建时间
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	// 服务器更新时间
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// MCPConfig 表示 MCP 的配置结构
type MCPConfig struct {
	// 配置名称
	Name string `json:"name" yaml:"name"`
	// 租户信息
	Tenant string `json:"tenant" yaml:"tenant"`
	// 配置创建时间
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	// 配置更新时间
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
	// 配置删除时间，非零值表示所有信息已被删除
	DeletedAt time.Time `json:"deletedAt,omitempty" yaml:"deletedAt,omitempty"`
	// 路由配置列表
	Routers []RouterConfig `json:"routers,omitempty" yaml:"routers,omitempty"`
	// 服务器配置列表
	Servers []ServerConfig `json:"servers,omitempty" yaml:"servers,omitempty"`
	// 工具配置列表
	Tools []ToolConfig `json:"tools,omitempty" yaml:"tools,omitempty"`
	// 提示配置列表
	Prompts []PromptConfig `json:"prompts,omitempty" yaml:"prompts,omitempty"`
	// 代理 MCP 服务器配置列表
	McpServers []MCPServerConfig `json:"mcpServers,omitempty" yaml:"mcpServers,omitempty"`
}

// RouterConfig 表示路由的配置结构
type RouterConfig struct {
	// 服务器地址
	Server string `json:"server" yaml:"server"`
	// 路由前缀
	Prefix string `json:"prefix" yaml:"prefix"`
	// SSE 路由前缀
	SSEPrefix string `json:"ssePrefix" yaml:"ssePrefix"`
	// CORS 配置
	CORS *CORSConfig `json:"cors,omitempty" yaml:"cors,omitempty"`
	// 认证配置
	Auth *Auth `json:"auth,omitempty" yaml:"auth,omitempty"`
}

// CORSConfig 表示 CORS（跨域资源共享）的配置结构
type CORSConfig struct {
	// 允许的源列表
	AllowOrigins []string `json:"allowOrigins,omitempty" yaml:"allowOrigins,omitempty"`
	// 允许的方法列表
	AllowMethods []string `json:"allowMethods,omitempty" yaml:"allowMethods,omitempty"`
	// 允许的请求头列表
	AllowHeaders []string `json:"allowHeaders,omitempty" yaml:"allowHeaders,omitempty"`
	// 暴露的响应头列表
	ExposeHeaders []string `json:"exposeHeaders,omitempty" yaml:"exposeHeaders,omitempty"`
	// 是否允许携带凭证
	AllowCredentials bool `json:"allowCredentials" yaml:"allowCredentials"`
}

// ProxyConfig 表示代理的配置结构
type ProxyConfig struct {
	// 代理主机地址
	Host string `json:"host" yaml:"host"`
	// 代理端口
	Port int `json:"port" yaml:"port"`
	// 代理类型，支持 http, https, socks5
	Type string `json:"type" yaml:"type"`
}

// ServerConfig 表示服务器的配置结构
type ServerConfig struct {
	// 服务器名称
	Name string `json:"name" yaml:"name"`
	// 服务器描述
	Description string `json:"description" yaml:"description"`
	// 允许的工具列表
	AllowedTools []string `json:"allowedTools,omitempty" yaml:"allowedTools,omitempty"`
	// 服务器配置键值对
	Config map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
}

// ToolConfig 表示工具的配置结构
type ToolConfig struct {
	// 工具名称
	Name string `json:"name" yaml:"name"`
	// 工具描述
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// 请求方法
	Method string `json:"method" yaml:"method"`
	// 请求端点
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// 代理配置
	Proxy *ProxyConfig `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	// 请求头键值对
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	// 参数配置列表
	Args []ArgConfig `json:"args,omitempty" yaml:"args,omitempty"`
	// 请求体内容
	RequestBody string `json:"requestBody"  yaml:"requestBody"`
	// 响应体内容
	ResponseBody string `json:"responseBody" yaml:"responseBody"`
	// 输入模式
	InputSchema map[string]any `json:"inputSchema,omitempty" yaml:"inputSchema,omitempty"`
	// 注解信息
	Annotations map[string]any `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// MCPServerConfig 表示 MCP 服务器的配置结构
type MCPServerConfig struct {
	// 服务器类型，支持 sse, stdio 和 streamable-http
	Type string `json:"type" yaml:"type"`
	// 服务器名称
	Name string `json:"name" yaml:"name"`
	// 用于 stdio 类型的命令
	Command string `json:"command,omitempty" yaml:"command,omitempty"`
	// 用于 stdio 类型的命令参数
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
	// 用于 stdio 类型的环境变量
	Env map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	// 用于 sse 和 streamable-http 类型的 URL
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
	// 启动策略，支持两种模式：onStart（启动时立即启动）和 onDemand（按需启动）
	Policy MCPStartupPolicy `json:"policy" yaml:"policy"`
	// 是否在 mcp-gateway 启动时安装此 MCP 服务器
	Preinstalled bool `json:"preinstalled" yaml:"preinstalled"`
}

// ArgConfig 表示参数的配置结构
type ArgConfig struct {
	// 参数名称
	Name string `json:"name" yaml:"name"`
	// 参数位置，支持 header, query, path, body
	Position string `json:"position" yaml:"position"`
	// 参数是否必填
	Required bool `json:"required" yaml:"required"`
	// 参数类型
	Type string `json:"type" yaml:"type"`
	// 参数描述
	Description string `json:"description" yaml:"description"`
	// 参数默认值
	Default string `json:"default" yaml:"default"`
	// 数组类型参数的子项配置
	Items ItemsConfig `json:"items,omitempty" yaml:"items,omitempty"`
}

// ItemsConfig 表示数组子项的配置结构
type ItemsConfig struct {
	// 子项类型
	Type string `json:"type" yaml:"type"`
	// 枚举值列表
	Enum []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	// 对象类型子项的属性
	Properties map[string]any `json:"properties,omitempty" yaml:"properties,omitempty"`
	// 嵌套子项配置
	Items *ItemsConfig `json:"items,omitempty" yaml:"items,omitempty"`
	// 必填属性列表
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`
}

// MCPConfigVersion 表示 MCP 配置的一个版本
type MCPConfigVersion struct {
	// 版本号
	Version int `json:"version" yaml:"version"`
	// 创建者
	CreatedBy string `json:"created_by" yaml:"created_by"`
	// 版本创建时间
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	// 操作类型，支持 Create, Update, Delete, Revert
	ActionType ActionType `json:"action_type" yaml:"action_type"`
	// 配置名称
	Name string `json:"name" yaml:"name"`
	// 租户信息
	Tenant string `json:"tenant" yaml:"tenant"`
	// 路由配置信息
	Routers string `json:"routers" yaml:"routers"`
	// 服务器配置信息
	Servers string `json:"servers" yaml:"servers"`
	// 工具配置信息
	Tools string `json:"tools" yaml:"tools"`
	// 提示配置信息
	Prompts string `json:"prompts" yaml:"prompts"`
	// 代理 MCP 服务器配置信息
	McpServers string `json:"mcp_servers" yaml:"mcp_servers"`
	// 指示此版本是否当前激活
	IsActive bool `json:"is_active" yaml:"is_active"`
	// 配置内容的哈希值
	Hash string `json:"hash" yaml:"hash"`
}

// Auth 表示认证的配置结构
type Auth struct {
	// 认证模式
	Mode AuthMode `json:"mode" yaml:"mode"`
}

// PromptConfig 表示提示的配置结构
type PromptConfig struct {
	// 提示名称
	Name string `json:"name" yaml:"name"`
	// 提示描述
	Description string `json:"description" yaml:"description"`
	// 提示参数列表
	Arguments []PromptArgument `json:"arguments" yaml:"arguments"`
	// 提示响应列表
	PromptResponse []PromptResponse `json:"promptResponse,omitempty" yaml:"promptResponse,omitempty"`
}

// PromptArgument 表示提示参数的配置结构
type PromptArgument struct {
	// 参数名称
	Name string `json:"name" yaml:"name"`
	// 参数描述
	Description string `json:"description" yaml:"description"`
	// 参数是否必填
	Required bool `json:"required" yaml:"required"`
}

// PromptResponse 表示提示响应的配置结构
type PromptResponse struct {
	// 角色信息
	Role string `json:"role" yaml:"role"`
	// 响应内容
	Content PromptResponseContent `json:"content" yaml:"content"`
}

// PromptResponseContent 表示提示响应内容的配置结构
type PromptResponseContent struct {
	// 内容类型
	Type string `json:"type" yaml:"type"`
	// 内容文本
	Text string `json:"text" yaml:"text"`
}

// ToToolSchema 将 ToolConfig 转换为 ToolSchema
func (t *ToolConfig) ToToolSchema() ToolSchema {
	// 创建输入模式的属性映射
	properties := make(map[string]any)
	// 存储必填参数名称列表
	required := make([]string, 0)
	// 遍历所有参数，构建属性映射
	for _, arg := range t.Args {
		property := map[string]any{
			"type":        arg.Type,
			"description": arg.Description,
		}

		// 如果参数类型为数组，处理数组子项配置
		if arg.Type == "array" {
			items := make(map[string]any)
			if len(arg.Items.Enum) > 0 {
				items["enum"] = tool.Union(arg.Items.Enum)
			} else {
				items["type"] = arg.Items.Type
				// 如果子项是对象类型，递归处理其属性
				if arg.Items.Properties != nil {
					items["properties"] = arg.Items.Properties
				}
			}
			property["items"] = items
		}

		properties[arg.Name] = property
		// 如果参数是必填的，将其名称添加到必填列表中
		if arg.Required {
			required = append(required, arg.Name)
		}
	}

	// 如果存在已有的输入模式，将其合并到属性映射中
	if t.InputSchema != nil {
		for k, v := range t.InputSchema {
			properties[k] = v
		}
	}
	var annotations *ToolAnnotations
	// 如果存在注解信息，解析注解
	if t.Annotations != nil {
		annotations = &ToolAnnotations{
			Title:           tool.GetString(t.Annotations, "title", ""),
			DestructiveHint: tool.GetBool(t.Annotations, "destructiveHint", true),
			IdempotentHint:  tool.GetBool(t.Annotations, "idempotentHint", false),
			OpenWorldHint:   tool.GetBool(t.Annotations, "openWorldHint", true),
			ReadOnlyHint:    tool.GetBool(t.Annotations, "readOnlyHint", false),
		}
	}

	// 返回转换后的 ToolSchema
	return ToolSchema{
		Name:        t.Name,
		Description: t.Description,
		InputSchema: ToolInputSchema{
			Type:       "object",
			Properties: properties,
			Required:   required,
		},
		Annotations: annotations,
	}
}

// ToPromptSchema 将 PromptConfig 转换为 PromptSchema
func (t *PromptConfig) ToPromptSchema() PromptSchema {
	// 初始化提示参数模式列表
	args := make([]PromptArgumentSchema, len(t.Arguments))
	// 遍历所有提示参数，构建提示参数模式列表
	for i, a := range t.Arguments {
		args[i] = PromptArgumentSchema{
			Name:        a.Name,
			Description: a.Description,
			Required:    a.Required,
		}
	}
	// 初始化提示响应模式列表
	var responses []PromptResponseSchema
	// 遍历所有提示响应，构建提示响应模式列表
	for _, r := range t.PromptResponse {
		responses = append(responses, PromptResponseSchema{
			Role: r.Role,
			Content: PromptResponseContentSchema{
				Type: r.Content.Type,
				Text: r.Content.Text,
			},
		})
	}
	// 返回转换后的 PromptSchema
	return PromptSchema{
		Name:           t.Name,
		Description:    t.Description,
		Arguments:      args,
		PromptResponse: responses,
	}
}
