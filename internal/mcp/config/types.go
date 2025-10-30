package config

import (
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"time"
)

const (
	SecurityTypeOr  string = "or"
	SecurityTypeAnd string = "and"
)

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
	AuthModeApiKey        AuthMode = "apiKey"
	AuthModeHttp          AuthMode = "http"
	AuthModeOauth2        AuthMode = "oauth2"
	AuthModeOpenIdConnect AuthMode = "openIdConnect"
	AuthModeEmpty         AuthMode = ""
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
	ID           uint              `json:"id"`
	Name         string            `json:"name" yaml:"name"`
	UUID         string            `json:"uuid" yaml:"uuid"`
	Description  string            `json:"description" yaml:"description"`
	Urls         []string          `json:"urls,omitempty" yaml:"urls,omitempty"`
	CreatedAt    time.Time         `json:"createdAt" yaml:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt" yaml:"updatedAt"`
	Config       map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
	Tools        []*ToolConfig     `json:"tools,omitempty" yaml:"tools,omitempty"`
	Version      string            `json:"version"`
	SecurityList []*Security       `json:"securityList"` //
	AllTools     []string          `json:"allTools"`
	//CORS        CORSConfig `json:"cors"` //暂时不考虑
}

// ToolConfig name：工具的唯一标识符
// title：用于显示目的的可选的、人类可读的工具名称。
// description：人类可读的功能描述
// inputSchema：定义预期参数的 JSON Schema
// outputSchema：可选 JSON Schema 定义预期输出结构
// annotations：描述工具行为的可选属性
// ToolConfig 表示工具的配置结构
type ToolConfig struct {
	ID uint `json:"id"`

	UUID string `json:"uuid"`
	// 工具名称,需要考虑如果没有的时候需要新构建一个
	Name string `json:"name"`
	// 工具描述
	Description string `json:"description,omitempty"`
	// 请求方法
	Method string `json:"method"`
	// 请求端点
	Endpoint string `json:"endpoint"`
	// 请求头键值对
	Headers map[string]string `json:"headers,omitempty"`
	// 参数配置列表
	Args []ArgConfig `json:"args,omitempty"`
	// 请求体内容
	RequestBody string `json:"requestBody"`
	// 响应模式, 目前只支持 application/json
	ContentType string `json:"contentType"`
	// 响应体内容
	ResponseBody string `json:"responseBody"`
	// 输入模式
	ToolSchema *protocol.InputSchema `json:"tool_schema,omitempty"`
	// 注解信息
	Annotations map[string]any `json:"annotations,omitempty"`
	//是否需要认证，为原有接口的鉴权
	Security *Security `json:"security"`
	// 认证模式
	SecurityMode AuthMode `json:"securityMode"`
	// 认证级别
	SecurityLevel SecurityLevel `json:"authLevel"` // 1 api，2 path 3 doc
	// 是否显示,如果不符合认证信息，则不显示
	IsShow bool `json:"isShow"` //true 显示，false不显示
}

type Security struct {
	SecurityKey  string       `json:"securityKey"` //必选. security scheme 的类型。有效值包括 "apiKey", "http", "oauth2", "openIdConnect",只兼容前两者
	Mode         AuthMode     `json:"mode"`        //Any http, apiKey, oauth2, openIdConnect
	Name         string       `json:"name"`
	Scheme       string       `json:"scheme"` //http ("bearer")
	In           AuthPosition `json:"in"`     //"query"、"header" 或 "cookie".
	Description  string       `json:"description"`
	BearerFormat string       `json:"bearerFormat"` // http ("bearer")
	//Flows        *Flows `json:"flows,omitempty"`
	//OpenIdConnectUrl string `json:"openIdConnectUrl`
}
type SecurityLevel int

const (
	SecurityLevelPublic SecurityLevel = 0 //没有授权
	SecurityLevelApi    SecurityLevel = 1 //api级别授权
	SecurityLevelDoc    SecurityLevel = 2 // doc级别授权
)

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
	// 枚举值列表
	Enum    []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Explode bool     `json:"explode" yaml:"explode"` // 是否展开参数
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

type (
	// ToolSchema represents a tool definition
	ToolSchema struct {
		// The name of the tool
		Name string `json:"name"`
		// A human-readable description of the tool
		Description string `json:"description"`
		// A JSON Schema object defining the expected parameters for the tool
		InputSchema protocol.InputSchema `json:"inputSchema"`
		// Annotations for the tool
		Annotations *ToolAnnotations `json:"annotations,omitempty"`
	}

	// https://github.com/modelcontextprotocol/modelcontextprotocol/blob/main/schema/2025-03-26/schema.json
	ToolAnnotations struct {
		DestructiveHint bool `json:"destructiveHint,omitempty"`
		IdempotentHint  bool `json:"idempotentHint,omitempty"`
		OpenWorldHint   bool `json:"openWorldHint,omitempty"`
		ReadOnlyHint    bool `json:"readOnlyHint,omitempty"`
		// A human-readable title for the tool.
		Title string `json:"title,omitempty"`
	}

	ToolInputSchema struct {
		Type       string         `json:"type"`
		Properties map[string]any `json:"properties"`
		Required   []string       `json:"required,omitempty"`
		Title      string         `json:"title"`
		Enum       []any          `json:"enum,omitempty"`
	}
)

