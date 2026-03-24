package interceptor

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"test.com/project-common/encrypts"
	"test.com/project-grpc/task"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/repo"
)

// TODO[缓存优化]: 需要修复的问题：
// 1. 缓存穿透：空结果也会每次查询数据库，需要添加空值缓存标记
// 2. 缓存一致性：数据更新/删除时没有主动清除缓存，可能返回脏数据
// 3. 建议按方法分组配置不同的缓存过期时间
//
// 适合缓存的方法特征：
// - 读多写少（读:写 > 10:1）
// - 数据变化频率低（配置类、字典类、列表类）
// - 允许几秒-几分钟延迟
// - 请求参数固定或有限
//
// 不适合缓存的方法：
// - GetUserById、GetOrderById 等实时性要求高的
// - Create/Update/Delete 等写入类接口
// CacheInterceptor 除了缓存拦截器 实现日志拦截器 打印参数内容值 请求的时间 等等的
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
	cacheMap["/task.service.v1.TaskService/TaskList"] = &task.TaskListResponse{}
	return &CacheInterceptor{cache: dao.Rc, cacheMap: cacheMap}
}

func (c *CacheInterceptor) Cache() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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
	})
}

func (c *CacheInterceptor) CacheInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		c = New()
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

func ClearTaskListCache() {
	con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 清除 TaskList 缓存
	keys, _ := dao.Rc.Keys(con, "/task.service.v1.TaskService/TaskList::*")
	if len(keys) > 0 {
		dao.Rc.Delete(con, keys)
		zap.L().Info("已删除所有任务缓存", zap.Int("count", len(keys)))
	}

	// TODO: 这里直接删除 FindMemInfoById 缓存不太合理，因为 member 信息变化不频繁且影响范围广
	// 但在以下特定场景需要双删（延迟双删）以保证一致性：
	// 1. 修改任务执行人后 - 任务列表需要显示最新的执行人信息
	// 2. 修改成员账号信息后（如头像、名称）- 需要同步更新所有引用该成员的地方
	// 当前实现：修改 account 或 task 时统一清除所有 FindMemInfoById 缓存
	// 理想实现：根据 memberCode 精确清除特定成员的缓存，避免影响其他无关成员
	userKeys, _ := dao.Rc.Keys(con, "/login.service.v1.LoginService/FindMemInfoById::*")
	if len(userKeys) > 0 {
		dao.Rc.Delete(con, userKeys)
		zap.L().Info("已删除所有member信息缓存(任务相关)", zap.Int("count", len(userKeys)))
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		con2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel2()
		keys2, _ := dao.Rc.Keys(con2, "/task.service.v1.TaskService/TaskList::*")
		if len(keys2) > 0 {
			dao.Rc.Delete(con2, keys2)
			zap.L().Info("延迟删除所有任务缓存", zap.Int("count", len(keys2)))
		}
		userKeys2, _ := dao.Rc.Keys(con2, "/login.service.v1.LoginService/FindMemInfoById::*")
		if len(userKeys2) > 0 {
			dao.Rc.Delete(con2, userKeys2)
			zap.L().Info("延迟删除所有member信息缓存(任务相关)", zap.Int("count", len(userKeys2)))
		}
	}()
}

func ClearAccountCache() {
	con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 清除 project-project 的 account 缓存
	keys, _ := dao.Rc.Keys(con, "/account.service.v1.AccountService/Account::*")
	if len(keys) > 0 {
		dao.Rc.Delete(con, keys)
		zap.L().Info("已删除所有账户缓存", zap.Int("count", len(keys)))
	}

	// 清除 project-user 的 FindMemInfoById 缓存（共享Redis）
	userKeys, _ := dao.Rc.Keys(con, "/login.service.v1.LoginService/FindMemInfoById::*")
	if len(userKeys) > 0 {
		dao.Rc.Delete(con, userKeys)
		zap.L().Info("已删除所有member信息缓存", zap.Int("count", len(userKeys)))
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		con2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel2()
		keys2, _ := dao.Rc.Keys(con2, "/account.service.v1.AccountService/Account::*")
		if len(keys2) > 0 {
			dao.Rc.Delete(con2, keys2)
			zap.L().Info("延迟删除所有账户缓存", zap.Int("count", len(keys2)))
		}
		userKeys2, _ := dao.Rc.Keys(con2, "/login.service.v1.LoginService/FindMemInfoById::*")
		if len(userKeys2) > 0 {
			dao.Rc.Delete(con2, userKeys2)
			zap.L().Info("延迟删除所有member信息缓存", zap.Int("count", len(userKeys2)))
		}
	}()
}
