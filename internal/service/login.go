package service

import (
	"flow-bridge-mcp/api/user"
	"flow-bridge-mcp/internal/conf"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginService struct {
	log *conf.Logger
}

func NewLoginService(log *conf.Logger) *LoginService {
	return &LoginService{
		log: log,
	}
}

func (l *LoginService) Login(ctx *gin.Context) {
	var json user.UserLoginReq
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//l.log.Info(json.Username + json.Password)
	l.log.WithContext(ctx).InfoF("username: %s, password: %s", json.Username, json.Password)

	ctx.JSON(http.StatusOK, json)
}
