package model

type McpServer struct {
	BaseModel
	UUID          string          `json:"uuid" gorm:"column:uuid;type:varchar(36);not null;default:'';comment:服务器唯一标识"`
	Name          string          `json:"name" gorm:"column:name;type:varchar(500);not null;default:'';comment:服务器名称"`
	Description   string          `json:"description" gorm:"column:description;type:text;default:'';comment:服务器详细描述"`
	Urls          string          `json:"urls" gorm:"column:urls;type:varchar(255);not null;default:'';comment:服务器访问地址"`
	AllowedTools  string          `json:"allowedTools" gorm:"column:allowed_tools;type:text;default:'';comment:允许使用的工具列表，JSON格式存储"`
	Version       string          `json:"version" gorm:"column:version;type:varchar(20);default:v1.0.0;comment:OpenAPI版本号"`
	HaveTools     HaveToolsStatus `json:"haveTools" gorm:"column:have_tools;type:SMALLINT;default:1;comment:服务器是否支持工具,0未知，1 不支持，2支持"`
	IsAuth        AuthStatus      `json:"isAuth" gorm:"column:is_auth;type:SMALLINT;default:1;comment:服务器是否需要认证,如果开启全局认证，则所有的方法都需要认证,0未知，1 不开启，2开启service授权，3开启平台授权，4开启所有的授权"`
	ServiceToken  string          `json:"serviceToken" gorm:"column:service_token;type:text;comment:服务认证Token，用于访问用户提供的接口"`
	PlatformToken string          `json:"platformToken" gorm:"column:platform_token;type:text;comment:平台认证Token，平台添加的认证令牌"`

	Security string `json:"security" gorm:"column:security;type:text;default:'';comment:认证信息"` //认证的原始信息

}

func (m *McpServer) TableName() string {
	return "mcp_server"
}