// ArgsToInputSchema 将 ToolConfig 的 Args 转换为 protocol.InputSchema
func (t *ToolConfig) ArgsToInputSchema() *protocol.InputSchema {
	// 如果没有参数，返回空的 InputSchema
	if len(t.Args) == 0 {
		return &protocol.InputSchema{
			Type:       protocol.Object,
			Properties: make(map[string]*protocol.Property),
			Required:   make([]string, 0),
		}
	}

	schema := &protocol.InputSchema{
		Type:       protocol.Object,
		Properties: make(map[string]*protocol.Property),
		Required:   make([]string, 0),
	}

	// 遍历所有参数，构建属性映射
	for _, arg := range t.Args {
		property := t.ConvertArgToProperty(arg)
		schema.Properties[arg.Name] = property

		// 如果参数是必填的，将其名称添加到必填列表中
		if arg.Required {
			schema.Required = append(schema.Required, arg.Name)
		}
	}

	return schema
}

// ConvertArgToProperty 将 ArgConfig 转换为 protocol.Property
func (t *ToolConfig) ConvertArgToProperty(arg ArgConfig) *protocol.Property {
	property := &protocol.Property{
		Type:        convertStringToDataType(arg.Type),
		Description: arg.Description,
		Enum:        arg.Enum,
	}

	// 处理数组类型参数
	if arg.Type == "array" {
		property.Items = t.convertItemsToProperty(arg.Items)
	}

	// 处理对象类型参数
	if arg.Type == "object" {
		property.Properties = t.convertProperties(arg.Items.Properties)
		if len(arg.Items.Required) > 0 {
			property.Required = arg.Items.Required
		}
	}

	return property
}

// convertItemsToProperty 将 ItemsConfig 转换为 protocol.Property
func (t *ToolConfig) convertItemsToProperty(items ItemsConfig) *protocol.Property {
	if items.Type == "" {
		return nil
	}

	property := &protocol.Property{
		Type: convertStringToDataType(items.Type),
		Enum: items.Enum,
	}

	// 递归处理嵌套数组
	if items.Items != nil {
		property.Items = t.convertItemsToProperty(*items.Items)
	}

	// 处理对象数组
	if items.Type == "object" {
		property.Properties = t.convertProperties(items.Properties)
		if len(items.Required) > 0 {
			property.Required = items.Required
		}
	}

	return property
}

// convertProperties 将 map[string]any 转换为 map[string]*protocol.Property
func (t *ToolConfig) convertProperties(properties map[string]any) map[string]*protocol.Property {
	if len(properties) == 0 {
		return nil
	}

	result := make(map[string]*protocol.Property)

	for key, value := range properties {
		if propMap, ok := value.(map[string]any); ok {
			prop := &protocol.Property{}

			// 转换类型
			if typeStr, ok := propMap["type"].(string); ok {
				prop.Type = convertStringToDataType(typeStr)
			}

			// 转换描述
			if desc, ok := propMap["description"].(string); ok {
				prop.Description = desc
			}

			// 转换枚举
			if enum, ok := propMap["enum"].([]string); ok {
				prop.Enum = enum
			} else if enum, ok := propMap["enum"].([]interface{}); ok {
				// 处理 []interface{} 类型的枚举
				enumStr := make([]string, len(enum))
				for i, v := range enum {
					enumStr[i] = fmt.Sprintf("%v", v)
				}
				prop.Enum = enumStr
			}

			// 递归处理嵌套属性
			if props, ok := propMap["properties"].(map[string]any); ok {
				prop.Properties = t.convertProperties(props)
			}

			// 处理 required 字段
			if required, ok := propMap["required"].([]string); ok {
				prop.Required = required
			} else if required, ok := propMap["required"].([]interface{}); ok {
				requiredStr := make([]string, len(required))
				for i, v := range required {
					requiredStr[i] = fmt.Sprintf("%v", v)
				}
				prop.Required = requiredStr
			}

			result[key] = prop
		}
	}

	return result
}

// convertStringToDataType 将字符串类型转换为 protocol.DataType
func convertStringToDataType(typeStr string) protocol.DataType {
	switch typeStr {
	case "string":
		return protocol.String
	case "number":
		return protocol.Number
	case "integer":
		return protocol.Integer
	case "boolean":
		return protocol.Boolean
	case "array":
		return protocol.Array
	case "object":
		return protocol.ObjectT
	case "null":
		return protocol.Null
	default:
		return protocol.String
	}
}

// ToProtocolTool 将 ToolConfig 转换为 protocol.Tool
func (t *ToolConfig) ToProtocolTool() (*protocol.Tool, error) {
	inputSchema := t.ArgsToInputSchema()

	tool := &protocol.Tool{
		Name:        t.Name,
		Description: t.Description,
		InputSchema: *inputSchema,
	}

	// 处理注解信息
	if t.Annotations != nil {
		annotations := &protocol.ToolAnnotations{}

		if title, ok := t.Annotations["title"].(string); ok {
			annotations.Title = title
		}

		if readOnly, ok := t.Annotations["readOnlyHint"].(bool); ok {
			annotations.ReadOnlyHint = &readOnly
		}

		if destructive, ok := t.Annotations["destructiveHint"].(bool); ok {
			annotations.DestructiveHint = &destructive
		}

		if idempotent, ok := t.Annotations["idempotentHint"].(bool); ok {
			annotations.IdempotentHint = &idempotent
		}

		if openWorld, ok := t.Annotations["openWorldHint"].(bool); ok {
			annotations.OpenWorldHint = &openWorld
		}

		tool.Annotations = annotations
	}

	return tool, nil
}
