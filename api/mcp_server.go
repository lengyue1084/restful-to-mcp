package api

type UpdateMcpServerByUUIDRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type UpdateMcpServerByUUIDResponse struct {
}
