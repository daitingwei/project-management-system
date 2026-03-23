package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	srv "test.com/project-common"
	"test.com/project-common/kk"
	"test.com/project-project/config"
	"test.com/project-project/router"
	"test.com/project-project/tracing"
)

func main() {
	// TODO: 微服务拆分计划
	// 当前 project-project 过于臃肿，包含项目/任务/菜单/权限/账户/部门等多个模块
	// 建议拆分为独立微服务：
	//   - project-service: 项目管理
	//   - task-service: 任务管理
	//   - menu-service: 菜单管理
	//   - auth-service: 权限管理
	//   - account-service: 账户管理
	//   - department-service: 部门管理
	// 当前架构更像是单体架构硬拆成3个服务，微服务划分不合理

	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider(config.C.JaegerConfig.Endpoints)
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// 路由
	router.InitRouter(r)
	// 初始化rpc调用
	router.InitUserRpc()
	// grpc服务注册
	gc := router.RegisterGrpc()
	// grpc服务注册到etcd
	router.RegisterEtcdServer()
	//初始化kafka
	kafkaWriter := config.InitKafkaWriter()

	// 初始化 Kafka 消费者，将日志转发到 Logstash
	kafkaConsumer := kk.NewKafkaConsumer(
		[]string{config.C.KafkaConfig.Broker},
		"msproject_log",
		"msproject_consumer_group",
		"localhost",
		50000,
	)
	ctx, cancel := context.WithCancel(context.Background())
	if err := kafkaConsumer.Start(ctx); err != nil {
		log.Printf("Failed to start kafka consumer: %v", err)
	} else {
		log.Printf("Kafka consumer started successfully")
	}

	// 初始化缓存读取器
	cacheReader := config.NewCacheReader()
	go cacheReader.DeleteCache()

	stop := func() {
		log.Printf("Shutting down services...")
		cancel()
		kafkaConsumer.Stop()
		gc.Stop()
		kafkaWriter()
		cacheReader.R.Close()
		log.Printf("All services stopped")
	}

	// 启动优雅关闭处理
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Printf("Received shutdown signal")
		stop()
		os.Exit(0)
	}()

	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
