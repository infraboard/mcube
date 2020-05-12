package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"github.com/infraboard/mcube/cache"
)

const (
	// AdapterName registry adapter name
	AdapterName = "redis"
)

var (
	// DefaultPrefix redis缓存的健的默认前缀
	DefaultPrefix = "osmp"
	// DefaultTTL 缓存默认的ttl时间
	DefaultTTL = 300 * time.Second
)

type adapter struct {
	address  string
	db       int
	password string
	prefix   string
	ttl      time.Duration

	client *redis.Client
	config map[string]string
}

// NewCache new an redis cache instance
func NewCache() cache.Cache {
	return &adapter{
		prefix: DefaultPrefix,
		ttl:    DefaultTTL,
	}
}

// config is like {"prefix":"collection key","address":"connection info","db":"0", "password": ""}
func (a *adapter) Config(config string) error {
	if err := a.validate(config); err != nil {
		return err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     a.address,
		Password: a.password,
		DB:       a.db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	a.client = client
	return nil
}

func (a *adapter) SetDefaultTTL(ttl time.Duration) {
	a.ttl = ttl
}

func (a *adapter) PutWithTTL(key string, val interface{}, ttl time.Duration) error {
	if a.client == nil {
		return errors.New("redis adapter not config or config failed")
	}

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if err := a.client.Set(key, b, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (a *adapter) Put(key string, val interface{}) error {
	if a.client == nil {
		return errors.New("redis adapter not config or config failed")
	}

	return a.PutWithTTL(key, val, a.ttl)
}

func (a *adapter) Get(key string, val interface{}) error {
	if a.client == nil {
		return errors.New("redis adapter not config or config failed")
	}

	v, err := a.client.Get(key).Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(v, val); err != nil {
		return err
	}

	return nil
}

func (a *adapter) Delete(key string) error {
	if a.client == nil {
		return errors.New("redis adapter not config or config failed")
	}

	if err := a.client.Del(key).Err(); err != nil {
		return err
	}

	return nil
}

func (a *adapter) IsExist(key string) bool {
	if ret := a.client.Exists(key).Val(); ret == 1 {
		return true
	}

	return false
}

func (a *adapter) ClearAll() error {
	return a.client.FlushDB().Err()
}

func (a *adapter) Decr(key string) error {
	return a.client.Decr(key).Err()
}

func (a *adapter) Incr(key string) error {
	return a.client.Incr(key).Err()
}

func (a *adapter) Close() error {
	if a.client == nil {
		return errors.New("redis adapter not config or config failed")
	}

	if err := a.client.Close(); err != nil {
		return err
	}

	a.client = nil
	return nil
}

func (a *adapter) validate(config string) error {
	var cf map[string]string
	if err := json.Unmarshal([]byte(config), &cf); err != nil {
		return fmt.Errorf("ummarshal redis adapter config error, %s", err)
	}

	v, ok := cf["address"]
	if !ok {
		return fmt.Errorf("redis adapter address not config")
	}
	a.address = v

	if v, ok := cf["db"]; !ok {
		a.db = 0
	} else {
		a.db, _ = strconv.Atoi(v)
	}

	if v, ok := cf["password"]; ok {
		a.password = v
	}

	return nil
}

func init() {
	cache.Register(AdapterName, NewCache)
}
