package biz

import (
	"flow-bridge-mcp/api"
	"flow-bridge-mcp/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	Login(ctx *gin.Context) (*api.LoginReplay, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *logger.Logger
}

func NewUserUseCase(repo UserRepo, log *logger.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log}
}

func (u *UserUseCase) UserLogin(ctx *gin.Context) (*api.LoginReplay, error) {
	u.log.Info("biz 示例")
	return u.repo.Login(ctx)
}
