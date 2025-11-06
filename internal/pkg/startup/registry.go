package startup

import (
	"flow-bridge-mcp/internal/conf"
	"flow-bridge-mcp/pkg/logger"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sync"
)

type Registry struct {
	conf *conf.Conf
	log  *logger.Logger
	sync.Mutex
	namingClient naming_client.INamingClient
}

func NewRegistry(conf *conf.Conf, log *logger.Logger) *Registry {

	// 创建动态配置客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				NamespaceId:         conf.Conf.GetString("registry.nacos.namespace"),
				TimeoutMs:           uint64(conf.Conf.GetInt("registry.nacos.timeout")),      //超时时间
				BeatInterval:        int64(conf.Conf.GetInt("registry.nacos.beatInterval")),  //心跳时间
				AppName:             conf.Conf.GetString("server.name"),                      //程序名称
				NotLoadCacheAtStart: conf.Conf.GetBool("registry.nacos.notLoadCacheAtStart"), //启动时不加载缓存
				//LogDir:              conf.Conf.GetString("registry.nacos.logDir"),            //日志
				//CacheDir:            conf.Conf.GetString("registry.nacos.cacheDir"),          //持久化Nacos服务信息的目录,                                     //持久化Nacos服务信息的目录
				LogLevel: conf.Conf.GetString("registry.nacos.logLevel"), //可选值：debug, info, warn, error
				Username: conf.Conf.GetString("registry.nacos.username"),
				Password: conf.Conf.GetString("registry.nacos.password"),
				//Endpoint:            "", // 非直连模式 ServerConfigs不配置，通过Endpoint动态获取
				//TLSCfg:              constant.TLSConfig{},
				//UpdateThreadNum:      0, //更新Nacos服务信息的goroutine数量默认值：20
				//ContextPath:      "", //Nacos服务器的上下文路径 默认值：/nacos Nacos服务的访问路径前缀
			},
			ServerConfigs: []constant.ServerConfig{
				{
					IpAddr:      conf.Conf.GetString("registry.nacos.ipAddr"), //nacos的服务地址
					ContextPath: conf.Conf.GetString("registry.nacos.contextPath"),
					Port:        uint64(conf.Conf.GetInt("registry.nacos.port")), //nacos的端口号
					Scheme:      conf.Conf.GetString("registry.nacos.scheme"),    //nacos的协议
				},
			},
		},
	)
	if err != nil {
		//panic(fmt.Errorf("NewConfigClient error: %v", err))
		log.Errorf("NewConfigClient error: %v", err)
	}
	return &Registry{
		conf:         conf,
		log:          log,
		namingClient: namingClient,
	}
}

func (r *Registry) register() {
	r.Lock()
	ok, err := r.namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          r.conf.Conf.GetString("server.http.ip"),
		Port:        uint64(r.conf.Conf.GetInt("server.http.port")),
		ServiceName: r.conf.Conf.GetString("server.name"),
		Weight:      r.conf.Conf.GetFloat64("registry.nacos.weight"),
		Enable:      r.conf.Conf.GetBool("registry.nacos.enable"),
		Healthy:     r.conf.Conf.GetBool("registry.nacos.healthy"),
		ClusterName: r.conf.Conf.GetString("registry.nacos.clusterName"), // 集群名称
		GroupName:   r.conf.Conf.GetString("registry.nacos.groupName"),   // 分组名称
		Ephemeral:   r.conf.Conf.GetBool("registry.nacos.ephemeral"),
		Metadata: map[string]string{ // 元数据
			"version": r.conf.Conf.GetString("server.http.version"),
			//"env": "production",
		},
	})
	if err != nil {
		panic(fmt.Errorf("RegisterInstance error: %v", err))
		return
	}
	if !ok {
		panic(fmt.Errorf("RegisterInstance failed"))
		return
	}
	defer r.Unlock()

}
func (r *Registry) Deregister() {
	r.Lock()
	defer r.Unlock()
	success, err := r.namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          r.conf.Conf.GetString("server.http.ip"),
		Port:        uint64(r.conf.Conf.GetInt("server.http.port")),
		ServiceName: r.conf.Conf.GetString("server.name"),
		Ephemeral:   true,
		//Cluster:     "DEFAULT", // 默认值DEFAULT
		//GroupName:   "默认值DEFAULT_GROUP",   // 默认值DEFAULT_GROUP
	})
	if err != nil {
		r.log.Errorf("DeregisterInstance error: %v", err)
	}
	if !success {
		r.log.Error("DeregisterInstance failed")
	}
}
