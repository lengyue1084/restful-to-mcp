package model

type McpServer struct {
	BaseModel
	UUID         string `json:"uuid" gorm:"column:uuid;type:varchar(36);not null;comment:服务器唯一标识"`
	Name         string `json:"name" gorm:"column:name;type:varchar(500);not null;comment:服务器名称"`
	Description  string `json:"description" gorm:"column:description;type:text;comment:服务器详细描述"`
	Url          string `json:"url" gorm:"column:url;type:varchar(255);not null;comment:服务器访问地址"`
	Port         int    `json:"port" gorm:"column:port;type:int;default:8080;comment:服务端口号"`
	AllowedTools string `json:"allowedTools" gorm:"column:allowed_tools;type:text;comment:允许使用的工具列表，JSON格式存储"`
	Version      string `json:"version" gorm:"column:version;type:varchar(20);default:v1.0.0;comment:OpenAPI版本号"`
}

func (m *McpServer) TableName() string {
	return "mcp_server"
}
