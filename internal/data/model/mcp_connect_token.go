package model

type McpConnectToken struct {
	BaseModel
	McpServerUUID string `json:"mcpServerUUID" gorm:"column:mcp_server_uuid;uniqueIndex:idx_mcp_tools_unique;type:varchar(36);not null;default:'';comment:工具唯一标识"`
	McpServerId   int64  `json:"mcp_server_id" gorm:"column:mcp_server_id;type:bigint;not null;comment:关联的MCP服务器ID"`
	McpServerName string `json:"mcp_server_name" gorm:"column:mcp_server_name;type:varchar(500);not null;default:'';comment:工具名称"`
	ConnectToken  string `json:"connect_token" gorm:"column:connect_token;type:varchar(36);not null;default:'';comment:连接token"`
}

func (m *McpConnectToken) TableName() string {
	return "mcp_connect_token"
}
