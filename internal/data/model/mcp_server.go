package model

type McpServer struct {
	BaseModel
	Names       string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Version     string `json:"version"`
}

func (m *McpServer) TableName() string {
	return "mcp_server"
}
