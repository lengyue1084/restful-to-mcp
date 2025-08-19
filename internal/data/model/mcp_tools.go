package model

type McpTools struct {
	BaseModel
	McpServerId   uint   `json:"McpServerId" gorm:"column:mcp_server_id;type:bigint;not null;comment:关联的MCP服务器ID"`
	Name          string `json:"name" gorm:"column:name;type:varchar(500);not null;comment:工具名称"`
	Description   string `json:"description" gorm:"column:description;type:text;comment:工具描述信息"`
	McpServerType uint8  `json:"McpServerType" gorm:"column:mcp_server_type;type:int8;default:1;comment:服务器类型 1:openapi 2:grpc"`
	Method        string `json:"method" gorm:"column:method;type:varchar(10);comment:HTTP请求方法(GET/POST/PUT/DELETE等)"`
	Endpoint      string `json:"endpoint" gorm:"column:endpoint;type:varchar(255);comment:API端点地址"`
	Headers       string `json:"headers" gorm:"column:headers;type:text;comment:请求头信息，JSON格式存储"`
	Args          string `json:"args" gorm:"column:args;type:text;comment:参数配置，JSON格式存储"`
	InputSchema   string `json:"inputSchema" gorm:"column:input_schema;type:text;comment:输入参数Schema定义"`
	RequestBody   string `json:"requestBody" gorm:"column:request_body;type:text;comment:请求体模板或内容"`
	ResponseBody  string `json:"responseBody" gorm:"column:response_body;type:text;comment:响应体模板或内容"`
	Annotations   string `json:"annotations" gorm:"column:annotations;type:text;comment:工具注解信息，JSON格式存储"`
}

func (m *McpTools) TableName() string {
	return "mcp_tools"
}
