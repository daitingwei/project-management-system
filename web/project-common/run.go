package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动HTTP服务
// 参数说明：
//   - r: gin.Engine 实例
//   - srvName: 服务名称
//   - addr: 监听地址，格式如 ":8080"
//   - stop: 服务停止时的回调函数
//
// TODO: 若需启用HTTPS，可改为调用 RunWithTLS() 方法
// 示例：
//   certFile := "./cert/server.crt"  // 证书文件路径
//   keyFile := "./cert/server.key"  // 私钥文件路径
//   RunWithTLS(r, srvName, addr, certFile, keyFile, stop)
func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	//保证下面的优雅启停
	go func() {
		log.Printf("%s running in %s \n", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	//SIGINT 用户发送INTR字符(Ctrl+C)触发 kill -2
	//SIGTERM 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down menu %s... \n", srvName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %v", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout....")
	}
	log.Printf("%s stop success... \n", srvName)
}
