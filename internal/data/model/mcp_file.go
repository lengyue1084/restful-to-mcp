package model

type McpFile struct {
	BaseModel
	Md5  string `json:"md5" gorm:"column:md5;type:varchar(36);not null;default:'';"`
	Name string `json:"name" gorm:"column:name;type:varchar(500);not null;default:'';"`
}

func (m *McpFile) TableName() string {
	return "mcp_file"
}
