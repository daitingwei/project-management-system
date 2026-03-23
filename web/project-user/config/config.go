package config

import (
	"bytes"
	"github.com/go-redis/redis/v8"
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
	MysqlConfig  *MysqlConfig
	JwtConfig    *JwtConfig
	JaegerConfig *JaegerConfig
}

type JaegerConfig struct {
	Endpoints string
}
type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name     string
	Addr    string
	Version  string
	Weight   int64
	EtcdAddr string
}

type EtcdConfig struct {
	Addrs []string
}

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Db       string
}

type JwtConfig struct {
	AccessExp     int64
	RefreshExp    int64
	AccessSecret  string
	RefreshSecret string
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
	c.ReadGrpcConfig()
	c.ReadEtcdConfig()
	c.InitMysqlConfig()
	c.InitJwtConfig()
	c.InitJaegerConfig()
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

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	gc.EtcdAddr = c.viper.GetString("grpc.etcdAddr")
	c.GC = gc
}

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
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
func (c *Config) InitMysqlConfig() {
	mc := &MysqlConfig{}
	mc.Username = c.viper.GetString("mysql.username")
	mc.Password = c.viper.GetString("mysql.password")
	mc.Host = c.viper.GetString("mysql.host")
	mc.Port = c.viper.GetInt("mysql.port")
	mc.Db = c.viper.GetString("mysql.db")
	c.MysqlConfig = mc
}
func (c *Config) InitJwtConfig() {
	mc := &JwtConfig{}
	mc.AccessSecret = c.viper.GetString("jwt.accessSecret")
	mc.AccessExp = c.viper.GetInt64("jwt.accessExp")
	mc.RefreshExp = c.viper.GetInt64("jwt.refreshExp")
	mc.RefreshSecret = c.viper.GetString("jwt.refreshSecret")
	c.JwtConfig = mc
}

func (c *Config) InitJaegerConfig() {
	mc := &JaegerConfig{
		Endpoints: c.viper.GetString("jaeger.endpoints"),
	}
	c.JaegerConfig = mc
}
