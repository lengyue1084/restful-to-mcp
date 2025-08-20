package transformer

import (
	"context"
	"flow-bridge-mcp/internal/mcp/config"
	"time"
)

// Metadata 转换器元数据（可选）router映射、等一些附加信息都可以存储，预留
type Metadata struct {
	RequestID      string         `json:"request_id"`
	ClientIP       string         `json:"client_ip"`
	UserAgent      string         `json:"user_agent"`
	ProcessingTime time.Duration  `json:"processing_time"`
	Data           map[string]any `json:"custom_tags"`
}

func NewMetadata() *Metadata {
	return &Metadata{}
}

// Transformer 定义了转换器的通用接口,后续扩展其他协议的时候需要实现该方法
// 所有具体的转换器实现都需要实现这个接口
type Transformer interface {
	VersionDetector
	Validator
	// Convert 将输入的规范数据转换为 MCP 配置
	// 参数 ctx 为上下文，用于控制超时和取消
	// 参数 data 为规范的字节数据
	// 返回 MCP 配置指针和可能出现的错误
	Convert(ctx context.Context, data []byte) (*config.MCPServer, error)
	// Metadata 可选
	// Metadata 获取转换过程的元数据信息
	Metadata(ctx context.Context) *Metadata // 修正方法名和返回类型

}

// VersionDetector 子接口定义
type VersionDetector interface {
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
