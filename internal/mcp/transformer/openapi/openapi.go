package openapi

// package openapi 定义了处理 OpenAPI 规范转换的包

// 导入必要的包
import (
	"context"
	"encoding/json"
	"flow-bridge-mcp/internal/mcp/config"
	"flow-bridge-mcp/internal/mcp/transformer"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
	"github.com/google/wire"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

var ProviderSet = wire.NewSet(
	NewConverter,
)

// Converter 结构体用于处理从 OpenAPI 规范到 MCP 配置的转换
// 目前结构体为空，可根据需求添加必要的字段
type Converter struct {
	// Add any necessary fields here
	log *logger.Logger
}

var _ transformer.Transformer = (*Converter)(nil)

// 定义 OpenAPI 版本常量
const (
	// OpenAPIVersion2 表示 OpenAPI 2.0 版本（也称为 Swagger 2.0）
	OpenAPIVersion2 = "2.0"
	// OpenAPIVersion3 表示 OpenAPI 3.0 版本
	OpenAPIVersion3 = "3.0"
	// OpenAPIVersion31 表示 OpenAPI 3.1 版本
	OpenAPIVersion31 = "3.1"
)

// NewConverter 创建一个新的 Converter 实例
// 返回一个指向 Converter 结构体的指针

func NewConverter(log *logger.Logger) transformer.Transformer {
	return &Converter{
		log: log,
	}
}
func (c *Converter) Metadata(ctx context.Context) *transformer.Metadata {

	return &transformer.Metadata{}
}

func (c *Converter) Validate(ctx context.Context, data []byte) error {

	//todo 分离验证方法
	return nil
}

// Convert 将 OpenAPI 规范数据转换为 MCP 配置
// 参数 specData 为 OpenAPI 规范的字节数据
// 返回 MCP 配置指针和可能出现的错误
func (c *Converter) Convert(ctx context.Context, specData []byte) (*config.MCPConfig, error) {

	var mcpUUID string
	UUID := ctx.Value("uuid")
	if UUID != nil {
		mcpUUID = UUID.(string)
	}

	// 检测 OpenAPI 版本
	version, err := c.DetectVersion(ctx, specData)
	if err != nil {
		return nil, err
	}

	// 根据 API 版本选择对应的处理逻辑
	if strings.HasPrefix(version, OpenAPIVersion2) {
		// 处理 Swagger 2.0 版本,统一转成3.0的版本处理
		return c.convertSwagger2(ctx, specData)
	}

	// 处理 OpenAPI 3.x 版本
	// 创建一个新的 OpenAPI 3.x 加载器
	loader := openapi3.NewLoader()

	// 如果是 OpenAPI 3.1 版本，允许外部引用
	if strings.HasPrefix(version, OpenAPIVersion31) {
		loader.IsExternalRefsAllowed = true
	}

	// 从字节数据加载 OpenAPI 文档
	doc, err := loader.LoadFromData(specData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI specification: %w", err)
	}

	// 如果是 OpenAPI 3.0 版本，验证文档的有效性
	if strings.HasPrefix(version, OpenAPIVersion3) {
		if err := doc.Validate(loader.Context); err != nil {
			return nil, fmt.Errorf("invalid OpenAPI specification: %w", err)
		}
	}

	// 生成一个 4 位的随机字符串
	//rs := tool.RandStringByLen(4)
	if mcpUUID == "" {
		mcpUUID = tool.RandStringByLen(4)
	}

	// 创建基础的 MCP 配置
	mcpConfig := &config.MCPConfig{
		//Name:      doc.Info.Title + "_" + rs,
		//Name:      doc.Info.Title + "_" + mcpUUID,
		Name:      doc.Info.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Routers:   make([]config.RouterConfig, 0),
		Servers:   make([]config.ServerConfig, 0),
		Tools:     make([]config.ToolConfig, 0),
	}

	// 创建服务器配置
	server := config.ServerConfig{
		Name:         mcpConfig.Name,
		Description:  doc.Info.Description,
		Config:       make(map[string]string),
		AllowedTools: make([]string, 0),
	}

	// 将服务器 URL 添加到配置中
	if len(doc.Servers) > 0 {
		// server默认服务器地址为第一个服务器地址
		server.Config["url"] = c.selectServer(doc.Servers)
	}

	// 为服务器创建一个默认的路由配置
	router := config.RouterConfig{
		Server: mcpConfig.Name,
		//Prefix: fmt.Sprintf("/mcp/%s", rs), // 为每个路由生成一个随机前缀
		Prefix: fmt.Sprintf("/mcp/%s", mcpUUID), // 为每个路由生成一个随机前缀,如果前端传来则使用前端的，如果没有则使用自己生成的
		CORS: &config.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization", "Mcp-Session-Id", "mcp-protocol-version"},
			ExposeHeaders:    []string{"Mcp-Session-Id", "mcp-protocol-version"},
			AllowCredentials: true,
		},
	}

	// 将文档中的路径转换为工具配置
	for path, pathItem := range doc.Paths.Map() {
		// 为每个 HTTP 方法创建一个工具配置
		for method, operation := range pathItem.Operations() {
			if method == "options" {
				continue // 跳过 CORS 预检请求
			}

			// 如果操作 ID 为空，则根据方法和路径生成一个操作 ID
			if operation.OperationID == "" {
				// 转换路径为操作 ID 格式，例如：/users/email/{email} -> users_email_argemail
				pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
				for i, part := range pathParts {
					if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
						pathParts[i] = "arg" + strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
					}
				}
				operation.OperationID = fmt.Sprintf("%s_%s", strings.ToLower(method), strings.Join(pathParts, "_"))
			}

			// 创建工具配置
			tool := config.ToolConfig{
				Name:         operation.OperationID,
				Description:  tool.FirstNonEmpty(operation.Description, operation.Summary),
				Method:       method,
				Endpoint:     fmt.Sprintf("{{.Config.url}}%s", path),
				Headers:      make(map[string]string),
				Args:         make([]config.ArgConfig, 0),
				ResponseBody: "{{.Response.Body}}", // 使用透传响应
			}

			// 添加默认请求头
			tool.Headers["Content-Type"] = "application/json"
			tool.Headers["Authorization"] = "{{.Request.Headers.Authorization}}"

			// 定义不同位置的参数切片
			var bodyArgs []config.ArgConfig
			var pathArgs []config.ArgConfig
			var queryArgs []config.ArgConfig
			var headerArgs []config.ArgConfig

			// 处理操作中的参数
			for _, param := range operation.Parameters {
				// 创建参数配置
				arg := config.ArgConfig{
					Name:        param.Value.Name,
					Position:    param.Value.In,
					Required:    param.Value.Required,
					Type:        "string", // 默认参数类型为字符串
					Description: param.Value.Description,
				}

				// 如果参数有 schema 定义，尝试获取参数类型和默认值
				if param.Value.Schema != nil && param.Value.Schema.Value != nil {
					if param.Value.Schema.Value.Type != nil {
						types := param.Value.Schema.Value.Type.Slice()
						if len(types) > 0 {
							arg.Type = types[0]
						}
					}
					if param.Value.Schema.Value.Default != nil {
						arg.Default = fmt.Sprintf("%v", param.Value.Schema.Value.Default)
					}
				}

				// 根据参数位置进行不同处理
				switch param.Value.In {
				case "path":
					// 路径参数总是必需的
					arg.Required = true
					pathArgs = append(pathArgs, arg)
					// 更新端点中的路径参数占位符
					tool.Endpoint = strings.ReplaceAll(tool.Endpoint, fmt.Sprintf("{%s}", arg.Name), fmt.Sprintf("{{.Args.%s}}", arg.Name))
				case "query":
					queryArgs = append(queryArgs, arg)
				case "header":
					tool.Headers[arg.Name] = fmt.Sprintf("{{.Args.%s}}", arg.Name)
					headerArgs = append(headerArgs, arg)
				}
			}

			// 处理请求体
			if operation.RequestBody != nil {
				// 获取请求体是否必需的标志,否则会400错误，如果为true可以给一个空的json {}
				requestBodyRequired := operation.RequestBody.Value.Required
				// 遍历请求体支持的内容类型
				for contentType, mediaType := range operation.RequestBody.Value.Content {
					if contentType == "application/json" {
						tool.RequestBody = contentType
						// 添加请求体参数
						if mediaType.Schema != nil {
							schema := mediaType.Schema.Value
							// 处理 schema 引用
							if mediaType.Schema.Ref != "" {
								refName := strings.TrimPrefix(mediaType.Schema.Ref, "#/components/schemas/")
								if refSchema, ok := doc.Components.Schemas[refName]; ok {
									schema = refSchema.Value
								}
							}

							// 如果 schema 有属性定义
							if schema.Properties != nil {
								for name, prop := range schema.Properties {
									// 跳过响应专用字段
									if strings.HasPrefix(name, "response") || name == "id" || name == "createdAt" {
										continue
									}

									// 创建请求体参数配置
									arg := config.ArgConfig{
										Name:        name,
										Position:    "body",
										Required:    requestBodyRequired || contains(schema.Required, name),
										Type:        "string", // 默认参数类型为字符串
										Description: prop.Value.Description,
									}

									// 如果属性有类型定义
									if prop.Value != nil && prop.Value.Type != nil {
										types := prop.Value.Type.Slice()
										if len(types) > 0 {
											arg.Type = types[0]
											// 如果是数组类型且有 items 定义
											if arg.Type == "array" && prop.Value.Items != nil && prop.Value.Items.Value != nil {
												arg.Items = buildNestedArg(prop.Value.Items.Value)
											}
										}
									}

									// 如果属性有默认值
									if prop.Value.Default != nil {
										arg.Default = fmt.Sprintf("%v", prop.Value.Default)
									}

									bodyArgs = append(bodyArgs, arg)
								}
							}
						}
					}
				}
			}

			// 合并所有参数
			tool.Args = append(tool.Args, pathArgs...)
			tool.Args = append(tool.Args, queryArgs...)
			tool.Args = append(tool.Args, bodyArgs...)
			tool.Args = append(tool.Args, headerArgs...)

			// 如果有请求体参数，构建请求体模板
			if len(bodyArgs) > 0 {
				var bodyTemplate strings.Builder
				bodyTemplate.WriteString("{\n")
				for i, arg := range bodyArgs {
					bodyTemplate.WriteString(fmt.Sprintf(`    "%s": {{ toJSON .Args.%s}}`, arg.Name, arg.Name))
					if i < len(bodyArgs)-1 {
						bodyTemplate.WriteString(",\n")
					} else {
						bodyTemplate.WriteString("\n")
					}
				}
				bodyTemplate.WriteString("}")
				tool.RequestBody = bodyTemplate.String()
			}

			// 将工具配置添加到 MCP 配置中
			mcpConfig.Tools = append(mcpConfig.Tools, tool)
			// 将工具名称添加到服务器允许的工具列表中
			server.AllowedTools = append(server.AllowedTools, tool.Name)
		}
	}

	// 将服务器配置添加到 MCP 配置中
	mcpConfig.Servers = append(mcpConfig.Servers, server)
	// 将路由配置添加到 MCP 配置中
	mcpConfig.Routers = append(mcpConfig.Routers, router)

	return mcpConfig, nil
}

