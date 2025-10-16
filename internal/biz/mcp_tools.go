package biz

import (
	"context"
	"encoding/json"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/config"
	"flow-bridge-mcp/internal/pkg/cache"
	_const "flow-bridge-mcp/pkg/const"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

type McpToolsRepo interface {
	Create(ctx context.Context, mcpToolInfo *model.McpTools) (err error)
	CreateMcpToolsBatch(ctx context.Context, mcpServerId int64, uuid string, allTools []string, mcpToolInfo []*model.McpTools) (err error)
	UpdateToolsForAuthWithTx(ctx context.Context, uuid string, tools []*api.Tools) (err error)
	GetMcpServerTools(ctx context.Context, uuid string) (mcpTools []*model.McpTools, err error)
	GetMcpServerToolsByUUID(ctx context.Context, uuid string) (mcpTools []*model.McpTools, err error)
	CreateMcpServerTool(ctx context.Context, mcpToolInfo *model.McpTools) (uuid string, err error)
	GetMcpServerToolByNameWithUUID(ctx context.Context, uuid string, name string) (tool *model.McpTools, err error)
	GetMcpServerToolByName(ctx context.Context, name string) (tool *model.McpTools, err error)
}

type McpToolsUserCase struct {
	mtRepo McpToolsRepo
	msRepo McpServerRepo
	log    *logger.Logger
	cache  *cache.MemoryCache
}

func NewMcpToolsUserCase(mtRepo McpToolsRepo, msRepo McpServerRepo, log *logger.Logger, cache *cache.MemoryCache) *McpToolsUserCase {
	return &McpToolsUserCase{
		mtRepo: mtRepo,
		msRepo: msRepo,
		log:    log,
		cache:  cache,
	}
}

func (m *McpToolsUserCase) GetMcpServerTools(ctx context.Context, uuid string) (resp *api.GetMcpServerToolsResponse, err error) {
	tools, err := m.mtRepo.GetMcpServerTools(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "查询工具列表失败,err:%+v", err)
		return nil, fmt.Errorf("查询工具列表失败,err:%+v", err)
	}

	var toolsList = make([]*protocol.Tool, len(tools))
	for _, tool := range tools {
		var annotations *protocol.ToolAnnotations
		if tool.Annotations != "" {
			err = json.Unmarshal([]byte(tool.Annotations), &annotations)
			if err != nil {
				m.log.ErrorWithContext(ctx, "tool.Annotations json转换错误，err:%+v", err)
				return nil, err
			}
		}
		var toolSchema = &protocol.InputSchema{}
		err = json.Unmarshal([]byte(tool.ToolSchema), toolSchema)
		if err != nil {
			m.log.ErrorWithContext(ctx, "tool.ToolSchema json转换错误，err:%+v", err)
			return nil, err
		}
		tmpTool := &protocol.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: *toolSchema,
			Annotations: annotations,
		}
		toolsList = append(toolsList, tmpTool)
	}
	resp = &api.GetMcpServerToolsResponse{
		Tools: toolsList,
	}

	return resp, nil
}

func (m *McpToolsUserCase) GetMcpServerToolsByUUID(ctx context.Context, uuid string) (resp *api.GetMcpServerToolsByUUIDResponse, err error) {

	tools, err := m.mtRepo.GetMcpServerToolsByUUID(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "查询工具列表失败,err:%+v", err)
		return nil, fmt.Errorf("查询工具列表失败,err:%+v", err)
	}
	var toolsList []*api.CommonToolItemInfo
	if err = tool.Copy(&toolsList, tools); err != nil {
		m.log.ErrorWithContext(ctx, "工具列表转换失败,err:%+v", err)
		return nil, fmt.Errorf("工具列表转换失败,err:%+v", err)
	}
	resp = &api.GetMcpServerToolsByUUIDResponse{
		Tools: toolsList,
	}

	return
}

