package biz

import (
	"context"
	"errors"
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/internal/mcp/transformer"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
	"fmt"
)

type OpenapiUseCase struct {
	log           *logger.Logger
	transformer   transformer.Transformer
	Tx            Transaction
	mcpServerRepo McpServerRepo
	mfUc          *McpFileUserCase
	mcpToolsRepo  McpToolsRepo
	conf          *conf.Conf
}

func NewOpenapiUserCase(
	transformer transformer.Transformer,
	log *logger.Logger,
	tx Transaction,
	mcpServerRepo McpServerRepo,
	mfUc *McpFileUserCase,
	mcpToolsRepo McpToolsRepo,
	conf *conf.Conf,
) *OpenapiUseCase {
	return &OpenapiUseCase{
		log:           log,
		transformer:   transformer,
		Tx:            tx,
		mcpServerRepo: mcpServerRepo,
		mfUc:          mfUc,
		mcpToolsRepo:  mcpToolsRepo,
		conf:          conf,
	}
}

func (o *OpenapiUseCase) Create(ctx context.Context, req *api.ServerInfoRequest) (resp *api.ServerInfoResponse, err error) {

	var decodeString []byte
	uuidStr := req.UUID
	if uuidStr == "" {
		o.log.ErrorWithContext(ctx, "UUID不能为空")
		return nil, fmt.Errorf("UUID不能为空")
	}
	ctx = context.WithValue(ctx, "uuid", uuidStr)
	cleanedContent := tool.CleanBase64String(req.FileContent)
	if err := tool.ValidateBase64String(cleanedContent); err != nil {
		o.log.ErrorWithContext(ctx, "Base64字符串验证失败: %+v", err)
		return nil, err
	}
	decodeString, err = tool.TryMultipleBase64Decodings(cleanedContent)
	if err != nil {
		o.log.ErrorWithContext(ctx, "Base64解析失败，err:%+v", err)
		return nil, fmt.Errorf("base64 decode failed: %w", err)
	}
	contentMd5Str := tool.MD5(req.FileContent)

	//创建McpFile记录
	mcpFileInfo, err := o.mfUc.GetMcpFileInfoByMd5(ctx, contentMd5Str)
	if err != nil {
		o.log.ErrorWithContext(ctx, "获取McpFile失败，err:%+v", err)
		return nil, errors.New("获取McpFile失败")
	}
	if mcpFileInfo.ID == 0 {
		fileName := tool.FileNameByUUid()
		// 创建McpFile文件，写入记录
		err := o.mfUc.CreateMcpFile(ctx, fileName, req.Name, req.Suffix, contentMd5Str, req.Description)
		if err != nil {
			o.log.ErrorWithContext(ctx, "创建McpFile失败，err:%+v", err)
			return nil, errors.New("创建McpFile失败")
		}

		path := o.conf.Conf.GetString("file.path")
		// todo 写入文件中,md5表明是否变化
		_, err = tool.WriteFile(path, fileName, decodeString)
		if err != nil {
			o.log.ErrorWithContext(ctx, "写入文件失败，err:%+v", err)
			return nil, err
		}
	} else {
		err = o.mfUc.UpdateMcpFile(ctx, req.Name, req.Suffix, contentMd5Str, req.Description)
		if err != nil {
			o.log.ErrorWithContext(ctx, "更新McpFile失败，err:%+v", err)
			return nil, errors.New("更新McpFile失败")
		}

	}

	mcpConfig, err := o.transformer.Convert(ctx, decodeString)
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

	return
}

func (o *OpenapiUseCase) ValidateBase64String(ctx context.Context, fileContent string) bool {
	// 清理Base64字符串
	cleanedContent := tool.CleanBase64String(fileContent)
	o.log.ErrorWithContext(ctx, "UUID不能为空")
	// 验证Base64字符串
	if err := tool.ValidateBase64String(cleanedContent); err != nil {
		o.log.ErrorWithContext(ctx, "Base64字符串验证失败: %+v", err)
		return false
	}
	return true

}
