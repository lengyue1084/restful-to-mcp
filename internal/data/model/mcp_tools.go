package model

type McpTools struct {
	BaseModel
	McpServerId   uint                `json:"McpServerId" gorm:"column:mcp_server_id;type:bigint;not null;comment:关联的MCP服务器ID"`
	Name          string              `json:"name" gorm:"column:name;type:varchar(500);not null;default:'';comment:工具名称"`
	Description   string              `json:"description" gorm:"column:description;type:text;default:'';comment:工具描述信息"`
	McpServerType McpServerTypeStatus `json:"McpServerType" gorm:"column:mcp_server_type;type:SMALLINT;default:1;comment:服务器类型 1:openapi 2:grpc"`
	Method        string              `json:"method" gorm:"column:method;type:varchar(10);default:'';comment:HTTP请求方法(GET/POST/PUT/DELETE等)"`
	Endpoint      string              `json:"endpoint" gorm:"column:endpoint;type:varchar(255);default:'';comment:API端点地址"`
	Headers       string              `json:"headers" gorm:"column:headers;type:text;default:'';comment:请求头信息，JSON格式存储"`
	Args          string              `json:"args" gorm:"column:args;type:text;default:'';comment:参数配置，JSON格式存储"`
	RequestBody   string              `json:"requestBody" gorm:"column:request_body;type:text;default:'';comment:请求体模板或内容"`
	ResponseBody  string              `json:"responseBody" gorm:"column:response_body;type:text;default:'';comment:响应体模板或内容"`

	InputSchema string `json:"inputSchema" gorm:"column:input_schema;type:text;default:'';comment:输入参数Schema定义"`
	Annotations string `json:"annotations" gorm:"column:annotations;type:text;default:'';comment:工具注解信息，JSON格式存储"`

	IsAuth   AuthStatus `json:"isAuth" gorm:"column:is_auth;type:SMALLINT;default:0;comment:是否需要认证"` // 是否需要认证 默认0不需要认证,1需要认证，这个认证为接口级别
	AuthType string     `json:"authType" gorm:"column:auth_type;type:varchar(30);default:'';comment:认证方式"`
}

func (m *McpTools) TableName() string {
	return "mcp_tools"
}