func (m *McpToolsUserCase) CreateMcpServerTool(ctx context.Context, req *api.CreateMcpServerToolRequest) (resp *api.CreateMcpServerToolResponse, err error) {
	if req.Path == "" {
		return nil, fmt.Errorf("请输入正确的路径")
	}
	path := tool.ConvertPathToArgsFormat(req.Path)
	if path[0] != '/' {
		path = fmt.Sprintf("/%s", path)
	}
	if req.IsAuth == _const.IsAuthYes && req.AuthMode == "" {
		return nil, fmt.Errorf("请输入鉴权参数")
	}
	mcpServerInfo, err := m.msRepo.GetMcpServerInfoByUUID(ctx, req.McpServerUUID)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return nil, err
	}
	if mcpServerInfo.ID == 0 {
		m.log.ErrorWithContext(ctx, "GetMcpServerInfoByUUID error: %v", err)
		return nil, fmt.Errorf("没有查询到该server信息")
	}

	toolInfo, err := m.mtRepo.GetMcpServerToolByNameWithUUID(ctx, req.McpServerUUID, req.Name)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerToolsByName error: %v", err)
		return nil, err
	}
	if toolInfo.ID != 0 {
		m.log.ErrorWithContext(ctx, "GetMcpServerToolsByName error: %v", err)
		return nil, fmt.Errorf("工具%s已存在", req.Name)
	}

	// 判断工具名是否重复，这里不区分 server
	toolInfoForName, err := m.mtRepo.GetMcpServerToolByName(ctx, req.Name)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetMcpServerToolsByName error: %v", err)
		return nil, err
	}
	isRepeat := _const.CommonStatusNo
	if toolInfoForName.ID != 0 {
		isRepeat = _const.CommonStatusYes
	}

	var security = &config.Security{
		SecurityKey: req.SecurityKey,
		Mode:        req.AuthMode,
		Name:        req.SecurityKey,
		Scheme:      req.Scheme,
		In:          req.Position,
		Description: "",
	}
	securityByte, err := json.Marshal(security)
	if err != nil {
		m.log.ErrorWithContext(ctx, "security json转换错误，err:%+v", err)
		return nil, err
	}

	var annotationsMap = make(map[string]string)
	annotationsMap["description"] = req.Description
	annotationsMap["title"] = req.Name
	annotationsJson, err := json.Marshal(annotationsMap)
	if err != nil {
		m.log.ErrorWithContext(ctx, "annotations json转换错误，err:%+v", err)
		return nil, err
	}

	tools := &model.McpTools{
		UUID:           tool.NewUUID(),
		Name:           req.Name,
		Description:    req.Description,
		Endpoint:       fmt.Sprintf("{{.Config.url}}%s", path),
		Method:         req.Method,
		McpServerUUID:  req.McpServerUUID,
		McpServerId:    int64(mcpServerInfo.ID),
		SerialNumber:   mcpServerInfo.SerialNumber,
		IsShow:         req.IsShow,
		IsPlatformAuth: req.IsPlatformAuth,
		IsAuth:         req.IsAuth,
		AuthMode:       req.AuthMode.String(),
		Security:       string(securityByte),
		IsRepeat:       isRepeat,
		ResponseBody:   "{{.Response.Body}}",
		Annotations:    string(annotationsJson),
	}

	uuid, err := m.mtRepo.CreateMcpServerTool(ctx, tools)
	if err != nil {
		m.log.ErrorWithContext(ctx, "创建工具失败,err:%+v", err)
		return nil, fmt.Errorf("创建工具失败,err:%+v", err)
	}
	resp = &api.CreateMcpServerToolResponse{
		UUID: uuid,
	}
	return
}

func (m *McpToolsUserCase) UpdateMcpServerTool(ctx context.Context, req *api.UpdateMcpServerToolRequest) (resp *api.UpdateMcpServerToolResponse, err error) {

	return

}
