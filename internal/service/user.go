package service

import (
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/conf"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserService struct {
	uc  *biz.UserUseCase
	log *conf.Logger
}

func NewUserService(uc *biz.UserUseCase, log *conf.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log,
	}
}

// 绑定 JSON
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func (u *UserService) Login(ctx *gin.Context) {
	u.log.Info("service示例")
	u.uc.UserLogin(ctx)
	var json Login
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
