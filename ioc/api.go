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
	return store.Namespace(API_NAMESPACE)
}

func ApiPathPrefix(pathPrefix string, obj Object) string {
	customPrefix := obj.Meta().CustomPathPrefix
	if customPrefix != "" {
		return customPrefix
	}

	pathPrefix = strings.TrimSuffix(pathPrefix, "/")
	return path.Join(pathPrefix, obj.Version(), obj.Name())
}

// LoadGinApi 装载所有的gin app
// func LoadGinApi(pathPrefix string, root gin.IRouter) {
// 	objects := store.Namespace(API_NAMESPACE)
// 	objects.ForEach(func(w *ObjectWrapper) {
// 		api, ok := w.Value.(GinApiObject)
// 		if !ok {
// 			return
// 		}
// 		api.Registry(root.Group(ApiPathPrefix(pathPrefix, api)))
// 	})
// }

// // LoadHttpApp 装载所有的http app
// func LoadGoRestfulApi(pathPrefix string, root *restful.Container) {
// 	objects := store.Namespace(API_NAMESPACE)
// 	objects.ForEach(func(w *ObjectWrapper) {
// 		api, ok := w.Value.(GoRestfulApiObject)
// 		if !ok {
// 			return
// 		}

// 		ws := new(restful.WebService)
// 		ws.
// 			Path(ApiPathPrefix(pathPrefix, api)).
// 			Consumes(restful.MIME_JSON, form.MIME_POST_FORM, form.MIME_MULTIPART_FORM, yaml.MIME_YAML, yamlk8s.MIME_YAML).
// 			Produces(restful.MIME_JSON, yaml.MIME_YAML, yamlk8s.MIME_YAML)
// 		api.Registry(ws)
// 		root.Add(ws)
// 	})
// }
