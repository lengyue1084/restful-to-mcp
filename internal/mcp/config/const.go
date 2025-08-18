package config

// JSPNRPCVersion Protocol versions
const (
	JSPNRPCVersion = "2.0"
)

// Methods
const (
	Initialize              = "initialize"
	NotificationInitialized = "notifications/initialized"
	NotificationMessage     = "notifications/message"
	Ping                    = "ping"
	ToolsList               = "tools/list"
	ToolsCall               = "tools/call"
	// PromptsList 如果不涉及 AI 提示词功能，可以不实现
	PromptsList = "prompts/list"
	PromptsGet  = "prompts/get"
)

// Response
const (
	Accepted = "Accepted"

	NotificationRootsListChanged    = "notifications/roots/list_changed"
	NotificationCancelled           = "notifications/cancelled"
	NotificationProgress            = "notifications/progress"
	NotificationResourceUpdated     = "notifications/resources/updated"
	NotificationResourceListChanged = "notifications/resources/list_changed"
	NotificationToolListChanged     = "notifications/tools/list_changed"
	NotificationPromptListChanged   = "notifications/prompts/list_changed"

	SamplingCreateMessage = "sampling/createMessage"
	LoggingSetLevel       = "logging/setLevel"

	// ResourcesList 如果不涉及文件或资源配置，可以不实现
	ResourcesList          = "resources/list"
	ResourcesTemplatesList = "resources/templates/list"
	ResourcesRead          = "resources/read"
)

// Error codes for MCP protocol
// Standard JSON-RPC error codes
const (
	ErrorCodeParseError     = -32700
	ErrorCodeInvalidRequest = -32600
	ErrorCodeMethodNotFound = -32601
	ErrorCodeInvalidParams  = -32602
	ErrorCodeInternalError  = -32603
)

// SDKs and applications error codes
const (
	ErrorCodeConnectionClosed = -32000
	ErrorCodeRequestTimeout   = -32001
)

const (
	HeaderMcpSessionID = "SessionId"
)

const (
	TextContentType  = "text"
	ImageContentType = "image"
	AudioContentType = "audio"
)
