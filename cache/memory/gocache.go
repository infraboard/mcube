package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bluele/gcache"

	"github.com/infraboard/mcube/cache"
)

const (
	// AdapterName registry adapter name
	AdapterName = "memory"
)

var (
	// DefaultTTL 缓存默认的ttl时间
	DefaultTTL = 300 * time.Second
	// DefaultSize 缓存容量
	DefaultSize = 100
)

type adapter struct {
	ttl  time.Duration
	size int
	gc   gcache.Cache

	config map[string]string
}

// NewCache new an redis cache instance
func NewCache() cache.Cache {
	return &adapter{
		ttl:  DefaultTTL,
		size: DefaultSize,
	}
}

// {"ttl": "300", "size": "1000"}
func (a *adapter) Config(config string) error {
	if err := a.validate(config); err != nil {
		return err
	}

	gc := gcache.New(a.size).
		LRU().
		Expiration(a.ttl).
		Build()

	a.gc = gc
	return nil
}

func (a *adapter) SetDefaultTTL(ttl time.Duration) {
	a.ttl = ttl
}

func (a *adapter) Put(key string, val interface{}) error {
	return a.PutWithTTL(key, val, a.ttl)
}

func (a *adapter) PutWithTTL(key string, val interface{}, ttl time.Duration) error {
	if a.gc == nil {
		return errors.New("memcache adapter not config or config failed")
	}

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return a.gc.SetWithExpire(key, b, a.ttl)
}

func (a *adapter) Get(key string, val interface{}) error {
	if a.gc == nil {
		return errors.New("memcache adapter not config or config failed")
	}

	v, err := a.gc.Get(key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(v.([]byte), val); err != nil {
		return err
	}

	return nil

}

func (a *adapter) Delete(key string) error {
	if a.gc == nil {
		return errors.New("memcache adapter not config or config failed")
	}

	if !a.gc.Remove(key) {
		return errors.New("delete cache error")
	}

	return nil
}

func (a *adapter) IsExist(key string) bool {
	if a.gc == nil {
		return false
	}

	return a.gc.Has(key)
}

func (a *adapter) ClearAll() error {
	if a.gc == nil {
		return errors.New("memcache adapter not config or config failed")
	}

	a.gc.Purge()
	return nil
}

func (a *adapter) Incr(key string) error {
	return errors.New("not impliment")
}

func (a *adapter) Decr(key string) error {
	return errors.New("not impliment")
}

func (a *adapter) Close() error {
	return nil
}

func (a *adapter) validate(config string) error {
	var cf map[string]string
	if err := json.Unmarshal([]byte(config), &cf); err != nil {
		return fmt.Errorf("ummarshal redis adapter config error, %s", err)
	}

	if v, ok := cf["size"]; !ok {
		a.size = DefaultSize
	} else {
		a.size, _ = strconv.Atoi(v)
		if a.size == 0 {
			a.size = DefaultSize
		}
	}

	if v, ok := cf["ttl"]; !ok {
		a.ttl = DefaultTTL
	} else {
		ttlInt, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if ttlInt != 0 {
			a.ttl = time.Duration(ttlInt) * time.Second
		} else {
			a.ttl = DefaultTTL
		}
	}

	return nil
}

func init() {
	cache.Register(AdapterName, NewCache)
}
