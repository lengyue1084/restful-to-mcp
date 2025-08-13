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
	flagConf    string
	hostname, _ = os.Hostname()
	logger      *conf.Logger
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs", "config path, eg: -conf config.yaml")
}
func main() {
	flag.Parse()
	config := conf.NewConf(flagConf)

	app, cleanup, err := initApp(config)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	// 启动服务
	err = app.Run(
		fmt.Sprintf("%s:%s", config.Conf.GetString("server.http.addr"),
			config.Conf.GetString("server.http.port")),
	)
	if err != nil {
		logger.Error("server run failed", zap.Error(err))
		panic(err)
	}
}
