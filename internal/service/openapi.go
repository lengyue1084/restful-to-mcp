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
	oaUc *biz.OpenapiUseCase
	log  *logger.Logger
}

func NewOpenapiService(oaUc *biz.OpenapiUseCase, log *logger.Logger) *OpenapiService {
	return &OpenapiService{
		oaUc: oaUc,
		log:  log,
	}
}

func (o *OpenapiService) Upload(c *gin.Context) {
	var req *api.OpenapiUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	if req.Suffix != "json" && req.Suffix != "yaml" && req.Suffix != "yml" {
		o.log.ErrorWithContext(c, "文件格式错误")
		response.Error(c, "文件格式错误,只支持json、yaml格式", nil)
		return
	}
	resp, err := o.oaUc.Create(c, req)
	if err != nil {
		o.log.ErrorWithContext(c, "创建失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("创建失败，err:%+v", err), nil)
		return
	}
	response.Success(c, "创建成功", resp)
}

func (o *OpenapiService) UpdateForAuth(c *gin.Context) {
	var req *api.OpenapiUpdateForAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, fmt.Sprintf("参数错误,err:%+v", err), nil)
		return
	}
	resp, err := o.oaUc.UpdateForAuth(c, req)
	if err != nil {
		o.log.ErrorWithContext(c, "更新失败,err:%+v", err)
		response.Error(c, fmt.Sprintf("更新失败，err:%+v", err), nil)
		return
	}
	response.Success(c, "更新成功", resp)

}