// 在 MCP 转换器中根据条件选择 server
func (c *Converter) selectServer(servers []*openapi3.Server) string {
	// 根据请求头选择环境,暂时不考虑吧，后期有需要再做
	//if env := ctx.Value("env"); env == "staging" {
	//	return servers[1].URL  // 返回预发布环境URL
	//}
	return servers[0].URL // 默认生产环境
}

// DetectVersion 从规范数据中检测 OpenAPI 版本
// 参数 specData 为 OpenAPI 规范的字节数据
// 返回检测到的版本字符串和可能出现的错误
func (c *Converter) DetectVersion(_ context.Context, data []byte) (string, error) {
	var spec map[string]interface{}

	// 尝试用 JSON 解析规范数据
	if err := json.Unmarshal(data, &spec); err != nil {
		// 如果 JSON 解析失败，尝试用 YAML 解析
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return "", fmt.Errorf("failed to parse specification: %w", err)
		}
	}

	// 检查是否为 OpenAPI 3.x 版本
	if openapi, ok := spec["openapi"].(string); ok {
		return openapi, nil
	}

	// 检查是否为 Swagger 2.0 版本
	if swagger, ok := spec["swagger"].(string); ok {
		return swagger, nil
	}

	return "", fmt.Errorf("could not determine OpenAPI version")
}

