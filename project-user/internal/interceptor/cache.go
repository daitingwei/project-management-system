package interceptor

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"test.com/project-common/encrypts"
	"test.com/project-grpc/user/login"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/repo"
)

// TODO[缓存优化]:
// 1. 缓存拦截器 Cache() 方法已实现但未在 router.go 中启用，需取消注释
// 2. 需要修复的缓存问题同 project-project 服务
// 3. 可扩展缓存方法：FindByEmail、ValidateCode 等查询接口
//
// 适合缓存的接口：
// - MyOrgList - 用户组织列表（用户级别缓存）
// - FindMemInfoById - 成员信息（变化少）
//
// 不适合缓存：
// - 登录验证（实时性要求高）
// - 验证码验证（一次性）

type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

type CacheRespOption struct {
	path   string
	typ    any
	expire time.Duration
}

func New() *CacheInterceptor {
	cacheMap := make(map[string]any)
	cacheMap["/login.service.v1.LoginService/MyOrgList"] = &login.OrgListResponse{}
	cacheMap["/login.service.v1.LoginService/FindMemInfoById"] = &login.MemberMessage{}
	return &CacheInterceptor{cache: dao.Rc, cacheMap: cacheMap}
}

// Cache gRPC缓存拦截器（返回grpc.ServerOption类型）
// 修改说明：新增此方法，返回grpc.ServerOption类型以适配grpc.NewServer()的拦截器参数
// 背景：原Cache()方法返回的是grpc.UnaryServerInterceptor类型，无法直接在grpc.NewServer()中使用
// 修改点：将方法返回值改为grpc.ServerOption，内部实现与CacheInterceptor()相同
// func (c *CacheInterceptor) Cache() grpc.ServerOption {
// 	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		respType := c.cacheMap[info.FullMethod]
// 		if respType == nil {
// 			return handler(ctx, req)
// 		}
// 		con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 		defer cancel()
// 		marshal, _ := json.Marshal(req)
// 		cacheKey := encrypts.Md5(string(marshal))
// 		respJson, _ := c.cache.Get(con, info.FullMethod+"::"+cacheKey)
// 		if respJson != "" {
// 			json.Unmarshal([]byte(respJson), &respType)
// 			zap.L().Info(info.FullMethod + " 走了缓存")
// 			return respType, nil
// 		}
// 		resp, err = handler(ctx, req)
// 		bytes, _ := json.Marshal(resp)
// 		c.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
// 		zap.L().Info(info.FullMethod + " 放入缓存")
// 		return
// 	})
// }

// CacheInterceptor gRPC缓存拦截器（返回拦截器函数类型）
// 修改说明：修复返回值问题
// 背景：原方法名Cache()与新添加的Cache()方法冲突，且返回值类型不匹配
// 修改点：重命名为CacheInterceptor()，保持原有功能
func (c *CacheInterceptor) CacheInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))
		respJson, _ := c.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + " 走了缓存")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		c.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		zap.L().Info(info.FullMethod + " 放入缓存")
		return
	}
}
