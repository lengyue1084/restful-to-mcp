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
	UpdateMcpServerTool(ctx context.Context, tool *model.McpTools, uuid string) (err error)
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
		if tool.ToolSchema != "" {
			err = json.Unmarshal([]byte(tool.ToolSchema), toolSchema)
			if err != nil {
				m.log.ErrorWithContext(ctx, "tool.ToolSchema json转换错误，err:%+v", err)
				return nil, err
			}
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
		m.log.ErrorWithContext(ctx, "GetMcpServerToolsByName the tool already exists, tool name: %v", err)
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

	var argsSlice []config.ArgConfig
	// 从path中提取参数
	args := tool.ExtractArgsFromPath(path)
	if len(args) > 0 {
		argsSlice = args
	}
	// 解析入参的InputArgs,直接从前端获取，py层负责转换，这里就直接获取即可
	var inputArgs = req.Args
	argsSlice = append(argsSlice, inputArgs...)
	argsJson, err := json.Marshal(argsSlice)
	if err != nil {
		m.log.ErrorWithContext(ctx, "args json转换错误，err:%+v", err)
		return nil, err
	}
	// 鉴权从server中提取，对于表单构建的方式鉴权，server中入库的为可用的鉴权，openapi可能为多种并存
	var securitySlice []config.Security
	var security = &config.Security{}
	if mcpServerInfo.Security != "" {
		err = json.Unmarshal([]byte(mcpServerInfo.Security), &securitySlice)
		if err != nil {
			return nil, err
		}
		if len(securitySlice) > 0 {
			security = &securitySlice[0]
		}
	}
	securityJson, err := json.Marshal(security)
	if err != nil {
		m.log.ErrorWithContext(ctx, "security json转换错误，err:%+v", err)
		return nil, err
	}
	// 是否从参数中解析 header部分填充到header中，不需要，请求端已经添加了
	var (
		headersMap  = make(map[string]string)
		headersJson []byte
	)
	if mcpServerInfo.Header != "" {
		err := json.Unmarshal([]byte(mcpServerInfo.Header), &headersMap)
		if err != nil {
			m.log.ErrorWithContext(ctx, "mcpServerInfo.Header json转换错误，err:%+v", err)
			return nil, err
		}
	}
	if len(headersMap) == 0 {
		if _, ok := headersMap["Content-Type"]; !ok {
			headersMap["Content-Type"] = "application/json"
		}
	}
	headersJson, err = json.Marshal(headersMap)
	if err != nil {
		m.log.ErrorWithContext(ctx, "headers json转换错误，err:%+v", err)
		return nil, err
	}

	var toolConfig = &config.ToolConfig{
		Args: args,
	}
	var inputSchema = toolConfig.ArgsToInputSchema()
	toolSchemaJson, err := json.Marshal(inputSchema)
	if err != nil {
		m.log.ErrorWithContext(ctx, "toolSchema json转换错误，err:%+v", err)
		return nil, err
	}

	tools := &model.McpTools{
		UUID:           tool.NewUUID(),
		Name:           req.Name,
		Description:    req.Description,
		Endpoint:       fmt.Sprintf("{{.Config.url}}%s", path),
		Method:         req.Method,
		Headers:        string(headersJson),
		McpServerUUID:  req.McpServerUUID,
		McpServerId:    int64(mcpServerInfo.ID),
		SerialNumber:   mcpServerInfo.SerialNumber,
		IsShow:         _const.StatusHidden,
		IsPlatformAuth: req.IsPlatformAuth,
		IsAuth:         req.IsAuth,
		AuthMode:       security.Mode.String(),
		Security:       string(securityJson),
		IsRepeat:       isRepeat,
		ToolSchema:     string(toolSchemaJson),
		RequestBody:    "{}",                 //暂时该字段没有使用，这里就不用构建了
		ResponseBody:   "{{.Response.Body}}", //暂时该字段没有使用，这里就不用构建了
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
	if req.Path != "" {
		path = tool.ConvertPathToArgsFormat(req.Path)
		if path != "" && path[0] != '/' {
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
	mcpServerInfo, err := m.msRepo.GetMcpServerInfoByUUID(ctx, toolInfo.McpServerUUID)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerInfoByUUID error: %v", err)
		return nil, err
	}
	if mcpServerInfo.ID == 0 {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerInfoByUUID error: %v", err)
		return nil, fmt.Errorf("没有查询到该server信息")
	}
	isRepeat := toolInfo.IsRepeat
	// 判断该mcp server 工具名是否重复,同一个mcp server 下不允许工具名重复
	if toolInfo.Name != req.Name {
		toolInfoNew, err := m.mtRepo.GetMcpServerToolByNameWithUUID(ctx, toolInfo.McpServerUUID, req.Name)
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerToolByNameWithUUID error: %v", err)
			return nil, err
		}
		if toolInfoNew.ID != 0 {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerToolByNameWithUUID error: %v", err)
			return nil, fmt.Errorf("该服务的工具%s已存在", req.Name)
		}
		isRepeat = _const.CommonStatusNo
		// 判断工具名是否重复，这里不区分 server
		toolInfoForName, err := m.mtRepo.GetMcpServerToolByName(ctx, req.Name)
		if err != nil {
			m.log.ErrorWithContext(ctx, "UpdateMcpServerTool GetMcpServerToolsByName error: %v", err)
			return nil, err
		}
		if toolInfoForName.ID != 0 {
			isRepeat = _const.CommonStatusYes
		}
	}

	// 组装参数，从path中提取参数,并合并前端入参的参数
	var argsSlice []config.ArgConfig
	args := tool.ExtractArgsFromPath(path)
	if len(args) > 0 {
		argsSlice = args
	}
	// 解析入参的InputArgs,直接从前端获取，py层负责转换，这里就直接获取即可
	var inputArgs = req.Args
	argsSlice = append(argsSlice, inputArgs...)
	argsJson, err := json.Marshal(argsSlice)
	if err != nil {
		m.log.ErrorWithContext(ctx, "args json转换错误，err:%+v", err)
		return nil, err
	}

	if req.IsPlatformAuth != _const.IsAuthNo && req.IsPlatformAuth != _const.IsAuthYes {
		return nil, fmt.Errorf("请选择正确的平台鉴权方式")
	}
	if req.IsAuth != _const.IsAuthNo && req.IsAuth != _const.IsAuthYes {
		return nil, fmt.Errorf("请选择正确接口鉴权的鉴权方式")
	}

	var annotationsMap = make(map[string]string)
	if req.Description != "" {
		annotationsMap["description"] = req.Description
	}
	if req.Name != "" {
		annotationsMap["title"] = req.Name
	}
	annotationsJson, err := json.Marshal(annotationsMap)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool annotations json转换错误，err:%+v", err)
		return nil, err
	}

	// 鉴权从server中提取，对于表单构建的方式鉴权，server中入库的为可用的鉴权，openapi可能为多种并存
	var securitySlice []config.Security
	var security = &config.Security{}
	if mcpServerInfo.Security != "" {
		err = json.Unmarshal([]byte(mcpServerInfo.Security), &securitySlice)
		if err != nil {
			return nil, err
		}
		if len(securitySlice) > 0 {
			security = &securitySlice[0]
		}
	}
	securityJson, err := json.Marshal(security)
	if err != nil {
		m.log.ErrorWithContext(ctx, "security json转换错误，err:%+v", err)
		return nil, err
	}

	// 是否从参数中解析 header部分填充到header中，不需要，请求端已经添加了
	var (
		headersMap  = make(map[string]string)
		headersJson []byte
	)
	if mcpServerInfo.Header != "" {
		err := json.Unmarshal([]byte(mcpServerInfo.Header), &headersMap)
		if err != nil {
			m.log.ErrorWithContext(ctx, "mcpServerInfo.Header json转换错误，err:%+v", err)
			return nil, err
		}
	}
	if len(headersMap) == 0 {
		if _, ok := headersMap["Content-Type"]; !ok {
			headersMap["Content-Type"] = "application/json"
		}
	}
	headersJson, err = json.Marshal(headersMap)
	if err != nil {
		m.log.ErrorWithContext(ctx, "headers json转换错误，err:%+v", err)
		return nil, err
	}

	var toolConfig = &config.ToolConfig{
		Args: args,
	}
	var inputSchema = toolConfig.ArgsToInputSchema()
	toolSchemaJson, err := json.Marshal(inputSchema)
	if err != nil {
		m.log.ErrorWithContext(ctx, "toolSchema json转换错误，err:%+v", err)
		return nil, err
	}

	var toolInfoModel = &model.McpTools{
		UUID:           req.UUID,
		Name:           req.Name,
		Description:    req.Description,
		Method:         req.Method,
		Endpoint:       fmt.Sprintf("{{.Config.url}}%s", path),
		Headers:        string(headersJson),
		Args:           string(argsJson),
		Security:       string(securityJson),
		ToolSchema:     string(toolSchemaJson),
		Annotations:    string(annotationsJson),
		IsShow:         _const.StatusHidden,
		IsPlatformAuth: req.IsPlatformAuth,
		IsAuth:         req.IsAuth,
		AuthMode:       security.Mode.String(),
		IsRepeat:       isRepeat, // 做个判断，如果前端有重复，则返回错误
	}

	err = m.mtRepo.UpdateMcpServerTool(ctx, toolInfoModel, req.UUID)
	if err != nil {
		m.log.ErrorWithContext(ctx, "UpdateMcpServerTool 创建工具失败,err:%+v", err)
		return nil, fmt.Errorf("创建工具失败,err:%+v", err)
	}
	resp = &api.UpdateMcpServerToolResponse{}
	return

}

func (m *McpToolsUserCase) GetToolsInfoByUUID(ctx context.Context, uuid string) (resp *api.GetToolsInfoByUUIDResponse, err error) {

	toolInfo, err := m.mtRepo.GetMcpServerToolInfoByUUID(ctx, uuid)
	if err != nil {
		m.log.ErrorWithContext(ctx, "GetToolsInfoByUUID error: %v", err)
		return nil, err
	}
	resp = &api.GetToolsInfoByUUIDResponse{
		ID:            toolInfo.ID,
		UUID:          toolInfo.UUID,
		CreatedAt:     toolInfo.CreatedAt.String(),
		UpdatedAt:     toolInfo.UpdatedAt.String(),
		McpServerId:   toolInfo.McpServerId,
		McpServerUUID: toolInfo.McpServerUUID,
		Name:          toolInfo.Name,
		Description:   toolInfo.Description,
		McpServerType: toolInfo.McpServerType,
		Method:        toolInfo.Method,
		Endpoint:      "",
		Headers:       "",
		Args:          nil,
		//Security:       ,
		IsAuth:       0,
		AuthMode:     "",
		IsShow:       0,
		SerialNumber: "",
		IsRepeat:     0,
	}
	return
}
