package apidoc

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-openapi/spec"
)

// DocBundle 一次构建产物（Swagger2 + OpenAPI3 + ETag）
type DocBundle struct {
	Swagger     *spec.Swagger
	SwaggerJSON []byte
	OpenAPIJSON []byte
	ETag        string
	BuiltAt     time.Time
}

// DocCache 进程内文档缓存
type DocCache struct {
	mu     sync.Mutex
	bundle *DocBundle
}

func (c *DocCache) Get(ttlSec int, build func() (*DocBundle, error)) (*DocBundle, error) {
	if c == nil {
		return build()
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttlSec >= 0 && c.bundle != nil {
		if ttlSec == 0 || time.Since(c.bundle.BuiltAt) < time.Duration(ttlSec)*time.Second {
			return c.bundle, nil
		}
	}
	b, err := build()
	if err != nil {
		return nil, err
	}
	c.bundle = b
	return b, nil
}

func (c *DocCache) Invalidate() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bundle = nil
}

func NewDocBundle(swagger *spec.Swagger, openapi map[string]any) (*DocBundle, error) {
	swJSON, err := json.Marshal(swagger)
	if err != nil {
		return nil, err
	}
	oaJSON, err := json.Marshal(openapi)
	if err != nil {
		return nil, err
	}
	sum := sha1.Sum(swJSON)
	return &DocBundle{
		Swagger:     swagger,
		SwaggerJSON: swJSON,
		OpenAPIJSON: oaJSON,
		ETag:        `"` + hex.EncodeToString(sum[:8]) + `"`,
		BuiltAt:     time.Now(),
	}, nil
}
