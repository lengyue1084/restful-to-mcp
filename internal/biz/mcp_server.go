package biz

import (
	"context"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
)

type McpServerRepo interface {
	Create(ctx context.Context, serverInfo *model.McpServer) (err error)
	GetMcpServerInfo(ctx context.Context, UUID string) (serverInfo *model.McpServer, err error)
	GetMcpServerInfoWithAllTools(ctx context.Context) (serverInfo []*model.McpServer, err error)
	CreateWithTx(ctx context.Context, serverInfo *model.McpServer) (id int64, err error)
	GetMcpServerInfoByID(ctx context.Context, id int64) (mcpServerInfo *model.McpServer, err error)
	GetMcpServerInfoByUUID(ctx context.Context, id string) (mcpServerInfo *model.McpServer, err error)
	UpdateMcpServerForAuthWithTx(ctx context.Context, uuid string, isAuth _const.AuthStatus, serviceToken, platformToken string) (err error)
	UpdateMcpServerByUUID(ctx context.Context, uuid, name, description string) (resp *api.UpdateMcpServerByUUIDResponse, err error)
	GetCountMcpServerInfoBySerialNumber(ctx context.Context, serialNumber string) (count int64, err error)
}

type McpServerUseCase struct {
	msRepo  McpServerRepo
	log     *logger.Logger
	mctRepo McpConnectTokenRepo
}

func NewMcpServerUseCase(msRepo McpServerRepo, mctRepo McpConnectTokenRepo, log *logger.Logger) *McpServerUseCase {
	return &McpServerUseCase{
		msRepo:  msRepo,
		log:     log,
		mctRepo: mctRepo,
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

func (m *McpServerUseCase) GetMcpConnectTokenByUUID(ctx context.Context, uuid string) (resp *api.GetMcpConnectTokenByUUIDResponse, err error) {

	mcpServerInfo, err := m.msRepo.GetMcpServerInfoByUUID(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return nil, err
	}
	if mcpServerInfo.ID == 0 {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return nil, fmt.Errorf("没有查询到该server信息")
	}
	connectToken := tool.RandStringWithLowercaseAndDigits(16)
	var mcpConnectToken = &model.McpConnectToken{
		McpServerUUID: uuid,
		McpServerId:   int64(mcpServerInfo.ID),
		McpServerName: mcpServerInfo.Name,
		ConnectToken:  connectToken,
	}
	err = m.mctRepo.Create(ctx, mcpConnectToken)
	if err != nil {
		m.log.ErrorWithContext(ctx, "create mcp connect token error: %v", err)
		return nil, err
	}
	resp = &api.GetMcpConnectTokenByUUIDResponse{
		ConnectToken: connectToken,
	}
	return
}

func (m *McpServerUseCase) GetMcpServerInfoWithAllTools(ctx context.Context) (mcpServerInfo []*model.McpServer, err error) {
	mcpServerInfo, err = m.msRepo.GetMcpServerInfoWithAllTools(ctx)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoWithTools error: %v", err)
		return nil, err
	}
	return
}
