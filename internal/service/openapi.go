package service

import (
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/pkg/response"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

type OpenapiService struct {
	uc  *biz.OpenapiUseCase
	log *logger.Logger
}

func NewOpenapiService(uc *biz.OpenapiUseCase, log *logger.Logger) *OpenapiService {
	return &OpenapiService{
		uc:  uc,
		log: log,
	}
}

func (o *OpenapiService) Upload(c *gin.Context) {
	var req *api.OpenapiUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "参数错误", nil)
		return
	}
	resp, err := o.uc.Create(c, req)
	if err != nil {
		o.log.ErrorWithContext(c, "创建失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("创建失败，err:%s", err), nil)
		return
	}
	response.Success(c, "创建成功", resp)
}
