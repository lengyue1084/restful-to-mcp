package biz

import (
	"context"
	"encoding/json"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/pkg/cache"
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
}

type McpToolsUserCase struct {
	mtRepo McpToolsRepo
	log    *logger.Logger
	cache  *cache.MemoryCache
}

func NewMcpToolsUserCase(mtRepo McpToolsRepo, log *logger.Logger, cache *cache.MemoryCache) *McpToolsUserCase {
	return &McpToolsUserCase{
		mtRepo: mtRepo,
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

	var toolsList []*protocol.Tool
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
