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
func (c *Converter) Convert(ctx context.Context, specData []byte) (*config.MCPServer, error) {

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
	if doc == nil {
		return nil, fmt.Errorf("failed to load OpenAPI specification,err:%s", "doc is nil")
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
	var (
		info            *openapi3.Info
		serverName      string
		servers         []*openapi3.Server
		paths           *openapi3.Paths
		securitySchemes openapi3.SecuritySchemes
		urls            []string
		components      *openapi3.Components
		toolsSlice      []*config.ToolConfig
		securityKey     string
	)
	info = doc.Info
	if info == nil {
		c.log.ErrorWithContext(ctx, "info is nil")
		return nil, fmt.Errorf("info is nil")
	}
	serverName = doc.Info.Title
	if serverName == "" {
		c.log.ErrorWithContext(ctx, "openapi info title must no empty")
		return nil, fmt.Errorf("openapi info title must no empty")
	}
	servers = doc.Servers
	if servers == nil || len(servers) == 0 {
		c.log.ErrorWithContext(ctx, "openapi servers is empty")
		return nil, fmt.Errorf("openapi servers is empty")
	}
	if len(servers) <= 0 {
		c.log.ErrorWithContext(ctx, "openapi servers is empty")
		return nil, fmt.Errorf("openapi servers url is needed")
	}
	for _, server := range servers {
		url := server.URL
		if url != "" && (strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https")) {
			urls = append(urls, server.URL)
		}
	}
	if urls == nil || len(urls) == 0 {
		c.log.ErrorWithContext(ctx, "openapi servers url is invalid")
		return nil, fmt.Errorf("openapi servers url is invalid")
	}
	paths = doc.Paths
	if paths == nil {
		c.log.ErrorWithContext(ctx, "paths is nil")
		return nil, fmt.Errorf("paths is nil")
	}
	components = doc.Components
	securitySchemes = components.SecuritySchemes
	var authSecuritySchemes []*config.Security
	for key, scheme := range securitySchemes {
		val := scheme.Value
		if val == nil {
			continue
		}
		if val.Type == config.AuthModeApiKey.String() && val.Name == "" {
			c.log.ErrorWithContext(ctx, "security scheme name is empty", "key", key)
			return nil, fmt.Errorf("security scheme name is empty")
		}
		authSecuritySchemes = append(authSecuritySchemes, &config.Security{
			SecurityKey:  key,
			Name:         val.Name,
			Mode:         config.AuthMode(val.Type),
			Description:  val.Description,
			In:           config.AuthPosition(val.In),
			Scheme:       val.Scheme,
			BearerFormat: val.BearerFormat,
		})

	}

	if doc.Security != nil && len(doc.Security) > 0 {
		//for _, sec := range doc.Security {
		// 只要有一种授权方式，不支持，则返回错误
		ok, err := c.isAllNotAllowed(doc.Security, authSecuritySchemes)
		if err != nil {
			return nil, fmt.Errorf("doc security isAllNotAllowed failed: %w", err)
		}
		if ok {
			return nil, fmt.Errorf("doc security only support:%s and%s  methods，", config.AuthModeHttp, config.AuthModeApiKey)
		}
		//对于文档级别的鉴权方式，如果有多个鉴权方式，只取第一个
		securityKey = c.getFirstSecurityName(doc.Security)
	}

	toolsSlice, err = c.PathsToTools(paths, components, securityKey, authSecuritySchemes)
	if err != nil {
		c.log.ErrorWithContext(ctx, "failed to convert OpenAPI paths to MCP tools: %w", err)
		return nil, err
	}
	var allTools []string
	for _, tool := range toolsSlice {
		allTools = append(allTools, tool.Name)
	}

	//serverTitle := strings.ReplaceAll(serverName, " ", "")
	// 创建服务器配置
	server := &config.MCPServer{
		Name:         info.Title,
		UUID:         mcpUUID,
		Description:  info.Description,
		Urls:         urls,
		CreatedAt:    time.Time{}, //目前没有创建则时间为零值
		UpdatedAt:    time.Time{}, //目前没有创建则时间为零值
		Config:       nil,
		Tools:        toolsSlice,
		Version:      info.Version,
		SecurityList: authSecuritySchemes,
		AllTools:     allTools,
	}

	return server, nil
}

// 只要有个认证方式不支持就返回true
func (c *Converter) getFirstSecurityName(security openapi3.SecurityRequirements) (securityName string) {

	for _, sec := range security {
		for name := range sec {
			securityName = name
			break
		}
	}
	return securityName
}
func (c *Converter) isAllNotAllowed(security openapi3.SecurityRequirements, securitySlice []*config.Security) (isAllNotAllowed bool, err error) {
	isAllNotAllowed = false
	//if len(security) > 1 {
	//	// or
	//	return true, errors.New("api only one security is allowed")
	//}
	var securityMap = make(map[string]*config.Security)
	for _, val := range securitySlice {
		securityMap[val.SecurityKey] = val
	}
	for _, sec := range security {
		// and
		//if len(sec) > 1 {
		//	return true, errors.New("api only one security is allowed")
		//}
		for name := range sec {
			if v, ok := securityMap[name]; ok {
				securityMap[name] = v
				if v.Mode != config.AuthModeHttp && v.Mode != config.AuthModeApiKey {
					isAllNotAllowed = true
					return isAllNotAllowed, fmt.Errorf("api only support:%s and%s  methods，", config.AuthModeHttp, config.AuthModeApiKey)
				}
			}

		}
	}
	return isAllNotAllowed, nil
}

func (c *Converter) PathsToTools(paths *openapi3.Paths, components *openapi3.Components, docSecuritykey string, authSecuritySchemes []*config.Security) (toolsSlice []*config.ToolConfig, err error) {
	var securityLevel = config.SecurityLevelPublic
	if docSecuritykey != "" {
		securityLevel = config.SecurityLevelDoc
	}
	var (
		pathSecurityLevel = securityLevel
		pathSecurityKey   = docSecuritykey
		//接口是否显示，对于不符合规则的接口，默认不显示，不会注册到mcp server中
		isShow = true
	)

	// 将文档中的路径转换为工具配置
	for path, pathItem := range paths.Map() {

		isShow = true
		// 为每个 HTTP 方法创建一个工具配置
		for method, operation := range pathItem.Operations() {
			if method == "options" {
				continue // 跳过 CORS 预检请求
			}

			//判断OperationID 是否为空，则根据方法和路径生成一个操作 OperationID
			if operation.OperationID == "" {
				// 如果OperationID 格式，例如：/user/order/{userid} -> users_email_arguserid
				pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
				for i, part := range pathParts {
					if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
						pathParts[i] = "arg" + strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}")
					}
				}
				// OperationID = get_users_email_arguserid
				operation.OperationID = fmt.Sprintf("%s_%s", strings.ToLower(method), strings.Join(pathParts, "_"))
			}
			var securityInfo *config.Security = nil
			var pathSecurityMode = config.AuthModeEmpty
			pathSecurityKey = docSecuritykey
			pathSecurityLevel = securityLevel
			security := operation.Security
			if security != nil {
				pathSecurityLevel = config.SecurityLevelApi
				ok, err := c.isAllNotAllowed(*security, authSecuritySchemes)
				if err != nil {
					return nil, err
				}
				if ok {
					return nil, fmt.Errorf("path:%s,security not allowed", path)
				}
				pathSecurityKey = c.getFirstSecurityName(*security)
			}

			for _, val := range authSecuritySchemes {
				if pathSecurityKey == val.SecurityKey {
					pathSecurityMode = val.Mode
					securityInfo = val
					if pathSecurityLevel != config.SecurityLevelPublic && (pathSecurityMode != config.AuthModeHttp && val.Mode != config.AuthModeApiKey) {
						isShow = false
					}
					break
				}
			}

			// 创建工具配置
			toolInfo := &config.ToolConfig{
				Name:          operation.OperationID,
				Description:   tool.FirstNonEmpty(operation.Description, operation.Summary),
				Method:        method,
				Endpoint:      fmt.Sprintf("{{.Config.url}}%s", path),
				Headers:       make(map[string]string),
				Args:          make([]config.ArgConfig, 0),
				ResponseBody:  "{{.Response.Body}}", // 使用透传响应
				SecurityMode:  pathSecurityMode,
				SecurityLevel: pathSecurityLevel,
				IsShow:        isShow,
				Security:      securityInfo,
			}
			//if toolInfo.Name == "addPet" {
			//	fmt.Printf("operation.OperationID is empty")
			//}

			// 添加默认请求头
			toolInfo.Headers["Content-Type"] = "application/json"

			// 定义不同位置的参数切片
			var (
				bodyArgs   []config.ArgConfig
				pathArgs   []config.ArgConfig
				queryArgs  []config.ArgConfig
				headerArgs []config.ArgConfig
			)

			// 处理操作中的参数
			for _, param := range operation.Parameters {
				// 创建参数配置
				arg := config.ArgConfig{
					Name:        param.Value.Name,
					Position:    param.Value.In,
					Required:    param.Value.Required,
					Type:        "string", // 默认参数类型为字符串
					Description: param.Value.Description,
					Explode:     false,
				}

				// 如果参数有 schema 定义，尝试获取参数类型和默认值
				if param.Value.Schema != nil && param.Value.Schema.Value != nil {
					schema := param.Value.Schema.Value

					// 处理参数类型
					if schema.Type != nil {
						types := schema.Type.Slice()
						if len(types) > 0 {
							arg.Type = types[0]
							//解析数组类型的 query参数
							// 对于数组类型，尝试获取元素类型
							if arg.Type == "array" && schema.Items != nil && schema.Items.Value != nil {
								arg.Explode = true
								if schema.Items.Value.Type != nil {
									itemTypes := schema.Items.Value.Type.Slice()
									if len(itemTypes) > 0 {
										// 可以记录数组元素类型，或者构建更复杂的类型描述
										arg.Items = buildNestedArg(schema.Items.Value)
									}
								}
							}
						}
					}
					if param.Value.Schema.Value.Default != nil {
						arg.Default = fmt.Sprintf("%v", param.Value.Schema.Value.Default)
					}
					// 处理其他有用的 schema 属性
					if schema.Description != "" && arg.Description == "" {
						arg.Description = schema.Description
					}

					// 处理枚举值
					if len(schema.Enum) > 0 {
						// 可以将枚举值存储在 arg 的扩展字段中
						enumValues := make([]string, len(schema.Enum))
						for i, enumVal := range schema.Enum {
							enumValues[i] = fmt.Sprintf("%v", enumVal)
						}
						// 如果有扩展字段可以存储枚举信息
						arg.Enum = enumValues
						arg.Description = fmt.Sprintf("%s,其中参数只能从 [%s]选取", arg.Description, strings.Join(enumValues, ", "))
					}
					if arg.Description == "" {
						arg.Description = "参数描述"
					}
					if schema.Min != nil {
						arg.Description = fmt.Sprintf("%s,参数最小值是%v", arg.Description, *schema.Min)
					}
					if schema.Max != nil {
						arg.Description = fmt.Sprintf("%s,参数最大值是%v", arg.Description, *schema.Max)
					}
					if schema.MinLength != 0 {
						arg.Description = fmt.Sprintf("%s,参数最小长度是%v", arg.Description, schema.MinLength)
					}
					if schema.MaxLength != nil {
						arg.Description = fmt.Sprintf("%s,参数最大长度是%v", arg.Description, schema.MaxLength)
					}
					if schema.Example != nil {
						arg.Description = fmt.Sprintf("%s,参数示例是%v", arg.Description, schema.Example)
					}
				}

				// 根据参数位置进行不同处理
				switch param.Value.In {
				case "path":
					// 路径参数总是必需的
					arg.Required = true
					pathArgs = append(pathArgs, arg)
					// 更新端点中的路径参数占位符
					toolInfo.Endpoint = strings.ReplaceAll(toolInfo.Endpoint, fmt.Sprintf("{%s}", arg.Name), fmt.Sprintf("{{.Args.%s}}", arg.Name))
				case "query":
					queryArgs = append(queryArgs, arg)
				case "header":
					toolInfo.Headers[arg.Name] = fmt.Sprintf("{{.Args.%s}}", arg.Name)
					headerArgs = append(headerArgs, arg)
				}
			}

			// 处理请求体
			if operation.RequestBody != nil {
				// 获取请求体是否必需的标志,否则会400错误，如果为true可以给一个空的json {}
				//requestBodyRequired := operation.RequestBody.Value.Required
				// 遍历请求体支持的内容类型
				for contentType, contentValue := range operation.RequestBody.Value.Content {
					if contentType == "application/json" { //只处理application/json，过滤其他类型包括二进制文件的类型
						toolInfo.RequestBody = contentType
						toolInfo.ContentType = contentType
						toolInfo.IsShow = true
						// 添加请求体参数
						if contentValue.Schema != nil {
							schema := contentValue.Schema.Value
							// 处理 schema 引用
							if contentValue.Schema.Ref != "" {
								refName := strings.TrimPrefix(contentValue.Schema.Ref, "#/components/schemas/")
								if refSchema, ok := components.Schemas[refName]; ok {
									schema = refSchema.Value
								}
								// 处理数组类型
							} else if schema.Type != nil && len(schema.Type.Slice()) > 0 && schema.Type.Slice()[0] == "array" {
								if schema.Items != nil && schema.Items.Value != nil {
									refName := strings.TrimPrefix(schema.Items.Ref, "#/components/schemas/")
									if refSchema, ok := components.Schemas[refName]; ok {
										schema = refSchema.Value
									}
								}
							}

							if schema.Properties != nil { // 如果 schema 有属性定义
								for name, prop := range schema.Properties {
									// 跳过响应专用字段
									if strings.HasPrefix(name, "response") || name == "createdAt" {
										continue
									}

									// 创建请求体参数配置
									arg := config.ArgConfig{
										Name:     name, //字段名
										Position: "body",
										//Required:    requestBodyRequired || contains(schema.Required, name), // 判断属性是否必需，如果全局没有设置，则取，则判断ref中required是否定义了字段
										Required:    contains(schema.Required, name), // 判断属性是否必需，如果全局没有设置，则取，则判断ref中required是否定义了字段
										Type:        "string",                        // 设置默认参数类型为字符串
										Description: prop.Value.Description,
									}
									// 处理 schema 引用
									if prop.Ref != "" {
										refName := strings.TrimPrefix(prop.Ref, "#/components/schemas/")
										if refSchema, ok := components.Schemas[refName]; ok {
											// 使用引用的 schema 信息
											if refSchema.Value != nil {
												prop = refSchema // 替换为引用的 schema
											}
										}
									}

									// 如果属性有类型定义
									if prop.Value != nil && prop.Value.Type != nil {
										types := prop.Value.Type.Slice()
										if len(types) > 0 {
											arg.Type = types[0] //重新赋值参数类型
											// 如果是数组类型且有 items 定义
											if arg.Type == "array" && prop.Value.Items != nil && prop.Value.Items.Value != nil {
												arg.Items = buildNestedArg(prop.Value.Items.Value)
											} else if arg.Type == "object" && prop.Value.Properties != nil {
												// 处理对象类型
												arg.Items = buildNestedArg(prop.Value)
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
						break
					} else { //其他类型的不显示，比如multipart/form-data、application/x-www-form-urlencoded、text/plain、application/octet-stream
						toolInfo.IsShow = false
						toolInfo.ContentType = contentType //如果非application/json，则取最后一次的类型
					}

				}
			}

			// 合并所有参数
			toolInfo.Args = append(toolInfo.Args, pathArgs...)
			toolInfo.Args = append(toolInfo.Args, queryArgs...)
			toolInfo.Args = append(toolInfo.Args, bodyArgs...)
			toolInfo.Args = append(toolInfo.Args, headerArgs...)

			// 如果有请求体参数，构建请求体模板,格式为json，如果参数为quantity，shipDate，status，complete，petId，则模板为：
			//{
			//	"quantity": {{ toJSON .Args.quantity}},
			//	"shipDate": {{ toJSON .Args.shipDate}},
			//	"status": {{ toJSON .Args.status}},
			//	"complete": {{ toJSON .Args.complete}},
			//	"petId": {{ toJSON .Args.petId}}
			//}
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
				toolInfo.RequestBody = bodyTemplate.String()
			}
			toolInfo.ToolSchema = toolInfo.ArgsToInputSchema()
			toolInfo.Annotations = map[string]any{
				"title":       toolInfo.Name,
				"description": toolInfo.Description,
			}
			toolsSlice = append(toolsSlice, toolInfo)
		}
	}
	return toolsSlice, nil
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
func (c *Converter) convertSwagger2(ctx context.Context, specData []byte) (*config.MCPServer, error) {
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
func (c *Converter) ConvertFromJSON(ctx context.Context, jsonData []byte) (*config.MCPServer, error) {
	return c.Convert(ctx, jsonData)
}

// ConvertFromYAML 将 YAML 格式的 OpenAPI 规范转换为 MCP 配置
// 参数 yamlData 为 YAML 格式的 OpenAPI 规范字节数据
// 返回 MCP 配置指针和可能出现的错误
func (c *Converter) ConvertFromYAML(ctx context.Context, yamlData []byte) (*config.MCPServer, error) {
	return c.Convert(ctx, yamlData)
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
