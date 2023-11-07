package application

import "net/http"

type RouterBuilder interface {
	Init() error
	GetRouter() http.Handler
}

var RouterBuilderStore = map[WEB_FRAMEWORK]RouterBuilder{
	WEB_FRAMEWORK_GO_RESTFUL: NewGoRestfulRouterBuilder(),
	WEB_FRAMEWORK_GIN:        NewGinRouterBuilder(),
}
