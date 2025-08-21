package model

type McpServer struct {
	BaseModel
	UUID         string          `json:"uuid" gorm:"column:uuid;type:varchar(36);not null;default:'';comment:服务器唯一标识"`
	Name         string          `json:"name" gorm:"column:name;type:varchar(500);not null;default:'';comment:服务器名称"`
	Description  string          `json:"description" gorm:"column:description;type:text;default:'';comment:服务器详细描述"`
	Urls         string          `json:"urls" gorm:"column:urls;type:varchar(255);not null;default:'';comment:服务器访问地址"`
	AllowedTools string          `json:"allowedTools" gorm:"column:allowed_tools;type:text;default:'';comment:允许使用的工具列表，JSON格式存储"`
	Version      string          `json:"version" gorm:"column:version;type:varchar(20);default:v1.0.0;comment:OpenAPI版本号"`
	HaveTools    HaveToolsStatus `json:"haveTools" gorm:"column:have_tools;type:SMALLINT;default:0;comment:服务器是否支持工具,,0 不支持，1支持"`
	IsAuth       AuthStatus      `json:"isAuth" gorm:"column:is_auth;type:SMALLINT;default:0;comment:服务器是否需要认证,如果开启全局认证，则所有的方法都需要认证,0 不开启，1开启"`
	Token        string          `json:"token" gorm:"column:token;type:varchar(255);comment:全局认证Token"`
	Auth         string          `json:"auth" gorm:"column:auth;type:text;default:'';comment:认证信息"` //认证的原始信息
}

func (m *McpServer) TableName() string {
	return "mcp_server"
}
