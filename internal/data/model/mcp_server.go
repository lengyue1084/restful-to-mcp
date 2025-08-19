package model

type McpServer struct {
	BaseModel
	UUID        string `json:"uuid" gorm:"uuid"`
	Name        string `json:"name" gorm:"name"`
	Description string `json:"description" gorm:"description"`
	Url         string `json:"url" gorm:"url"`
	Version     string `json:"version" gorm:"version"`
}

func (m *McpServer) TableName() string {
	return "mcp_server"
}
