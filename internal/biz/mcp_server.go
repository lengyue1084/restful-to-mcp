package biz

import (
	"context"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
)

type McpServerRepo interface {
	Create(ctx context.Context, serverInfo *model.McpServer) (err error)
	GetMcpServerInfo(ctx context.Context, UUID string) (serverInfo *model.McpServer, err error)
	CreateWithTx(ctx context.Context, serverInfo *model.McpServer) (id int64, err error)
	GetMcpServerInfoByID(ctx context.Context, id int64) (mcpServerInfo *model.McpServer, err error)
	GetMcpServerInfoByUUID(ctx context.Context, id string) (mcpServerInfo *model.McpServer, err error)
	UpdateMcpServerForAuthWithTx(ctx context.Context, uuid string, isAuth _const.AuthStatus, serviceToken, platformToken string) (err error)
	UpdateMcpServerByUUID(ctx context.Context, uuid, name, description string) (resp *api.UpdateMcpServerByUUIDResponse, err error)
}

type McpServerUseCase struct {
	msRepo McpServerRepo
	log    *logger.Logger
}

func NewMcpServerUseCase(msRepo McpServerRepo, log *logger.Logger) *McpServerUseCase {
	return &McpServerUseCase{
		msRepo: msRepo,
		log:    log,
	}

}

func (m *McpServerUseCase) UpdateMcpServerByUUID(ctx context.Context, uuid, name, description string) (resp *api.UpdateMcpServerByUUIDResponse, err error) {

	resp, err = m.msRepo.UpdateMcpServerByUUID(ctx, uuid, name, description)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerByUUID error: %v", err)
		return
	}
	return
}
