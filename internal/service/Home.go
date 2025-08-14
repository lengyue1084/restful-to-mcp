package service

import (
	"flow-bridge-mcp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type HomeService struct {
	log *logger.Logger
}

func NewHome(log *logger.Logger) *HomeService {
	return &HomeService{
		log: log,
	}
}

func (h *HomeService) Index(ctx *gin.Context) {
	h.log.Info("welcome Home page")
}
