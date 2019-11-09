package register

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Register 服务注册接口
type Register interface {
	Registe(service *ServiceInstance) (<-chan HeatbeatResonse, error)
	UnRegiste() error
}

// HeatbeatResonse 心态的返回
type HeatbeatResonse interface {
	TTL() int64
}

// ServiceType 服务类型
type ServiceType string

const (
	// API 提供API访问的服务
	API = ServiceType("api")
	// Worker 后台作业服务
	Worker = ServiceType("worker")
)

// ServiceInstance todo
type ServiceInstance struct {
	Region       string      `json:"region,omitempty"`
	InstanceName string      `json:"instanceName,omitempty"`
	ServiceName  string      `json:"serviceName,omitempty"`
	Type         ServiceType `json:"serviceType,omitempty"`
	Address      string      `json:"address,omitempty"`
	Version      string      `json:"version,omitempty"`
	GitBranch    string      `json:"gitBranch,omitempty"`
	GitCommit    string      `json:"gitCommit,omitempty"`
	BuildEnv     string      `json:"buildEnv,omitempty"`
	BuildAt      string      `json:"buildAt,omitempty"`
	Online       int64       `json:"online,omitempty"` // 毫秒时间戳

	Meta map[string]interface{} `json:"meta,omitempty"`

	Prefix   string        `json:"-"`
	Interval time.Duration `json:"-"`
	TTL      int64         `json:"-"`
}

// Validate 服务实例注册参数校验
func (s *ServiceInstance) Validate() error {
	if s.InstanceName == "" && s.ServiceName == "" || s.Type == "" {
		return errors.New("service instance name or service_name or type not config")
	}

	return nil
}

// Name 实例显示名称
func (s *ServiceInstance) Name() string {
	return fmt.Sprintf("%s.%s", s.ServiceName, s.InstanceName)
}

// MakeRegistryKey 构建etcd对应的key
func (s *ServiceInstance) MakeRegistryKey() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", s.Prefix, s.ServiceName, s.Type, s.Region, s.InstanceName)
}

// ParseInstanceKey 解析key中的服务信息
func ParseInstanceKey(key string) (serviceName, region, instanceName string, serviceType ServiceType, err error) {
	kl := strings.Split(key, "/")

	if len(kl) != 5 {
		err = errors.New("key format error, must like <prefix>/<service_name>/<service_type>/<namespace>/<instance_name>")
		return
	}

	serviceName, serviceType, region, instanceName = kl[1], ServiceType(kl[2]), kl[3], kl[4]
	return
}
