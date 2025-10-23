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
	GetMcpServerToolsByServerUUID(ctx context.Context, uuid string) (mcpTools []*model.McpTools, err error)
	CreateMcpServerTool(ctx context.Context, mcpToolInfo *model.McpTools) (uuid string, err error)
	GetMcpServerToolByNameWithUUID(ctx context.Context, uuid string, name string) (tool *model.McpTools, err error)
	GetMcpServerToolByName(ctx context.Context, name string) (tool *model.McpTools, err error)
	GetMcpServerToolInfoByUUID(ctx context.Context, uuid string) (tool *model.McpTools, err error)
	UpdateMcpServerTool(ctx context.Context, tool map[string]interface{}, uuid string) (err error)
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
	tools, err := m.mtRepo.GetMcpServerToolsByServerUUID(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "查询工具列表失败,err:%+v", err)
		return nil, fmt.Errorf("查询工具列表失败,err:%+v", err)
	}

	var toolsList = make([]*protocol.Tool, 0)
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

	tools, err := m.mtRepo.GetMcpServerToolsByServerUUID(ctx, uuid)
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

	var annotationsMap = make(map[string]string)
	annotationsMap["description"] = req.Description
	annotationsMap["title"] = req.Name
	annotationsJson, err := json.Marshal(annotationsMap)
	if err != nil {
		m.log.ErrorWithContext(ctx, "annotations json转换错误，err:%+v", err)
		return nil, err
	}

	// 从path中提取参数
	args := tool.ExtractArgsFromPath(path)
	argsJson, err := json.Marshal(args)
	if err != nil {
		m.log.ErrorWithContext(ctx, "args json转换错误，err:%+v", err)
		return nil, err
	}
	var authMode config.Security
	err = json.Unmarshal([]byte(mcpServerInfo.Security), &authMode)
	if err != nil {
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
		IsShow:         _const.StatusHidden,
		IsPlatformAuth: req.IsPlatformAuth,
		IsAuth:         req.IsAuth,
		AuthMode:       authMode.Mode.String(),
		Security:       mcpServerInfo.Security,
		IsRepeat:       isRepeat,
		ResponseBody:   "{{.Response.Body}}",
		Annotations:    string(annotationsJson),
		Args:           string(argsJson),
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

	var path string
	if req.Path != nil && *req.Path != "" {
		path = tool.ConvertPathToArgsFormat(*req.Path)
		if path[0] != '/' {
			path = fmt.Sprintf("/%s", path)
		}
	}

	toolInfo, err := m.mtRepo.GetMcpServerToolInfoByUUID(ctx, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool error: %v", err)
		return nil, err
	}
	if toolInfo.ID == 0 {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool error: %v,tool uuid：%v", "工具不存在", req.UUID)
		return nil, fmt.Errorf("工具不存在")
	}
	isRepeat := toolInfo.IsRepeat
	// 判断该mcp server 工具名是否重复,同一个mcp server 下不允许工具名重复
	if req.Name != nil && toolInfo.Name != *req.Name {
		toolInfoNew, err := m.mtRepo.GetMcpServerToolByNameWithUUID(ctx, toolInfo.McpServerUUID, *req.Name)
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerToolByNameWithUUID error: %v", err)
			return nil, err
		}
		if toolInfoNew.ID != 0 {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerToolByNameWithUUID error: %v", err)
			return nil, fmt.Errorf("该服务的工具%s已存在", *req.Name)
		}
		isRepeat = _const.CommonStatusNo
		// 判断工具名是否重复，这里不区分 server
		toolInfoForName, err := m.mtRepo.GetMcpServerToolByName(ctx, *req.Name)
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerToolsByName error: %v", err)
			return nil, err
		}
		if toolInfoForName.ID != 0 {
			isRepeat = _const.CommonStatusYes
		}
	}

	var dataMap = make(map[string]interface{})
	if req.IsAuth != nil {
		dataMap["is_auth"] = *req.IsAuth
	}

	var annotationsMap = make(map[string]string)
	if req.Description != nil {
		annotationsMap["description"] = *req.Description
	}
	if req.Name != nil {
		annotationsMap["title"] = *req.Name
	}
	annotationsJson, err := json.Marshal(annotationsMap)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool annotations json转换错误，err:%+v", err)
		return nil, err
	}

	if req.Name != nil {
		dataMap["name"] = req.Name
	}
	if req.Description != nil {
		dataMap["description"] = req.Description
	}
	if req.Path != nil {
		dataMap["endpoint"] = fmt.Sprintf("{{.Config.url}}%s", path)
	}
	if req.Method != nil {
		dataMap["method"] = req.Method
	}
	if req.IsShow != nil {
		dataMap["is_show"] = req.IsShow
	}
	if req.IsPlatformAuth != nil {
		dataMap["is_platform_auth"] = req.IsPlatformAuth
	}
	if req.IsAuth != nil {
		dataMap["is_auth"] = req.IsAuth
	}

	dataMap["is_repeat"] = isRepeat
	dataMap["response_body"] = "{{.Response.Body}}"
	dataMap["annotations"] = string(annotationsJson)
	var args []*config.ArgConfig
	if path != "" {
		// 从path中提取参数
		args = tool.ExtractArgsFromPath(path)

	}
	argsJson, err := json.Marshal(args)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool args json转换错误，err:%+v", err)
		return nil, err
	}
	dataMap["args"] = string(argsJson)

	//[{"name":"lastName","position":"body","required":false,"type":"string","description":"","default":"","items":{"type":""},"explode":false},{"name":"password","position":"body","required":false,"type":"string","description":"","default":"","items":{"type":""},"explode":false},{"name":"phone","position":"body","required":false,"type":"string","description":"","default":"","items":{"type":""},"explode":false},{"name":"userStatus","position":"body","required":false,"type":"integer","description":"User Status","default":"","items":{"type":""},"explode":false},{"name":"username","position":"body","required":false,"type":"string","description":"","default":"","items":{"type":""},"explode":false},{"name":"email","position":"body","required":false,"type":"string","description":"","default":"","items":{"type":""},"explode":false},{"name":"firstName","position":"body","required":false,"type":"string","description":"","default":"","items":{"type":""},"explode":false},{"name":"id","position":"body","required":false,"type":"integer","description":"","default":"","items":{"type":""},"explode":false}]

	err = m.mtRepo.UpdateMcpServerTool(ctx, dataMap, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool 创建工具失败,err:%+v", err)
		return nil, fmt.Errorf("创建工具失败,err:%+v", err)
	}
	resp = &api.UpdateMcpServerToolResponse{}
	return

}
