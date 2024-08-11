package cache

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

type ObjectFinder func(ctx context.Context, objectId string) (any, error)

func NewGetter(ctx context.Context, f ObjectFinder) *Getter {
	return &Getter{
		ctx: ctx,
		f:   f,
		l:   log.Sub("cache"),
	}
}

type Getter struct {
	ctx          context.Context
	f            ObjectFinder
	ttl          int64
	namespace    string
	resourceType string
	l            *zerolog.Logger
}

func (g *Getter) GetKey(key string) string {
	keys := []string{}
	if g.namespace != "" {
		keys = append(keys, g.namespace)
	}
	if g.resourceType != "" {
		keys = append(keys, g.resourceType)
	}
	keys = append(keys, key)
	return strings.Join(keys, ".")
}

func (g *Getter) WithTTL(ttl int64) *Getter {
	g.ttl = ttl
	return g
}

func (g *Getter) WithNamespace(namespace string) *Getter {
	g.namespace = namespace
	return g
}

func (g *Getter) WithResourceType(resourceType string) *Getter {
	g.resourceType = resourceType
	return g
}

func (g *Getter) Get(key string, value any) error {
	err := C().Get(g.ctx, g.GetKey(key), value)
	if err == nil {
		g.l.Info().Msgf("get object %s from cache", g.GetKey(key))
		return nil
	}

	if err == ErrKeyNotFound {
		goto DO_CACHE
	}

	return err

DO_CACHE:
	v, err := g.f(g.ctx, key)
	if err != nil {
		return err
	}

	if err := C().Set(g.ctx, g.GetKey(key), v, WithExpiration(g.ttl)); err != nil {
		g.l.Warn().Msgf("set cache error, %s", err)
	} else {
		g.l.Info().Msgf("set cache object %s ttl: %d second", g.GetKey(key), g.ttl)
	}

	//
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.TypeOf(v).Kind() == reflect.Ptr {
			reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v).Elem())
		} else {
			reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v))
		}
	} else {
		return fmt.Errorf("value must be ptr")
	}

	return nil
}
