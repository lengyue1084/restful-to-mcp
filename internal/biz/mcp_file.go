package biz

import (
	"context"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
	"flow-bridge-mcp/pkg/tool"
)

type McpFileRepo interface {
	Create(ctx context.Context, McpFileInfo *model.McpFile) (err error)
	GetMcpFileInfoByMd5(ctx context.Context, md5 string) (mcpFileInfo *model.McpFile, err error)
	UpdateMcpFileById(ctx context.Context, id int64, serverInfo *model.McpFile) (err error)
}
type McpFileUserCase struct {
	mfRepo McpFileRepo
	log    *logger.Logger
}

func NewMcpFileUserCase(mfRepo McpFileRepo, log *logger.Logger) *McpFileUserCase {
	return &McpFileUserCase{
		mfRepo: mfRepo,
		log:    log,
	}
}
func (m *McpFileUserCase) GetMcpFileInfoByMd5(ctx context.Context, fileMd5 string) (McpFileInfo *model.McpFile, err error) {
	mcpFileInfo, err := m.mfRepo.GetMcpFileInfoByMd5(ctx, fileMd5)
	if err != nil {
		m.log.ErrorWithContext(ctx, "获取McpFile失败，err:%+v", err)
		return nil, err
	}
	return mcpFileInfo, nil

}

func (m *McpFileUserCase) CreateMcpFile(ctx context.Context, fileName, sourceName, suffix, contentMd5Str, description string) (err error) {

	if fileName == "" {
		fileName = tool.FileNameByUUid()
	}
	mcpFile := model.McpFile{
		Name:        fileName,
		SourceName:  sourceName,
		Md5:         contentMd5Str,
		Description: description,
		Suffix:      suffix,
	}
	err = m.mfRepo.Create(ctx, &mcpFile)
	if err != nil {
		m.log.ErrorWithContext(ctx, "创建McpFile失败，err:%+v", err)
		return err
	}
	return
}

func (m *McpFileUserCase) UpdateMcpFile(ctx context.Context, id int64, fileName, suffix, contentMd5Str, description string) (err error) {

	mcpFile := &model.McpFile{
		SourceName:  fileName,
		Md5:         contentMd5Str,
		Description: description,
		Suffix:      suffix,
	}
	err = m.mfRepo.UpdateMcpFileById(ctx, id, mcpFile)
	if err != nil {
		m.log.ErrorWithContext(ctx, "创建McpFile失败，err:%+v", err)
		return err
	}
	return
}