// convertSwagger2 将 Swagger 2.0 规范转换为 OpenAPI 3.0 规范，然后再转换为 MCP 配置
// 参数 specData 为 Swagger 2.0 规范的字节数据
// 返回 MCP 配置指针和可能出现的错误
func (c *Converter) convertSwagger2(ctx context.Context, specData []byte) (*config.MCPConfig, error) {
	var swagger2Doc openapi2.T
	// 尝试用 JSON 解析 Swagger 2.0 文档
	if err := json.Unmarshal(specData, &swagger2Doc); err != nil {
		// 如果 JSON 解析失败，尝试用 YAML 解析
		if err := yaml.Unmarshal(specData, &swagger2Doc); err != nil {
			return nil, fmt.Errorf("failed to parse Swagger 2.0 specification: %w", err)
		}
	}

	// 将 Swagger 2.0 文档转换为 OpenAPI 3.0 文档
	openapi3Doc, err := openapi2conv.ToV3(&swagger2Doc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Swagger 2.0 to OpenAPI 3.0: %w", err)
	}

	// 将 OpenAPI 3.0 文档序列化为 JSON 字节数据
	openapi3Data, err := json.Marshal(openapi3Doc)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize OpenAPI 3.0 document: %w", err)
	}

	// 递归调用 Convert 方法处理 OpenAPI 3.0 数据
	return c.Convert(ctx, openapi3Data)
}

