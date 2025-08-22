package config

import (
	"time"
)

// MCPStartupPolicy represents the startup policy for MCP servers

type (
	MCPStartupPolicy string
	AuthMode         string
	AuthPosition     string
)

const (
	PolicyOnStart  MCPStartupPolicy = "onStart"
	PolicyOnDemand MCPStartupPolicy = "onDemand"
)

func (m MCPStartupPolicy) Sting() string {
	switch m {
	case PolicyOnStart:
		return "onStart"
	case PolicyOnDemand:
		return "onDemand"
	default:
		return "unknown"
	}
}

const (
	// AuthModeApiKey 必选. security scheme 的类型。有效值包括 "apiKey", "http", "oauth2", "openIdConnect".
	AuthModeApiKey AuthMode = "apiKey"
	AuthModeHttp   AuthMode = "http"
)

func (a AuthMode) String() string {
	switch a {
	case AuthModeApiKey:
		return "apiKey"
	case AuthModeHttp:
		return "http"
	default:
		return "unknown"
	}
}

var (
	AuthPositionHeader AuthPosition = "header"
	AuthPositionQuery  AuthPosition = "query"
	AuthPositionCookie AuthPosition = "cookie"
)

func (a AuthPosition) String() string {
	switch a {
	case AuthPositionHeader:
		return "header"
	case AuthPositionQuery:
		return "query"
	case AuthPositionCookie:
		return "cookie"
	default:
		return "unknown"
	}
}

// MCPServer 表示 MCP 服务器的数据结构
type MCPServer struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name" yaml:"name"`
	UUID        string            `json:"uuid" yaml:"uuid"`
	Description string            `json:"description" yaml:"description"`
	Urls        []string          `json:"urls,omitempty" yaml:"urls,omitempty"`
	CreatedAt   time.Time         `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" yaml:"updatedAt"`
	Config      map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
	Auth        []*Auth           `json:"auth"`
	Tools       []*ToolConfig     `json:"tools,omitempty" yaml:"tools,omitempty"`
	Version     string            `json:"version"`
	//CORS        CORSConfig `json:"cors"` //暂时不考虑
}

// CORSConfig 表示 CORS（跨域资源共享）的配置结构
//type CORSConfig struct {
//	AllowOrigins     []string `json:"allowOrigins,omitempty" yaml:"allowOrigins,omitempty"`
//	AllowMethods     []string `json:"allowMethods,omitempty" yaml:"allowMethods,omitempty"`
//	AllowHeaders     []string `json:"allowHeaders,omitempty" yaml:"allowHeaders,omitempty"`
//	ExposeHeaders    []string `json:"exposeHeaders,omitempty" yaml:"exposeHeaders,omitempty"`
//	AllowCredentials bool     `json:"allowCredentials" yaml:"allowCredentials"`
//}

// ToolConfig name：工具的唯一标识符
// title：用于显示目的的可选的、人类可读的工具名称。
// description：人类可读的功能描述
// inputSchema：定义预期参数的 JSON Schema
// outputSchema：可选 JSON Schema 定义预期输出结构
// annotations：描述工具行为的可选属性
// ToolConfig 表示工具的配置结构
type ToolConfig struct {
	// 工具名称,需要考虑如果没有的时候需要新构建一个
	Name string `json:"name" yaml:"name"`
	// 工具描述
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// 请求方法
	Method string `json:"method" yaml:"method"`
	// 请求端点
	Endpoint string `json:"endpoint" yaml:"endpoint"`
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

// Auth 表示认证的配置结构
type Auth struct {
	// 认证模式
	Mode         AuthMode     `json:"mode"`         //Any http, apiKey, oauth2, openIdConnect
	Description  string       `json:"description"`  //Any
	Name         string       `json:"name"`         //apiKey
	In           AuthPosition `json:"in"`           //apiKey
	Scheme       string       `json:"scheme"`       //http
	BearerFormat string       `json:"bearerFormat"` //http ("bearer")
	//Flows        *Flows `json:"flows,omitempty"`
	//OpenIdConnectUrl string `json:"openIdConnectUrl`
}

//// ToToolSchema 将 ToolConfig 转换为 ToolSchema
//func (t *ToolConfig) ToToolSchema() ToolSchema {
//	// 创建输入模式的属性映射
//	properties := make(map[string]any)
//	// 存储必填参数名称列表
//	required := make([]string, 0)
//	// 遍历所有参数，构建属性映射
//	for _, arg := range t.Args {
//		property := map[string]any{
//			"type":        arg.Type,
//			"description": arg.Description,
//		}
//
//		// 如果参数类型为数组，处理数组子项配置
//		if arg.Type == "array" {
//			items := make(map[string]any)
//			if len(arg.Items.Enum) > 0 {
//				items["enum"] = tool.Union(arg.Items.Enum)
//			} else {
//				items["type"] = arg.Items.Type
//				// 如果子项是对象类型，递归处理其属性
//				if arg.Items.Properties != nil {
//					items["properties"] = arg.Items.Properties
//				}
//			}
//			property["items"] = items
//		}
//
//		properties[arg.Name] = property
//		// 如果参数是必填的，将其名称添加到必填列表中
//		if arg.Required {
//			required = append(required, arg.Name)
//		}
//	}
//
//	// 如果存在已有的输入模式，将其合并到属性映射中
//	if t.InputSchema != nil {
//		for k, v := range t.InputSchema {
//			properties[k] = v
//		}
//	}
//	var annotations *ToolAnnotations
//	// 如果存在注解信息，解析注解
//	if t.Annotations != nil {
//		annotations = &ToolAnnotations{
//			Title:           tool.GetString(t.Annotations, "title", ""),
//			DestructiveHint: tool.GetBool(t.Annotations, "destructiveHint", true),
//			IdempotentHint:  tool.GetBool(t.Annotations, "idempotentHint", false),
//			OpenWorldHint:   tool.GetBool(t.Annotations, "openWorldHint", true),
//			ReadOnlyHint:    tool.GetBool(t.Annotations, "readOnlyHint", false),
//		}
//	}
//
//	// 返回转换后的 ToolSchema
//	return ToolSchema{
//		Name:        t.Name,
//		Description: t.Description,
//		InputSchema: ToolInputSchema{
//			Type:       "object",
//			Properties: properties,
//			Required:   required,
//		},
//		Annotations: annotations,
//	}
//}
