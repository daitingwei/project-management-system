package config

import (
	"bytes"
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/project-common/logs"
	"test.com/project-common/nacos"
)

var C = InitConfig()

type Config struct {
	viper        *viper.Viper
	SC           *ServerConfig
	GC           *GrpcConfig
	EtcdConfig   *EtcdConfig
	JaegerConfig *JaegerConfig
	MinioConfig  *MinioConfig
}

type JaegerConfig struct {
	Endpoints string
}

type MinioConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Bucket     string
	UseSSL     bool
}
type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name string
	Addr string
}

type EtcdConfig struct {
	Addrs []string
}

func InitConfig() *Config {
	// InitConfig 初始化配置
	// 修改说明：使用新的nacos客户端封装，支持配置热更新
	// 背景：
	// 1. 使用project-common模块中重构的nacos封装，提供更简洁的API
	// 2. 支持配置热更新：nacos配置变更时自动重新加载
	// 3. 优先从nacos读取配置，读取失败则回退到本地配置文件
	conf := &Config{viper: viper.New()}
	//先从nacos读取配置，如果读取不到 在本地读取
	bootConf := nacos.InitBootstrap()
	nacosClient := nacos.NewConfigClient(bootConf)
	configYaml, err2 := nacosClient.GetConfig()
	if err2 != nil {
		log.Fatalln(err2)
	}
	err2 = nacosClient.ListenConfig(func(namespace, group, dataId, data string) {
		log.Printf("load nacos config changed %s \n", data)
		err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(data)))
		if err != nil {
			log.Printf("load nacos config changed err : %s \n", err.Error())
		}
		conf.ReLoadAllConfig()
	})
	if err2 != nil {
		log.Fatalln(err2)
	}
	conf.viper.SetConfigType("yaml")
	if configYaml != "" {
		err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(configYaml)))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		workDir, _ := os.Getwd()
		conf.viper.SetConfigName("config")
		conf.viper.AddConfigPath(workDir + "/config")
		err := conf.viper.ReadInConfig()
		if err != nil {
			log.Fatalln(err)
		}
	}
	conf.ReLoadAllConfig()
	return conf
}

func (c *Config) ReLoadAllConfig() {
	c.ReadServerConfig()
	c.InitZapLog()
	c.ReadEtcdConfig()
	c.InitJaegerConfig()
	c.ReadMinioConfig()
}

func (c *Config) InitZapLog() {
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}

func (c *Config) InitJaegerConfig() {
	mc := &JaegerConfig{
		Endpoints: c.viper.GetString("jaeger.endpoints"),
	}
	c.JaegerConfig = mc
}

func (c *Config) ReadMinioConfig() {
	mc := &MinioConfig{
		Endpoint:  c.viper.GetString("minio.endpoint"),
		AccessKey: c.viper.GetString("minio.accessKey"),
		SecretKey: c.viper.GetString("minio.secretKey"),
		Bucket:    c.viper.GetString("minio.bucket"),
		UseSSL:    c.viper.GetBool("minio.useSSL"),
	}
	c.MinioConfig = mc
}
