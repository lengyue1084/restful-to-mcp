package database

import (
	"context"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPgClient 创建 PostgreSQL 数据库连接
func NewPgClient(config *conf.Conf, log *conf.Logger) (*gorm.DB, func(), error) {
	dsn := config.Conf.GetString("data.database.pg_source")

	// PostgresSQL 配置
	pgConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // 禁用简单协议，使用扩展查询协议
	}

	// 连接数据库
	dbOpen, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{})
	if err != nil {
		log.Error("open postgresql error:%+v", zap.Error(err))
		return nil, nil, err
	}

	// 获取底层数据库连接
	db, err := dbOpen.DB()
	if err != nil {
		log.Error("get postgresql db instance error:%+v", zap.Error(err))
		return nil, nil, err
	}

	// 设置连接池参数
	db.SetMaxIdleConns(config.Conf.GetInt("data.database.pg_max_idle_conn"))
	db.SetMaxOpenConns(config.Conf.GetInt("data.database.pg_max_open_conn"))

	// 连接池最大生命周期
	if maxLifetime := config.Conf.GetDuration("data.database.pg_conn_max_lifetime"); maxLifetime > 0 {
		db.SetConnMaxLifetime(maxLifetime)
	}

	// 连接池最大空闲时间
	if maxIdleTime := config.Conf.GetDuration("data.database.pg_conn_max_idle_time"); maxIdleTime > 0 {
		db.SetConnMaxIdleTime(maxIdleTime)
	}

	// 测试连接
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Error("ping postgresql error:%+v", zap.Error(err))
		return nil, nil, err
	}

	// 自动迁移 schema
	if err := dbOpen.AutoMigrate(model.User{}); err != nil {
		log.Error("auto migrate postgresql error:%+v", zap.Error(err))
		return nil, nil, err
	}

	// 清理函数
	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Error("close postgresql error:%+v", zap.Error(err))
		} else {
			log.Info("close postgresql success")
		}
	}
	log.Info("open postgresql success")
	return dbOpen, cleanup, nil
}
