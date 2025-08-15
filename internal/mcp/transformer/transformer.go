package transformer

import (
	"context"
	"flow-bridge-mcp/internal/mcp/config"
)

// Transformer 定义了转换器的通用接口,后续扩展其他协议的时候需要实现该方法
// 所有具体的转换器实现都需要实现这个接口
type Transformer interface {
	VersionedTransformer
	Validator
	// Convert 将输入的规范数据转换为 MCP 配置
	// 参数 ctx 为上下文，用于控制超时和取消
	// 参数 data 为规范的字节数据
	// 返回 MCP 配置指针和可能出现的错误
	Convert(ctx context.Context, data []byte) (*config.MCPConfig, error)
	// ConvertWithOptions ConvertFromYAML ConvertToJSON 将 MCP 配置转换为 JSON 格式
	// 输入参数 mcpConfig 为 MCP 配置指针
	// 返回 JSON 字节数据和可能出现的错误
	//ConvertToJSON(ctx context.Context, mcpConfig *config.MCPConfig) ([]byte, error)
	// ConvertFromYAML 将 YAML 格式的规范转换为 MCP 配置
	// 参数 yamlData 为 YAML 格式的规范字节数据
	// 返回 MCP 配置指针和可能出现的错误
	//ConvertFromYAML(ctx context.Context, yamlData []byte) (*config.MCPConfig, error)
	// ConvertWithOptions 将规范转换为 MCP 配置，可指定租户和前缀
	// 参数 data 为规范的字节数据
	// 参数 tenant 为租户名称
	// 参数 prefix 为前缀
	// 返回 MCP 配置指针和可能出现的错误
	//ConvertWithOptions(ctx context.Context, data []byte, tenant, prefix string) (*config.MCPConfig, error)

}

// VersionedTransformer 定义了支持版本检测的转换器接口
// 继承自 Transformer 接口，增加了版本检测功能
type VersionedTransformer interface {
	// DetectVersion 从规范数据中检测版本
	// 参数 data 为规范的字节数据
	// 返回检测到的版本字符串和可能出现的错误
	DetectVersion(ctx context.Context, data []byte) (string, error)
}

// Validator 定义了验证功能的接口
type Validator interface {
	// Validate 验证规范数据的有效性
	// 参数 data 为规范的字节数据
	// 返回验证结果和可能出现的错误
	Validate(ctx context.Context, data []byte) error
}
