package main

import (
	"flag"
	"flow-bridge-mcp/internal/conf"
	"fmt"
	"go.uber.org/zap"
	"os"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	flagConf string
	id, _    = os.Hostname()
	//logger   *zap.Logger
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs", "config path, eg: -conf config.yaml")
}
func main() {
	flag.Parse()
	config := conf.NewConf(flagConf)
	logger := conf.NewZapLogger(config)

	app, cleanup, err := initApp(config, logger)
	if err != nil {
		logger.Error("init app failed", zap.Error(err))
		panic(err)
	}
	defer cleanup()
	logger.Info("start http server")
	// 启动服务
	if err := app.Run(
		fmt.Sprintf("%s:%s", config.Conf.GetString("server.http.addr"),
			config.Conf.GetString("server.http.port")),
	); err != nil {
		logger.Error("server run failed", zap.Error(err))
		panic(err)
	}
}
