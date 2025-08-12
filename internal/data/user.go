package data

import (
	"flow-bridge-mcp/api/user"
	"flow-bridge-mcp/internal/biz"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type userRepo struct {
	UserData *UserData
	log      *conf.Logger
}

func NewUserRepo(d *UserData, log *conf.Logger) biz.UserRepo {
	return &userRepo{
		UserData: d,
		log:      log,
	}
}

func (u *userRepo) Login(ctx *gin.Context) (*user.LoginReplay, error) {
	userMode := &model.User{
		Name: "WQ",
		Age:  18,
	}
	u.UserData.log.Info("data 日志示例")
	_ = u.UserData.db.Create(&userMode)
	var userList []model.User
	_ = u.UserData.db.Table("users").Find(&userList)
	u.UserData.log.Info("gorm示例",
		zap.String("第一条数据的name值", userList[0].Name),
		zap.Uint8("第二条数据的age值", userList[1].Age),
	)
	var list = user.LoginReplay{
		Code:    0,
		Message: "success",
	}
	//err := rdb.Set(ctx, "key", "value", 0).Err()
	u.UserData.redis.Set(ctx, "flow-bridge-gateway:user01", "WQredis", 10*time.Second)
	val := u.UserData.redis.Get(ctx, "flow-bridge-gateway:user01").Val()
	u.UserData.log.Info("redis示例",
		zap.String("flow-bridge-gateway:user01：", val),
	)

	return &list, nil
}
