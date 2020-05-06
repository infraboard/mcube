package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/register"
)

type etcd struct {
	leaseID              clientv3.LeaseID
	client               *clientv3.Client
	requestTimeout       time.Duration
	headbeatResponseChan chan register.HeatbeatResonse

	isStopped     bool
	instanceKey   string
	stopInstance  chan struct{}
	keepAliveStop context.CancelFunc

	logger.Logger
}

type headbeatResponse struct {
	ttl int64
}

func (h *headbeatResponse) TTL() int64 {
	return h.ttl
}

// NewEtcdRegister 初始化一个基于etcd的实例注册器
func NewEtcdRegister(endpoints []string, username, password string, l logger.Logger) (register.Register, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(5) * time.Second,
		Username:    username,
		Password:    password,
	})

	if err != nil {
		return nil, fmt.Errorf("connect etcd error, %s", err)
	}

	etcdR := new(etcd)
	etcdR.Logger = l
	etcdR.client = client
	etcdR.stopInstance = make(chan struct{}, 1)
	etcdR.requestTimeout = time.Duration(5) * time.Second
	etcdR.headbeatResponseChan = make(chan register.HeatbeatResonse, 3)

	return etcdR, nil
}

// Register use to registe serice endpoint to etcd. when etcd is down,
// register can retry to registe util the etcd up.
//
// name is service name, use to discovery service address, eg. keyauth
// host and port is service endpoint, eg. 127.0.0.0:50000
// target is etcd addr, eg. 127.0.0.0:2379
// interval is service refresh interval, eg. 10s
// ttl is service ttl, eg. 15
// TODO: 判断服务是否已经被其他人注册了, 如果注册了 则需要更换名称才能注册
func (e *etcd) Registe(service *register.ServiceInstance) (<-chan register.HeatbeatResonse, error) {
	if err := service.Validate(); err != nil {
		return nil, err
	}

	// 注册服务的key
	sjson, err := json.Marshal(service)
	if err != nil {
		e.Errorf("marshal service object to json error, %s", err)
	}
	serviceValue := string(sjson)

	// 注册服务到etcd
	serviceKey := service.MakeRegistryKey()

	// 后台续约
	// 并没有直接使用KeepAlive, 因为存在偶然断开, 就不续约的情况
	ctx, cancel := context.WithCancel(context.Background())
	e.keepAliveStop = cancel
	go e.keepAlive(ctx, serviceKey, serviceValue, service.TTL, e.headbeatResponseChan)

	return e.headbeatResponseChan, nil
}

func (e *etcd) getLeaseID(ttl int64) (clientv3.LeaseID, error) {
	resp, err := e.client.Lease.Grant(context.TODO(), ttl)
	if err != nil {
		return 0, err
	}
	e.leaseID = resp.ID
	return e.leaseID, nil
}

func (e *etcd) addOnce(key, value string, ttl int64) error {
	// 获取leaseID
	resp, err := e.client.Lease.Grant(context.TODO(), ttl)
	if err != nil {
		return fmt.Errorf("get etcd lease id error, %s", err)
	}
	e.leaseID = resp.ID

	// 写入key
	if _, err := e.client.Put(context.Background(), key, value, clientv3.WithLease(e.leaseID)); err != nil {
		return fmt.Errorf("registe service '%s' with ttl to etcd3 failed: %s", key, err.Error())
	}
	e.instanceKey = key

	return nil
}

func (e *etcd) keepAlive(ctx context.Context, key, value string, ttl int64, respChan chan register.HeatbeatResonse) {
	// 初始化注册
	if err := e.addOnce(key, value, ttl); err != nil {
		e.Errorf("registry error, %s", err)
		return
	}
	e.Infof("服务实例(%s)注册成功", key)

	// 不停的续约
	interval := ttl / 5
	tk := time.NewTicker(time.Duration(interval) * time.Second)
	defer tk.Stop()

	for {
		select {
		case <-ctx.Done():
			e.Infof("keepalive goroutine exit")
			return
		case <-tk.C:
			Opctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			resp, err := e.client.Lease.KeepAliveOnce(Opctx, e.leaseID)
			if err != nil {
				if strings.Contains(err.Error(), "requested lease not found") {
					// 避免程序卡顿造成leaseID失效(比如mac 电脑休眠))
					if err := e.addOnce(key, value, ttl); err != nil {
						e.Errorf("refresh registry error, %s", err)
					} else {
						e.Warn("refresh registry success")
					}
				}
				e.Errorf("lease keep alive error, %s", err)
			} else {
				respChan <- &headbeatResponse{ttl: resp.TTL}
			}
		}
	}
}

// UnRegiste delete registered service from etcd, if etcd is down
// unregister while timeout.
func (e *etcd) UnRegiste() error {
	if e.isStopped {
		return errors.New("the instance has unregisted")
	}

	// delete instance key
	e.stopInstance <- struct{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if resp, err := e.client.Delete(ctx, e.instanceKey); err != nil {
		e.Warnf("unregister '%s' failed: connect to etcd server timeout, %s", e.instanceKey, err.Error())
	} else {
		if resp.Deleted == 0 {
			e.Infof("unregister '%s' failed, the key not exist", e.instanceKey)
		} else {
			e.Infof("服务实例(%s)注销成功", e.instanceKey)
		}
	}

	// revoke lease
	_, err := e.client.Lease.Revoke(context.TODO(), e.leaseID)
	if err != nil {
		e.Warnf("revoke lease error, %s", err)
		return err
	}
	e.isStopped = true

	// 停止续约心态
	e.keepAliveStop()

	return nil
}
