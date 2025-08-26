package model

type McpFile struct {
	BaseModel
	Md5         string `json:"md5" gorm:"column:md5;type:varchar(36);not null;default:'';"`
	Name        string `json:"name" gorm:"column:name;type:varchar(500);not null;default:'';"`
	SourceName  string `json:"sourceName" gorm:"column:source_name;type:varchar(500);not null;default:'';"`
	Description string `json:"description" gorm:"column:description;type:text;default:'';"`
	Suffix      string `json:"suffix" gorm:"column:suffix;type:varchar(255);not null;default:'yaml';"`
}

func (m *McpFile) TableName() string {
	return "mcp_file"
}
