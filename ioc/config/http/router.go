package http

import "net/http"

type RouterHook func(http.Handler)

type RouterBuilder interface {
	BeforeLoadHooks(...RouterHook)
	AfterLoadHooks(...RouterHook)
	Build() (http.Handler, error)
}

func NewRouterBuilderHooks() *RouterBuilderHooks {
	return &RouterBuilderHooks{
		Before: []RouterHook{},
		After:  []RouterHook{},
	}
}

type RouterBuilderHooks struct {
	Before []RouterHook
	After  []RouterHook
}

func (b *RouterBuilderHooks) BeforeLoadHooks(hooks ...RouterHook) {
	b.Before = append(b.Before, hooks...)
}

func (b *RouterBuilderHooks) AfterLoadHooks(hooks ...RouterHook) {
	b.After = append(b.After, hooks...)
}
