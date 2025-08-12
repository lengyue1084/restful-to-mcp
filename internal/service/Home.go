package service

import (
	"flow-bridge-mcp/internal/conf"
	"github.com/gin-gonic/gin"
)

type HomeService struct {
	log *conf.Logger
}

func NewHome(log *conf.Logger) *HomeService {
	return &HomeService{
		log: log,
	}
}

func (h *HomeService) Index(ctx *gin.Context) {
	h.log.Info("welcome Home page")
}
