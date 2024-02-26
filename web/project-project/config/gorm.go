package config

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"test.com/project-project/internal/database/gorms"
)

var _db *gorm.DB

func (c *Config) ReConnMysql() {
	// 打印数据库配置
	// fmt.Printf("=== DB Config ===\n")
	// fmt.Printf("Separation: %v\n", c.DbConfig.Separation)
	// fmt.Printf("Master: %s:%s@%s:%d/%s\n",
	// 	c.DbConfig.Master.Username,
	// 	c.DbConfig.Master.Password,
	// 	c.DbConfig.Master.Host,
	// 	c.DbConfig.Master.Port,
	// 	c.DbConfig.Master.Db)
	// for i, s := range c.DbConfig.Slave {
	// 	fmt.Printf("Slave%d: %s:%s@%s:%d/%s\n",
	// 		i+1, s.Username, s.Password, s.Host, s.Port, s.Db)
	// }
	// fmt.Printf("==================\n")
	if c.DbConfig.Separation {
		//读写分离配置
		username := c.DbConfig.Master.Username //账号
		password := c.DbConfig.Master.Password //密码
		host := c.DbConfig.Master.Host         //数据库地址，可以是Ip或者域名
		port := c.DbConfig.Master.Port         //数据库端口
		Dbname := c.DbConfig.Master.Db         //数据库名
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
		var err error
		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			zap.L().Error("连接数据库失败 ", zap.Error(err))
			return
		}
		replicas := []gorm.Dialector{}
		for _, v := range c.DbConfig.Slave {
			username := v.Username //账号
			password := v.Password //密码
			host := v.Host         //数据库地址，可以是Ip或者域名
			port := v.Port         //数据库端口
			Dbname := v.Db         //数据库名
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
			cfg := mysql.Config{
				DSN: dsn,
			}
			replicas = append(replicas, mysql.New(cfg))
		}
		err = _db.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: dsn,
			})},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(10).
			SetMaxOpenConns(200))
		if err != nil {
			zap.L().Error("Use slave err ", zap.Error(err))
			return
		}
	} else {
		//配置MySQL连接参数
		username := c.MysqlConfig.Username //账号
		password := c.MysqlConfig.Password //密码
		host := c.MysqlConfig.Host         //数据库地址，可以是Ip或者域名
		port := c.MysqlConfig.Port         //数据库端口
		Dbname := c.MysqlConfig.Db         //数据库名
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
		var err error
		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("连接数据库失败, error=" + err.Error())
		}
	}
	gorms.SetDB(_db)
}
