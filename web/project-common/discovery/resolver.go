package discovery

import (
	"context"
	"net/url"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

const (
	schema = "etcd"
)

// Resolver for grpc client
type Resolver struct {
	schema      string   // 解析器的 schema，例如 "etcd"
	EtcdAddrs   []string // etcd 服务器地址列表
	DialTimeout int      // 连接超时时间，单位秒

	closeCh      chan struct{}      // 关闭通道，用于通知关闭解析器
	watchCh      clientv3.WatchChan // 监听通道，用于接收 etcd 事件
	cli          *clientv3.Client   // etcd 客户端
	keyPrifix    string             // 服务注册路径前缀
	srvAddrsList []resolver.Address // 服务地址列表

	cc     resolver.ClientConn // gRPC 客户端连接
	logger *zap.Logger         // 日志记录器
}

// NewResolver create a new resolver.Builder base on etcd
// NewResolver 创建一个基于 etcd 的解析器
func NewResolver(etcdAddrs []string, logger *zap.Logger) *Resolver {
	return &Resolver{
		schema:      schema,
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Scheme returns the scheme supported by this resolver.
// Scheme 返回解析器支持的 scheme，例如 "etcd"
func (r *Resolver) Scheme() string {
	return r.schema
}

// Build creates a new resolver.Resolver for the given target
// Build 创建一个新的解析器实例
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	// 使用新版本 gRPC API 解析 target
	// 格式：etcd://authority/serviceName/version
	var serviceName, version string
	serviceName, version = parseTargetURL(target.URL)
	// 构建服务注册路径前缀
	r.keyPrifix = BuildPrefix(Server{Name: serviceName, Version: version})
	// 启动解析器
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// parseTargetURL 解析新版 gRPC target URL 格式
// 格式：etcd://authority/serviceName/version
func parseTargetURL(targetURL url.URL) (serviceName, version string) {
	// 移除路径开头的斜杠
	path := strings.TrimPrefix(targetURL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) >= 1 {
		serviceName = parts[0]
	}
	if len(parts) >= 2 {
		version = parts[1]
	}

	// 如果没有版本号，尝试从 authority 获取
	if version == "" && targetURL.Host != "" {
		version = targetURL.Host
	}

	return serviceName, version
}

// ResolveNow resolver.Resolver interface
// ResolveNow 触发解析器立即解析目标
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close 关闭解析器
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// start 启动解析器
func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	// 创建 etcd 客户端
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	// 注册解析器
	resolver.Register(r)
	// 创建关闭通道
	r.closeCh = make(chan struct{})
	// 同步获取所有地址信息
	if err = r.sync(); err != nil {
		return nil, err
	}
	// 启动监听通道
	go r.watch()

	return r.closeCh, nil
}

// watch update events
// watch 监听 etcd 事件
func (r *Resolver) watch() {
	// 创建定时器，每分钟同步一次地址信息
	ticker := time.NewTicker(time.Minute)
	// 创建监听通道，监听指定路径前缀的事件
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrifix, clientv3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", zap.Error(err))
			}
		}
	}
}

// update
// update 更新服务地址列表
func (r *Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		var info Server
		var err error

		switch ev.Type {
		// 处理 PUT 事件，更新服务地址
		case mvccpb.PUT:
			info, err = ParseValue(ev.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
			if !Exist(r.srvAddrsList, addr) {
				r.srvAddrsList = append(r.srvAddrsList, addr)
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
			}
			// 更新服务权重
		case mvccpb.DELETE:
			info, err = SplitPath(string(ev.Kv.Key))
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr}
			if s, ok := Remove(r.srvAddrsList, addr); ok {
				r.srvAddrsList = s
				r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
			}
		}
	}
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 同步获取 etcd 中所有服务地址
	res, err := r.cli.Get(ctx, r.keyPrifix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	r.srvAddrsList = []resolver.Address{}

	for _, v := range res.Kvs {
		info, err := ParseValue(v.Value)
		if err != nil {
			continue
		}
		addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
		r.srvAddrsList = append(r.srvAddrsList, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
	return nil
}