// ConvertFromJSON 将 JSON 格式的 OpenAPI 规范转换为 MCP 配置
// 参数 jsonData 为 JSON 格式的 OpenAPI 规范字节数据
// 返回 MCP 配置指针和可能出现的错误
func (c *Converter) ConvertFromJSON(ctx context.Context, jsonData []byte) (*config.MCPConfig, error) {
	return c.Convert(ctx, jsonData)
}

// ConvertFromYAML 将 YAML 格式的 OpenAPI 规范转换为 MCP 配置
// 参数 yamlData 为 YAML 格式的 OpenAPI 规范字节数据
// 返回 MCP 配置指针和可能出现的错误
func (c *Converter) ConvertFromYAML(ctx context.Context, yamlData []byte) (*config.MCPConfig, error) {
	return c.Convert(ctx, yamlData)
}

// ConvertWithOptions 将 OpenAPI 规范转换为 MCP 配置，可指定租户和前缀
// 参数 specData 为 OpenAPI 规范的字节数据
// 参数 tenant 为租户名称
// 参数 prefix 为前缀
// 返回 MCP 配置指针和可能出现的错误
func (c *Converter) ConvertWithOptions(ctx context.Context, specData []byte, tenant, prefix string) (*config.MCPConfig, error) {
	config, err := c.Convert(ctx, specData)
	if err != nil {
		return nil, err
	}
	// 去除前缀前的斜杠
	cleanPrefix := strings.TrimPrefix(prefix, "/")
	if tenant != "" && prefix != "" {
		if len(config.Routers) > 0 {
			// 生成一个 4 位的随机字符串
			rs := tool.RandStringByLen(4)
			config.Routers[0].Prefix = "/" + cleanPrefix + "/" + rs
		}
	} else if tenant != "" {
		if len(config.Routers) > 0 {
			// 自动生成前缀，逻辑与默认逻辑相同
			rs := tool.RandStringByLen(4)
			config.Routers[0].Prefix = "/" + rs
		}
	}
	return config, nil
}

// contains 检查字符串是否在字符串切片中
// 参数 slice 为字符串切片
// 参数 str 为要检查的字符串
// 返回布尔值，表示是否包含
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// buildNestedArg 根据 OpenAPI 3.x 的 schema 构建嵌套的参数配置
// 参数 schema 为 OpenAPI 3.x 的 schema 指针
// 返回 ItemsConfig 结构体
func buildNestedArg(schema *openapi3.Schema) config.ItemsConfig {
	var items config.ItemsConfig

	// 如果 schema 有类型定义
	if schema.Type != nil {
		types := schema.Type.Slice()
		if len(types) > 0 {
			schemaType := types[0]
			items.Type = schemaType

			// 如果是对象类型且有属性定义
			if schemaType == "object" && schema.Properties != nil {
				properties := make(map[string]any)
				for childName, childProp := range schema.Properties {
					propType := "string"
					if childProp.Value.Type != nil && len(childProp.Value.Type.Slice()) > 0 {
						propType = childProp.Value.Type.Slice()[0]
					}
					propSchema := map[string]any{
						"type":        propType,
						"description": childProp.Value.Description,
					}
					// 递归处理嵌套的 object/array 类型
					if propType == "object" {
						nested := buildNestedArg(childProp.Value)
						if len(nested.Properties) > 0 {
							propSchema["properties"] = nested.Properties
						}
					} else if propType == "array" && childProp.Value.Items != nil && childProp.Value.Items.Value != nil {
						nested := buildNestedArg(childProp.Value.Items.Value)
						propSchema["items"] = nested
					}
					properties[childName] = propSchema
				}
				items.Properties = properties
				items.Required = schema.Required
			} else if schemaType == "array" && schema.Items != nil && schema.Items.Value != nil {
				nested := buildNestedArg(schema.Items.Value)
				items.Type = "array"
				items.Properties = nil
				items.Items = &nested
			}
		}
	}
	return items
}
