package model

import "flow-bridge-mcp/pkg/const"

type McpTools struct {
	BaseModel
	McpServerId   int64                      `json:"mcpServerId" gorm:"column:mcp_server_id;uniqueIndex:idx_mcp_tools_unique;type:bigint;not null;comment:关联的MCP服务器ID"`
	McpServerUUID string                     `json:"mcpServerUUID" gorm:"column:mcp_server_uuid;uniqueIndex:idx_mcp_tools_unique;type:varchar(36);not null;default:'';comment:工具唯一标识"`
	Name          string                     `json:"name" gorm:"column:name;uniqueIndex:idx_mcp_tools_unique;type:varchar(500);not null;default:'';comment:工具名称"`
	Description   string                     `json:"description" gorm:"column:description;type:text;default:'';comment:工具描述信息"`
	McpServerType _const.McpServerTypeStatus `json:"mcpServerType" gorm:"column:mcp_server_type;type:SMALLINT;default:1;comment:服务器类型 1:openapi 2:grpc"`
	Method        string                     `json:"method" gorm:"column:method;type:varchar(10);default:'';comment:HTTP请求方法(GET/POST/PUT/DELETE等)"`
	Endpoint      string                     `json:"endpoint" gorm:"column:endpoint;type:varchar(255);default:'';comment:API端点地址"`
	Headers       string                     `json:"headers" gorm:"column:headers;type:text;default:'';comment:请求头信息，JSON格式存储"`
	Args          string                     `json:"args" gorm:"column:args;type:text;default:'';comment:参数配置，JSON格式存储"`
	RequestBody   string                     `json:"requestBody" gorm:"column:request_body;type:text;default:'';comment:请求体模板或内容"`
	ResponseBody  string                     `json:"responseBody" gorm:"column:response_body;type:text;default:'';comment:响应体模板或内容"`

	ToolSchema  string `json:"toolSchema" gorm:"column:tool_schema;type:text;default:'';comment:工具参数Schema定义"`
	Annotations string `json:"annotations" gorm:"column:annotations;type:text;default:'';comment:工具注解信息，JSON格式存储"`
	Security    string `json:"security" gorm:"column:security;type:text;default:'';comment:认证信息"`

	IsAuth         _const.AuthStatus `json:"isAuth" gorm:"column:is_auth;type:SMALLINT;default:1;comment:是否需要认证"` // 是否需要认证 0未知，默认1不需要认证,2需要认证，这个认证为接口级别
	AuthMode       string            `json:"authMode" gorm:"column:auth_mode;type:varchar(20);default:'';comment:认证模式"`
	IsPlatformAuth _const.AuthStatus `json:"isPlatformAuth" gorm:"column:is_platform_auth;type:SMALLINT;default:1;comment:是否平台认证"`
	IsShow         _const.Status     `json:"isShow" gorm:"column:is_show;type:SMALLINT;default:1;comment:是否显示"`
	SerialNumber   string            `json:"serialNumber" gorm:"column:serial_number;type:varchar(36);default:'';comment:服务序列号"` //来源于server 的serialNumber
	//是否重复
	IsRepeat _const.Status `json:"isRepeat" gorm:"column:is_repeat;type:SMALLINT;default:1;comment:是否重复"`
}

func (m *McpTools) TableName() string {
	return "mcp_tools"
}
