package biz

import (
	"context"
	"encoding/json"
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
	DeleteMcpServerByUUID(ctx context.Context, uuid string) (err error)
	GetCountMcpServerInfoBySerialNumber(ctx context.Context, serialNumber string) (count int64, err error)
	CreateMcpServerByForm(ctx context.Context, serverInfo *model.McpServer) (mcpServer *model.McpServer, err error)
	UpdateMcpServerByForm(ctx context.Context, serverInfo *model.McpServer) (mcpServer *model.McpServer, err error)
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
func (m *McpServerUseCase) GetMcpServerInfoByUUID(ctx context.Context, uuid string) (resp *api.GetMcpServerInfoByUUIDResponse, err error) {

	mcpServerInfo, err := m.msRepo.GetMcpServerInfoByUUID(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return
	}
	if mcpServerInfo.ID == 0 {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return nil, fmt.Errorf("没有查询到该server信息")
	}
	var urls []string
	err = json.Unmarshal([]byte(mcpServerInfo.Urls), &urls)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return
	}
	resp = &api.GetMcpServerInfoByUUIDResponse{
		ID:        mcpServerInfo.ID,
		CreatedAt: mcpServerInfo.CreatedAt.String(),
		UpdatedAt: mcpServerInfo.UpdatedAt.String(),
		CommonMcpServerByForm: api.CommonMcpServerByForm{
			UUID:          mcpServerInfo.UUID,
			Name:          mcpServerInfo.Name,
			Description:   mcpServerInfo.Description,
			Urls:          urls,
			Version:       mcpServerInfo.Version,
			IsAuth:        int8(mcpServerInfo.IsAuth),
			PlatformToken: mcpServerInfo.PlatformToken,
		},
	}

	return
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

func (m *McpServerUseCase) DeleteMcpServerByUUID(ctx context.Context, uuid string) (err error) {

	err = m.msRepo.DeleteMcpServerByUUID(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "DeleteMcpServerByUUID error: %v", err)
		return err
	}
	return
}

func (m *McpServerUseCase) CreateMcpServerByForm(ctx context.Context, req *api.CreateMcpServerByFormRequest) (resp *api.CreateMcpServerByFormResponse, err error) {

	if len(req.Urls) == 0 {
		return nil, fmt.Errorf("至少填写一个url")
	}
	urlsStr, err := json.Marshal(req.Urls)
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpServerByForm error: %+v", err)
		return nil, err
	}
	var (
		serialNumber = ""
	)
	for i := 0; i < _const.CommonRetryTimes; i++ {
		serialNumber = tool.RandStringWithLowercaseAndDigits(6)
		count, err := m.msRepo.GetCountMcpServerInfoBySerialNumber(ctx, serialNumber)
		if err != nil {
			m.log.ErrorWithContext(ctx, "CreateMcpServerByForm error: %v", err)
			return nil, err
		}
		if count == 0 {
			break
		}
		if i == _const.CommonRetryTimes-1 {
			return nil, fmt.Errorf("生成序列号失败")
		}
	}

	var mcpServerInfo = &model.McpServer{
		UUID:          req.UUID,
		Name:          req.Name,
		Description:   req.Description,
		Urls:          string(urlsStr),
		AllTools:      "[]",
		Version:       req.Version,
		McpServerType: _const.McpServerTypeOpenapi,
		HaveTools:     _const.HaveToolsNo,
		IsAuth:        _const.AuthStatus(req.IsAuth),
		ServiceToken:  req.ServiceToken,
		PlatformToken: req.PlatformToken,
		Security:      "",
		Status:        _const.ServerHadSetToken,
		SerialNumber:  serialNumber,
		Source:        _const.SourceTypeForm,
	}
	mcpServerInfo, err = m.msRepo.CreateMcpServerByForm(ctx, mcpServerInfo)
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpServerByForm error: %v", err)
		return nil, err
	}
	resp = &api.CreateMcpServerByFormResponse{
		ID:        mcpServerInfo.ID,
		CreatedAt: mcpServerInfo.CreatedAt.String(),
		CommonMcpServerByForm: api.CommonMcpServerByForm{
			UUID:          mcpServerInfo.UUID,
			Name:          mcpServerInfo.Name,
			Description:   mcpServerInfo.Description,
			Urls:          nil,
			Version:       mcpServerInfo.Version,
			IsAuth:        int8(mcpServerInfo.IsAuth),
			PlatformToken: mcpServerInfo.PlatformToken,
		},
	}
	var urls []string
	if err = json.Unmarshal([]byte(mcpServerInfo.Urls), &urls); err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpServerByForm error: %v", err)
		return nil, err
	}
	resp.Urls = urls
	return
}

func (m *McpServerUseCase) UpdateMcpServerByForm(ctx context.Context, req *api.UpdateMcpServerByFormRequest) (resp *api.UpdateMcpServerByFormResponse, err error) {

	urlsStr, err := json.Marshal(req.Urls)
	if err != nil {
		return nil, err
	}
	var mcpServerInfo = &model.McpServer{
		UUID:          req.UUID,
		Name:          req.Name,
		Description:   req.Description,
		Urls:          string(urlsStr),
		Version:       req.Version,
		IsAuth:        _const.AuthStatus(req.IsAuth),
		PlatformToken: req.PlatformToken,
	}
	mcpServerInfo, err = m.msRepo.UpdateMcpServerByForm(ctx, mcpServerInfo)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerByForm error: %v", err)
		return nil, err
	}
	resp = &api.UpdateMcpServerByFormResponse{
		ID:        mcpServerInfo.ID,
		CreatedAt: mcpServerInfo.CreatedAt.String(),
		UpdatedAt: mcpServerInfo.UpdatedAt.String(),
		CommonMcpServerByForm: api.CommonMcpServerByForm{
			UUID:          mcpServerInfo.UUID,
			Name:          mcpServerInfo.Name,
			Description:   mcpServerInfo.Description,
			Urls:          req.Urls,
			Version:       mcpServerInfo.Version,
			IsAuth:        int8(mcpServerInfo.IsAuth),
			PlatformToken: mcpServerInfo.PlatformToken,
		},
	}
	return
}
