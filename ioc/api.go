package ioc

import (
	"path"
	"strings"
)

const (
	API_NAMESPACE = "apis"
)

// 用于托管RestApi对象的Ioc空间, 最后初始化
func Api() StoreUser {
	return DefaultStore.Namespace(API_NAMESPACE)
}

func ApiPathPrefix(pathPrefix string, obj Object) string {
	customPrefix := obj.Meta().CustomPathPrefix
	if customPrefix != "" {
		return customPrefix
	}

	pathPrefix = strings.TrimSuffix(pathPrefix, "/")
	return path.Join(pathPrefix, obj.Version(), obj.Name())
}
