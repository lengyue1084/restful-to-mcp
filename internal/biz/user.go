package biz

import (
	"flow-bridge-mcp/api/user"
	"flow-bridge-mcp/internal/conf"
	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	Login(ctx *gin.Context) (*user.LoginReplay, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *conf.Logger
}

func NewUserUseCase(repo UserRepo, log *conf.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log}
}

func (u *UserUseCase) UserLogin(ctx *gin.Context) (*user.LoginReplay, error) {
	u.log.Info("biz 示例")
	return u.repo.Login(ctx)
}
