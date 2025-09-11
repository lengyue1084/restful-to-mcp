package conf

import (
	"flow-bridge-mcp/pkg/tool"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

type Conf struct {
	Conf *viper.Viper
}

func NewConf(confFile string) *Conf {
	if ok, err := tool.FileExists(confFile); err != nil || !ok {
		if ok, err = tool.FileExists(fmt.Sprintf("%s%s", confFile, "/config.yaml")); err != nil || !ok {
			panic("请检查配置文件" + confFile + "是否存在")
		}
	}

	configViper := viper.New()
	configViper.SetConfigName("config") // name of config file (without extension)
	configViper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	configViper.AddConfigPath(confFile) // path to look for the config file in
	//vip.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	configViper.AddConfigPath(".") // optionally look for config in the working directory

	// 设置环境变量前缀
	configViper.SetEnvPrefix("MCP") // 环境变量前缀为MCP_
	configViper.AutomaticEnv()

	// 替换环境变量中的分隔符
	configViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//如：nacos.server.port  → MCP_NACOS_SERVER_PORT

	err := configViper.ReadInConfig() // Find and read the config file
	if err != nil {                   // Handle response reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	configViper.WatchConfig()
	configViper.OnConfigChange(func(e fsnotify.Event) {

		fmt.Println("Config file changed:", e.Name)
		//todo
		// 记录变更日志
		//if logger != nil {
		//	logger.Info("Config file changed",
		//		zap.String("file", e.Name),
		//		zap.String("operation", e.Op.String()))
		//}

		// viper 会自动更新配置值，但你可能需要：
		// 1. 通知相关组件配置已变更
		// 2. 重新初始化依赖配置的服务
		// 3. 更新全局变量等

		// 示例：通知配置变更
		//notifyConfigChange()
	})
	return &Conf{
		Conf: configViper,
	}
}
