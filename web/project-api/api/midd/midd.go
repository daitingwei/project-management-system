package midd

// TODO[Token验证缓存优化]:
// TokenVerify() 每次都调用 rpc.LoginServiceClient.TokenVerify 进行 gRPC 远程调用
// 这种方式在高并发场景下性能较差，建议添加本地缓存：
//
// 优化方案：
// 1. 使用本地缓存（如 bigcache / go-cache / caffeine）缓存 token 验证结果
// 2. 缓存 key: token 值
// 3. 缓存时间: token 剩余有效期的 80%（或固定 5-10 分钟）
// 4. 缓存命中时直接返回，无需调用 gRPC 服务
// 5. Token 失效/主动退出时主动清除缓存
//
// 适合缓存的原因：
// - Token 在有效期内验证结果基本不变
// - 验签操作频繁，本地缓存可大幅降低延迟
// - 减少 gRPC 服务压力
//
// 注意事项：
// - 需处理 token 失效时的缓存清除
// - 考虑分布式场景下的缓存一致性问题

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
)

// GetIp 获取ip函数
func GetIp(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}

func TokenVerify() func(*gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		// 1.从header中获取token
		token := c.GetHeader("Authorization")
		// 2.调佣user服务进行token认证
		ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelFunc()
		ip := GetIp(c)
		// 先去查询node表 确认不使用登录控制的接口，不做登录认证了

		response, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token, Ip: ip})
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort()
			return
		}
		// 3.处理结果，认证通过 将信息放入gin的上下文 失败返回未登录
		c.Set("memberId", response.Member.Id)
		c.Set("memberName", response.Member.Name)
		c.Set("organizationCode", response.Member.OrganizationCode)
		c.Next()
	}
}
