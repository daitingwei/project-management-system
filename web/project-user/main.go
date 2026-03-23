package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	srv "test.com/project-common"
	"test.com/project-user/config"
	"test.com/project-user/internal/rpc"
	"test.com/project-user/router"
	"test.com/project-user/tracing"
)

func main() {
	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider(config.C.JaegerConfig.Endpoints)
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// 初始化RPC客户端
	// 修改说明：新增此调用，用于初始化project-project服务的RPC客户端
	// 背景：用户注册时需要调用project-project的AccountService创建账户记录
	// 实现：通过rpc.InitRpcProjectClient()建立与project-project服务的gRPC连接
	rpc.InitRpcProjectClient()
	// 路由
	router.InitRouter(r)
	// grpc服务注册
	gc := router.RegisterGrpc()
	// grpc服务注册到etcd
	router.RegisterEtcdServer()

	stop := func() {
		gc.Stop()
	}
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
