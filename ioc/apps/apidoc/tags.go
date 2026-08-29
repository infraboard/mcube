package apidoc

import (
	"sort"
	"strings"
	"sync"

	"github.com/go-openapi/spec"
)

var (
	tagDescMu sync.RWMutex
	tagDescs  = map[string]string{}
	postBuild []func(*spec.Swagger)
)

// RegisterTagDescription 注册某一 OpenAPI Tag 的 Markdown 说明（模块导读）。
// 在服务 Init/init 阶段调用；BuildSwagger PostBuild 时写入 swo.Tags。
func RegisterTagDescription(name, markdown string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	tagDescMu.Lock()
	defer tagDescMu.Unlock()
	tagDescs[name] = strings.TrimSpace(markdown)
}

// RegisterTagDescriptions 批量注册。
func RegisterTagDescriptions(m map[string]string) {
	for k, v := range m {
		RegisterTagDescription(k, v)
	}
}

// RegisterPostBuild 追加自定义 PostBuild（在 Sanitize / Enrich / Tag 注入之后执行）。
func RegisterPostBuild(fn func(*spec.Swagger)) {
	if fn == nil {
		return
	}
	tagDescMu.Lock()
	defer tagDescMu.Unlock()
	postBuild = append(postBuild, fn)
}

// ApplyTagDescriptions 把已注册的 Tag 说明与路径上出现的 tag 合并进 swagger.Tags。
func ApplyTagDescriptions(swo *spec.Swagger) {
	if swo == nil {
		return
	}
	used := collectTagNames(swo)
	tagDescMu.RLock()
	defer tagDescMu.RUnlock()

	byName := map[string]spec.Tag{}
	for _, t := range swo.Tags {
		byName[t.Name] = t
	}
	for name := range used {
		t, ok := byName[name]
		if !ok {
			t = spec.Tag{TagProps: spec.TagProps{Name: name}}
		}
		if desc, ok := tagDescs[name]; ok && desc != "" {
			t.Description = desc
		}
		byName[name] = t
	}
	// 仅注册了说明、尚无路由的 tag 也挂上（便于文档预览）
	for name, desc := range tagDescs {
		if desc == "" {
			continue
		}
		if _, ok := byName[name]; ok {
			continue
		}
		byName[name] = spec.Tag{TagProps: spec.TagProps{Name: name, Description: desc}}
	}

	names := make([]string, 0, len(byName))
	for n := range byName {
		names = append(names, n)
	}
	sort.Strings(names)
	out := make([]spec.Tag, 0, len(names))
	for _, n := range names {
		out = append(out, byName[n])
	}
	swo.Tags = out
}

// RunPostBuildHooks 执行业务侧注册的额外 PostBuild。
func RunPostBuildHooks(swo *spec.Swagger) {
	tagDescMu.RLock()
	fns := append([]func(*spec.Swagger){}, postBuild...)
	tagDescMu.RUnlock()
	for _, fn := range fns {
		fn(swo)
	}
}

func collectTagNames(swo *spec.Swagger) map[string]struct{} {
	used := map[string]struct{}{}
	if swo.Paths == nil {
		return used
	}
	addOp := func(op *spec.Operation) {
		if op == nil {
			return
		}
		for _, t := range op.Tags {
			t = strings.TrimSpace(t)
			if t != "" {
				used[t] = struct{}{}
			}
		}
	}
	for _, item := range swo.Paths.Paths {
		addOp(item.Get)
		addOp(item.Post)
		addOp(item.Put)
		addOp(item.Patch)
		addOp(item.Delete)
		addOp(item.Head)
		addOp(item.Options)
	}
	return used
}
