package rpc

import (
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"test.com/project-api/config"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	"test.com/project-grpc/account"
	"test.com/project-grpc/auth"
	"test.com/project-grpc/department"
	"test.com/project-grpc/menu"
	"test.com/project-grpc/project"
	// 修改说明：使用带版本号的包名，以适配新的gRPC服务定义
	// 背景：project-grpc模块重构后，task服务包名从task改为task_service_v1
	task_service_v1 "test.com/project-grpc/task"
)

var ProjectServiceClient project.ProjectServiceClient
// TaskServiceClient 修改说明：变量类型从task.TaskServiceClient改为task_service_v1.TaskServiceClient
var TaskServiceClient task_service_v1.TaskServiceClient
var AccountServiceClient account.AccountServiceClient
var DepartmentServiceClient department.DepartmentServiceClient
var AuthServiceClient auth.AuthServiceClient
var MenuServiceClient menu.MenuServiceClient

func InitRpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.Dial(
		"etcd:///project",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ProjectServiceClient = project.NewProjectServiceClient(conn)
	TaskServiceClient = task_service_v1.NewTaskServiceClient(conn)
	AccountServiceClient = account.NewAccountServiceClient(conn)
	DepartmentServiceClient = department.NewDepartmentServiceClient(conn)
	AuthServiceClient = auth.NewAuthServiceClient(conn)
	MenuServiceClient = menu.NewMenuServiceClient(conn)

}
