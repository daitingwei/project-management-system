package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

type ConfigClient struct {
	confClient config_client.IConfigClient
	group      string
	dataId     string
}

func NewConfigClient(bootConf *BootConfig) *ConfigClient {
	clientConfig := constant.ClientConfig{
		NamespaceId:         bootConf.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      bootConf.NacosConfig.IpAddr,
			ContextPath: bootConf.NacosConfig.ContextPath,
			Port:        uint64(bootConf.NacosConfig.Port),
			Scheme:      bootConf.NacosConfig.Scheme,
		},
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	nc := &ConfigClient{
		confClient: configClient,
		group:      bootConf.NacosConfig.Group,
		dataId:     bootConf.NacosConfig.DataId,
	}
	return nc
}

func (c *ConfigClient) GetConfig() (string, error) {
	return c.confClient.GetConfig(vo.ConfigParam{
		DataId: c.dataId,
		Group:  c.group,
	})
}

func (c *ConfigClient) ListenConfig(onChange func(namespace, group, dataId, data string)) error {
	return c.confClient.ListenConfig(vo.ConfigParam{
		DataId: c.dataId,
		Group:  c.group,
		OnChange: onChange,
	})
}
