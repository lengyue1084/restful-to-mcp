package database

import (
	"context"
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"flow-bridge-mcp/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPgClient 创建 PostgreSQL 数据库连接
func NewPgClient(config *conf.Conf, log2 *logger.GormLogger, log *logger.Logger) (*gorm.DB, func(), error) {
	dsn := config.Conf.GetString("data.database.pg_source")

	// PostgresSQL 配置
	pgConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false, // 禁用简单协议，使用扩展查询协议
	}

	// GORM 配置 - 启用日志
	gormConfig := &gorm.Config{}

	// 根据配置决定是否启用 SQL 日志
	if config.Conf.GetBool("data.database.log_sql") {
		gormConfig.Logger = log2 // 自定义日志记录器
	}

	// 连接数据库
	dbOpen, err := gorm.Open(postgres.New(pgConfig), gormConfig)
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
	if err := dbOpen.AutoMigrate(
		model.McpServer{},
		model.McpTools{},
		model.McpFile{},
	); err != nil {
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
