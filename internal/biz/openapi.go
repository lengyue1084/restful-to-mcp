package biz

import (
	"context"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/transformer"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
)

type OpenapiRepo interface {
	Create(ctx context.Context, serverInfo *model.McpServer) (err error)
}
type OpenapiUseCase struct {
	OpenapiRepo OpenapiRepo
	log         *logger.Logger
	Transformer transformer.Transformer
}

func NewOpenapiUserCase(repo OpenapiRepo, Transformer transformer.Transformer, log *logger.Logger) *OpenapiUseCase {
	return &OpenapiUseCase{
		OpenapiRepo: repo,
		log:         log,
		Transformer: Transformer,
	}
}

func (o *OpenapiUseCase) Create(ctx context.Context, req *api.ServerInfoRequest) (resp *api.ServerInfoResponse, err error) {
	// 清理Base64字符串
	cleanedContent := tool.CleanBase64String(req.FileContent)
	uuidStr := req.UUID
	if uuidStr == "" {
		return nil, fmt.Errorf("UUID不能为空")
	}
	//
	ctx = context.WithValue(ctx, "uuid", uuidStr)
	o.log.ErrorWithContext(ctx, "UUID不能为空")
	// 验证Base64字符串
	if err := tool.ValidateBase64String(cleanedContent); err != nil {
		o.log.ErrorWithContext(ctx, "Base64字符串验证失败: %+v", err)
		return nil, err
	}

	// 尝试多种解码方式
	decodeString, err := tool.TryMultipleBase64Decodings(cleanedContent)
	if err != nil {
		o.log.ErrorWithContext(ctx, "Base64解析失败，err:%+v", err)
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	//req.FileContent = string(decodeString)

	//converter := openapi.NewConverter()
	mcpConfig, err := o.Transformer.Convert(ctx, decodeString)
	if err != nil {
		o.log.ErrorWithContext(ctx, "数据转换错误，err:%+v", err)
		return nil, err
	}
	fmt.Printf("%+v", mcpConfig)

	var serverInfo *model.McpServer
	if err := tool.Copy(&serverInfo, req); err != nil {
		o.log.ErrorWithContext(ctx, "数据转换错误，err:%+v", err)
		return nil, err
	}

	err = o.OpenapiRepo.Create(ctx, serverInfo)
	return
}
