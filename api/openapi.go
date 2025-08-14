package api

import "time"

// ServerInfoRequest 创建服务器信息请求
type ServerInfoRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	FileContent string `json:"file_content" binding:"required"`
}

// ServerInfoResponse 服务器信息响应
type ServerInfoResponse struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	FileContent string    `json:"file_content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
