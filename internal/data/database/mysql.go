package database

import (
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/internal/data/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysqlClient 创建 Mysql 数据库连接
func NewMysqlClient(config *conf.Conf, log *conf.Logger) (*gorm.DB, func(), error) {
	dsn := config.Conf.GetString("data.database.source")
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	dbOpen, err := gorm.Open(mysql.New(mysqlConfig))
	if err != nil {
		log.Errorf("open mysql error:%+v", zap.Error(err))
		panic(err)
	}
	db, err := dbOpen.DB()
	if err != nil {
		log.Errorf("open mysql error:%+v", zap.Error(err))
		panic(err)
	}
	db.SetMaxIdleConns(config.Conf.GetInt("data.database.max_idle_conn"))
	db.SetMaxOpenConns(config.Conf.GetInt("data.database.max_open_conn"))
	// 迁移 schema
	if err := dbOpen.AutoMigrate(model.User{}); err != nil {
		panic(err)
	}
	cleanup := func() {
		if err := db.Close(); err != nil {
			//log.Println("close mysql err:", err)
			log.Errorf("close mysql err:%+v", zap.Error(err))
		}
		//log.Println("close mysql success")
		log.Info("close mysql success")
	}
	//log.Println("open mysql success")
	log.Info("open mysql success")
	return dbOpen, cleanup, nil
}
