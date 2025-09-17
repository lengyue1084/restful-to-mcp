package api

type UpdateMcpServerByUUIDRequest struct {
	UUID        string `json:"uuid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type UpdateMcpServerByUUIDResponse struct {
}

type GetMcpConnectTokenByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type GetMcpConnectTokenByUUIDResponse struct {
	ConnectToken string `json:"connectToken" binding:"required"`
}

type DeleteMcpServerByUUIDRequest struct {
	UUID string `json:"uuid" binding:"required"`
}
type DeleteMcpServerByUUIDResponse struct {
}

type CommonMcpServerByForm struct {
	UUID          string `json:"uuid" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Url           string `json:"url" binding:"required"`
	Version       string `json:"version" binding:"required"`
	IsAuth        int8   `json:"isAuth" binding:"required"`
	PlatformToken string `json:"platformToken"`
}

type CreateMcpServerByFormRequest struct {
	CommonMcpServerByForm
}

type CreateMcpServerByFormResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	CommonMcpServerByForm
}

type UpdateMcpServerByFormRequest struct {
	CommonMcpServerByForm
}

type UpdateMcpServerByFormResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CommonMcpServerByForm
}
