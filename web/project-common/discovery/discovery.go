package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// Register for grpc server
type Register struct {
	EtcdAddrs   []string // etcd数组地址
	DialTimeout int      // 连接超时时间

	closeCh     chan struct{}                           // 关闭返回的通道
	leasesID    clientv3.LeaseID                        // 租约id
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 租约保持通道

	srvInfo Server           // 服务信息
	srvTTL  int64            // 服务过期时间
	cli     *clientv3.Client // etcd客户端
	logger  *zap.Logger      // 日志
}

// NewRegister create a register base on etcd
// 创建etcd实例
func NewRegister(etcdAddrs []string, logger *zap.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register a service to etcd
// 注册服务到etcd
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	// 检查服务地址是否合法
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}
	log.Printf("[Discovery] Connecting to etcd: %v", r.EtcdAddrs)
	// 创建etcd客户端
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		log.Printf("[Discovery] Failed to connect to etcd: %v", err)
		return nil, err
	}
	log.Printf("[Discovery] Connected to etcd successfully")
	// 设置服务信息和过期时间
	r.srvInfo = srvInfo
	r.srvTTL = ttl
	// 注册服务
	if err = r.register(); err != nil {
		log.Printf("[Discovery] Failed to register service: %v", err)
		return nil, err
	}
	log.Printf("[Discovery] Service registered successfully: %s at %s", srvInfo.Name, srvInfo.Addr)
	// 创建关闭通道
	r.closeCh = make(chan struct{})
	// 启动租约保持
	go r.keepAlive()
	// 返回关闭通道
	return r.closeCh, nil
}

// Stop stop register
// 停止注册服务
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 注册节点
func (r *Register) register() error {
	// 创建租约上下文
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()
	// 创建租约
	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	// 设置租约id
	r.leasesID = leaseResp.ID
	// 创建租约保持通道
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	// 删除节点
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive
// 保持租约
func (r *Register) keepAlive() {
	// 创建租约保持定时器
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		// 关闭通道
		case <-r.closeCh:
			// 注销节点
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
			// 撤销租约
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		// 租约保持通道
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		// 租约保持定时器
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
}

// UpdateHandler return http handler
// UpdateHandler 更新服务权重
func (r *Register) UpdateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("update server weight success"))
	})
}

// GetServerInfo 获取服务信息
func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}
