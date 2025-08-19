package model

type McpFile struct {
	BaseModel
	Md5  string `json:"md5" gorm:"md5"`
	Name string `json:"name" gorm:"name"`
}

func (m *McpFile) TableName() string {
	return "mcp_file"
}
