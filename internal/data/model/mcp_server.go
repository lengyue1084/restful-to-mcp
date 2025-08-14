package model

type McpServer struct {
	ID        int    `json:"id"`
	User      string `json:"user"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (m *McpServer) TableName() string {
	return "mcp_server"
}
