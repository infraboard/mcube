package ioc

import (
	"path"
	"strings"
)

const (
	API_NAMESPACE = "apis"
)

type ApiCustomPathExtender interface {
	ApiPath() string
}

// 用于托管RestApi对象的Ioc空间, 最后初始化
func Api() StoreUser {
	return DefaultStore.Namespace(API_NAMESPACE)
}

func ApiPathPrefix(pathPrefix string, obj Object) string {
	customPrefix := obj.Meta().CustomPathPrefix
	if customPrefix != "" {
		return customPrefix
	}

	if extender, ok := obj.(ApiCustomPathExtender); ok {
		return extender.ApiPath()
	}

	pathPrefix = strings.TrimSuffix(pathPrefix, "/")
	return path.Join(pathPrefix, obj.Version(), obj.Name())
}
