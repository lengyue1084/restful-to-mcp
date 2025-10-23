package biz

import (
	"context"
	"encoding/json"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/config"
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
	var header []map[string]string
	if mcpServerInfo.Header != "" {
		err = json.Unmarshal([]byte(mcpServerInfo.Header), &header)
		if err != nil {
			return nil, err
		}
	}

	var security config.Security
	if err := json.Unmarshal([]byte(mcpServerInfo.Security), &security); err != nil {
		// 处理错误，例如记录日志或返回错误
		m.log.ErrorWithContext(ctx, "Failed to unmarshal security info: %v", err)
		return nil, err
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
			IsAuth:        _const.AuthTypeStatus(mcpServerInfo.IsAuth),
			PlatformToken: mcpServerInfo.PlatformToken,
			ServiceToken:  mcpServerInfo.ServiceToken,
			Header:        header,
			Security:      security,
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
	// 如果为api服务鉴权，则必须填写鉴权mode
	if req.IsAuth == _const.IsAuthServiceAuth && req.Security.Mode == "" {
		return nil, fmt.Errorf("api鉴权模式，缺少必要的鉴权参数")
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
	var headerJson string
	if req.Header != nil && len(req.Header) > 0 {
		headerBytes, err := json.Marshal(req.Header)
		if err != nil {
			return nil, fmt.Errorf("header json error: %v", err)
		}
		headerJson = string(headerBytes)
	}

	securityStr := "{}"
	if req.IsAuth == _const.IsAuthServiceAuth {
		var security = &config.Security{}
		security.SecurityKey = req.Security.SecurityKey
		if req.Security.Mode == "" {
			return nil, fmt.Errorf("请输入鉴权参数")
		}
		security.Mode = req.Security.Mode
		switch req.Security.Mode {
		case config.AuthModeHttp:
			// http 模式下 scheme 必填
			if req.Security.Scheme == "" {
				return nil, fmt.Errorf("%s模式，scheme为必填项", config.AuthModeHttp)
			}
			security.Scheme = req.Security.Scheme
		case config.AuthModeApiKey:
			// apiKey 模式下 name、in 必填
			if req.Name == "" {
				return nil, fmt.Errorf("%s模式，name为必填项", config.AuthModeApiKey)
			}
			security.Name = req.Name
			if req.Security.In == "" {
				return nil, fmt.Errorf("%s模式，position为必填项", config.AuthModeApiKey)
			}
			security.In = req.Security.In
		}

		securityByte, err := json.Marshal(security)
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool security json转换错误，err:%+v", err)
			return nil, err
		}
		if req.Security.SecurityKey != "" {
			securityStr = string(securityByte)
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
		Security:      securityStr,
		Status:        _const.ServerHadSetToken,
		SerialNumber:  serialNumber,
		Source:        _const.SourceTypeForm,
		Header:        headerJson,
	}
	mcpServerInfo, err = m.msRepo.CreateMcpServerByForm(ctx, mcpServerInfo)
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpServerByForm error: %v", err)
		return nil, err
	}
	resp = &api.CreateMcpServerByFormResponse{
		ID:        mcpServerInfo.ID,
		UUID:      mcpServerInfo.UUID,
		CreatedAt: mcpServerInfo.CreatedAt.String(),
	}
	return
}

func (m *McpServerUseCase) UpdateMcpServerByForm(ctx context.Context, req *api.UpdateMcpServerByFormRequest) (resp *api.UpdateMcpServerByFormResponse, err error) {

	if len(req.Urls) == 0 {
		return nil, fmt.Errorf("至少填写一个url")
	}
	// 如果为api服务鉴权，则必须填写鉴权mode
	if req.IsAuth == _const.IsAuthServiceAuth && req.Security.Mode == "" {
		return nil, fmt.Errorf("api鉴权模式，缺少必要的鉴权参数")
	}
	urlsStr, err := json.Marshal(req.Urls)
	if err != nil {
		m.log.ErrorWithContext(ctx, "CreateMcpServerByForm error: %+v", err)
		return nil, err
	}
	var headerJson string
	if req.Header != nil && len(req.Header) > 0 {
		headerBytes, err := json.Marshal(req.Header)
		if err != nil {
			return nil, fmt.Errorf("header json error: %v", err)
		}
		headerJson = string(headerBytes)
	}
	securityStr := "{}"
	if req.IsAuth == _const.IsAuthServiceAuth {
		var security = &config.Security{}
		security.SecurityKey = req.Security.SecurityKey
		if req.Security.Mode == "" {
			return nil, fmt.Errorf("请输入鉴权参数")
		}
		security.Mode = req.Security.Mode
		switch req.Security.Mode {
		case config.AuthModeHttp:
			// http 模式下 scheme 必填
			if req.Security.Scheme == "" {
				return nil, fmt.Errorf("%s模式，scheme为必填项", config.AuthModeHttp)
			}
			security.Scheme = req.Security.Scheme
		case config.AuthModeApiKey:
			// apiKey 模式下 name、in 必填
			if req.Name == "" {
				return nil, fmt.Errorf("%s模式，name为必填项", config.AuthModeApiKey)
			}
			security.Name = req.Name
			if req.Security.In == "" {
				return nil, fmt.Errorf("%s模式，position为必填项", config.AuthModeApiKey)
			}
			security.In = req.Security.In
		}

		securityByte, err := json.Marshal(security)
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool security json转换错误，err:%+v", err)
			return nil, err
		}
		if req.Security.SecurityKey != "" {
			securityStr = string(securityByte)
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
		Header:        headerJson,
		Security:      securityStr,
	}
	mcpServerInfo, err = m.msRepo.UpdateMcpServerByForm(ctx, mcpServerInfo)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerByForm error: %v", err)
		return nil, err
	}
	resp = &api.UpdateMcpServerByFormResponse{
		ID:        mcpServerInfo.ID,
		UUID:      mcpServerInfo.UUID,
		CreatedAt: mcpServerInfo.CreatedAt.String(),
	}
	return
}
